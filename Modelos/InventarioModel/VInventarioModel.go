package inventariomodel

import "gopkg.in/mgo.v2/bson"

//Producto -> Estructura para guardar y recibir los datos de los productos
type Producto struct {
	ID         bson.ObjectId   `bson:"_id,omitempty"`
	Nombre     string          `bson:"Nombre"`
	PrecioPub  float64         `bson:"PrecioPub"`
	PrecioUti  float64         `bson:"PrecioUti"`
	Categoria  string          `bson:"Categoria"`
	Imagen     bson.ObjectId   `bson:"Imagen"`
	Productos  []bson.ObjectId `bson:"Productos"`
	Cantidades []int           `bson:"Cantidades"`
}

//Almacen -> Estructura para manejar los almacenes y los productos almacenados en ellos
type Almacen struct {
	ID         bson.ObjectId   `bson:"_id,omitempty"`
	Nombre     string          `bson:"Nombre"`
	Productos  []bson.ObjectId `bson:"Productos"`
	Existencia []int           `bson:"Existencia"`
}

//VentaDiaria -> Recupera todos los productos y cantidades
type VentaDiaria struct {
	Productos  []bson.ObjectId
	Existencia []int
}
