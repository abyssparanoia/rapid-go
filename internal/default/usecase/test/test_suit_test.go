package test

import (
	"database/sql"
	"testing"

	"github.com/abyssparanoia/rapid-go/internal/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testDB *sql.DB

var _ = BeforeSuite(func() {
	testDB = testutil.PrepareDefaultTestDB()
})

var _ = AfterSuite(func() {
	testDB.Close()
})

var _ = AfterEach(func() {
	testutil.TruncatDefaultTestDB(testDB)
})

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Usecase Test Suite")
}
