package databases

import (
	"testing"

	"github.com/nambelaas/payroll-system-go/pkg"
	"github.com/stretchr/testify/assert"
)

func TestConnectTestDb(t *testing.T) {
	db := pkg.ConnectTestDB()
	assert.NotNil(t, db)

	sqlDb, err := db.DB()
	assert.NoError(t, err)
	err = sqlDb.Ping()
	assert.NoError(t, err)
}
