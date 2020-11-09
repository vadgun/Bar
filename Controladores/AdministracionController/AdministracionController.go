package administracioncontroller

import (
	"fmt"
	"strconv"

	"github.com/kataras/iris/v12"
	sessioncontroller "github.com/vadgun/Bar/Controladores/SessionController"
	indexmodel "github.com/vadgun/Bar/Modelos/IndexModel"
	ventamodel "github.com/vadgun/Bar/Modelos/VentaModel"
)

//Administracion -> Devuelve la vista de la Administracion
func Administracion(ctx iris.Context) {
	fmt.Println("Prueba de administracion")

	userOn := indexmodel.GetUserOn(sessioncontroller.Sess.Start(ctx).GetString("UserID"))
	ctx.ViewData("Usuario", userOn)

	MesasDisponibles := ventamodel.ExtraeMesas()

	ctx.ViewData("MesasDisponibles", MesasDisponibles)
	if err := ctx.View("Administracion.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

//AdministracionFondo -> Cambia el fondo de pantalla
func AdministracionFondo(ctx iris.Context) {
	//saber el ultimo fondo usado
	var fondo int
	fondo = ventamodel.ExtraeFondo()
	fondostring := strconv.Itoa(fondo)
	ctx.HTML(fondostring)
}

//AdministracionFondoCambiar -> Guarda el ultimo fondo seleccionado hasta 33
func AdministracionFondoCambiar(ctx iris.Context) {
	fondo, _ := ctx.PostValueInt("data")
	ventamodel.ActualizaFondo(fondo)
}
