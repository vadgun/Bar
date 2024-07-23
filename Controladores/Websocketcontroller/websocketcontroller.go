package websocketcontroller

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// 	"github.com/kataras/iris/v12/context"
// )

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func WebsocketHandler(ctx *context.Context) {
// 	conn, err := upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	for {
// 		messageType, p, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		if err := conn.WriteMessage(messageType, p); err != nil {
// 			log.Println(err)
// 			return
// 		}
// 	}
// }

import (
	another "context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/context"
	conexiones "github.com/vadgun/Bar/Conexiones"
	db "github.com/vadgun/Bar/Modelos/Db"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// Global variables to hold active connections
// Variables globales para retener las conexiones activas que llegan del websocket
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Definimos un objeto
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// Mutex para sincronizar el acceso al mapa de clientes
var mutex = &sync.Mutex{}

func WebsocketHandler(ctx *context.Context) {
	conn, err := upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Registrar un nuevo cliente
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	// Dar de baja el cliente cuando se desconecte
	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
	}()

	for {
		var msg Message
		// Leer un mensaje del cliente
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("error:", err)
			break
		}
		// Envia el mensaje recibido al canal de transmision
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("error:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func MongoSupervisor() {

	client, _ := db.ConectarMongoDB()
	mesasdiarias := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_MD)

	defer client.Disconnect(another.Background())
	//Implementar una pipeline de mongodb
	pipeline := mongo.Pipeline{}

	changeStream, err := mesasdiarias.Watch(another.TODO(), pipeline)
	if err != nil {
		fmt.Println("No se pudo crear el stream de mongodb", err)
	}

	defer changeStream.Close(another.TODO())
	fmt.Println("Esperando cambios en la coleccion mesasdiarias...")
	for changeStream.Next(another.TODO()) {
		var changeEvent bson.M
		if err := changeStream.Decode(&changeEvent); err != nil {
			log.Fatal(err)
		}

		fmt.Println(changeEvent)
		// Send it out to every client that is currently connected
		mutex.Lock()
		var newmessage Message
		newmessage.Username = "Jose Roberto"
		newmessage.Message = "Camacho Christy"
		for client := range clients {
			err := client.WriteJSON(newmessage)
			if err != nil {
				log.Println("error:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock() // TODO: Send update to WebSocket clients
	}

	if err := changeStream.Err(); err != nil {
		log.Fatal(err)
	}
}
