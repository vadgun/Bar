package inventariomodel

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	conexiones "github.com/vadgun/Bar/Conexiones"
	db "github.com/vadgun/Bar/Modelos/Db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func productos(client *mongo.Client) *mongo.Collection {
	var productos *mongo.Collection
	productos = client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_P)
	return productos
}

// GuardaProducto -> Guarda el producto en la Base de datos
func GuardaProducto(producto Producto) bool {

	client, _ := db.ConectarMongoDB()
	productos := productos(client)

	res, err := productos.InsertOne(context.TODO(), producto)
	if err != nil {
		fmt.Println(err, "Error Agregando el producto")
		return false
	}
	fmt.Printf("Producto Agregado : %v, %v\n", producto.Nombre, res.InsertedID)
	return true

}

func UploadImageToMongo(path string, namefile string) primitive.ObjectID {

	client, _ := db.ConectarMongoDB()
	bucket, err := gridfs.NewBucket(client.Database("Bar"))
	if err != nil {
		fmt.Println("Error creating GridFS bucket")
	}

	file2, err := os.Open(path + "/" + namefile)
	check(err, "Error al abrir el archivo o el archivo no existe")
	defer file2.Close()
	stat, err := file2.Stat()
	check(err, "Error al leer el archivo")
	bs := make([]byte, stat.Size()) // read the file
	_, err = file2.Read(bs)
	check(err, "Error al crear objeto que contendrá el archivo")

	fileID, err := bucket.UploadFromStream(namefile, bytes.NewBuffer(bs))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("new file created with ID %s", fileID)
	return fileID
}

func UpdateImgArt(idimg primitive.ObjectID, id primitive.ObjectID) {

	client, _ := db.ConectarMongoDB()
	productos := productos(client)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"Imagen": idimg}}
	opts2 := options.Update().SetUpsert(false)
	result, err := productos.UpdateOne(context.TODO(), filter, update, opts2)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Imagen no encontrada")
		}
	}

	if result.MatchedCount != 0 {
		fmt.Println("Id de imagen actualizada: ")
		return
	}

}

func check(err error, mensaje string) {
	if err != nil {
		fmt.Println("Error en :", mensaje)
		panic(err)
	}
}

func TraerImagenActa(idimg primitive.ObjectID) string {
	client, err := db.ConectarMongoDB()
	check(err, "Error conectando a MongoDB")
	bucket, err := gridfs.NewBucket(client.Database("Bar"))
	check(err, "Error creando el bucket de GridFS")
	downloadStream, err := bucket.OpenDownloadStream(idimg)
	if err != nil {
		log.Fatal(err)
	}
	defer downloadStream.Close()

	buf := make([]byte, downloadStream.GetFile().Length)
	_, err = downloadStream.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	// Codificar el archivo en Base64
	base64Image := base64.StdEncoding.EncodeToString(buf)

	return base64Image
}

// TodosLosProductos -> Regresa todos los productos para una tabla la cual tendra accines para agregar al almacen
func TodosLosProductos() []Producto {

	var productosAll []Producto
	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	filter := bson.M{}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := productoscol.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err, "Error en Buscar todos los productos")
	}
	if err = cursor.All(context.TODO(), &productosAll); err != nil {
		fmt.Println(err, "Error en Buscar Todos los productos")
	}

	return productosAll
}

// ExtraeBotellas -> Regresa unicamente las botellas
func ExtraeBotellas() []Producto {

	var botellas []Producto
	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	filter := bson.M{"Categoria": "botella"}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := productoscol.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err, "Error en Buscar Botellas")
	}
	if err = cursor.All(context.TODO(), &botellas); err != nil {
		fmt.Println(err, "Error en Buscar Botellas")
	}

	return botellas

}

// ExtraeCervezas -> Regresa unicamente las cervezas
func ExtraeCervezas() []Producto {
	var cervezas []Producto
	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	filter := bson.M{"Categoria": "cerveza"}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := productoscol.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err, "Error en Buscar Cervezas")
	}
	if err = cursor.All(context.TODO(), &cervezas); err != nil {
		fmt.Println(err, "Error en Buscar Cervezas")
	}

	return cervezas

}

