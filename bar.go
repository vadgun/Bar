package main

import (
	"github.com/kataras/iris/v12"
	administracioncontroller "github.com/vadgun/Bar/Controladores/AdministracionController"
	indexcontroller "github.com/vadgun/Bar/Controladores/IndexController"
	inventariocontroller "github.com/vadgun/Bar/Controladores/InventarioController"
	logincontroller "github.com/vadgun/Bar/Controladores/LoginController"
	ventacontroller "github.com/vadgun/Bar/Controladores/VentaController"
	websocketcontroller "github.com/vadgun/Bar/Controladores/Websocketcontroller"
)

func main() {
	app := iris.New()
	app.HandleDir("/Recursos", "./Recursos")
	app.Favicon("./Recursos/Imagenes/favicon.ico")
	app.RegisterView(iris.HTML("./Vistas", ".html").Reload(true))

	app.Get("/", logincontroller.Getlogin)
	app.Get("/login", logincontroller.Getlogin)
	app.Post("/login", logincontroller.Getlogin)

	app.Get("/logout", logincontroller.Getlogout)

	app.Post("/index", indexcontroller.Index)
	app.Get("/index", indexcontroller.Index)

	app.Get("/inventario", inventariocontroller.Inventario)
	app.Post("/inventario", inventariocontroller.Inventario)

	app.Get("/venta", ventacontroller.Venta)
	app.Post("/venta", ventacontroller.Venta)

	app.Get("/ventadiaria", ventacontroller.VentaDiaria)
	app.Post("/ventadiaria", ventacontroller.VentaDiaria)

	app.Get("/administracion", administracioncontroller.Administracion)
	app.Post("/administracion", administracioncontroller.Administracion)

	//Formularios
	app.Post("/altaform", inventariocontroller.Altaform)
	app.Post("/altaproducto", inventariocontroller.GuardarProducto)
	app.Post("/editarproducto", inventariocontroller.EditarProducto)
	app.Post("/editandoproducto", inventariocontroller.EditandoProducto)
	app.Post("/eliminarproducto", inventariocontroller.EliminarProducto)
	app.Post("/editaralmacen", inventariocontroller.EditarAlmacen)
	app.Post("/actualizaralmacen", inventariocontroller.ActualizarAlmacen)
	app.Post("/productosASelect", inventariocontroller.ProductosASelect)

	//Venta
	app.Post("/configurarmesasdiaras", ventacontroller.ConfigurarMesas)
	app.Post("/verificarventa", ventacontroller.VerificarVenta)
	app.Post("/abrirventadiaria", ventacontroller.AbrirVentaDiaria)
	app.Post("/cerrarventademesa", ventacontroller.CerrarVentaDeMesa)
	app.Post("/productosEnModal", ventacontroller.ProductoDeModal)
	app.Post("/productosVendidos", ventacontroller.ProductosVendidos)

	app.Post("/productosAgregadosEnModal", ventacontroller.ProductosAgregadosEnModal)
	app.Post("/productosAgregadosEnModalDesdePromo", ventacontroller.ProductosAgregadosEnModalDesdePromo)

	app.Post("/clasificaventa", ventacontroller.ClasificaVenta)

	//Traspasos
	app.Post("/enviarproductoabodega", ventacontroller.EnviarProductoABodega)
	app.Post("/enviarproductoarefrigerador", ventacontroller.EnviarProductoARefrigerador)
	app.Post("/eliminarproductodebodega", ventacontroller.EliminarProductoDeBodega)
	app.Post("/eliminarproductoderefrigerador", ventacontroller.EliminarProductoDeRefrigerador)

	//Imprmir PDF
	app.Get("/imprimirventadiaria", ventacontroller.ImprimirVentaDiaria)
	app.Post("/imprimirventadiaria", ventacontroller.ImprimirVentaDiaria)
	app.Post("/imprimeventa", ventacontroller.ImprimirVentaDiariaEnPDF2)

	//Eliminar Colecciones
	app.Post("/vaciarbasededatos", ventacontroller.VaciarColeccion)

	//Cambiar Fondo de pantalla y Recordarlo
	app.Post("/ultimofondousado", administracioncontroller.AdministracionFondo)
	app.Post("/cambiarfondo", administracioncontroller.AdministracionFondoCambiar)
	app.Post("/configurarmensajes", administracioncontroller.AdministracionMensajesCambiar)

	// WebSocket
	app.Get("/ws", websocketcontroller.WebsocketHandler)
	go websocketcontroller.HandleMessages()
	go websocketcontroller.MongoSupervisor()

	app.Run(iris.Addr(":8080"))
}
