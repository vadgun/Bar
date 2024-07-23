package ventamodel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mesa -> Controla la mesa diaria por dia.
type Mesa struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Mesa       int                  `bson:"Mesa"`
	Estatus    bool                 `bson:"Estatus"`
	Abierta    bool                 `bson:"Abierta"`
	Cerrada    bool                 `bson:"Cerrada"`
	Fecha      time.Time            `bson:"Fecha"`
	Productos  []primitive.ObjectID `bson:"Productos"`
	Cantidades []int                `bson:"Cantidades"`
	GranTotal  float64              `bson:"GranTotal"`
}

// ConfigurarMesas -> Controla la estructura para mesas diarias.
type ConfigurarMesas struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Configuracion string             `bson:"Configuracion"`
	Disponibles   int                `bson:"Disponibles"`
}

type Fondo struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Configuracion string             `bson:"Configuracion"`
	Disponibles   int                `bson:"Disponibles"`
}
