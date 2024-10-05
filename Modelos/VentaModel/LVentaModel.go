package ventamodel

import (
	"context"
	"errors"
	"fmt"
	"time"

	conexiones "github.com/vadgun/Bar/Conexiones"
	db "github.com/vadgun/Bar/Modelos/Db"
	inventariomodel "github.com/vadgun/Bar/Modelos/InventarioModel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log"
)

// ExtraeMesas -> Extrae las mesas dadas de alta
func ExtraeMesas() int {
	var mesas ConfigurarMesas
	var err error

	client, _ := db.ConectarMongoDB()
	configuraciones := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_CFG)
	opts := options.FindOne().SetSort(bson.M{"Configuracion": 1})
	err = configuraciones.FindOne(context.TODO(), bson.M{"Configuracion": "Mesas"}, opts).Decode(&mesas)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err)
		}
	}

	return mesas.Disponibles
}

// ActualizarMesasDiarias -> Actualizara el numero de mesas para crear diariamente
func ActualizarMesasDiarias(numero int) {
	client, _ := db.ConectarMongoDB()

	configuraciones := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_CFG)
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.M{"Configuracion": "Mesas"}
	update := bson.M{"$set": bson.M{"Disponibles": numero}}
	var updatedDocument bson.M
	err := configuraciones.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		opts,
	).Decode(&updatedDocument)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return
		}
		log.Println(err)
		fmt.Println("Mesas no actualizadas")
	}
	fmt.Println("Mesas actualizadas")

}

// GuardaMesaDiaria -> Crea la mesa diaria con fecha del dia, numero de mesa y status desocupado
func GuardaMesaDiaria(mesa Mesa) {
	mesa.ID = primitive.NewObjectID()
	client, _ := db.ConectarMongoDB()
	mesasdiarias := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_MD)
	res, err := mesasdiarias.InsertOne(context.TODO(), mesa)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Mesa diaria insertada ID %v\n", res.InsertedID)
}

// BuscarVentasDiarias -> Busca por dia las mesas diarias
func BuscarVentasDiarias(fecha time.Time) []Mesa {
	var mesasdiarias []Mesa

	client, _ := db.ConectarMongoDB()
	mesasdiariascol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_MD)
	opts := options.Find().SetSort(bson.M{"Mesa": 1})
	cursor, err := mesasdiariascol.Find(context.TODO(), bson.M{"Fecha": fecha}, opts)
	if err != nil {
		fmt.Println(err, "Error en BuscarVentasDiarias")
	}
	if err = cursor.All(context.TODO(), &mesasdiarias); err != nil {
		fmt.Println(err, "Error en BuscarVentasDiarias")
	}

	return mesasdiarias
}

// ExtraeMesa -> Regresa la mesa solicitada
func ExtraeMesa(idmesa string) Mesa {
	idobjmesa, err := primitive.ObjectIDFromHex(idmesa)
	if err != nil {
		fmt.Println(err)
	}

	var mesa Mesa

	client, _ := db.ConectarMongoDB()
	mesasdiarias := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_MD)
	opts := options.FindOne().SetSort(bson.M{"Mesa": 1})
	err = mesasdiarias.FindOne(context.TODO(), bson.M{"_id": idobjmesa}, opts).Decode(&mesa)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Mesa no encontrada")
		}
	}

	return mesa
}
func tiempoTranscurrido(inicio, fin time.Time) string {
	// Calculamos la diferencia entre las dos fechas
	diferencia := fin.Sub(inicio)

	// Convertimos la diferencia a horas, minutos y segundos
	horas := int(diferencia.Hours())
	minutos := int(diferencia.Minutes()) % 60
	segundos := int(diferencia.Seconds()) % 60

	// Formateamos el resultado en "hh:mm:ss"
	return fmt.Sprintf("%02d:%02d:%02d", horas, minutos, segundos)
}

