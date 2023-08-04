{{- $short := (shortname .Name "err" "res" "sqlstr" "db" "YOLog") -}}
{{- $table := (.Table.TableName) -}}
// {{ .Name }} represents a row from '{{ $table }}'.
type {{ .Name }} struct {
{{- range .Fields }}
{{- if eq (.Col.DataType) (.Col.ColumnName) }}
	{{ .Name }} string `spanner:"{{ .Col.ColumnName }}" json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }} enum
{{- else if .CustomType }}
	{{ .Name }} {{ retype .CustomType }} `spanner:"{{ .Col.ColumnName }}" json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }}
{{- else }}
	{{ .Name }} {{ .Type }} `spanner:"{{ .Col.ColumnName }}" json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }}
{{- end }}
{{- end }}
}

type {{ .Name }}Slice []*{{ .Name }}

func {{ .Name }}TableName() string {
	return "{{ .Table.TableName }}"
}

{{ if .PrimaryKey }}
func {{ .Name }}PrimaryKeys() []string {
     return []string{
{{- range .PrimaryKeyFields }}
		"{{ colname .Col }}",
{{- end }}
	}
}
{{- end }}

func {{ .Name }}Columns() []string {
	return []string{
{{- range .Fields }}
		"{{ colname .Col }}",
{{- end }}
	}
}

func {{ .Name }}WritableColumns() []string {
	return []string{
{{- range .Fields }}
	{{- if not .Col.IsGenerated }}
		"{{ colname .Col }}",
	{{- end }}
{{- end }}
	}
}

func ({{ $short }} *{{ .Name }}) columnsToPtrs(cols []string, customPtrs map[string]interface{}) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		if val, ok := customPtrs[col]; ok {
			ret = append(ret, val)
			continue
		}

		switch col {
{{- range .Fields }}
		case "{{ colname .Col }}":
			ret = append(ret, &{{ $short }}.{{ .Name }})
{{- end }}
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}
	return ret, nil
}

func ({{ $short }} *{{ .Name }}) columnsToValues(cols []string) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		switch col {
{{- range .Fields }}
		case "{{ colname .Col }}":
			{{- if .CustomType }}
			ret = append(ret, {{ .Type }}({{ $short }}.{{ .Name }}))
			{{- else }}
			ret = append(ret, {{ $short }}.{{ .Name }})
			{{- end }}
{{- end }}
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}

	return ret, nil
}

// new{{ .Name }}_Decoder returns a decoder which reads a row from *spanner.Row
// into {{ .Name }}. The decoder is not goroutine-safe. Don't use it concurrently.
func new{{ .Name }}_Decoder(cols []string) func(*spanner.Row) (*{{ .Name }}, error) {
	{{- range .Fields }}
		{{- if .CustomType }}
			var {{ customtypeparam .Name }} {{ .Type }}
		{{- end }}
	{{- end }}
	customPtrs := map[string]interface{}{
		{{- range .Fields }}
			{{- if .CustomType }}
				"{{ colname .Col }}": &{{ customtypeparam .Name }},
			{{- end }}
	{{- end }}
	}

	return func(row *spanner.Row) (*{{ .Name }}, error) {
        var {{ $short }} {{ .Name }}
        ptrs, err := {{ $short }}.columnsToPtrs(cols, customPtrs)
        if err != nil {
            return nil, err
        }

        if err := row.Columns(ptrs...); err != nil {
            return nil, err
        }
        {{- range .Fields }}
            {{- if .CustomType }}
                {{ $short }}.{{ .Name }} = {{ retype .CustomType }}({{ customtypeparam .Name }})
            {{- end }}
        {{- end }}


		return &{{ $short }}, nil
	}
}

func ({{ $short }} *{{ .Name }}) Insert(ctx context.Context) error {
	params := make(map[string]interface{})
	{{- range .Fields }}
		params[fmt.Sprintf("{{ .Name }}")] = {{ $short }}.{{ .Name }}
	{{- end }}

	values := []string{
		{{- range .Fields }}
			fmt.Sprintf("@{{ .Name }}"),
		{{- end }}
	}
	rowValue := fmt.Sprintf("(%s)", strings.Join(values, ","))

	sql := fmt.Sprintf(`
    INSERT INTO {{ $table }}
        ({{ colnames .Fields }})
    VALUES
        %s
    `, rowValue)

	err := GetSpannerTransaction(ctx).ExecContext(ctx, sql, params)
	if err != nil {
		return err
	}

	return nil
}

func ({{ $short }}Slice {{ .Name }}Slice) InsertAll(ctx context.Context) error {
	if len({{ $short }}Slice) == 0 {
        return nil
    }

    params := make(map[string]interface{})
    valueStmts := make([]string, 0, len({{ $short }}Slice))
    for i, m := range {{ $short }}Slice {
        {{- range .Fields }}
			params[fmt.Sprintf("{{ .Name }}%d", i)] = m.{{ .Name }}
        {{- end }}


        values := []string{
            {{- range .Fields }}
            fmt.Sprintf("@{{ .Name }}%d", i),
            {{- end }}
        }
        rowValue := fmt.Sprintf("(%s)", strings.Join(values, ","))
        valueStmts = append(valueStmts, rowValue)
    }

    sql := fmt.Sprintf(`
    INSERT INTO {{ $table }}
        ({{ colnames .Fields }})
    VALUES
        %s
    `, strings.Join(valueStmts, ","))

	err := GetSpannerTransaction(ctx).ExecContext(ctx, sql, params)
	if err != nil {
		return err
	}

	return nil
}

{{ if ne (fieldnames .Fields $short .PrimaryKeyFields) "" }}
// Update the {{ .Name }}
func ({{ $short }} *{{ .Name }}) Update(ctx context.Context) error {
	{{ $primaryKeyFields := .PrimaryKeyFields -}}
	updateColumns := []string{}
	{{ range $i, $field := .Fields -}}
		{{- $include := false }}
		{{- range $j,$p := $primaryKeyFields -}}
			{{- if eq $field.Name $p.Name }}
				{{ $include = true }}
			{{- end }}
		{{- end }}
		{{- if eq $include false }}
			updateColumns = append(updateColumns, "{{$field.Name}} = @param_{{$field.Name}}")
		{{- end }}
	{{- end }}

	sql := fmt.Sprintf(`
	UPDATE {{ $table }}
	SET
		%s
    WHERE
        {{- range $i,$v := .PrimaryKeyFields }}
          {{- if eq $i 0 }}
            {{ $v.Name }} = @update_params{{ $i }}
          {{- else }}
            AND {{ $v.Name }} = @update_params{{ $i }}
          {{- end }}
        {{- end }}
	`, strings.Join(updateColumns, ","))

	setParams := map[string]interface{}{
	{{- range $i, $field := .Fields -}}
		{{ $include := false }}
		{{- range $j,$p := $primaryKeyFields }}
			{{- if eq $field.Name $p.Name }}
				{{ $include = true }}
			{{- end }}
		{{- end }}
		{{- if eq $include false }}
		"param_{{$field.Name}}": {{ $short }}.{{.Name}},
		{{- end }}
	{{- end }}
	}

	whereParams := map[string]interface{}{
	{{- range $i, $field := .PrimaryKeyFields }}
		"update_params{{$i}}": {{ $short }}.{{ $field.Name }},
	{{- end }}
	}

	params := make(map[string]interface{})
	for key, value := range setParams {
	    params[key] = value
	}
	for key, value := range whereParams {
        params[key] = value
    }

	err := GetSpannerTransaction(ctx).ExecContext(ctx, sql, params)
	if err != nil {
		return err
	}

	return nil
}
{{ end }}

// Delete the {{ .Name }} from the database.
func ({{ $short }} *{{ .Name }}) Delete(ctx context.Context) error {
	sql := fmt.Sprintf(`
        	DELETE FROM {{ $table }}
        	WHERE
        	    %s
        	`,
        	fmt.Sprintf("({{ colnamesquery .PrimaryKeyFields " AND " }})"),
		)
	
	params := map[string]interface{}{
	{{- range $i, $field := .PrimaryKeyFields }}
		"param{{$i}}": {{ $short }}.{{ $field.Name }},
	{{- end }}
	}

	if err := GetSpannerTransaction(ctx).ExecContext(ctx, sql, params); err != nil {
		return err
	}
	return nil
}
