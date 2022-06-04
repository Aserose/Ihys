package repository

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/dto"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRepository(t *testing.T) {
	log := customLogger.NewLogger()
	rep := newRepo(log)
	key := "keykeykey"
	user := dto.TGUser{
		int64(11111),
		int64(22222),
	}

	convey.Convey("init", t, func() {

		convey.Convey("auth", func() {
			rep.testAuth(user, key)
		})

	})

}

type repo struct {
	rep Repository
	log customLogger.Logger
}

func newRepo(log customLogger.Logger) repo {
	return repo{
		rep: NewRepository(log, config.NewCfg(log).Repository),
		log: log,
	}
}

func (r repo) testAuth(user dto.TGUser, key string) {
	defer r.clean(user)

	r.testAuthVk(user, key)

}

func (r repo) clean(user dto.TGUser) {
	r.rep.Auth.Vk().DeleteKey(user)
}

func (r repo) testAuthVk(user dto.TGUser, key string) {
	authVk := r.rep.Auth.Vk()

	authVk.PutKey(user, key)
	convey.So(authVk.GetKey(user), convey.ShouldEqual, key)
}
