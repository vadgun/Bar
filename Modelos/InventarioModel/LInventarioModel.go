package inventariomodel

import (
	"bytes"
	"encoding/base64"
	"fmt"
	conexiones "github.com/vadgun/Bar/Conexiones"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"image/jpeg"
	"log"
	"os"
	"strconv"
)

//GuardaProducto -> Guarda el producto en la Base de datos
func GuardaProducto(producto Producto) bool {

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		panic(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.Insert(producto)
	if err1 != nil {
		fmt.Println(err1)
		return false
	} else {
		fmt.Println(" -> ", producto.Nombre, " Agregado")
		return true
	}

}

func UploadImageToMongo(path string, namefile string) bson.ObjectId {
	db, err := mgo.Dial(conexiones.MONGO_SERVER)
	check(err, "Error al conectar con mongo")
	base := db.DB(conexiones.MONGO_DB)
	file2, err := os.Open(path + "/" + namefile)
	check(err, "Error al abrir el archivo o el archivo no existe")
	defer file2.Close()
	stat, err := file2.Stat()
	check(err, "Error al leer el archivo")
	bs := make([]byte, stat.Size()) // read the file
	_, err = file2.Read(bs)
	check(err, "Error al crear objeto que contendrá el archivo")
	img, err := base.GridFS("ImagenesProducto").Create(namefile)
	ids_img := img.Id()
	check(err, "error al crear archivo en mongo")
	_, err = img.Write(bs)
	check(err, "error al escribir archivo en mongo")
	fmt.Println("File uploaded successfully to mongo ")
	err = img.Close()
	check(err, "error al cerrar img de mongo")
	db.Close()
	idimg := getObjectIdToInterface(ids_img)
	return idimg
}

func getObjectIdToInterface(i interface{}) bson.ObjectId {
	var v = i.(bson.ObjectId)
	return v
}

func UpdateImgArt(idimg bson.ObjectId, id bson.ObjectId) {
	c, s := conectMgoPersonal()
	defer s.Close()
	// err := c.Update(bson.M{"_id": id}, bson.M{"$push": bson.M{"ActaNacimiento": idimg}})
	err := c.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"Imagen": idimg}})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Imagen de Producto Updated  :)")
	}
}

func conectMgoPersonal() (*mgo.Collection, *mgo.Session) {
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	if err != nil {
		panic(err)
	}
	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	return c, session
}

func check(err error, mensaje string) {
	if err != nil {
		panic(err)
	}
}

func TraerImagenActa(idimg bson.ObjectId) string {
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()

	check(err, "Error al conectar con MongoDB")
	Base := session.DB(conexiones.MONGO_DB)

	img, err1 := Base.GridFS("ImagenesProducto").OpenId(idimg)
	check(err1, "Error al leer Imagen de MONGO: "+idimg.Hex())

	b := make([]byte, img.Size())
	n, err := img.Read(b)
	check(err, "Error al crear mapa de bytes")

	fmt.Println("N -> ", n)
	imagen, err := jpeg.Decode(bytes.NewReader(b))
	check(err, "Error al decodificar")
	buffer := new(bytes.Buffer)

	err2 := jpeg.Encode(buffer, imagen, nil)
	check(err2, "Error al codificar.")

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())

	defer img.Close()

	return str

}

//TodosLosProductos -> Regresa todos los productos para una tabla la cual tendra accines para agregar al almacen
func TodosLosProductos() []Producto {

	var productos []Producto
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.Find(bson.M{}).Sort("Nombre").All(&productos)
	if err1 != nil {
		fmt.Println(err1)
	}

	return productos
}

//ExtraeBotellas -> Regresa unicamente las botellas
func ExtraeBotellas() []Producto {
	var botellas []Producto
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.Find(bson.M{"Categoria": "botella"}).All(&botellas)
	if err1 != nil {
		fmt.Println(err1)
	}
	return botellas

}

//ExtraeCervezas -> Regresa unicamente las cervezas
func ExtraeCervezas() []Producto {
	var cervezas []Producto
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.Find(bson.M{"Categoria": "cerveza"}).All(&cervezas)
	if err1 != nil {
		fmt.Println(err1)
	}
	return cervezas

}

//ExtraeBotanas -> Regresa unicamente las botanas
func ExtraeBotanas() []Producto {
	var botanas []Producto
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.Find(bson.M{"Categoria": "botana"}).All(&botanas)
	if err1 != nil {
		fmt.Println(err1)
	}
	return botanas

}

//ExtraePromos -> Regresa unicamente las promociones
func ExtraePromos() []Producto {
	var promos []Producto
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.Find(bson.M{"Categoria": "promo"}).All(&promos)
	if err1 != nil {
		fmt.Println(err1)
	}
	return promos

}