// ExtraeBotanas -> Regresa unicamente las botanas
func ExtraeBotanas() []Producto {
	var botanas []Producto
	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	filter := bson.M{"Categoria": "botana"}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := productoscol.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err, "Error en Buscar Botanas")
	}
	if err = cursor.All(context.TODO(), &botanas); err != nil {
		fmt.Println(err, "Error en Buscar Botanas")
	}

	return botanas

}

// ExtraePromos -> Regresa unicamente las promociones
func ExtraePromos() []Producto {
	var promos []Producto
	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	filter := bson.M{"Categoria": "promo"}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := productoscol.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err, "Error en Buscar Promos")
	}
	if err = cursor.All(context.TODO(), &promos); err != nil {
		fmt.Println(err, "Error en Buscar Promos")
	}
	return promos

}

// ExtraeCancion -> Regresa unicamente las canciones
func ExtraeCancion() []Producto {
	var canciones []Producto
	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	filter := bson.M{"Categoria": "cancion"}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := productoscol.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err, "Error en Buscar Cancion")
	}
	if err = cursor.All(context.TODO(), &canciones); err != nil {
		fmt.Println(err, "Error en Buscar Cancion")
	}
	return canciones

}

// ExtraeFichas -> Regresa unicamente las fichas
func ExtraeFichas() []Producto {
	var fichas []Producto
	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	filter := bson.M{"Categoria": "ficha"}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := productoscol.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err, "Error en Buscar Ficha")
	}
	if err = cursor.All(context.TODO(), &fichas); err != nil {
		fmt.Println(err, "Error en Buscar Ficha")
	}
	return fichas

}

// ExtraeCigarros -> Regresa unicamente los cigarros
func ExtraeCigarros() []Producto {
	var cigarros []Producto
	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	filter := bson.M{"Categoria": "cigarros"}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := productoscol.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err, "Error en Buscar Ficha")
	}
	if err = cursor.All(context.TODO(), &cigarros); err != nil {
		fmt.Println(err, "Error en Buscar Ficha")
	}
	return cigarros

}

// ExtraeCopas -> Regresa unicamente las copas
func ExtraeCopas() []Producto {
	var copas []Producto
	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	filter := bson.M{"Categoria": "copa"}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := productoscol.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err, "Error en Buscar Ficha")
	}
	if err = cursor.All(context.TODO(), &copas); err != nil {
		fmt.Println(err, "Error en Buscar Ficha")
	}
	return copas

}

func ExtraerAlmacenes() []Almacen {
	var almacenes []Almacen
	client, _ := db.ConectarMongoDB()
	almacenescol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_AL)
	opts := options.Find().SetSort(bson.M{"Nombre": 1})
	filter := bson.M{}

	cursor, err := almacenescol.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err, "Error en Buscar Almacen")
	}
	if err = cursor.All(context.TODO(), &almacenes); err != nil {
		fmt.Println(err, "Error en Buscar Almacen")
	}
	return almacenes

}

// ExtraerAlmacen -> Extrae el almacen unico a editar
func ExtraerAlmacen(info string) *Almacen {

	idobj, err := primitive.ObjectIDFromHex(info)
	if err != nil {
		fmt.Println(err, "Error convirtiendo id de almacen a objectID")
	}

	var almacen Almacen
	client, _ := db.ConectarMongoDB()

	almacenescol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_AL)
	err = almacenescol.FindOne(context.TODO(), bson.M{"_id": idobj}).Decode(&almacen)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Almacen no encontrado")
		}
	}

	return &almacen

}

// ExtraeBodegaID -> Devuelve el ID para el almacen de "Bodega" en MongoDB
func ExtraeBodegaID() string {

	var almacen Almacen
	var err error
	client, _ := db.ConectarMongoDB()
	almacenescol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_AL)
	opts := options.FindOne().SetSort(bson.M{"Mesa": 1})
	err = almacenescol.FindOne(context.TODO(), bson.M{"Nombre": "Bodega"}, opts).Decode(&almacen)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Bodega no encontrado")
		}
	}
	return almacen.ID.Hex()
}

