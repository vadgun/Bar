package indexcontroller

import (
	"fmt"
	"github.com/kataras/iris/v12"
	sessioncontroller "github.com/vadgun/Bar/Controladores/SessionController"
	indexmodel "github.com/vadgun/Bar/Modelos/IndexModel"
)

//Index -> Regresa la pagina de inicio
func Index(ctx iris.Context) {
	var usuario indexmodel.MongoUser
	var autorizado bool
	autorizado2, _ := sessioncontroller.Sess.Start(ctx).GetBoolean("Autorizado")

	if autorizado2 == false {
		usuario.Key = ctx.PostValue("pass")
		usuario.Usuario = ctx.PostValue("usuario")
		autorizado, usuario = indexmodel.VerificarUsuario(usuario)
		if autorizado {
			sessioncontroller.Sess.Start(ctx).Set("Autorizado", true)
			sessioncontroller.Sess.Start(ctx).Set("UserID", usuario.ID.Hex())
		}
	}

	if autorizado || autorizado2 {
		userOn := indexmodel.GetUserOn(sessioncontroller.Sess.Start(ctx).GetString("UserID"))
		ctx.ViewData("Usuario", userOn)
		// ctx.ViewData("Permisos", permisos)

		// Extraer el total de los presupuestos 2020 -
		// Saber que numeros economicos faltan y cuales ya fueron capturados
		// Crear un resumen general para la pagina principal

		// ac := accounting.Accounting{Symbol: "$", Precision: 2}
		// total, totales, totalorden := herramientasmodel.ExtraeTotalPresupuesto(2020)
		fmt.Println("Bienvenido ", userOn.Nombre)
		// totalformat := ac.FormatMoney(total)
		// ctx.ViewData("Total", totalformat)
		// ctx.ViewData("Totales", totales)
		// ctx.ViewData("Totalorden", totalorden)

		if err := ctx.View("Index.html"); err != nil {
			ctx.Application().Logger().Infof(err.Error())
		}
	} else {
		ctx.Redirect("/login", iris.StatusSeeOther)
	}
}
