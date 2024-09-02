add-global-variants = true
add-panic-variants = true
no-tests = true
no-auto-timestamps = true
no-hooks = true
output = "internal/infrastructure/database/internal/dbmodel"
pkgname = "dbmodel"
wipe = true
templates = [
  "{{GOPATH}}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.16.2/templates/main",
  "{{GOPATH}}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.16.2/templates/test",
  "db/main/templates",
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