// CierraMesa -> Cambia los status de la mesa sin borrar el registro para la venta diaria y el gran total que se esta buscando
func CierraMesa(mesaX Mesa) {

	var err error

	client, _ := db.ConectarMongoDB()
	mesasdiarias := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_MD)
	mesaX.FechaTermino = time.Now()

	filter := bson.M{"_id": mesaX.ID}
	update := bson.M{"$set": bson.M{"Estatus": false, "Abierta": false, "Cerrada": true, "FechaTermino": mesaX.FechaTermino, "Ocupacion": tiempoTranscurrido(mesaX.FechaInicio, mesaX.FechaTermino)}}
	opts2 := options.Update().SetUpsert(false)
	result, err := mesasdiarias.UpdateOne(context.TODO(), filter, update, opts2)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Mesa no encontrada")
		}
	}

	if result.MatchedCount != 0 {
		fmt.Println("Mesa actualizada: ", mesaX.Mesa)
	}

	var cantidades []int
	var objectsids []primitive.ObjectID

	var nuevamesa Mesa
	nuevamesa.ID = primitive.NewObjectID()
	nuevamesa.Fecha = mesaX.Fecha
	nuevamesa.Mesa = mesaX.Mesa
	nuevamesa.Abierta = true
	nuevamesa.Mesero = ""
	nuevamesa.Cerrada = false
	nuevamesa.Estatus = false
	nuevamesa.Productos = objectsids
	nuevamesa.Cantidades = cantidades
	nuevamesa.GranTotal = 0

	res, err := mesasdiarias.InsertOne(context.TODO(), nuevamesa)
	if err != nil {
		fmt.Println(err, "Error insertando mesa nueva")
	}

	fmt.Printf("Mesa nueva insertada %v\n", res.InsertedID)

}

// EliminarColeccionVentasDiarias -> Elimina la Base de Datos de la coleccion de mesas diarias.
func EliminarColeccionVentasDiarias() {
	client, _ := db.ConectarMongoDB()
	ventadiaria := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_MD)
	// if err := ventadiaria.Drop(context.TODO()); err != nil {
	// 	fmt.Println("Error eliminando la coleccion")
	// }

	// Ahora en ves de eliminar toda la coleccion, vamos a introducir una fecha del dia para eliminar las mesas del dia unicamente
	// posteriormente se tiene que realizar los cambios al formulario correspondiente pero sera muy parecido al de imprimir venta
	// o venta diaria. esto se hace para prevenir que la funcion de mongosupervisor deje de apuntar a la coleccion mesasdiarias,
	// ya que al eliminar la coleccion por completo y crearla al iniciar la venta del dia, el puntero ya no inicia el watch de nuevo.
	// por eso se implementa la consistencia de datos y solo se eliminaran documentos.

	filter := bson.M{}
	result, err := ventadiaria.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	if result.DeletedCount != 0 {
		fmt.Printf("Documentos eliminados %v", result.DeletedCount)
	}

}

// ActualizarMesaDiaria -> Actualiza la mesa diaria
func ActualizarMesaDiaria(mesa Mesa) {
	var sumatotal float64
	for k, v := range mesa.Cantidades {
		precio := PrecioProducto(mesa.Productos[k])
		sumatotal += (float64(v) * precio)
	}
	var creatorsDate time.Time

	if mesa.FechaInicio == creatorsDate {
		mesa.FechaInicio = time.Now()
	}

	client, _ := db.ConectarMongoDB()
	mesasdiarias := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_MD)
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": mesa.ID}
	update := bson.M{"$set": bson.M{"Estatus": true, "GranTotal": sumatotal, "Productos": mesa.Productos, "Cantidades": mesa.Cantidades, "Mesero": mesa.Mesero, "FechaInicio": mesa.FechaInicio, "Ocupacion": tiempoTranscurrido(mesa.FechaInicio, time.Now())}}
	result, err := mesasdiarias.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount != 0 {
		fmt.Println("Mesa actualizada: ", mesa.Mesa)
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("Se inserto una nueva mesa %v\n", result.UpsertedID)
	}

}

