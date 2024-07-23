package inventariomodel

import "go.mongodb.org/mongo-driver/bson/primitive"

//Producto -> Estructura para guardar y recibir los datos de los productos
type Producto struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Nombre     string               `bson:"Nombre"`
	PrecioPub  float64              `bson:"PrecioPub"`
	PrecioUti  float64              `bson:"PrecioUti"`
	Categoria  string               `bson:"Categoria"`
	Imagen     primitive.ObjectID   `bson:"Imagen"`
	Productos  []primitive.ObjectID `bson:"Productos"`
	Cantidades []int                `bson:"Cantidades"`
}

//Almacen -> Estructura para manejar los almacenes y los productos almacenados en ellos
type Almacen struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Nombre     string               `bson:"Nombre"`
	Productos  []primitive.ObjectID `bson:"Productos"`
	Existencia []int                `bson:"Existencia"`
}

//VentaDiaria -> Recupera todos los productos y cantidades
type VentaDiaria struct {
	Productos  []primitive.ObjectID
	Existencia []int
}
