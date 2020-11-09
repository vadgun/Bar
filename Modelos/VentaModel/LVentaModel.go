package ventamodel

import (
	"fmt"
	"time"

	conexiones "github.com/vadgun/Bar/Conexiones"
	inventariomodel "github.com/vadgun/Bar/Modelos/InventarioModel"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"log"
)

//ExtraeMesas -> Extrae las mesas dadas de alta
func ExtraeMesas() int {
	var mesas ConfigurarMesas
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_CFG)
	err1 := c.Find(bson.M{"Configuracion": "Mesas"}).One(&mesas)
	if err1 != nil {
		fmt.Println(err1)
	}
	return mesas.Disponibles

}

//ActualizarMesasDiarias -> Actualizara el numero de mesas para crear diariamente
func ActualizarMesasDiarias(numero int) {

	var mesas ConfigurarMesas
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_CFG)
	err1 := c.Find(bson.M{"Configuracion": "Mesas"}).One(&mesas)
	if err1 != nil {
		fmt.Println(err1)
	}
	mesas.Disponibles = numero

	errx := c.UpdateId(mesas.ID, mesas)
	if errx != nil {
		fmt.Println(errx)
	}

}

//GuardaMesaDiaria -> Crea la mesa diaria con fecha del dia, numero de mesa y status desocupado
func GuardaMesaDiaria(mesa Mesa) {

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	mesa.ID = bson.NewObjectId()

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_MD)
	errx := c.Insert(mesa)
	if errx != nil {
		fmt.Println(errx)
	}

}

//BuscarVentasDiarias -> Busca por dia las mesas diarias
func BuscarVentasDiarias(fecha time.Time) []Mesa {
	var mesasdiarias []Mesa

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_MD)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	errx := c.Find(bson.M{"Fecha": fecha}).Sort("Mesa").All(&mesasdiarias)
	if errx != nil {
		fmt.Println(errx)
	}

	return mesasdiarias
}

//ExtraeMesa -> Regresa la mesa solicitada
func ExtraeMesa(idmesa string) Mesa {

	idobjmesa := bson.ObjectIdHex(idmesa)

	var mesa Mesa

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_MD)
	err1 := c.FindId(idobjmesa).One(&mesa)
	if err1 != nil {
		fmt.Println(err1)
	}

	return mesa

}

//CierraMesa -> Cambia los status de la mesa sin borrar el registro para la venta diaria y el gran total que se esta buscando
func CierraMesa(mesaX Mesa) {

	var mesa Mesa

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_MD)
	err1 := c.FindId(mesaX.ID).One(&mesa)
	if err1 != nil {
		fmt.Println(err1)
	}

	mesa.Estatus = false
	mesa.Abierta = false
	mesa.Cerrada = true

	err2 := c.UpdateId(mesa.ID, mesa)
	if err2 != nil {
		fmt.Println(err2)
	}

	var cantidades []int
	var objectsids []bson.ObjectId

	var nuevamesa Mesa
	nuevamesa.ID = bson.NewObjectId()
	nuevamesa.Fecha = mesa.Fecha
	nuevamesa.Mesa = mesa.Mesa
	nuevamesa.Abierta = true
	nuevamesa.Cerrada = false
	nuevamesa.Estatus = false
	nuevamesa.Productos = objectsids
	nuevamesa.Cantidades = cantidades
	nuevamesa.GranTotal = 0

	err3 := c.Insert(nuevamesa)
	if err3 != nil {
		fmt.Println(err3)
	}

}

//EliminarColeccionVentasDiarias -> Elimina la Base de Datos de la coleccion de mesas diarias.
func EliminarColeccionVentasDiarias() {
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_MD)
	err1 := c.DropCollection()
	if err1 != nil {
		fmt.Println(err1)
	}

}

