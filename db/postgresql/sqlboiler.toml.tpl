add-global-variants = true
add-panic-variants = true
no-tests = true
no-auto-timestamps = true
no-hooks = true
output = "internal/infrastructure/postgresql/internal/dbmodel"
pkgname = "dbmodel"
wipe = true
templates = [
  "{{GOPATH}}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.19.1/templates/main",
  "{{GOPATH}}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.19.1/templates/test",
  "db/postgresql/templates",
]

[struct-tag-cases]
toml = "snake"
yaml = "snake"
json = "snake"
boil = "snake"

[psql]
  dbname = "maindb"
  host   = "localhost"
  port   = 5432
  user   = "postgres"
  pass   = "postgres"
  sslmode = "disable"
  blacklist = [
   "goose_db_version",
   "content_types", 
   "asset_types", 
   "staff_roles"
  ]

[[types]]
 [types.match]
    nullable = false
    db_type = "date"

 [types.replace]
    type = "custom_types.Date"

 [types.imports]
    third_party = ['"github.com/abyssparanoia/rapid-go/db/postgresql/custom_types"']

[[types]]
 [types.match]
    nullable = true
    db_type = "date"

 [types.replace]
    type = "custom_types.NullDate"

 [types.imports]
    third_party = ['"github.com/abyssparanoia/rapid-go/db/postgresql/custom_types"']
