add-global-variants = true
add-panic-variants = true
no-driver-templates = true
no-tests = true
output = "internal/infrastructure/database/internal/dbmodel"
pkgname = "dbmodel"
wipe = true
templates = [
  "{{GOPATH}}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.14.2/templates/main",
  "{{GOPATH}}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.14.2/templates/test",
  "db/main/templates",
]

[mysql]
  dbname  = "maindb"
  host    = "localhost"
  port    = 3306
  user    = "root"
  pass    = "password"
  sslmode = "false"