//ExtraeCancion -> Regresa unicamente las canciones
func ExtraeCancion() []Producto {
	var canciones []Producto
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.Find(bson.M{"Categoria": "cancion"}).All(&canciones)
	if err1 != nil {
		fmt.Println(err1)
	}
	return canciones

}

//ExtraeFichas -> Regresa unicamente las fichas
func ExtraeFichas() []Producto {
	var fichas []Producto
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.Find(bson.M{"Categoria": "ficha"}).All(&fichas)
	if err1 != nil {
		fmt.Println(err1)
	}
	return fichas

}

//ExtraeCigarros -> Regresa unicamente los cigarros
func ExtraeCigarros() []Producto {
	var cigarros []Producto
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.Find(bson.M{"Categoria": "cigarros"}).All(&cigarros)
	if err1 != nil {
		fmt.Println(err1)
	}
	return cigarros

}

//ExtraeCopas -> Regresa unicamente las copas
func ExtraeCopas() []Producto {
	var copas []Producto
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.Find(bson.M{"Categoria": "copa"}).All(&copas)
	if err1 != nil {
		fmt.Println(err1)
	}
	return copas

}

func ExtraerAlmacenes() []Almacen {
	var almacenes []Almacen
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_AL)
	err1 := c.Find(bson.M{}).All(&almacenes)
	if err1 != nil {
		fmt.Println(err1)
	}
	return almacenes

}

//ExtraerAlmacen -> Extrae el almacen unico a editar
func ExtraerAlmacen(info string) Almacen {

	idobj := bson.ObjectIdHex(info)
	var almacen Almacen
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_AL)
	err1 := c.FindId(idobj).One(&almacen)
	if err1 != nil {
		fmt.Println(err1)
	}
	return almacen

}

//ExtraeBodegaID -> Devuelve el ID para el almacen de "Bodega" en MongoDB
func ExtraeBodegaID() string {

	var almacen Almacen
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_AL)
	err1 := c.Find(bson.M{"Nombre": "Bodega"}).Select(bson.M{"_id": 1}).One(&almacen)
	if err1 != nil {
		fmt.Println(err1)
	}
	return almacen.ID.Hex()
}

//AgregarAAlmacen -> Verifica la existencia del producto en el almacen y si no existe lo agrega
func AgregarAAlmacen(producto, almacenid string) bool {

	//Saber si ya existe en el almacen
	idobj := bson.ObjectIdHex(almacenid)
	idobjprod := bson.ObjectIdHex(producto)
	var almacen Almacen
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_AL)
	err1 := c.FindId(idobj).One(&almacen)
	if err1 != nil {
		fmt.Println(err1)
	}

	encontrado := false
	for _, v := range almacen.Productos {
		if v == idobjprod {
			encontrado = true
		}
	}

	if encontrado {
		return false
	} else {

		almacen.Productos = append(almacen.Productos, idobjprod)
		almacen.Existencia = append(almacen.Existencia, 0)
		err2 := c.UpdateId(almacen.ID, almacen)
		if err1 != nil {
			fmt.Println(err2)
		}
		return true
	}

}

//EliminarDeAlmacen -> Elimina el producto del almacen seleccionado
func EliminarDeAlmacen(producto, almacenid string) bool {
	//Saber si ya existe en el almacen
	idobj := bson.ObjectIdHex(almacenid)
	idobjprod := bson.ObjectIdHex(producto)
	var almacen Almacen

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_AL)
	err1 := c.FindId(idobj).One(&almacen)
	if err1 != nil {
		fmt.Println(err1)
	}

	encontrado := false

	contador := 0
	for k, v := range almacen.Productos {
		if v == idobjprod {
			encontrado = true
			contador = k

		}
	}

	if encontrado {

		if almacen.Existencia[contador] > 0 {
			return false
		} else {

			almacen.Existencia = RemoveIndex(almacen.Existencia, contador)

			err3 := c.UpdateId(almacen.ID, almacen)
			if err3 != nil {
				fmt.Println(err3)
			}
			err2 := c.Update(bson.M{"_id": almacen.ID}, bson.M{"$pull": bson.M{"Productos": idobjprod}})
			if err2 != nil {
				fmt.Println(err2)
			}

			return true
		}

	} else {
		return false
	}
}

//RemoveIndex -> Remueve el indice seleccionado del arreglo, haciendo un slice de 2 y lo vuelve a unir Version para Enteros
func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

//ExtraeRefrigeradorID -> Devuelve el ID para el almacen de "Refrigerador" en MongoDB
func ExtraeRefrigeradorID() string {

	var almacen Almacen
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_AL)
	err1 := c.Find(bson.M{"Nombre": "Refrigerador"}).Select(bson.M{"_id": 1}).One(&almacen)
	if err1 != nil {
		fmt.Println(err1)
	}
	return almacen.ID.Hex()
}

