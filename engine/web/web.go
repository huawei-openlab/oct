package web

import (
	"github.com/Unknwon/macaron"

	"github.com/huawei-openlab/oct/engine/middleware"
	"github.com/huawei-openlab/oct/engine/router"
	"github.com/huawei-openlab/oct/engine/setting"
)

func SetOctMacaron(m *macaron.Macaron) {
	//Setting
	setting.SetConfig("conf/engine.conf")

	//Setting Middleware
	middleware.SetMiddlewares(m)

	//Setting Router
	router.SetRouters(m)
}
