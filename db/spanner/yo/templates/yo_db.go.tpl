// YODB is the common interface for database operations.
type YODB interface {
	YORODB
}

// YORODB is the common interface for database operations.
type YORODB interface {
	ReadRow(ctx context.Context, table string, key spanner.Key, columns []string) (*spanner.Row, error)
	Read(ctx context.Context, table string, keys spanner.KeySet, columns []string) *spanner.RowIterator
	ReadUsingIndex(ctx context.Context, table, index string, keys spanner.KeySet, columns []string) (ri *spanner.RowIterator)
	Query(ctx context.Context, statement spanner.Statement) *spanner.RowIterator
}

// YOLog provides the log func used by generated queries.
var YOLog = func(context.Context, string, ...interface{}) { }

func newError(method, table string, err error) error {
	code := spanner.ErrCode(err)
	return newErrorWithCode(code, method, table, err)
}

func newErrorWithCode(code codes.Code, method, table string, err error) error {
	return &yoError{
		method: method,
		table:  table,
		err:    err,
		code:   code,
	}
}

type yoError struct {
	err    error
	method string
	table  string
	code   codes.Code
}

func (e yoError) Error() string {
	return fmt.Sprintf("yo error in %s(%s): %v", e.method, e.table, e.err)
}

func (e yoError) Unwrap() error {
	return e.err
}

func (e yoError) DBTableName() string {
	return e.table
}

// GRPCStatus implements a conversion to a gRPC status using `status.Convert(error)`.
// If the error is originated from the Spanner library, this returns a gRPC status of
// the original error. It may contain details of the status such as RetryInfo.
func (e yoError) GRPCStatus() *status.Status {
	var ae *apierror.APIError
	if errors.As(e.err, &ae) {
		return status.Convert(ae)
	}

	return status.New(e.code, e.Error())
}

func (e yoError) Timeout() bool { return e.code == codes.DeadlineExceeded }
func (e yoError) Temporary() bool { return e.code == codes.DeadlineExceeded }
func (e yoError) NotFound() bool { return e.code == codes.NotFound }


var (
	ctxTKey = struct{}{}
)

func GetSpannerTransaction(ctx context.Context) *SpannerTransaction {
	if txn, ok := ctx.Value(&ctxTKey).(*SpannerTransaction); ok {
		return txn
	}
	panic("error GetSpannerTransaction: transaction not found")
}

type SpannerTransactable struct {
	db *spanner.Client
}

// read only transaction
func (t *SpannerTransactable) ROTx(ctx context.Context, fn func(ctx context.Context) error) error {

	err := (func() error {
		ro := t.db.ReadOnlyTransaction()
		defer func() { ro.Close() }()
		txn := &SpannerTransaction{ro: ro}
		ctxWithTx := context.WithValue(ctx, &ctxTKey, txn)
		if err := fn(ctxWithTx); err != nil {
			return err
		}
		return nil
	})()

	if err != nil {
		return err
	}

	return nil
}

// read and write transaction
func (t *SpannerTransactable) RWTx(ctx context.Context, fn func(ctx context.Context) error) error {
	_, err := t.db.ReadWriteTransaction(ctx, func(ctx context.Context, rw *spanner.ReadWriteTransaction) error {

		txn := &SpannerTransaction{rw: rw}
		ctxWithTx := context.WithValue(ctx, &ctxTKey, txn)
		// トランザクション内で処理を実行
		if err := fn(ctxWithTx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func NewTransactable(
	db *spanner.Client,
) *SpannerTransactable {
	return &SpannerTransactable{db}
}

type SpannerTransaction struct {
	ro *spanner.ReadOnlyTransaction
	rw *spanner.ReadWriteTransaction
}

// isRo はReadOnlyTransactionかどうかを返す
func (t *SpannerTransaction) isRo() bool {
	return t.ro != nil
}

// isRw はReadWriteTransactionかどうかを返す
func (t *SpannerTransaction) isRw() bool {
	return t.rw != nil
}

func (t *SpannerTransaction) QueryContext(ctx context.Context, query string, params map[string]interface{}) (*SpannerRows, error) {
	if t.isRo() {
		iter := t.ro.Query(ctx, spanner.Statement{
			SQL:    query,
			Params: params,
		})
		return &SpannerRows{iter: iter, query: query}, nil
	}
	if t.isRw() {
		iter := t.rw.Query(ctx, spanner.Statement{
			SQL:    query,
			Params: params,
		})
		return &SpannerRows{iter: iter, query: query}, nil
	}
	return nil, fmt.Errorf("error QueryContext: empty transaction")
}

func (t *SpannerTransaction) ExecContext(ctx context.Context, query string, params map[string]interface{}) error {
	if t.isRo() || !t.isRw() {
		return fmt.Errorf("error ExecContext is executed in read only transaction")
	}

	_, err := t.rw.Update(ctx, spanner.Statement{
		SQL:    query,
		Params: params,
	})
	if err != nil {
		return err
	}
	return nil
}

type SpannerRows struct {
	iter    *spanner.RowIterator
	query   string
	nextRow *spanner.Row
}

func (r *SpannerRows) Close() error {
	r.iter.Stop()
	return nil
}

func (r *SpannerRows) Next() (ok bool, err error) {
	row, err := r.iter.Next()
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	r.nextRow = row
	return true, nil
}

func (r *SpannerRows) Scan(ptrs ...interface{}) error {
	if r.nextRow == nil {
		return fmt.Errorf("SpannerRows.Scan(): next row is nil")
	}
	return r.nextRow.Columns(ptrs...)
}

func (r *SpannerRows) ToStruct(p interface{}) error {
	return r.nextRow.ToStruct(p)
}
