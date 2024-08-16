package administracioncontroller

import (
	"fmt"

	"github.com/kataras/iris/v12"
	sessioncontroller "github.com/vadgun/Bar/Controladores/SessionController"
	indexmodel "github.com/vadgun/Bar/Modelos/IndexModel"
	ventamodel "github.com/vadgun/Bar/Modelos/VentaModel"
)

// Administracion -> Devuelve la vista de la Administracion
func Administracion(ctx iris.Context) {
	fmt.Println("Prueba de administracion")

	userOn := indexmodel.GetUserOn(sessioncontroller.Sess.Start(ctx).GetString("UserID"))
	ctx.ViewData("Usuario", userOn)

	MesasDisponibles := ventamodel.ExtraeMesas()
	fondo := ventamodel.ExtraeFondo()

	Mensaje1 := fondo.Mensajes[0]
	Mensaje2 := fondo.Mensajes[1]
	Mensaje3 := fondo.Mensajes[2]

	ctx.ViewData("MesasDisponibles", MesasDisponibles)
	ctx.ViewData("Mensaje1", Mensaje1)
	ctx.ViewData("Mensaje2", Mensaje2)
	ctx.ViewData("Mensaje3", Mensaje3)
	if err := ctx.View("Administracion.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

// AdministracionFondo -> Cambia el fondo de pantalla
func AdministracionFondo(ctx iris.Context) {
	//saber el ultimo fondo usado
	config := ventamodel.ExtraeFondo()

	ctx.JSON(config)

}

// AdministracionFondoCambiar -> Guarda el ultimo fondo seleccionado hasta 33
func AdministracionFondoCambiar(ctx iris.Context) {
	fondo, _ := ctx.PostValueInt("data")
	ventamodel.ActualizaFondo(fondo)
}

// AdministracionMensajesCambiar -> Guarda el ultimo fondo seleccionado hasta 33
func AdministracionMensajesCambiar(ctx iris.Context) {
	var mensajes = []string{ctx.PostValue("mensaje1"), ctx.PostValue("mensaje2"), ctx.PostValue("mensaje3")}
	fmt.Println(mensajes)

	ventamodel.ActualizaMensajes(mensajes)

	htmlcode := fmt.Sprintf(`
	<script>
		alert("Mensajes Actualizados");
		location.replace("/administracion");
	</script>
	`)

	ctx.HTML(htmlcode)
}