// AgregarAAlmacen -> Verifica la existencia del producto en el almacen y si no existe lo agrega
func AgregarAAlmacen(producto, almacenid string) bool {
	idobjprod, err := primitive.ObjectIDFromHex(producto)
	if err != nil {
		fmt.Println(err, "Error convirtiendo id de producto a objectID")
	}

	almacen := ExtraerAlmacen(almacenid)

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

		client, _ := db.ConectarMongoDB()
		almacenescol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_AL)

		fmt.Println(almacen)

		filter := bson.M{"_id": almacen.ID}
		update := bson.M{"$set": bson.M{"Productos": almacen.Productos, "Existencia": almacen.Existencia}}
		result, err := almacenescol.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				fmt.Println(err, "Id de Almacen no encontrado")
			}
		}
		fmt.Println(result)
		return true
	}

}

// EliminarDeAlmacen -> Elimina el producto del almacen seleccionado
func EliminarDeAlmacen(producto, almacenid string) bool {
	idobjprod, err := primitive.ObjectIDFromHex(producto)
	if err != nil {
		fmt.Println(err, "Error convirtiendo id de producto a objectID")
	}
	almacen := ExtraerAlmacen(almacenid)

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

			almacen.Existencia = RemoveInts(almacen.Existencia, contador)
			almacen.Productos = RemovePrimitives(almacen.Productos, contador)

			client, _ := db.ConectarMongoDB()
			almacenescol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_AL)

			filter := bson.M{"_id": almacen.ID}
			opts2 := options.Update().SetUpsert(false)
			result, err := almacenescol.UpdateOne(context.TODO(), filter, almacen, opts2)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					fmt.Println(err, "Id de Almacen no encontrado")
				}
			}

			if result.MatchedCount != 0 {
				fmt.Println("Producto Agregado al Almacen")
			}

			return true
		}

	} else {
		return false
	}
}

// RemoveInts -> Remueve el indice seleccionado del arreglo, haciendo un slice de 2 y lo vuelve a unir Version para Enteros
func RemoveInts(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

// RemovePrimitives -> Remueve el indice seleccionado del arreglo, haciendo un slice de 2 y lo vuelve a unir Version para Enteros
func RemovePrimitives(s []primitive.ObjectID, index int) []primitive.ObjectID {
	return append(s[:index], s[index+1:]...)
}

// ExtraeRefrigeradorID -> Devuelve el ID para el almacen de "Refrigerador" en MongoDB
func ExtraeRefrigeradorID() string {

	var almacen Almacen
	var err error

	client, _ := db.ConectarMongoDB()
	// Extrae el Almacen
	almacenes := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_AL)
	opts := options.FindOne().SetSort(bson.M{"Nombre": 1})
	filter := bson.M{"Nombre": "Refrigerador"}
	err = almacenes.FindOne(context.TODO(), filter, opts).Decode(&almacen)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Almacen no encontrado")
		}
	}
	return almacen.ID.Hex()
}

// ExtraeNombreProducto -> Regresa el nombre del producto
func ExtraeNombreProducto(id primitive.ObjectID) string {
	producto := ExtraeProducto(id.Hex())
	return producto.Nombre
}

// GuardaProductoEditado -> Edita el producto seleccionado
func GuardaProductoEditado(producto Producto) {

	var productoold *Producto
	var err error

	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	filter := bson.M{"_id": producto.ID}
	err = productoscol.FindOne(context.TODO(), filter).Decode(&productoold)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err)
		}
	}

	update := bson.M{"$set": bson.M{"Nombre": producto.Nombre, "PrecioPub": producto.PrecioPub, "PrecioUti": producto.PrecioUti, "Categoria": producto.Categoria, "Imagen": productoold.Imagen}}
	result, err := productoscol.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Producto no encontrado")
		}
	}

	if result.MatchedCount != 0 {
		fmt.Println("Producto actualizado: ", producto.Nombre)
	}

}

