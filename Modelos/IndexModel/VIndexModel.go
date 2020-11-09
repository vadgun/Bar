package indexmodel

import (
	"gopkg.in/mgo.v2/bson"
)

//MongoUser Controla datos de Usuario del sistema
type MongoUser struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Nombre       string        `bson:"Nombre"`
	Apellidos    string        `bson:"Apellidos"`
	Edad         int           `bson:"Edad"`
	Usuario      string        `bson:"Usuario"`
	Telefono     string        `bson:"Telefono"`
	Puesto       string        `bson:"Puesto"`
	Key          string        `bson:"Key"`
	Nombre2      string        `bson:"Nombre2"`
	Admin        bool          `bson:"Admin"`
	Presupuesto  bool          `bson:"Presupuesto"`
	Herramientas bool          `bson:"Herramientas"`
	Reportes     bool          `bson:"Reportes"`
	Produccion   bool          `bson:"Produccion"`
}
