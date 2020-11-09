package ventamodel

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Mesa -> Controla la mesa diaria por dia.
type Mesa struct {
	ID         bson.ObjectId   `bson:"_id,omitempty"`
	Mesa       int             `bson:"Mesa"`
	Estatus    bool            `bson:"Estatus"`
	Abierta    bool            `bson:"Abierta"`
	Cerrada    bool            `bson:"Cerrada"`
	Fecha      time.Time       `bson:"Fecha"`
	Productos  []bson.ObjectId `bson:"Productos"`
	Cantidades []int           `bson:"Cantidades"`
	GranTotal  float64         `bson:"GranTotal"`
}

//ConfigurarMesas -> Controla la estructura para mesas diarias.
type ConfigurarMesas struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Configuracion string        `bson:"Configuracion"`
	Disponibles   int           `bson:"Disponibles"`
}

type Fondo struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Configuracion string        `bson:"Configuracion"`
	Disponibles   int           `bson:"Disponibles"`
}