// ExtraeProducto -> Devuelve un producto mediante su id
func ExtraeProducto(idProducto string) Producto {

	idobjprod, err := primitive.ObjectIDFromHex(idProducto)
	if err != nil {
		fmt.Println(err, "Error convirtiendo id de producto a objectID")
	}

	var producto Producto
	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	err = productoscol.FindOne(context.TODO(), bson.M{"_id": idobjprod}).Decode(&producto)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Producto no encontrado")
		}
	}
	return producto
}

// ExtraeProductosSinPromo -> Devuelve los producto mediante su id
func ExtraeProductosSinPromo() []Producto {
	var productosr []Producto
	client, _ := db.ConectarMongoDB()
	productoscol := productos(client)
	filter := bson.M{"Categoria": bson.M{"$ne": "promo"}}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})
	cursor, err := productoscol.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err, "Error en Buscar Cancion")
	}
	if err = cursor.All(context.TODO(), &productosr); err != nil {
		fmt.Println(err, "Error en Buscar Cancion")
	}

	return productosr
}

// EliminarProducto -> Elimina un producto segun su id
func EliminarProducto(idProducto string) bool {

	idobjprod, err := primitive.ObjectIDFromHex(idProducto)
	if err != nil {
		fmt.Println(err, "Error convirtiendo id de producto a objectID")
	}

	var bodega Almacen
	var refrigerador Almacen

	client, _ := db.ConectarMongoDB()

	almacenescol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_AL)
	productoscol := productos(client)

	//Existe en Refrigerador con mas de 0
	err = almacenescol.FindOne(context.TODO(), bson.M{"Nombre": "Bodega"}).Decode(&bodega)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Almacen Bodega no encontrado")
		}
	}

	var encontradoenbodega bool
	encontradoenbodega = false
	for _, v := range bodega.Productos {

		if v == idobjprod {
			encontradoenbodega = true
		}

	}

	//Existe en Refrigerador con mas de 0
	err = almacenescol.FindOne(context.TODO(), bson.M{"Nombre": "Refrigerador"}).Decode(&refrigerador)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Almacen Refrigerador no encontrado")
		}
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

		res, err := productoscol.DeleteOne(context.TODO(), bson.M{"_id": idobjprod})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Producto %v eliminado de bodega y refrigerador\n", res.DeletedCount)

		return true
	}

	return false
}

// ActualizarExistenciasEnAlmacen -> Actualiza las existencias del almacen, segun los productos agregados al almacen desde un inicio
func ActualizarExistenciasEnAlmacen(almacen string, productos, existencias []string) bool {
	var existenciasint []int

	var almacenFs *Almacen
	almacenFs = ExtraerAlmacen(almacen)
	client, _ := db.ConectarMongoDB()
	almacenescol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_AL)

	for _, vv := range existencias {
		i, _ := strconv.Atoi(vv)
		existenciasint = append(existenciasint, i)
	}

	for k, v := range almacenFs.Productos {
		primt, _ := primitive.ObjectIDFromHex(productos[k])
		if v == primt {
			almacenFs.Existencia[k] = existenciasint[k]
		}
	}

	filter := bson.M{"_id": almacenFs.ID}
	update := bson.M{"$set": bson.M{"Productos": almacenFs.Productos, "Existencia": almacenFs.Existencia}}
	result, err := almacenescol.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Almacen no encontrado")
		}
	}

	if result.MatchedCount != 0 {
		fmt.Println("Almacen Actualizado")
	}

	return true

}

// ExtraeExistencias -> Extrae la Existencia del almacen seleccionado (Refrigerador)
func ExtraeExistencias(id primitive.ObjectID) int {

	var existencia int

	var almacenFs Almacen
	var err error

	client, _ := db.ConectarMongoDB()
	almacenescol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_AL)
	filter := bson.M{"Nombre": "Refrigerador"}
	err = almacenescol.FindOne(context.TODO(), filter).Decode(&almacenFs)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Refrigerador no encontrado")
		}
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
