package repository

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/dto"
	"IhysBestowal/pkg/customLogger"
	"github.com/jmoiron/sqlx"
	"github.com/smartystreets/goconvey/convey"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testPostgres := map[string]string{
		"PSQL_USER":     "postgres",
		"PSQL_PASSWORD": "postgres",
		"PSQL_PORT":     "5432",
		"PSQL_HOST":     "localhost",
		"PSQL_NAME":     "postgres",
		"PSQL_SSLMODE":  "disable",
	}
	for k := range testPostgres {
		if err := os.Setenv(k, testPostgres[k]); err != nil {
			log.Print(err.Error())
		}
	}

	m.Run()
}

func TestRepository(t *testing.T) {
	log := customLogger.NewLogger()
	db := newPsql(log, config.New(log).Repository.Psql)
	auth := newTestAuth(log, db)
	key := "keykeykey"
	user := dto.TGUser{
		int64(11111),
		int64(22222),
	}

	convey.Convey("init", t, func() {

		convey.Convey("auth", func() { auth.testAuth(user, key) })

	})

}

type testAuth struct {
	auth psqlAuth
	log  customLogger.Logger
}

func newTestAuth(log customLogger.Logger, db *sqlx.DB) testAuth {
	return testAuth{
		auth: newPsqlAuth(log, db),
		log:  log,
	}
}

func (a testAuth) testAuth(user dto.TGUser, key string) {
	defer a.clean(user)

	a.testAuthVk(user, key)

}

func (a testAuth) clean(user dto.TGUser) {
	a.auth.Vk().Delete(user)
}

func (a testAuth) testAuthVk(user dto.TGUser, key string) {
	authVk := a.auth.Vk()

	authVk.Create(user, key)
	convey.So(authVk.Get(user), convey.ShouldEqual, key)
}
