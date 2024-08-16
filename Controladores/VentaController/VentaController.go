package ventacontroller

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/kataras/iris/v12"
	"github.com/leekchan/accounting"
	"go.mongodb.org/mongo-driver/bson/primitive"

	sessioncontroller "github.com/vadgun/Bar/Controladores/SessionController"
	indexmodel "github.com/vadgun/Bar/Modelos/IndexModel"
	inventariomodel "github.com/vadgun/Bar/Modelos/InventarioModel"
	ventamodel "github.com/vadgun/Bar/Modelos/VentaModel"
)

// Venta -> Devuelve la vista de la Venta
func Venta(ctx iris.Context) {
	userOn := indexmodel.GetUserOn(sessioncontroller.Sess.Start(ctx).GetString("UserID"))
	ctx.ViewData("Usuario", userOn)

	//Enviar informacion de Numero de Mesas y la fecha del dia.

	// b := ventamodel.ExtraeMesas()

	// ctx.ViewData("Mesas", b)

	if err := ctx.View("Venta.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

// VentaDiaria -> Devuelve la vista de la Venta Diaria
func VentaDiaria(ctx iris.Context) {
	userOn := indexmodel.GetUserOn(sessioncontroller.Sess.Start(ctx).GetString("UserID"))
	ctx.ViewData("Usuario", userOn)

	//Enviar informacion de Numero de Mesas y la fecha del dia.

	// b := ventamodel.ExtraeMesas()

	// ctx.ViewData("Mesas", b)
	if err := ctx.View("VentaDiaria.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

// ConfigurarMesas -> Configura la cantidad de mesas para abrir la venta del dia
func ConfigurarMesas(ctx iris.Context) {

	var htmlcode string

	mesasdiarias, _ := ctx.PostValueInt("mesasdiarias")

	ventamodel.ActualizarMesasDiarias(mesasdiarias)

	htmlcode += fmt.Sprintf(`
	<script>
		alert("Mesas Actualizadas");
		location.replace("/administracion");
	</script>
	`)

	ctx.HTML(htmlcode)

}

// VerificarVenta -> Toma en cuenta la fecha para devolver las mesas diarias correspondientes a los dias, si no existe una cuenta diaria la crea por fecha y crea la mesa diara para la venta
func VerificarVenta(ctx iris.Context) {

	fechaVenta := ctx.PostValue("data")
	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	MesasDisponibles := ventamodel.ExtraeMesas()

	userOn := indexmodel.GetUserOn(sessioncontroller.Sess.Start(ctx).GetString("UserID"))

	layout := "2006-01-02"
	location, _ := time.LoadLocation("America/Mexico_City")

	t, _ := time.ParseInLocation(layout, fechaVenta, location)

	mesasdiarias := ventamodel.BuscarVentasDiarias(t)

	var htmlcode string

	if len(mesasdiarias) == 0 {

		htmlcode += fmt.Sprintf(`
		<div class="centrado">
			<h3>No se ha encontrado Venta para %v </h3>
		</div>`, fechaVenta)

		if userOn.Admin {
			htmlcode += fmt.Sprintf(`
		<div class="centrado">
			<a id="abrirventadiariabutton" class="btn btn-info btn-large padd" href="Javascript:AbrirVentaDelDia('%v');" role="button">Iniciar venta para %v </a>&nbsp;
		</div>

		<div class="centrado">
			<h3>Se abriran %v Mesas</h3>
		</div>

		`, fechaVenta, fechaVenta, MesasDisponibles)

		}
	} else {
		htmlcode += fmt.Sprintf(`<div class="container">
		<div class="row">
			<br>	`)

		for _, v := range mesasdiarias {

			if v.Estatus == false && v.Abierta == true {
				htmlcode += fmt.Sprintf(`
			<div class="col-6 col-sm-6 col-md-4 col-lg-4 col-xl-2" style="text-align: center; padding: 1PX; ">
				<div class="mesasinocupar">
					<H1 class="bignumber" >%v</H1>`, v.Mesa)

				htmlcode += fmt.Sprintf(`
					<br>DESOCUPADA<br>					
					<a class="btn-sm btn-warning" href="#" role="button" data-mesa="%v" data-id="%v" data-toggle="modal" data-target="#ModalProductos">Agregar Productos</a>
				</div>
		</div>`, v.Mesa, v.ID.Hex())

			} else if v.Estatus == true && v.Cerrada == false {
				htmlcode += fmt.Sprintf(`
						<div class="col-6 col-sm-6 col-md-4 col-lg-4 col-xl-2" style="text-align: center; padding: 1PX; ">
							<a href="Javascript:Ocupada('%v');" style="color: black;">
							<div class="mesasconclientes">
								<H2>%v</H2>`, v.Mesero, v.Mesa)

				htmlcode += fmt.Sprintf(`
								<h4 class="highgreen">%v</h4>
								<h4 class="highgreen">%v</h4>
								<a class="btn btn-primary btn" href="#" role="button" data-mesa="%v" data-id="%v" data-toggle="modal" data-target="#ModalProductos">Agregar Productos</a>
								`, ac.FormatMoney(v.GranTotal), v.Mesero, v.Mesa, v.ID.Hex())

				if userOn.Admin {
					htmlcode += fmt.Sprintf(`<a class="btn btn-danger btn" href="Javascript:CerrarVenta('%v');" role="button">Cerrar Venta</a>
					</div>
					</a>
				</div>`, v.ID.Hex())
				} else {
					htmlcode += fmt.Sprintf(`</div></a></div>`)
				}

			}

		}

		htmlcode += fmt.Sprintf(`
			</div>
		</div>
		`)

	}

	ctx.HTML(htmlcode)

}

// AbrirVentaDiaria -> Algoritmo que crea la venta del dia y la regresa al main contanier donde fue hecha la peticion
func AbrirVentaDiaria(ctx iris.Context) {

	fechaVenta := ctx.PostValue("data")
	var htmlcode string
	MesasDisponibles := ventamodel.ExtraeMesas()
	htmlcode += fmt.Sprintf(``)
	layout := "2006-01-02"
	location, _ := time.LoadLocation("America/Mexico_City")
	t, _ := time.ParseInLocation(layout, fechaVenta, location)

	//Crear las mesas diarias de acuerdo con la numeracion de mesas a abrir, por el momento solo tienen que llevar numero de mesa y status.
	var productos []primitive.ObjectID
	var cantidades []int

	for i := 1; i <= MesasDisponibles; i++ {

		var nuevamesa ventamodel.Mesa
		nuevamesa.Estatus = false
		nuevamesa.Mesa = i
		nuevamesa.Fecha = t
		nuevamesa.GranTotal = 0
		nuevamesa.Productos = productos
		nuevamesa.Cantidades = cantidades
		nuevamesa.Abierta = true
		nuevamesa.Cerrada = false

		ventamodel.GuardaMesaDiaria(nuevamesa)

	}

	htmlcode += fmt.Sprintf(`
		<script>
			AlertaCreada("Venta diaria creada");
		</script>
	`)

	ctx.HTML(htmlcode)

	//Crear el arreglo de mesas diarias con cuenta en 0 - las cuentas 0 seran la inicial de cada mesa si es diferente de 0 q sume 1
	//Las mesas seran abiertas del registro de configuracion de mesas diarias
	//Se creara un modelo a la base de datos que simule crear las mesas para abrirlas todas
	//Mesa Desocupada Status : False    Mesa Ocupada Status : true
	//Cada mesa tendra un arreglo de productos y existencia de productos similar al de los almacenes
	//Cada mesa tendra un gran total el cual se vera reflejado en la interfaz de las mesas
	//Para seguir vendiendo en dicha mesa despues de cerrar una venta habra que crear otra mesa con el mismo numero de mesa pero con otra cuenta
	//Las cuentas son las que seran indefinidas a lo largo del dia
	//Pueden haber 20 cuentas en una sola mesa
	//Cuenta 1-2-3-4-5-5-6-7-8
	//Una mesa puede tener muchas cuentas
	//Una cuenta solo puede tener una mesa
	//Cada mesa podra cerrar una venta y cambiar a estatus false cuando este al estar ocupada seria true
	//Cada mesa alimentara la venta diaria la cual alimentara la columna de utilidad al final del dia.
}

// EnviarProductoABodega -> Envia el producto a Bodega para dar de alta su existencia
func EnviarProductoABodega(ctx iris.Context) {

	var htmlcode string

	prodYAlmacen := ctx.PostValue("data")
	arreglo := strings.Split(prodYAlmacen, ":")
	producto := arreglo[0]
	almacen := arreglo[1]

	if inventariomodel.AgregarAAlmacen(producto, almacen) {

		htmlcode += fmt.Sprintf(`
		<script>
			alert("Producto Enviado a Bodega");
			$("#subcontainer").html('');
		</script>
		`)

	} else {
		htmlcode += fmt.Sprintf(`
		<script>
			alert("Este producto ya esta en bodega");
			$("#subcontainer").html('');
		</script>
		`)
	}

	ctx.HTML(htmlcode)

}

// EnviarProductoARefrigerador -> Envia el producto a Refrigerador para dar de alta su existencia
func EnviarProductoARefrigerador(ctx iris.Context) {

	var htmlcode string

	prodYAlmacen := ctx.PostValue("data")
	arreglo := strings.Split(prodYAlmacen, ":")
	producto := arreglo[0]
	almacen := arreglo[1]

	//Hacer una funcion en el modelo que saque el producto de una coleccion y lo agregue al array correspondiente,
	//solo se puede agregar una vez, asi q verificar el id q no se repita

	if inventariomodel.AgregarAAlmacen(producto, almacen) {

		htmlcode += fmt.Sprintf(`
		<script>
			alert("Producto Enviado a refrigerador");
			$("#subcontainer").html('');
		</script>
		`)

	} else {
		htmlcode += fmt.Sprintf(`
		<script>
			alert("Este producto ya esta en refrigerador");
			$("#subcontainer").html('');
		</script>
		`)
	}

	ctx.HTML(htmlcode)

}

// EliminarProductoDeBodega -> Eliminar de Bodega
func EliminarProductoDeBodega(ctx iris.Context) {
	var htmlcode string

	prodYAlmacen := ctx.PostValue("data")
	arreglo := strings.Split(prodYAlmacen, ":")
	producto := arreglo[0]
	almacen := arreglo[1]

	//Hacer una funcion en el modelo que saque el producto de una coleccion y lo agregue al array correspondiente,
	//solo se puede agregar una vez, asi q verificar el id q no se repita

	if inventariomodel.EliminarDeAlmacen(producto, almacen) {

		htmlcode += fmt.Sprintf(`
		<script>
			alert("Producto Eliminado de bodega");
			$("#subcontainer").html('');
		</script>
		`)

	} else {
		htmlcode += fmt.Sprintf(`
		<script>
			alert("Ya eliminado ó se tiene existencia de bodega");
			$("#subcontainer").html('');
		</script>
		`)
	}

	ctx.HTML(htmlcode)
}

// EliminarProductoDeRefrigerador -> Eliminar de Refrigerador
func EliminarProductoDeRefrigerador(ctx iris.Context) {
	var htmlcode string

	prodYAlmacen := ctx.PostValue("data")
	arreglo := strings.Split(prodYAlmacen, ":")
	producto := arreglo[0]
	almacen := arreglo[1]

	//Hacer una funcion en el modelo que saque el producto de una coleccion y lo agregue al array correspondiente,
	//solo se puede agregar una vez, asi q verificar el id q no se repita

	if inventariomodel.EliminarDeAlmacen(producto, almacen) {

		htmlcode += fmt.Sprintf(`
		<script>
			alert("Producto Eliminado de refrigerador");
			$("#subcontainer").html('');
		</script>
		`)

	} else {
		htmlcode += fmt.Sprintf(`
		<script>
			alert("Ya eliminardo ó Se tiene existencia en refrigerador");
			$("#subcontainer").html('');
		</script>
		`)
	}

	ctx.HTML(htmlcode)
}

// CerrarVentaDeMesa -> Cierra la venta diaria de una mesa y devuelve su estatus a false
func CerrarVentaDeMesa(ctx iris.Context) {
	idmesa := ctx.PostValue("data")
	cadenas := strings.Split(idmesa, ":")
	mesa := ventamodel.ExtraeMesa(cadenas[0])
	ventamodel.CierraMesa(mesa)
	ctx.HTML("Cerrando Mesas")
}

// AgregarProductosHTML -> Regresa el codigo HTML aplicando condiciones, for y switch
func AgregarProductosHTML(productos []inventariomodel.Producto) string {
	var htmlcode string

	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	for k, v := range productos {

		existenciaRefri := inventariomodel.ExtraeExistencias(v.ID)

		switch {
		case existenciaRefri == 0:
			htmlcode += fmt.Sprintf(`
			<tr>
				<th scope="row">%v</th>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
				<td>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">1</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">2</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">5</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">10</a>
				</td>
			</tr>
			`, k+1, v.Nombre, ac.FormatMoney(v.PrecioPub), existenciaRefri)
			break

		case existenciaRefri > 0 && existenciaRefri <= 1:
			htmlcode += fmt.Sprintf(`
			<tr>
				<th scope="row">%v</th>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
				<td>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:1:%v');" role="button">1</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">2</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">5</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">10</a>
				</td>
			</tr>
			`, k+1, v.Nombre, ac.FormatMoney(v.PrecioPub), existenciaRefri, v.ID.Hex(), v.Categoria)
			break
		case existenciaRefri > 0 && existenciaRefri < 3:
			htmlcode += fmt.Sprintf(`
			<tr>
				<th scope="row">%v</th>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
				<td>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:1:%v');" role="button">1</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:2:%v');" role="button">2</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">5</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">10</a>
				</td>
			</tr>
			`, k+1, v.Nombre, ac.FormatMoney(v.PrecioPub), existenciaRefri, v.ID.Hex(), v.Categoria, v.ID.Hex(), v.Categoria)
			break
		case existenciaRefri > 2 && existenciaRefri < 5:
			htmlcode += fmt.Sprintf(`
			<tr>
				<th scope="row">%v</th>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
				<td>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:1:%v');" role="button">1</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:2:%v');" role="button">2</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">5</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">10</a>

				</td>
			</tr>
			`, k+1, v.Nombre, ac.FormatMoney(v.PrecioPub), existenciaRefri, v.ID.Hex(), v.Categoria, v.ID.Hex(), v.Categoria)
			break
		case existenciaRefri >= 5 && existenciaRefri < 10:
			htmlcode += fmt.Sprintf(`
			<tr>
				<th scope="row">%v</th>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
				<td>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:1:%v');" role="button">1</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:2:%v');" role="button">2</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:5:%v');" role="button">5</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:alert('Insuficiente, Surtir Almacen.');" role="button">10</a>
				</td>
			</tr>
			`, k+1, v.Nombre, ac.FormatMoney(v.PrecioPub), existenciaRefri, v.ID.Hex(), v.Categoria, v.ID.Hex(), v.Categoria, v.ID.Hex(), v.Categoria)

			break
		case existenciaRefri >= 10:
			htmlcode += fmt.Sprintf(`
			<tr>
				<th scope="row">%v</th>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
				<td>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:1:%v');" role="button">1</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:2:%v');" role="button">2</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:5:%v');" role="button">5</a>
				<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesa('%v:10:%v');" role="button">10</a>
				</td>
			</tr>
			`, k+1, v.Nombre, ac.FormatMoney(v.PrecioPub), existenciaRefri, v.ID.Hex(), v.Categoria, v.ID.Hex(), v.Categoria, v.ID.Hex(), v.Categoria, v.ID.Hex(), v.Categoria)
			break

		}
	}

	return htmlcode
}

// ProductoDeModal -> Devuelve los productos selecionados a la ventana del modal
func ProductoDeModal(ctx iris.Context) {

	categoriaproductosAModal := ctx.PostValue("data")

	var productos []inventariomodel.Producto

	//Extraer Catergoria
	var htmlcode string

	switch categoriaproductosAModal {
	case "botana":

		productos = inventariomodel.ExtraeBotanas()

		lenbotanas := len(productos)

		if lenbotanas == 0 {
			htmlcode += fmt.Sprintf(`<div class="container">
			<div class="centrado">
				<h1 class="display-4">No hay botanas dadas de alta</h1>
			</div>
		</div>`)
		} else {
			htmlcode += fmt.Sprintf(`<table class="table table-hover" style="margin: auto; width: 65%s !important; font-size:14px; text-align: center;">
			<input type="hidden" id="tablaactiva" value="botana">
			<thead>
			  <tr>
				<th scope="col">#</th>
				<th scope="col">Producto</th>
				<th scope="col">Precio</th>
				<th scope="col">En Refrigerador</th>
				<th scope="col">Agregar</th>
			  </tr>
			</thead>
				<tbody>
					`, "%%")

			htmlcode += fmt.Sprintf(AgregarProductosHTML(productos))

			htmlcode += fmt.Sprintf(`	
				</tbody>
			</table>`)
		}

		break

	case "botella":

		productos = inventariomodel.ExtraeBotellas()
		lenbotellas := len(productos)

		if lenbotellas == 0 {
			htmlcode += fmt.Sprintf(`<div class="container">
			<div class="centrado">
				<h1 class="display-4">No hay botellas dadas de alta</h1>
			</div>
		</div>`)
		} else {
			htmlcode += fmt.Sprintf(`<table class="table table-hover" style="margin: auto; width: 65%s !important; font-size:14px; text-align: center;">
			<input type="hidden" id="tablaactiva" value="botella">
			<thead>
			  <tr>
				<th scope="col">#</th>
				<th scope="col">Producto</th>
				<th scope="col">Precio</th>
				<th scope="col">En Refrigerador</th>
				<th scope="col">Agregar</th>
			  </tr>
			</thead>
				<tbody>
					`, "%%")

			htmlcode += fmt.Sprintf(AgregarProductosHTML(productos))

			htmlcode += fmt.Sprintf(`	
				</tbody>
			</table>`)
		}

		break
	case "cigarros":
		productos = inventariomodel.ExtraeCigarros()
		lencigarros := len(productos)

		if lencigarros == 0 {
			htmlcode += fmt.Sprintf(`<div class="container">
			<div class="centrado">
				<h1 class="display-4">No hay cigarros dados de alta</h1>
			</div>
		</div>`)
		} else {
			htmlcode += fmt.Sprintf(`<table class="table table-hover" style="margin: auto; width: 65%s !important; font-size:14px; text-align: center;">
			<input type="hidden" id="tablaactiva" value="cigarros">
			<thead>
			  <tr>
				<th scope="col">#</th>
				<th scope="col">Producto</th>
				<th scope="col">Precio</th>
				<th scope="col">En Refrigerador</th>
				<th scope="col">Agregar</th>
			  </tr>
			</thead>
				<tbody>
					`, "%%")

			htmlcode += fmt.Sprintf(AgregarProductosHTML(productos))

			htmlcode += fmt.Sprintf(`	
				</tbody>
			</table>`)
		}
		break
	case "cerveza":
		productos = inventariomodel.ExtraeCervezas()
		lencervezas := len(productos)

		if lencervezas == 0 {
			htmlcode += fmt.Sprintf(`<div class="container">
			<div class="centrado">
				<h1 class="display-4">No hay cervezas dadas de alta</h1>
			</div>
		</div>`)
		} else {
			htmlcode += fmt.Sprintf(`<table class="table table-hover" style="margin: auto; width: 65%s !important; font-size:14px; text-align: center;">
			<input type="hidden" id="tablaactiva" value="cerveza">
			<thead>
			  <tr>
				<th scope="col">#</th>
				<th scope="col">Producto</th>
				<th scope="col">Precio</th>
				<th scope="col">En Refrigerador</th>
				<th scope="col">Agregar</th>
			  </tr>
			</thead>
				<tbody>
					`, "%%")

			htmlcode += fmt.Sprintf(AgregarProductosHTML(productos))

			htmlcode += fmt.Sprintf(`	
				</tbody>
			</table>`)
		}
		break
	case "copa":
		productos = inventariomodel.ExtraeCopas()
		lencopa := len(productos)

		if lencopa == 0 {
			htmlcode += fmt.Sprintf(`<div class="container">
			<div class="centrado">
				<h1 class="display-4">No hay copas dadas de alta</h1>
			</div>
		</div>`)
		} else {
			htmlcode += fmt.Sprintf(`<table class="table table-hover" style="margin: auto; width: 65%s !important; font-size:14px; text-align: center;">
			<input type="hidden" id="tablaactiva" value="copa">
			<thead>
			  <tr>
				<th scope="col">#</th>
				<th scope="col">Producto</th>
				<th scope="col">Precio</th>
				<th scope="col">En Refrigerador</th>
				<th scope="col">Agregar</th>
			  </tr>
			</thead>
				<tbody>
					`, "%%")

			htmlcode += fmt.Sprintf(AgregarProductosHTML(productos))

			htmlcode += fmt.Sprintf(`	
				</tbody>
			</table>`)
		}
		break
	case "ficha":
		productos = inventariomodel.ExtraeFichas()
		lenfichas := len(productos)

		if lenfichas == 0 {
			htmlcode += fmt.Sprintf(`<div class="container">
			<div class="centrado">
				<h1 class="display-4">No hay fichas dadas de alta</h1>
			</div>
		</div>`)
		} else {
			htmlcode += fmt.Sprintf(`<table class="table table-hover" style="margin: auto; width: 65%s !important; font-size:14px; text-align: center;">
			<input type="hidden" id="tablaactiva" value="ficha">
			<thead>
			  <tr>
				<th scope="col">#</th>
				<th scope="col">Producto</th>
				<th scope="col">Precio</th>
				<th scope="col">En Refrigerador</th>
				<th scope="col">Agregar</th>
			  </tr>
			</thead>
				<tbody>
					`, "%%")

			htmlcode += fmt.Sprintf(AgregarProductosHTML(productos))

			htmlcode += fmt.Sprintf(`	
				</tbody>
			</table>`)
		}
		break
	case "cancion":
		productos = inventariomodel.ExtraeCancion()
		lencanciones := len(productos)

		if lencanciones == 0 {
			htmlcode += fmt.Sprintf(`<div class="container">
			<div class="centrado">
				<h1 class="display-4">No hay canciones dadas de alta</h1>
			</div>
		</div>`)
		} else {
			htmlcode += fmt.Sprintf(`<table class="table table-hover" style="margin: auto; width: 65%s !important; font-size:14px; text-align: center;">
			<input type="hidden" id="tablaactiva" value="cancion">
			<thead>
			  <tr>
				<th scope="col">#</th>
				<th scope="col">Producto</th>
				<th scope="col">Precio</th>
				<th scope="col">En Refrigerador</th>
				<th scope="col">Agregar</th>
			  </tr>
			</thead>
				<tbody>
					`, "%%")
			htmlcode += fmt.Sprintf(AgregarProductosHTML(productos))

			htmlcode += fmt.Sprintf(`	
				</tbody>
			</table>`)
		}
		break
	case "promo":
		productos = inventariomodel.ExtraePromos()
		lenpromos := len(productos)

		if lenpromos == 0 {
			htmlcode += fmt.Sprintf(`<div class="container">
			<div class="centrado">
				<h1 class="display-4">No hay promociones dadas de alta</h1>
			</div>
		</div>`)
		} else {
			// htmlcode += fmt.Sprintf(`<table class="table table-hover" style="margin: auto; width: 65%s !important; font-size:14px; text-align: center;">
			// <thead>
			//   <tr>
			// 	<th scope="col">#</th>
			// 	<th scope="col">Producto</th>
			// 	<th scope="col">Precio</th>
			// 	<th scope="col">En Refrigerador</th>
			// 	<th scope="col">Agregar</th>
			//   </tr>
			// </thead>
			// 	<tbody>
			// 		`, "%%")
			// htmlcode += fmt.Sprintf(AgregarProductosHTML(productos))

			// htmlcode += fmt.Sprintf(`
			// 	</tbody>
			// </table>`)

			htmlcode += fmt.Sprintf(`
			<table class="table table-hover" style="margin: auto; width: 65%s !important; font-size:14px; text-align: center;">
			<input type="hidden" id="tablaactiva" value="promo">
			<thead>
				<tr>
					<th scope="col">#</th>
					<th scope="col">Producto</th>
					<th scope="col">Precio</th>
					<th scope="col">Subproductos</th>
					<th scope="col">Refrigerador</th>
					<th scope="col">Agregar</th>
				</tr>
			</thead>
			<tbody>`, "%%")

			for ar, v := range productos {
				var existenciasdepromo []int
				var nombredeproductosdepromo []string
				for _, vv := range v.Productos {
					existenciaRefri := inventariomodel.ExtraeExistencias(vv)
					nombredeproducto := inventariomodel.ExtraeNombreProducto(vv)
					nombredeproductosdepromo = append(nombredeproductosdepromo, nombredeproducto)
					existenciasdepromo = append(existenciasdepromo, existenciaRefri)
				}

				htmlcode += fmt.Sprintf(`
				<tr>
					<th scope="row">%v</th>
					<td>%v</td>
					<td>%v</td>`, ar+1, v.Nombre, v.PrecioPub)

				htmlcode += fmt.Sprintf(`<td>`)

				for ks, kv := range v.Productos {
					nombredeproducto := inventariomodel.ExtraeNombreProducto(kv)
					htmlcode += fmt.Sprintf(`
					%v,%v <br>`, nombredeproducto, v.Cantidades[ks])

				}

				htmlcode += fmt.Sprintf(`</td>`)

				htmlcode += fmt.Sprintf(`<td>`)

				for kd, vd := range nombredeproductosdepromo {
					htmlcode += fmt.Sprintf(`
						%v,%v <br>
					`, vd, existenciasdepromo[kd])
				}

				htmlcode += fmt.Sprintf(`</td>`)

				var sinexistencia bool

				for sr, sd := range v.Productos {
					existenciaRefri := inventariomodel.ExtraeExistencias(sd)
					if existenciaRefri < v.Cantidades[sr] {
						sinexistencia = true
					}

				}

				if sinexistencia {
					htmlcode += fmt.Sprintf(`
					<td>
						<a class="btn btn-danger btn-sm padd" href="Javascript:alert('Sin Existencias');" role="button">Sin Existencia</a>
					</td>
				</tr>`)
				} else {
					htmlcode += fmt.Sprintf(`
					<td>
						<a class="btn btn-success btn-sm padd" href="Javascript:ObtenDatosDeMesaDesdePromocion('%v:1:%v');" role="button">Agregar promo</a>
					</td>
				</tr>`, v.ID.Hex(), v.Categoria)

				}

			}

			htmlcode += fmt.Sprintf(`	
			</tbody>
		</table>
			`)
		}
	case "traspasos":
		htmlcode += fmt.Sprintf(`Traspasos`)

		break
	}

	ctx.HTML(htmlcode)

}

// ProductosAgregadosEnModalDesdePromo -> Evalua la mesa seleecionada y regresa al modal la cantidad de productos agregados
func ProductosAgregadosEnModalDesdePromo(ctx iris.Context) {

	idss := ctx.PostValue("data")
	cadenas := strings.Split(idss, ":")
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	mesa := ventamodel.ExtraeMesa(cadenas[0])
	primt, _ := primitive.ObjectIDFromHex(cadenas[1])
	producto := primt
	cantidad, _ := strconv.Atoi(cadenas[2])

	//tenemos la mesa , la promo y la cantidad la cual sera multiplicada por 1 por la cantidad de la promo si tiene 10 familiares seria 1*10 y eso se descontara del refrigerador
	var htmlcode string
	encontrado := false
	contador := 0
	for k, v := range mesa.Productos {
		if v == producto {
			encontrado = true
			contador = k
		}
	}

	if !encontrado {
		mesa.Productos = append(mesa.Productos, producto)
		mesa.Cantidades = append(mesa.Cantidades, cantidad)
	} else {
		mesa.Cantidades[contador] = mesa.Cantidades[contador] + cantidad
	}

	htmlcode += fmt.Sprintf(`
	<table class="table table-hover" style="margin: auto; width: 85%s !important; font-size:14px; text-align: center;">
	<input type="hidden" id="mesaenmodal" value="%v">
	<thead>
	  <tr>
		<th scope="col">#</th>
		<th scope="col">Producto</th>
		<th scope="col">Precio</th>
		<th scope="col">Cantidad</th>
		<th scope="col">Total</th>
	  </tr>
	</thead>
	<tbody>`, "%%", mesa.ID.Hex())

	totalsuma := 0.0
	for k, v := range mesa.Productos {

		productoenMesa := inventariomodel.ExtraeProducto(v.Hex())

		htmlcode += fmt.Sprintf(`
		
		<tr>
			<th scope="row">%v</th>
			<td>%v</td>
			<td>%v</td>
			<td>%v</td>
			<td>%v</td>
		</tr>`, k+1, productoenMesa.Nombre, ac.FormatMoney(productoenMesa.PrecioPub), mesa.Cantidades[k], ac.FormatMoney((productoenMesa.PrecioPub * float64(mesa.Cantidades[k]))))
		totalsuma += (productoenMesa.PrecioPub * float64(mesa.Cantidades[k]))
	}

	htmlcode += fmt.Sprintf(`
		<tr>
			<td colspan="4"></td>
			<td style="border-bottom: 2px solid #000000; text-align: center;">%v</td>
		</tr>
		</tbody>
	  </table>`, ac.FormatMoney(totalsuma))

	ventamodel.ActualizarMesaDiaria(mesa)

	ventamodel.ActuailzaAlmacenDesdeModalPromo(producto, cantidad)

	// 	htmlcode += fmt.Sprintf(`
	// 	<a class="btn btn-warning btn-large padd" href="Javascript:Traspasos('%v');" role="button">Traspasos</a>&nbsp;
	// 	<a class="btn btn-danger btn-large padd" href="Javascript:Cancelacion('%v');" role="button">Cancelar Mesa</a>&nbsp;
	// `, mesa.ID.Hex(), mesa.ID.Hex())

	ctx.HTML(htmlcode)

}

// ProductosVendidos -> Productos vendidos al momento de abrir el modal
func ProductosVendidos(ctx iris.Context) {

	idmesa := ctx.PostValue("data")
	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	var htmlcode string

	mesa := ventamodel.ExtraeMesa(idmesa)

	if len(mesa.Productos) == 0 {
		htmlcode += fmt.Sprintf(`
	<div class="container">
    	<div class="centrado">
        	<h1 class="display-4">Agrega productos a la mesa</h1>
			<input type="hidden" id="mesaenmodal" value="%v">
    	</div>
	</div>`, mesa.ID.Hex())

	} else {
		htmlcode += fmt.Sprintf(`
		<table class="table table-hover" style="margin: auto; width: 85%s !important; font-size:14px; text-align: center;">
		<input type="hidden" id="mesaenmodal" value="%v">
		<thead>
		  <tr>
			<th scope="col">#</th>
			<th scope="col">Producto</th>
			<th scope="col">Precio</th>
			<th scope="col">Cantidad</th>
			<th scope="col">Total</th>
		  </tr>
		</thead>
		<tbody>`, "%%", mesa.ID.Hex())

		totalsuma := 0.0
		for k, v := range mesa.Productos {

			productoenMesa := inventariomodel.ExtraeProducto(v.Hex())

			htmlcode += fmt.Sprintf(`
		
		<tr>
			<th scope="row">%v</th>
			<td>%v</td>
			<td>%v</td>
			<td>%v</td>
			<td>%v</td>
		</tr>`, k+1, productoenMesa.Nombre, ac.FormatMoney(productoenMesa.PrecioPub), mesa.Cantidades[k], ac.FormatMoney((productoenMesa.PrecioPub * float64(mesa.Cantidades[k]))))
			totalsuma += (productoenMesa.PrecioPub * float64(mesa.Cantidades[k]))
		}

		htmlcode += fmt.Sprintf(`
		<tr>
			<td colspan="4"></td>
			<td style="border-bottom: 2px solid #000000; text-align: center;">%v</td>
		</tr>
		</tbody>
	  </table>`, ac.FormatMoney(totalsuma))

		// 	htmlcode += fmt.Sprintf(`
		//   <div class="container centrado">
		//   <a class="btn btn-warning btn-large padd" href="Javascript:Traspasos('%v');" role="button">Traspasos</a>&nbsp;
		//   <a class="btn btn-danger btn-large padd" href="Javascript:Cancelacion('%v');" role="button">Cancelar Mesa</a>&nbsp;
		//   </div>`, mesa.ID.Hex(), mesa.ID.Hex())

	}

	ctx.HTML(htmlcode)

}

// ProductosAgregadosEnModal -> Evalua la Mesa seleccionada y regresa al modal la cantidad de productos agregados
func ProductosAgregadosEnModal(ctx iris.Context) {

	userOn := indexmodel.GetUserOn(sessioncontroller.Sess.Start(ctx).GetString("UserID"))
	log.Println("Usuario agregando producto", userOn.Nombre)
	idss := ctx.PostValue("data")
	cadenas := strings.Split(idss, ":")
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	mesa := ventamodel.ExtraeMesa(cadenas[0])
	primt, _ := primitive.ObjectIDFromHex(cadenas[1])
	producto := primt

	cantidad, _ := strconv.Atoi(cadenas[2])
	var htmlcode string
	encontrado := false
	contador := 0
	mesa.Mesero = userOn.Nombre

	for k, v := range mesa.Productos {
		if v == producto {
			encontrado = true
			contador = k
		}
	}

	if !encontrado {
		mesa.Productos = append(mesa.Productos, producto)
		mesa.Cantidades = append(mesa.Cantidades, cantidad)
	} else {
		mesa.Cantidades[contador] = mesa.Cantidades[contador] + cantidad
	}

	htmlcode += fmt.Sprintf(`
		<table class="table table-hover" style="margin: auto; width: 85%s !important; font-size:14px; text-align: center;">
		<input type="hidden" id="mesaenmodal" value="%v">
		<thead>
		  <tr>
			<th scope="col">#</th>
			<th scope="col">Producto</th>
			<th scope="col">Precio</th>
			<th scope="col">Cantidad</th>
			<th scope="col">Total</th>
		  </tr>
		</thead>
		<tbody>`, "%%", mesa.ID.Hex())

	totalsuma := 0.0
	for k, v := range mesa.Productos {

		productoenMesa := inventariomodel.ExtraeProducto(v.Hex())

		htmlcode += fmt.Sprintf(`
		
		<tr>
			<th scope="row">%v</th>
			<td>%v</td>
			<td>%v</td>
			<td>%v</td>
			<td>%v</td>
		</tr>`, k+1, productoenMesa.Nombre, ac.FormatMoney(productoenMesa.PrecioPub), mesa.Cantidades[k], ac.FormatMoney((productoenMesa.PrecioPub * float64(mesa.Cantidades[k]))))
		totalsuma += (productoenMesa.PrecioPub * float64(mesa.Cantidades[k]))
	}
	htmlcode += fmt.Sprintf(`
		<tr>
			<td colspan="4"></td>
			<td style="border-bottom: 2px solid #000000; text-align: center;">%v</td>
		</tr>
		</tbody>
	  </table>`, ac.FormatMoney(totalsuma))

	ventamodel.ActualizarMesaDiaria(mesa)
	ventamodel.ActuailzaAlmacenDesdeModal(producto, cantidad)
	// htmlcode += fmt.Sprintf(`
	// <div class="container centrado">
	// <a class="btn btn-warning btn-large padd" href="Javascript:Traspasos('%v');" role="button">Traspasos</a>&nbsp;
	// <a class="btn btn-danger btn-large padd" href="Javascript:Cancelacion('%v');" role="button">Cancelar Mesa</a>&nbsp;
	// </div>`, mesa.ID.Hex(), mesa.ID.Hex())

	ctx.HTML(htmlcode)

}

// FechaHTML Devuele la fecha formateada para html
func FechaHTML(fecha time.Time) string {
	fechahtml := fecha.Format("2006-01-02")
	// fechahtml = strings.Replace(fechahtml, "-", "/", -1)
	return fechahtml

}

// ImprimirVentaDiariaEnPDF2 -> Imprime el dia seleccionado en pdf listo para su impresion
func ImprimirVentaDiariaEnPDF2(ctx iris.Context) {

	var ventadiaria inventariomodel.VentaDiaria

	fecha := ctx.PostValue("data")

	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	layout := "2006-01-02"
	location, _ := time.LoadLocation("America/Mexico_City")

	t, _ := time.ParseInLocation(layout, fecha, location)

	mesasdiarias := ventamodel.BuscarVentasDiarias(t)

	var htmlcode string

	horafecha := time.Now()
	dia := horafecha.Day()
	mess := horafecha.Month().String()
	mes := MesEspanol(mess)
	anio := horafecha.Year()

	var grantotal float64

	for _, v := range mesasdiarias {

		if len(v.Productos) > 0 {

			for kk, vv := range v.Productos {
				ventadiaria.Productos = append(ventadiaria.Productos, vv)
				ventadiaria.Existencia = append(ventadiaria.Existencia, v.Cantidades[kk])
			}
			grantotal = +grantotal + v.GranTotal

		}

	}
	//Agrupar

	var agrupados []string
	var cantidades []int

	for k, v := range ventadiaria.Productos {

		if len(agrupados) == 0 {
			agrupados = append(agrupados, v.Hex())
			cantidades = append(cantidades, ventadiaria.Existencia[k])
		} else {

			encontrado := false
			contador := 0
			for kk, vv := range agrupados {
				if vv == v.Hex() {
					encontrado = true
					contador = kk
				}
			}

			if encontrado == false {
				agrupados = append(agrupados, v.Hex())
				cantidades = append(cantidades, ventadiaria.Existencia[k])
			} else {

				cantidades[contador] = cantidades[contador] + ventadiaria.Existencia[k]

			}

		}

	}

	diastring := strconv.Itoa(dia)
	aniostring := strconv.Itoa(anio)

	nArchivo := ctx.PostValue("data")

	pdf := gofpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 16)
	pdf.Cell(40, 7, tr("Bar La 49 - Venta del dia "+diastring+" de "+mes+" de "+aniostring+"."))
	pdf.Ln(2)
	pdf.SetFont("Times", "B", 14)
	pdf.Cell(40, 10, tr("------------------------------------------------------------------------------------------------------------------"))
	pdf.Ln(-1)
	pdf.SetFont("Courier", "B", 9)
	pdf.Cell(10, 10, tr("#"))
	pdf.Cell(90, 10, tr("Producto"))
	pdf.Cell(20, 10, tr("Cantd"))
	pdf.Cell(20, 10, tr("Precio"))
	pdf.Cell(25, 10, tr("Utild"))
	pdf.Cell(25, 10, tr("Total"))
	pdf.Ln(5)

	for k, v := range agrupados {
		product := inventariomodel.ExtraeProducto(v)
		kstring := strconv.Itoa(k + 1)
		cantidadesstring := strconv.Itoa(cantidades[k])
		totalpub := product.PrecioPub * float64(cantidades[k])
		totaluti := product.PrecioUti * float64(cantidades[k])

		pdf.Ln(5)
		pdf.Cell(10, 10, tr(kstring))
		pdf.Cell(90, 10, tr(product.Nombre))
		pdf.Cell(20, 10, tr(cantidadesstring))
		pdf.Cell(20, 10, tr(ac.FormatMoney(product.PrecioPub)))
		pdf.Cell(25, 10, tr(ac.FormatMoney(totaluti)))
		pdf.Cell(25, 10, tr(ac.FormatMoney(totalpub)))
	}
	pdf.Ln(10)
	pdf.Cell(10, 10, tr(""))
	pdf.Cell(90, 10, tr(""))
	pdf.Cell(20, 10, tr(""))
	pdf.Cell(20, 10, tr(""))
	pdf.Cell(25, 10, tr(""))
	pdf.Cell(25, 10, tr(ac.FormatMoney(grantotal)))

	fileee := `Recursos\Archivos\` + nArchivo + `.pdf`

	if len(agrupados) == 0 {

		htmlcode += fmt.Sprintf(`<div class="container centrado">
		<h1> No se ha encontrado venta para el día %v </h1>
		</div>`, nArchivo)
		fmt.Println("Agrupador esta en 0 -> ", len(agrupados))

	} else {
		err4 := pdf.OutputFileAndClose(fileee)
		if err4 != nil {
			fmt.Println(err4)
		}
		htmlcode += fmt.Sprintf(`<div class="container centrado">
		<h1> Archivo PDF creado con éxito</h1>
		</div>
		<div class="container centrado"><a class="btn btn-warning btn-large padd" href="Javascript:Descargar('%v');" role="button">Descargar Archivo %v</a></div>
		</div>`, nArchivo, nArchivo)

	}

	ctx.HTML(htmlcode)

}

// VaciarColeccion -> Elimina Colecciones de MongoDB
func VaciarColeccion(ctx iris.Context) {

	coleccion := ctx.PostValue("data")

	var htmlcode string

	switch coleccion {
	case "GG":

		htmlcode += fmt.Sprintf(`
		<script>
			alert("Base De DATOS ELIMINADA");
			location.replace("/administracion");
		</script>
		`)
		ventamodel.EliminarColeccionVentasDiarias()
		break

	}

	ctx.HTML(htmlcode)

}

// ImprimirVentaDiaria -> Devuelve la vista de la Venta Diaria
func ImprimirVentaDiaria(ctx iris.Context) {
	userOn := indexmodel.GetUserOn(sessioncontroller.Sess.Start(ctx).GetString("UserID"))
	ctx.ViewData("Usuario", userOn)

	//Enviar informacion de Numero de Mesas y la fecha del dia.

	// b := ventamodel.ExtraeMesas()

	// ctx.ViewData("Mesas", b)
	if err := ctx.View("ImprimirVentaDiaria.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

// ClasificaVenta -> Devuelve la clasifiacion de la venta diaria asi como los totales vendidos y la diferencia en precios
func ClasificaVenta(ctx iris.Context) {

	fecha := ctx.PostValue("data")
	var htmlcode string

	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	layout := "2006-01-02"
	location, _ := time.LoadLocation("America/Mexico_City")

	t, _ := time.ParseInLocation(layout, fecha, location)

	mesasdiarias := ventamodel.BuscarVentasDiarias(t)

	//Saber la venta del dia con una tabla

	var ventadiaria inventariomodel.VentaDiaria

	if len(mesasdiarias) == 0 {
		htmlcode += fmt.Sprintf(`<div class="container">
			<div class="centrado">
				<br>
				<h1 class="display-4">No se tiene ningun registro de %v</h1>
			</div>
		</div>`, t.Format(layout))
		ctx.HTML(htmlcode)
		return
	}

	htmlcode += fmt.Sprintf(`
		<table class="table table-hover" style="margin: auto; width: 85%s !important; font-size:14px; text-align: center;">
		<thead>
		  <tr>
			<th scope="col">#</th>
			<th scope="col">Producto</th>
			<th scope="col">Precio Pub</th>
			<th scope="col">Precio Uti</th>
			<th scope="col">Cantidad</th>
			<th scope="col">Utilidad</th>
			<th scope="col">Total</th>
		  </tr>
		</thead>
		<tbody>`, "%%")

	var grantotal float64

	fmt.Println("Gran total ->", grantotal)

	for _, v := range mesasdiarias {

		if len(v.Productos) > 0 {
			for kk, vv := range v.Productos {
				ventadiaria.Productos = append(ventadiaria.Productos, vv)
				ventadiaria.Existencia = append(ventadiaria.Existencia, v.Cantidades[kk])
			}
			grantotal += v.GranTotal
		}

	}

	//Agrupar
	var agrupados []string
	var cantidades []int

	for k, v := range ventadiaria.Productos {

		if len(agrupados) == 0 {
			agrupados = append(agrupados, v.Hex())
			cantidades = append(cantidades, ventadiaria.Existencia[k])
		} else {

			encontrado := false
			contador := 0
			for kk, vv := range agrupados {
				if vv == v.Hex() {
					encontrado = true
					contador = kk
				}
			}

			if encontrado == false {
				agrupados = append(agrupados, v.Hex())
				cantidades = append(cantidades, ventadiaria.Existencia[k])
			} else {

				cantidades[contador] = cantidades[contador] + ventadiaria.Existencia[k]

			}

		}

	}

	for k, v := range agrupados {

		product := inventariomodel.ExtraeProducto(v)

		htmlcode += fmt.Sprintf(`
		<tr>
			<th scope="row">%v</th>
			<td>%v</td>
			
			<td>%v</td>
			
			<td>%v</td>
			
			<td>%v</td>
			
			<td>%v</td>
			
			<td>%v</td>
		 </tr>`, k+1, product.Nombre, ac.FormatMoney(product.PrecioPub), ac.FormatMoney(product.PrecioUti), cantidades[k], ac.FormatMoney(product.PrecioUti*float64(cantidades[k])), ac.FormatMoney(product.PrecioPub*float64(cantidades[k])))

	}

	htmlcode += fmt.Sprintf(`
		<tr>
			<td colspan="6"></td>
			<td style="border-bottom: 2px solid #000000; text-align: center;">%v</td>
		</tr>
		</tbody>
	  </table>`, ac.FormatMoney(grantotal))

	ctx.HTML(htmlcode)

}

// MesEspanol Regresa el mes en español.
func MesEspanol(mes string) string {
	var mess string
	switch mes {
	case "January":
		mess = "Enero"
		break
	case "February":
		mess = "Febrero"
		break
	case "March":
		mess = "Marzo"
		break
	case "April":
		mess = "Abril"
		break
	case "May":
		mess = "Mayo"
		break
	case "June":
		mess = "Junio"
		break
	case "July":
		mess = "Julio"
		break
	case "August":
		mess = "Agosto"
		break
	case "September":
		mess = "Septiembre"
		break

	case "October":
		mess = "Octubre"
		break
	case "November":
		mess = "Noviembre"
		break
	case "December":
		mess = "Diciembre"
		break
	}
	return mess
}
