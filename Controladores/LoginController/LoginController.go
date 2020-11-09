package logincontroller

import (
	"fmt"

	"github.com/kataras/iris/v12"

	sessioncontroller "github.com/vadgun/Bar/Controladores/SessionController"
	indexmodel "github.com/vadgun/Bar/Modelos/IndexModel"
)

//Getlogin -> Funcion que regresa el GET de Login
func Getlogin(ctx iris.Context) {
	autorizado, _ := sessioncontroller.Sess.Start(ctx).GetBoolean("Autorizado")

	if autorizado {
		userOn := indexmodel.GetUserOn(sessioncontroller.Sess.Start(ctx).GetString("UserID"))
		ctx.ViewData("Usuario", userOn)
		fmt.Println("Vista -- Index.html")
		if err := ctx.View("Index.html"); err != nil {
			ctx.Application().Logger().Infof(err.Error())
		}
	} else {
		fmt.Println("Vista -- Login.html")
		if err := ctx.View("Login.html"); err != nil {
			ctx.Application().Logger().Infof(err.Error())
		}
	}
}

// Getlogout Cierra una sesion y destruye la sesion
func Getlogout(ctx iris.Context) {
	fmt.Println("Cerrando Sesion, Hasta Luego")
	sessioncontroller.Sess.Start(ctx).Set("Autorizado", false)
	sessioncontroller.Sess.Destroy(ctx)
	ctx.Redirect("/login", iris.StatusSeeOther)
}