//ExtraeNombreProducto -> Regresa el nombre del producto
func ExtraeNombreProducto(id bson.ObjectId) string {

	var producto Producto

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.FindId(id).One(&producto)
	if err1 != nil {
		fmt.Println(err1)
	}
	return producto.Nombre

}

//GuardaProductoEditado -> Edita el producto seleccionado
func GuardaProductoEditado(producto Producto) {

	var productoold Producto

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.FindId(producto.ID).One(&productoold)
	if err1 != nil {
		fmt.Println("err1", err1)
	}

	producto.Imagen = productoold.Imagen

	err2 := c.UpdateId(productoold.ID, producto)
	if err2 != nil {
		fmt.Println("err2", err2)
	}

}

//ExtraeProducto -> Devuelve un producto mediante su id
func ExtraeProducto(idProducto string) Producto {

	idobjprod := bson.ObjectIdHex(idProducto)

	var producto Producto

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.FindId(idobjprod).One(&producto)
	if err1 != nil {
		fmt.Println(err1)
	}
	return producto
}

//ExtraeProductosSinPromo -> Devuelve los producto mediante su id
func ExtraeProductosSinPromo() []Producto {
	var productos []Producto
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}
	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)
	err1 := c.Find(bson.M{"Categoria": bson.M{"$ne": "promo"}}).Sort("Nombre").All(&productos)
	if err1 != nil {
		fmt.Println(err1)
	}
	return productos
}

//EliminarProducto -> Elimina un producto segun su id
func EliminarProducto(idProducto string) bool {

	idobjprod := bson.ObjectIdHex(idProducto)

	var bodega Almacen
	var refrigerador Almacen

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	//Existe en Bodega con mas de 0
	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_AL)
	d := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_P)

	err1 := c.Find(bson.M{"Nombre": "Bodega"}).One(&bodega)
	if err1 != nil {
		fmt.Println(err1)
	}

	var encontradoenbodega bool
	encontradoenbodega = false
	for _, v := range bodega.Productos {

		if v == idobjprod {
			encontradoenbodega = true
		}

	}

	//Existe en Refrigerador con mas de 0
	err2 := c.Find(bson.M{"Nombre": "Refrigerador"}).One(&refrigerador)
	if err2 != nil {
		fmt.Println(err2)
	}

	var encontradoenrefrigerador bool
	encontradoenrefrigerador = false
	for _, vv := range bodega.Productos {

		if vv == idobjprod {
			encontradoenrefrigerador = true
		}

	}

	if encontradoenrefrigerador == true || encontradoenbodega == true {
		return false
	}

	if encontradoenrefrigerador == false && encontradoenbodega == false {

		err3 := d.RemoveId(idobjprod)
		if err3 != nil {
			fmt.Println(err3)
		}

		return true
	}

	return false
}

//ActualizarExistenciasEnAlmacen -> Actualiza las existencias del almacen, segun los productos agregados al almacen desde un inicio
func ActualizarExistenciasEnAlmacen(almacen string, productos, existencias []string) bool {

	almacenid := bson.ObjectIdHex(almacen)
	var existenciasint []int

	var almacenFs Almacen

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Alamcen Inicial ", almacenFs)

	for _, vv := range existencias {
		i, _ := strconv.Atoi(vv)
		existenciasint = append(existenciasint, i)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_AL)
	err1 := c.FindId(almacenid).One(&almacenFs)
	if err1 != nil {
		fmt.Println(err1)
	}

	for k, v := range almacenFs.Productos {
		if v == bson.ObjectIdHex(productos[k]) {
			almacenFs.Existencia[k] = existenciasint[k]
		}
	}

	err2 := c.UpdateId(almacenFs.ID, almacenFs)
	if err2 != nil {
		fmt.Println(err2)
	}

	fmt.Println("Alamcen Actualizado ", almacenFs)

	return true

}

//ExtraeExistencias -> Extrae la Existencia del almacen seleccionado (Refrigerador)
func ExtraeExistencias(id bson.ObjectId) int {

	var existencia int

	var almacenFs Almacen

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_AL)
	err1 := c.Find(bson.M{"Nombre": "Refrigerador"}).One(&almacenFs)
	if err1 != nil {
		fmt.Println(err1)
	}

	//Buscar el producto en el almacen

	encontrado := false
	contador := 0
	for k, v := range almacenFs.Productos {
		if id == v {
			encontrado = true
			contador = k
		}
	}

	if encontrado {
		existencia = almacenFs.Existencia[contador]
	}

	return existencia

}
