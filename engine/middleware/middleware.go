package middleware

import (
	"github.com/Unknwon/macaron"
)

func SetMiddlewares(m *macaron.Macaron) {
	//Set static file directory,static file access without log output
	m.Use(macaron.Static("static", macaron.StaticOptions{
		Expires: func() string { return "max-age=0" },
	}))

	//Set global Logger
	m.Map(Log)
	//Set logger handler function, deal with all the Request log output
	m.Use(logger("dev"))

	m.Use(macaron.Recovery())
}
