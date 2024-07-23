package indexmodel

import (
	"context"
	"errors"
	"fmt"

	conexiones "github.com/vadgun/Bar/Conexiones"
	db "github.com/vadgun/Bar/Modelos/Db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// VerificarUsuario Autentifica al usuario en la base de datos
func VerificarUsuario(usuario MongoUser) (bool, MongoUser) {
	var encontrado bool
	var err error

	client, _ := db.ConectarMongoDB()
	usuarios := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_U)

	opts := options.FindOne().SetSort(bson.M{"Nombre": 1})
	err = usuarios.FindOne(context.TODO(), bson.M{"Usuario": usuario.Usuario, "Key": usuario.Key}, opts).Decode(&usuario)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Usuario no encontrado")
		}
	}
	encontrado = usuario.Nombre != ""

	return encontrado, usuario
}

// GetUserOn Se extrae el usuario logeado
func GetUserOn(user string) MongoUser {

	var usuarioOn MongoUser
	usrobjid, err := primitive.ObjectIDFromHex(user)
	if err != nil {
		fmt.Println(err)
	}
	client, _ := db.ConectarMongoDB()
	usuarios := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_U)
	opts := options.FindOne().SetSort(bson.M{"Nombre": 1})
	err = usuarios.FindOne(context.TODO(), bson.M{"_id": usrobjid}, opts).Decode(&usuarioOn)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err, "Id no encontrada")
		}
	}
	return usuarioOn
}
