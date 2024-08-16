package inventariocontroller

import (
	"fmt"
	"io"
	"os"

	// "strings"
	"strconv"

	"github.com/kataras/iris/v12"
	sessioncontroller "github.com/vadgun/Bar/Controladores/SessionController"
	indexmodel "github.com/vadgun/Bar/Modelos/IndexModel"
	inventariomodel "github.com/vadgun/Bar/Modelos/InventarioModel"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Inventario -> Devuelve la vista del Inventario
func Inventario(ctx iris.Context) {
	fmt.Println("Prueba de invetario")
	userOn := indexmodel.GetUserOn(sessioncontroller.Sess.Start(ctx).GetString("UserID"))
	ctx.ViewData("Usuario", userOn)
	if err := ctx.View("Inventario.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

// Altaform -> Regresa la peticion del formulario de alta de productos
func Altaform(ctx iris.Context) {
	var htmlcode string

	info := ctx.PostValue("data")

	switch info {
	case "Almacenes":

		almacenesDisponibles := inventariomodel.ExtraerAlmacenes()

		// for k, v := range almacenesDisponibles {
		// 	fmt.Println("-> ", k)
		// 	htmlcode += fmt.Sprintf(`
		// 	%v
		// 	`, v.Nombre)
		// }

		htmlcode += fmt.Sprintf(`<hr>`)
		htmlcode += fmt.Sprintf(`<div class="container centrado">`)
		htmlcode += fmt.Sprintf(`<table class="table table-hover table-sm table-striped textito">`)
		htmlcode += fmt.Sprintf(`<tr>`)
		htmlcode += fmt.Sprintf(`
			
			<th class="textocentrado2">Nombre</th>
			<th class="textocentrado2">Acciones</th>`)
		htmlcode += fmt.Sprintf(`</tr>`)

		for _, v := range almacenesDisponibles {

			htmlcode += fmt.Sprintf(`<tr>`)
			htmlcode += fmt.Sprintf(`
			
			<td class="textocentrado2">%v</td>`, v.Nombre)

			htmlcode += fmt.Sprintf(`
			<td class="textocentrado2">
		<button class="btn-sm" title="Editar" onclick="EditarAlmacen('%v');">
			<img src="Recursos/Generales/Plugins/icons/build/svg/pencil.svg" height="15" alt="Editar"/>
		</button>
		</td>`, v.ID.Hex())

			htmlcode += fmt.Sprintf(`</tr>`)
		}

		htmlcode += fmt.Sprintf(`</div>`)

		break
	case "Alta":

		htmlcode += fmt.Sprintf(`
		<form method="POST" enctype="multipart/form-data" action="/altaproducto" name="altaproducto" id="altaproducto">

        <div class="col-12">
            <h6 class="border-bottoms-c"> Datos del producto: </h6>
            <div class="form-group row">
                <label for="nombre" class="col-sm-2 col-form-label negrita"> Nombre del producto : </label>
                <div class="col-sm-4">
                    <input type="text" class="form-control" id="nombre" name="nombre" value="" placeholder="Nombre del producto" required>
                </div>
            </div>

            <div class="form-group row">
                <label for="precio" class="col-sm-2 col-form-label negrita"> Precio al publico: </label>
                <div class="col-sm-4">
                    <input type="number" class="form-control" id="precio" name="precio" value="" placeholder="Precio del producto" required>
                </div>
            </div>

            <div class="form-group row">
                <label for="preciouti" class="col-sm-2 col-form-label negrita"> Precio utilidad: </label>
                <div class="col-sm-4">
                    <input type="number" class="form-control" id="preciouti" name="preciouti" value="" placeholder="Precio de utilidad" required>
                </div>
            </div>


            <div class="form-group row">
                <label for="categoria" class="col-sm-2 col-form-label negrita"> Categoria: </label>
                <div class="col-sm-4">
                    <select id="categoria" name="categoria" class="form-control" onchange="Javascript:ConvierteAPromo(this.value);">
                    <option value="">Selecciona Categoria</option>
                    <option value="botana">Botana</option>
					<option value="botella">Botella</option>
					<option value="cancion">Cancion</option>
                    <option value="cerveza">Cerveza</option>
					<option value="cigarros">Cigarros</option>
					<option value="copa">Copa</option>
                    <option value="ficha">Ficha</option>
                    <option value="promo">Promo</option>
                </select>
                </div>
            </div>
            <div class="form-group row" id="conversionpromo">

            </div>

            <div class="form-group row" id="promoproductos">

            </div>

            <div class="form-group row">
                <label for="imagenproducto" class="col-sm-2 col-form-label negrita"> Imagen del producto: </label>
                <div class="col-sm-4">
                    <input type="file" class="form-control" id="imagenproducto" name="imagenproducto">
                    <input type="hidden" value="" name="articulosagregados" id="articulosagregados">
                </div>
            </div>

            <div class="form-group row">
                <div class="col-sm-4">
                    <button type="submit" class="btn btn-success">Guardar</button>
                </div>
            </div>
        </div>
    </form>
		`)

		break
	case "botellas":

		var botellas []inventariomodel.Producto

		htmlcode += fmt.Sprintf(`
		<hr>
		<div class="container">
		<div class="row">
		`)

		botellas = inventariomodel.ExtraeBotellas()

		if len(botellas) == 0 {
			htmlcode += fmt.Sprintf(`<div class="centrado">
				<h1 class="display-4">No hay botellas dadas de alta</h1>
			</div>`)
		} else {

			for _, v := range botellas {

				imagenproducto := inventariomodel.TraerImagenActa(v.Imagen)

				htmlcode += fmt.Sprintf(`
			<div class=" col-sm-3 ">
				<img src="data:image/jpg;base64,%v" alt="... " class="img-thumbnail " width='50%%'>
				<h5>%v</h5>
				<h6>%v pesos</h6>
			</div>`, imagenproducto, v.Nombre, v.PrecioPub)
			}
		}

		htmlcode += fmt.Sprintf(`
			</div>
		</div>`)

		break
	case "cervezas":
		var cervezas []inventariomodel.Producto

		htmlcode += fmt.Sprintf(`
		<hr>
		<div class="container">
		<div class="row">
		`)

		cervezas = inventariomodel.ExtraeCervezas()

		if len(cervezas) == 0 {
			htmlcode += fmt.Sprintf(`<div class="centrado">
				<h1 class="display-4">No hay cervezas dadas de alta</h1>
			</div>`)
		} else {

			for _, v := range cervezas {
				imagenproducto := inventariomodel.TraerImagenActa(v.Imagen)

				htmlcode += fmt.Sprintf(`
			<div class=" col-sm-3 ">
				<img src="data:image/jpg;base64,%v" alt="... " class="img-thumbnail " width='50%%'>
				<h5>%v</h5>
				<h6>%v pesos</h6>
			</div>`, imagenproducto, v.Nombre, v.PrecioPub)
			}
		}

		htmlcode += fmt.Sprintf(`
			</div>
		</div>`)

		break
	case "botanas":
		var botanas []inventariomodel.Producto

		htmlcode += fmt.Sprintf(`
		<hr>
		<div class="container">
		<div class="row">
		`)

		botanas = inventariomodel.ExtraeBotanas()
		if len(botanas) == 0 {
			htmlcode += fmt.Sprintf(`<div class="centrado">
				<h1 class="display-4">No hay botanas dadas de alta</h1>
			</div>`)
		} else {

			for _, v := range botanas {

				imagenproducto := inventariomodel.TraerImagenActa(v.Imagen)

				htmlcode += fmt.Sprintf(`
			<div class=" col-sm-3 ">
				<img src="data:image/jpg;base64,%v" alt="... " class="img-thumbnail " width='50%%'>
				<h5>%v</h5>
				<h6>%v pesos</h6>
			</div>`, imagenproducto, v.Nombre, v.PrecioPub)
			}
		}
		htmlcode += fmt.Sprintf(`
			</div>
		</div>`)

		break
	case "cancion":
		var cancion []inventariomodel.Producto

		htmlcode += fmt.Sprintf(`
		<hr>
		<div class="container">
		<div class="row">
		`)

		cancion = inventariomodel.ExtraeCancion()

		if len(cancion) == 0 {
			htmlcode += fmt.Sprintf(`<div class="centrado">
				<h1 class="display-4">No hay canciones dadas de alta</h1>
			</div>`)
		} else {

			for _, v := range cancion {

				imagenproducto := inventariomodel.TraerImagenActa(v.Imagen)

				htmlcode += fmt.Sprintf(`
			<div class=" col-sm-3 ">
				<img src="data:image/jpg;base64,%v" alt="... " class="img-thumbnail " width='50%%'>
				<h5>%v</h5>
				<h6>%v pesos</h6>
			</div>`, imagenproducto, v.Nombre, v.PrecioPub)
			}
		}

		htmlcode += fmt.Sprintf(`
			</div>
		</div>`)
		break

	case "ficha":
		var fichas []inventariomodel.Producto

		htmlcode += fmt.Sprintf(`
	<hr>
	<div class="container">
	<div class="row">
	`)

		fichas = inventariomodel.ExtraeFichas()
		if len(fichas) == 0 {
			htmlcode += fmt.Sprintf(`<div class="centrado">
				<h1 class="display-4">No hay fichas dadas de alta</h1>
			</div>`)
		} else {

			for _, v := range fichas {

				imagenproducto := inventariomodel.TraerImagenActa(v.Imagen)

				htmlcode += fmt.Sprintf(`
		<div class=" col-sm-3 ">
			<img src="data:image/jpg;base64,%v" alt="... " class="img-thumbnail " width='50%%'>
			<h5>%v</h5>
			<h6>%v pesos</h6>
		</div>`, imagenproducto, v.Nombre, v.PrecioPub)
			}
		}

		htmlcode += fmt.Sprintf(`
		</div>
	</div>`)
		break
	case "promo":
		var promos []inventariomodel.Producto

		htmlcode += fmt.Sprintf(`
	<hr>
	<div class="container">
	<div class="row">
	`)

		promos = inventariomodel.ExtraePromos()
		if len(promos) == 0 {
			htmlcode += fmt.Sprintf(`<div class="centrado">
				<h1 class="display-4">No hay promociones dadas de alta</h1>
			</div>`)
		} else {

			for _, v := range promos {

				imagenproducto := inventariomodel.TraerImagenActa(v.Imagen)

				htmlcode += fmt.Sprintf(`
		<div class=" col-sm-3 ">
			<img src="data:image/jpg;base64,%v" alt="... " class="img-thumbnail " width='50%%'>
			<h5>%v</h5>
			<h6>%v pesos</h6>
		</div>`, imagenproducto, v.Nombre, v.PrecioPub)
			}
		}

		htmlcode += fmt.Sprintf(`
		</div>
	</div>`)
		break

	case "copa":
		var copas []inventariomodel.Producto

		htmlcode += fmt.Sprintf(`
	<hr>
	<div class="container">
	<div class="row">
	`)

		copas = inventariomodel.ExtraeCopas()
		if len(copas) == 0 {
			htmlcode += fmt.Sprintf(`<div class="centrado">
				<h1 class="display-4">No hay copas dadas de alta</h1>
			</div>`)
		} else {

			for _, v := range copas {

				imagenproducto := inventariomodel.TraerImagenActa(v.Imagen)

				htmlcode += fmt.Sprintf(`
		<div class=" col-sm-3 ">
			<img src="data:image/jpg;base64,%v" alt="... " class="img-thumbnail " width='50%%'>
			<h5>%v</h5>
			<h6>%v pesos</h6>
		</div>`, imagenproducto, v.Nombre, v.PrecioPub)
			}
		}

		htmlcode += fmt.Sprintf(`
		</div>
	</div>`)

		break

	case "cigarros":
		var cigarros []inventariomodel.Producto

		htmlcode += fmt.Sprintf(`
	<hr>
	<div class="container">
	<div class="row">
	`)

		cigarros = inventariomodel.ExtraeCigarros()

		if len(cigarros) == 0 {
			htmlcode += fmt.Sprintf(`<div class="centrado">
				<h1 class="display-4">No hay cigarros dados de alta</h1>
			</div>`)
		} else {

			for _, v := range cigarros {

				imagenproducto := inventariomodel.TraerImagenActa(v.Imagen)

				htmlcode += fmt.Sprintf(`
		<div class=" col-sm-3 ">
			<img src="data:image/jpg;base64,%v" alt="... " class="img-thumbnail " width='50%%'>
			<h5>%v</h5>
			<h6>%v pesos</h6>
		</div>`, imagenproducto, v.Nombre, v.PrecioPub)
			}
		}

		htmlcode += fmt.Sprintf(`
		</div>
	</div>`)
		break
	case "Traspasos":
		var todosLosProductos []inventariomodel.Producto

		todosLosProductos = inventariomodel.TodosLosProductos()

		// htmlcode += fmt.Sprintf(`Tabla AQUI`)

		htmlcode += fmt.Sprintf(`
		<br>
		<hr>
		<table class="table table-hover table-bordered table-lg" style="margin: auto; width: 50%s !important; font-size:14px;">
 	     <thead>
    	    <th class="textocentrado">
        	  Nombre
        	</th>
        	<th class="textocentrado">
          	  Precio Publico
        	</th>
        	<th class="textocentrado">
          	  Precio Utilidad
			</th>
			<th class="textocentrado">
			Imagen
		  	</th>	
        	<th class="textocentrado">
              Acciones
        	</th>
 	  	 </thead>
		  <tbody>`, "%%")

		var bodega string
		var refrigerador string

		bodega = inventariomodel.ExtraeBodegaID()
		refrigerador = inventariomodel.ExtraeRefrigeradorID()

		for _, v := range todosLosProductos {
			imagenproducto := inventariomodel.TraerImagenActa(v.Imagen)

			htmlcode += fmt.Sprintf(`
				<tr>
				<td class="textocentrado">
					%v
				</td>
				<td class="textocentrado">
					%v
				</td>
				<td class="textocentrado">
					%v
				</td>
				<td class="textocentrado">
					<img src="data:image/jpg;base64,%v" alt="... " class="img-thumbnail " width='100%%'>
				</td>`, v.Nombre, v.PrecioPub, v.PrecioUti, imagenproducto)

			htmlcode += fmt.Sprintf(`
				<td class="textocentrado">
					<button class="btn-sm" title="Editar Producto" onclick="EditarProducto('%v');">
						<img src="Recursos/Generales/Plugins/icons/build/svg/pencil.svg" height="15" alt="Editar Producto"/>
					</button>
	
					<button class="btn-sm" title="Eliminar Producto" onclick="EliminarProducto('%v');">
					<img src="Recursos/Generales/Plugins/icons/build/svg/trashcan.svg" height="15" alt="Eliminar Producto"/>
				</button>`, v.ID.Hex(), v.ID.Hex())

			if v.Categoria != "promo" {
				htmlcode += fmt.Sprintf(`
					<button class="btn-sm" title="Enviar a Bodega" onclick="EnviarBodega('%v:%v');">
					<img src="Recursos/Generales/Plugins/icons/build/svg/rocket.svg" height="15" alt="Enviar a Bodega"/>
				</button>
	
				<button class="btn-sm" title="Enviar a Refrigerador" onclick="EnviarRefrigerador('%v:%v');">
					<img src="Recursos/Generales/Plugins/icons/build/svg/rocket.svg" height="15" alt="Enviar a Refrigerador"/>
				</button>
	
				<button class="btn-sm" title="Eliminar de Bodega" onclick="EliminarDeBodega('%v:%v');">
				<img src="Recursos/Generales/Plugins/icons/build/svg/trashcan.svg" height="15" alt="Enviar a Bodega"/>
			</button>
	
			<button class="btn-sm" title="Eliminar de Refrigerador" onclick="EliminarDeRefrigerador('%v:%v');">
			<img src="Recursos/Generales/Plugins/icons/build/svg/trashcan.svg" height="15" alt="Enviar a Refrigerador"/>
		</button>
	
	
				</td>
				  
				</tr>`, v.ID.Hex(), bodega, v.ID.Hex(), refrigerador, v.ID.Hex(), bodega, v.ID.Hex(), refrigerador)

			} else {
				htmlcode += fmt.Sprintf(`
				</td>
				</tr>
				`)

			}

		}

		htmlcode += fmt.Sprintf(`
		</tbody>
		</table>
	  `)

		break
	}

	ctx.HTML(htmlcode)
}

func check(err error, mensaje string) {
	if err != nil {
		panic(err)
	}
}

// GuardarProducto -> recibe el
func GuardarProducto(ctx iris.Context) {

	var producto inventariomodel.Producto

	var htmlcode string

	producto.ID = primitive.NewObjectID()

	producto.Nombre = ctx.PostValue("nombre")
	producto.PrecioPub, _ = ctx.PostValueFloat64("precio")
	producto.PrecioUti, _ = ctx.PostValueFloat64("preciouti")
	producto.Categoria = ctx.PostValue("categoria")
	producto.Imagen = primitive.NewObjectID()

	// var productosstring []string
	// var cantidadstring []string

	if ctx.PostValue("categoria") == "promo" {
		numarticulos, _ := ctx.PostValueInt("articulosagregados")
		for i := 0; i < numarticulos; i++ {
			istring := strconv.Itoa(i)
			cantidadint, _ := ctx.PostValueInt("cantidadprod" + istring)
			prodid := ctx.PostValue("productoid" + istring)
			primt, _ := primitive.ObjectIDFromHex(prodid)
			producto.Productos = append(producto.Productos, primt)
			producto.Cantidades = append(producto.Cantidades, cantidadint)
		}
	}

	imagenacta, header, err := ctx.FormFile("imagenproducto")
	nombrearchivo := header.Filename
	// fmt.Println("Que trae el valor POST NAME--", nombrearchivo, "--")
	// fmt.Println("Que trae el valor POST ACTA--", imagenacta, "--")
	check(err, "Error al seleccionar la imagen")

	dirPath := "./Recursos/Imagenes/Productos"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Println("el directorio no existe")
		os.MkdirAll(dirPath, 0777)
	} else {
		fmt.Println("el directorio ya existe")
	}
	out, err := os.Create("./Recursos/Imagenes/Productos/" + nombrearchivo)
	check(err, "Unable to create the file for writing. Check your write access privilege")
	defer out.Close()
	_, err = io.Copy(out, imagenacta)
	check(err, "Error al escribir la imagen al directorio")

	//Atrapar los Datos Enviados por el formulario de alta para capturar la imagen

	if inventariomodel.GuardaProducto(producto) {

		idsImg := inventariomodel.UploadImageToMongo(dirPath, nombrearchivo)
		inventariomodel.UpdateImgArt(idsImg, producto.ID)

		htmlcode += fmt.Sprintf(`
		<script>
			alert("Producto Guardado");
			location.replace("/inventario");
		</script>
		`)
	} else {
		htmlcode += fmt.Sprintf(`
		<script>
			alert("Producto No Guardado");
			location.replace("/inventario");
		</script>
		`)
	}

	ctx.HTML(htmlcode)

}

// EditarAlmacen -> Regresa la informacion del almacen seleccionado para agregar o editar sus productos y/o existencias
func EditarAlmacen(ctx iris.Context) {
	var htmlcode string

	info := ctx.PostValue("data")

	//Crear una interfaz para editar el almacen puede ser un formulario de captura de productos

	//traer el almacen correspondiente☺

	almacenEditar := inventariomodel.ExtraerAlmacen(info)

	htmlcode += fmt.Sprintf(`
		<h4> %v </h4>
	`, almacenEditar.Nombre)

	htmlcode += fmt.Sprintf(`
	<form method="POST" action="/actualizaralmacen">
			<div class="col-12">
					<div class="form-group row">
							<label for="nombrealmacen" class="col-sm-2 col-form-label negrita">Nombre:</label>
							<div class="col-sm-2">
								<input type="text" class="form-control" id="nombrealmacen" name="nombrealmacen" value="%v" readonly>
							</div>
					</div>
	`, almacenEditar.Nombre)

	//Se iteran los productos almacenados en el almacen, pero esto solo serviria para los productos ya agregados, como agregar mas?!!!
	//	htmlcode += fmt.Sprintf(``)

	htmlcode += fmt.Sprintf(`
			<div class="form-group row">
				<label for="productos" class="col-sm-2 col-form-label negrita">Productos:</label>
				

		`)

	for k, v := range almacenEditar.Productos {

		nombredelproducto := inventariomodel.ExtraeNombreProducto(v)

		if k == 0 {
			htmlcode += fmt.Sprintf(`
			<div class="col-sm-4">
				<input type="text" class="form-control" id="producto%v" name="producto%v" value="%v" readonly>
				<input type="hidden" value="%v" name="id%v">
			</div>
			`, k, k, nombredelproducto, v.Hex(), k)

			htmlcode += fmt.Sprintf(`
			<div class="col-sm-4">
				<input type="text" class="form-control" id="existencias%v" name="existencias%v" value="%v" >
			</div>
			`, k, k, almacenEditar.Existencia[k])
		} else {
			htmlcode += fmt.Sprintf(`
			<div class="col-sm-4 offset-sm-2">
				<input type="text" class="form-control" id="producto%v" name="producto%v" value="%v" readonly>
				<input type="hidden" value="%v" name="id%v">
			</div>
			`, k, k, nombredelproducto, v.Hex(), k)

			htmlcode += fmt.Sprintf(`
			<div class="col-sm-4">
				<input type="text" class="form-control" id="existencias%v" name="existencias%v" value="%v" >
			</div>
			`, k, k, almacenEditar.Existencia[k])

		}

	}

	htmlcode += fmt.Sprintf(`
	</div>`)

	htmlcode += fmt.Sprintf(`
					<div class="centrado">
						<div class="form-group row">
							<input type="hidden" value="%v" name="longitud">
							<input type="hidden" value="%v" name="almaceneditar">
							<button type="submit"class="btn btn-success">Actualizar Almacen</button>					
						</div>
					</div>
			</div>
		</form>`, len(almacenEditar.Productos), almacenEditar.ID.Hex())

	ctx.HTML(htmlcode)

}

// EditandoProducto -> Edita el producto seleccionado
func EditandoProducto(ctx iris.Context) {

	var htmlcode string

	idprod := ctx.PostValue("idprod")
	nombreproducto := ctx.PostValue("nombre")
	preciopub, _ := ctx.PostValueFloat64("precio")
	preciouti, _ := ctx.PostValueFloat64("preciouti")
	categoria := ctx.PostValue("categoria")

	var producto inventariomodel.Producto

	primt, _ := primitive.ObjectIDFromHex(idprod)
	producto.ID = primt
	producto.Nombre = nombreproducto
	producto.PrecioPub = preciopub
	producto.PrecioUti = preciouti
	producto.Categoria = categoria

	if categoria == "promo" {
		promo := inventariomodel.ExtraeProducto(idprod)
		producto.Productos = promo.Productos
		producto.Cantidades = promo.Cantidades
	}

	inventariomodel.GuardaProductoEditado(producto)

	htmlcode += fmt.Sprintf(`
	<script>
		alert("Producto Actualizado");
		location.replace("/inventario");
	</script>
	`)

	ctx.HTML(htmlcode)

}

// EditarProducto -> Edita el producto seleccionado
func EditarProducto(ctx iris.Context) {

	var htmlcode string

	productoidhex := ctx.PostValue("data")

	producto := inventariomodel.ExtraeProducto(productoidhex)

	htmlcode += fmt.Sprintf(`
	<form method="POST" enctype="multipart/form-data" action="/editandoproducto" name="altaproducto" id="altaproducto">

	<div class="col-12">
		<h6 class="border-bottoms-c"> Datos del producto: </h6>
		<div class="form-group row">
			<label for="nombre" class="col-sm-2 col-form-label negrita"> Nombre del producto : </label>
			<div class="col-sm-4">
				<input type="text" class="form-control" id="nombre" name="nombre" value="%v" placeholder="Nombre del producto" required>
				<input type="hidden" value="%v" name="idprod">
			</div>
		</div>`, producto.Nombre, producto.ID.Hex())

	htmlcode += fmt.Sprintf(`
		<div class="form-group row">
			<label for="precio" class="col-sm-2 col-form-label negrita"> Precio al publico: </label>
			<div class="col-sm-4">
				<input type="number" class="form-control" id="precio" name="precio" value="%v" placeholder="Precio del producto" required>
			</div>
		</div>`, producto.PrecioPub)

	htmlcode += fmt.Sprintf(`
		<div class="form-group row">
		<label for="preciouti" class="col-sm-2 col-form-label negrita"> Precio utilidad: </label>
		<div class="col-sm-4">
			<input type="number" class="form-control" id="preciouti" name="preciouti" value="%v" placeholder="Precio de utilidad" required>
		</div>
		</div>`, producto.PrecioUti)

	htmlcode += fmt.Sprintf(`
		<div class="form-group row">
			<label for="categoria" class="col-sm-2 col-form-label negrita"> Categoria: </label>
			<div class="col-sm-4">
				<select id="categoria" name="categoria" class="form-control">`)

	switch producto.Categoria {
	case "botana":
		htmlcode += fmt.Sprintf(`
			<option value="">Selecciona Categoria</option>
			<option value="botana" selected>Botana</option>
			<option value="botella">Botella</option>
			<option value="cancion">Cancion</option>
			<option value="cerveza">Cerveza</option>
			<option value="cigarros">Cigarros</option>
			<option value="copa">Copa</option>
			<option value="ficha">Ficha</option>`)
		break
	case "botella":
		htmlcode += fmt.Sprintf(`
			<option value="">Selecciona Categoria</option>
			<option value="botana" >Botana</option>
			<option value="botella" selected>Botella</option>
			<option value="cancion">Cancion</option>
			<option value="cerveza">Cerveza</option>
			<option value="cigarros">Cigarros</option>
			<option value="copa">Copa</option>
			<option value="ficha">Ficha</option>`)
		break
	case "cancion":
		htmlcode += fmt.Sprintf(`
			<option value="">Selecciona Categoria</option>
			<option value="botana" >Botana</option>
			<option value="botella" >Botella</option>
			<option value="cancion" selected>Cancion</option>
			<option value="cerveza">Cerveza</option>
			<option value="cigarros">Cigarros</option>
			<option value="copa">Copa</option>
			<option value="ficha">Ficha</option>`)
		break
	case "cerveza":
		htmlcode += fmt.Sprintf(`
			<option value="">Selecciona Categoria</option>
			<option value="botana">Botana</option>
			<option value="botella">Botella</option>
			<option value="cancion">Cancion</option>
			<option value="cerveza" selected>Cerveza</option>
			<option value="cigarros">Cigarros</option>
			<option value="copa">Copa</option>
			<option value="ficha">Ficha</option>`)
		break
	case "cigarros":
		htmlcode += fmt.Sprintf(`
			<option value="">Selecciona Categoria</option>
			<option value="botana">Botana</option>
			<option value="botella">Botella</option>
			<option value="cancion">Cancion</option>
			<option value="cerveza">Cerveza</option>
			<option value="cigarros" selected>Cigarros</option>
			<option value="copa">Copa</option>
			<option value="ficha">Ficha</option>`)
		break
	case "copa":
		htmlcode += fmt.Sprintf(`
			<option value="">Selecciona Categoria</option>
			<option value="botana">Botana</option>
			<option value="botella">Botella</option>
			<option value="cancion">Cancion</option>
			<option value="cerveza">Cerveza</option>
			<option value="cigarros">Cigarros</option>
			<option value="copa" selected>Copa</option>
			<option value="ficha">Ficha</option>`)
		break
	case "ficha":
		htmlcode += fmt.Sprintf(`
			<option value="">Selecciona Categoria</option>
			<option value="botana">Botana</option>
			<option value="botella">Botella</option>
			<option value="cancion">Cancion</option>
			<option value="cerveza">Cerveza</option>
			<option value="cigarros">Cigarros</option>
			<option value="copa">Copa</option>
			<option value="ficha" selected>Ficha</option>`)
		break
	case "promo":
		htmlcode += fmt.Sprintf(`
		<option value="">Selecciona Categoria</option>
		<option value="botana">Botana</option>
		<option value="botella">Botella</option>
		<option value="cancion">Cancion</option>
		<option value="cerveza">Cerveza</option>
		<option value="cigarros">Cigarros</option>
		<option value="copa">Copa</option>
		<option value="ficha">Ficha</option>
		<option value="promo" selected>Promo</option>`)
		break
	}

	htmlcode += fmt.Sprintf(`
			</select>
			</div>
		</div>`)

	// 	<div class="form-group row">
	// 		<label for="imagenproducto" class="col-sm-2 col-form-label negrita"> Imagen del producto: </label>
	// 		<div class="col-sm-4">
	// 			<input type="file" class="form-control" id="imagenproducto" name="imagenproducto">
	// 		</div>
	// 	</div>`)

	htmlcode += fmt.Sprintf(`
		<div class="form-group row">
			<div class="col-sm-4">
				<button type="submit" class="btn btn-success">Actualizar Producto</button>
			</div>
		</div>
	</div>
</form>
	`)

	ctx.HTML(htmlcode)

}

// EliminarProducto -> Elimina el producto siempre y cuando no exista en los almacenes, para eliminarlo debera eliminarlo primero de los almacenes correspondientes
func EliminarProducto(ctx iris.Context) {

	//Verificar si el producto tiene existencia en almacen, si tiene no se puede eliminar, pedir al usuario q elimine la existencia antes de eliminar

	var htmlcode string

	producto := ctx.PostValue("data")

	if inventariomodel.EliminarProducto(producto) {
		htmlcode += fmt.Sprintf(`
		<script>
			alert("Producto Eliminado");
			location.replace("/inventario");
		</script>
		`)

	} else {
		htmlcode += fmt.Sprintf(`
		<script>
			alert("Producto encontrado en almacen");
			location.replace("/inventario");
		</script>
		`)
	}

	ctx.HTML(htmlcode)

}

// ActualizarAlmacen -> Recibe la peticion del arreglo para modificar el almacen
func ActualizarAlmacen(ctx iris.Context) {

	longitud, _ := ctx.PostValueInt("longitud")
	almaceneditar := ctx.PostValue("almaceneditar")

	var htmlcode string
	var ids []string
	var existencias []string

	for i := 0; i <= longitud; i++ {

		istring := strconv.Itoa(i)
		ids = append(ids, ctx.PostValue("id"+istring))
		existencias = append(existencias, ctx.PostValue("existencias"+istring))

	}

	//Almacenar correctamente los productos y las existencias en el almacen correspondiente

	if inventariomodel.ActualizarExistenciasEnAlmacen(almaceneditar, ids, existencias) {
		htmlcode += fmt.Sprintf(`
		<script>
			alert("Almacen Actualizado");
			location.replace("/inventario");
		</script>
		`)
	}

	ctx.HTML(htmlcode)

}

// ProductosASelect -> Regresa una lista de productos
func ProductosASelect(ctx iris.Context) {
	var htmlcode string

	productos := inventariomodel.ExtraeProductosSinPromo()

	htmlcode += fmt.Sprintf(`<option value="">Selecciona Articulo</option>`)

	for _, v := range productos {
		htmlcode += fmt.Sprintf(`
		<option value="%v">%v</option>
		`, v.ID.Hex(), v.Nombre)
	}

	ctx.HTML(htmlcode)

}