// ActuailzaAlmacenDesdeModalPromo -> Actualiza la cantidad de existencia en el Almacen Refrigerador
func ActuailzaAlmacenDesdeModalPromo(producto primitive.ObjectID, cantidad int) {

	var almacen inventariomodel.Almacen
	var promo inventariomodel.Producto
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

	//Extrae la Promo
	productos := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_P)
	err = productos.FindOne(context.TODO(), bson.M{"_id": producto}, opts).Decode(&promo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Promo no encontrada")
		}
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

	// Actualizar Almacen
	filter = bson.M{"_id": almacen.ID}
	update := bson.M{"$set": bson.M{"Existencia": almacen.Existencia}}
	result, err := almacenes.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Almacen Modal no encontrada")
		}
	}

	if result.MatchedCount != 0 {
		fmt.Println("Almacen Modal actualizado: ", almacen.Nombre)
	}
}
func GetNextOrderID() int {
	filter := bson.M{"_id": "orderid"}
	update := bson.M{"$inc": bson.M{"sequence_value": 1}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
	var result struct {
		SequenceValue int `bson:"sequence_value"`
	}
	client, _ := db.ConectarMongoDB()
	collection := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_c)
	err := collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)
	if err != nil {
		return 0
	}

	return result.SequenceValue

}

// ActuailzaAlmacenDesdeModal -> Actualiza la cantidad de existencia en el Almacen Refrigerador
func ActuailzaAlmacenDesdeModal(producto primitive.ObjectID, cantidad int) {

	var almacen inventariomodel.Almacen
	var err error

	client, _ := db.ConectarMongoDB()
	almacenes := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_AL)
	opts := options.FindOne().SetSort(bson.M{"Nombre": 1})
	filter := bson.M{"Nombre": "Refrigerador"}
	err = almacenes.FindOne(context.TODO(), filter, opts).Decode(&almacen)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Almacen no encontrado")
		}
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

	filter = bson.M{"_id": almacen.ID}
	update := bson.M{"$set": bson.M{"Existencia": almacen.Existencia}}
	result, err := almacenes.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id de Almacen no encontrado")
		}
	}

	if result.MatchedCount != 0 {
		fmt.Println("Almacen actualizado: ", almacen.Nombre)
	}
}

// PrecioProducto -> Hace una llamada a productos y regresa su precio al publico mediante su id
func PrecioProducto(producto primitive.ObjectID) float64 {

	prod := inventariomodel.ExtraeProducto(producto.Hex())

	return prod.PrecioPub

}

// ExtraeFondo -> Extrae la configuracion del fondo de pantalla
func ExtraeFondo() Fondo {
	var fondo Fondo
	var err error

	client, _ := db.ConectarMongoDB()
	configuraciones := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_CFG)
	opts := options.FindOne().SetSort(bson.M{"Configuracion": 1})
	err = configuraciones.FindOne(context.TODO(), bson.M{"Configuracion": "Fondo"}, opts).Decode(&fondo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err)
		}
	}

	return fondo
}

// ActualizaFondo -> Actualiza la configuracion del fondo de pantalla
func ActualizaFondo(fondint int) {
	client, _ := db.ConectarMongoDB()

	configuraciones := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_CFG)
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.M{"Configuracion": "Fondo"}
	update := bson.M{"$set": bson.M{"Disponibles": fondint}}
	var updatedDocument bson.M
	err := configuraciones.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		opts,
	).Decode(&updatedDocument)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
			fmt.Println("Fondo no actualizado")
		}
	}
	fmt.Println("Fondo actualizado")

}

// ActualizaMensajes -> Actualiza la configuracion del fondo de pantalla
func ActualizaMensajes(mensajes []string) {
	client, _ := db.ConectarMongoDB()

	configuraciones := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_CFG)
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.M{"Configuracion": "Fondo"}
	update := bson.M{"$set": bson.M{"Mensajes": mensajes}}
	var updatedDocument bson.M
	err := configuraciones.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		opts,
	).Decode(&updatedDocument)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
			fmt.Println("Mensajes no actualizados")
		}
	}
	fmt.Println("Mensajes actualizados")
}
