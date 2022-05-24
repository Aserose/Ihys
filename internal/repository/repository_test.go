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

	convey.Convey("init", t, func() {
		user := dto.TGUser{
			int64(11111),
			int64(22222),
		}
		key := "f2132f23"
		rep := NewRepository(log, config.NewCfg(log).Repository)

		convey.Convey("put", func() {
			defer rep.Vk().DeleteKey(user)

			rep.Vk().PutKey(user, key)

			convey.So(rep.Vk().GetKey(user), convey.ShouldEqual, key)

		})

	})

}
