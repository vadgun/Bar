package indexmodel

import (
	"fmt"
	"log"

	conexiones "github.com/vadgun/Bar/Conexiones"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//VerificarUsuario Autentifica al usuario en la base de datos
func VerificarUsuario(usuario MongoUser) (bool, MongoUser) {
	var encontrado bool
	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}
	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_U)
	err1 := c.Find(bson.M{"Usuario": usuario.Usuario, "Key": usuario.Key}).One(&usuario)
	if err1 != nil {
		fmt.Println(err1)
	}

	if usuario.Nombre == "" {
		encontrado = false
	} else {
		encontrado = true
	}

	return encontrado, usuario
}

//GetUserOn Se extrae el usuario logeado
func GetUserOn(user string) MongoUser {

	var usuarioOn MongoUser
	usrobjid := bson.ObjectIdHex(user)

	session, err := mgo.Dial(conexiones.MONGO_SERVER)
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}
	c := session.DB(conexiones.MONGO_DB).C(conexiones.MONGO_DB_U)
	err1 := c.FindId(usrobjid).One(&usuarioOn)
	if err1 != nil {
		fmt.Println(err1)
	}

	return usuarioOn
}