//ActualizarMesaDiaria -> Actualiza la mesa diaria
func ActualizarMesaDiaria(mesa Mesa) {
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	var sumatotal float64
	for k, v := range mesa.Cantidades {

		precio := PrecioProducto(mesa.Productos[k])
		sumatotal += (float64(v) * precio)

		fmt.Println("Suma ->", sumatotal, "   precio ->", precio, " v ->", v)
	}

	mesa.Estatus = true
	mesa.GranTotal = sumatotal

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_MD)
	err1 := c.UpdateId(mesa.ID, mesa)
	if err1 != nil {
		fmt.Println(err1)
	}
}

//ActuailzaAlmacenDesdeModalPromo -> Actualiza la cantidad de existencia en el Almacen Refrigerador
func ActuailzaAlmacenDesdeModalPromo(producto bson.ObjectId, cantidad int) {

	var almacen inventariomodel.Almacen
	var promo inventariomodel.Producto

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_AL)
	errx := c.Find(bson.M{"Nombre": "Refrigerador"}).One(&almacen)
	if errx != nil {
		fmt.Println(errx)
	}

	//Extrae la Promo

	d := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	errd := d.FindId(producto).One(&promo)
	if errd != nil {
		fmt.Println(errx)
	}

	var encontrado bool
	var contador int
	var contador2 int
	for k, v := range almacen.Productos {
		encontrado = false
		contador = 0
		contador2 = 0
		for kk, vv := range promo.Productos {

			if vv == v {
				encontrado = true
				contador = k
				contador2 = promo.Cantidades[kk]

				if encontrado {
					almacen.Existencia[contador] = almacen.Existencia[contador] - (cantidad * contador2)
				}
			}
		}

	}

	// fmt.Println("encontrado ->", encontrado)
	// fmt.Println("almacen ->", almacen)
	// fmt.Println("promo ->", promo)
	// fmt.Println("contador ->", contador)
	// fmt.Println("contador2 ->", contador2)
	// fmt.Println("cantidad ->", cantidad)

	// if encontrado {
	// 	almacen.Existencia[contador] = almacen.Existencia[contador] - (cantidad * contador2)
	// }

	err1 := c.UpdateId(almacen.ID, almacen)
	if err1 != nil {
		fmt.Println(err1)
	}

}

//ActuailzaAlmacenDesdeModal -> Actualiza la cantidad de existencia en el Almacen Refrigerador
func ActuailzaAlmacenDesdeModal(producto bson.ObjectId, cantidad int) {

	var almacen inventariomodel.Almacen

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_AL)
	errx := c.Find(bson.M{"Nombre": "Refrigerador"}).One(&almacen)
	if errx != nil {
		fmt.Println(errx)
	}

	encontrado := false
	contador := 0
	for k, v := range almacen.Productos {

		if v == producto {
			encontrado = true
			contador = k
		}

	}

	if encontrado {
		almacen.Existencia[contador] = almacen.Existencia[contador] - cantidad
	}

	err1 := c.UpdateId(almacen.ID, almacen)
	if err1 != nil {
		fmt.Println(err1)
	}

}

//PrecioProducto -> Hace una llamada a productos y regresa su precio al publico mediante su id
func PrecioProducto(producto bson.ObjectId) float64 {

	prod := inventariomodel.ExtraeProducto(producto.Hex())

	return prod.PrecioPub

}

//ExtraeFondo -> Extrae la configuracion del fondo de pantalla
func ExtraeFondo() int {

	var fondo Fondo

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_CFG)
	errx := c.Find(bson.M{"Configuracion": "Fondo"}).One(&fondo)
	if errx != nil {
		fmt.Println(errx)
	}

	return fondo.Disponibles

}

//ActualizaFondo -> Actualiza la configuracion del fondo de pantalla
func ActualizaFondo(fondint int) {
	var fondo Fondo

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_CFG)
	errx := c.Find(bson.M{"Configuracion": "Fondo"}).One(&fondo)
	if errx != nil {
		fmt.Println(errx)
	}

	fondo.Disponibles = fondint

	err1 := c.UpdateId(fondo.ID, fondo)
	if err1 != nil {
		fmt.Println(err1)
	}

}
