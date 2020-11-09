package sessioncontroller

import (
	"github.com/kataras/iris/v12/sessions"
)

var (
	cokkiename = "BarLa49"
	//Sess -> Variable que apunta a la sesion
	Sess = sessions.New(sessions.Config{Cookie: cokkiename})
)
