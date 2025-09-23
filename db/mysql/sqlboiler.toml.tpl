add-global-variants = true
add-panic-variants = true
no-tests = true
no-auto-timestamps = true
no-hooks = true
output = "internal/infrastructure/mysql/internal/dbmodel"
pkgname = "dbmodel"
wipe = true
templates = [
  "{{GOPATH}}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.19.1/templates/main",
  "{{GOPATH}}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.19.1/templates/test",
  "db/mysql/templates",
]

[struct-tag-cases]
toml = "snake"
yaml = "snake"
json = "snake"
boil = "snake"

[mysql]
  dbname  = "maindb"
  host    = "localhost"
  port    = 3306
  user    = "root"
  pass    = "password"
  sslmode = "false"

[[types]]
 [types.match]
    nullable = false
    db_type = "date"

 [types.replace]
    type = "custom_types.Date"

 [types.imports]
    third_party = ['"github.com/abyssparanoia/rapid-go/db/mysql/custom_types"']

[[types]]
 [types.match]
    nullable = true
    db_type = "date"

 [types.replace]
    type = "custom_types.NullDate"

 [types.imports]
    third_party = ['"github.com/abyssparanoia/rapid-go/db/mysql/custom_types"']
