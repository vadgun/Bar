package db

import (
	"context"
	"fmt"
	"time"

	conexiones "github.com/vadgun/Bar/Conexiones"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConectarMongoDB Crea un *mongo.Client que devuelve una conexcion a la base de datos
func ConectarMongoDB() (*mongo.Client, error) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conexiones.MONGO_SERVER))
	if err != nil {
		fmt.Println("Error conectando al servidor de Mongo DB: ->", err)
	}

	return client, nil
}
