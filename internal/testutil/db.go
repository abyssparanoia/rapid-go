package testutil

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/abyssparanoia/rapid-go/internal/pkg/gluemysql"
	env "github.com/caarlos0/env/v6"
	"github.com/go-testfixtures/testfixtures/v3"
)

var (
	fixtures *testfixtures.Loader
	envObj   *environment
)

func init() {
	envObj = &environment{}
	if err := env.Parse(envObj); err != nil {
		panic(err)
	}
}

// TruncatDefaultTestDB :
func TruncatDefaultTestDB(db *sql.DB) {
	sqlPath := filepath.Join(GetTestDataPath(), "truncate/default.sql")
	truncatePsqlTestDatabase(db, sqlPath)
}

// SetupDefaultTestDB :
func SetupDefaultTestDB() func() {
	db := PrepareDefaultTestDB()
	return func() {
		defer db.Close()
		TruncatDefaultTestDB(db)
	}
}

func truncatePsqlTestDatabase(db *sql.DB, path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(string(bytes))
	if err != nil {
		panic(err)
	}
}

// PrepareDefaultTestDB :
func PrepareDefaultTestDB() *sql.DB {
	return prepareDatabaseTest(
		envObj.DefaultDBHost,
		envObj.DefaultDBUser,
		envObj.DefaultDBDatabase,
		envObj.DefaultDBPassword,
	)
}

func prepareDatabaseTest(
	host,
	user,
	dbname,
	pass string,
) *sql.DB {
	envObj := &environment{}
	if err := env.Parse(envObj); err != nil {
		panic(err)
	}

	dbConn := gluemysql.NewClient(host, user, pass, dbname)

	return dbConn.DB
}

// GetTestDataPath :
func GetTestDataPath() string {
	_, b, _, _ := runtime.Caller(0)
	testutilPath := filepath.Dir(b)
	projectRoot := filepath.Dir(testutilPath)
	return filepath.Join(projectRoot, "testdata")
}
