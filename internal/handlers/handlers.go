package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConneciton]string)

// definig html template directory
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

// upgrades a http-request to websocket connection
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// handles websockets

// renders home page
func Home(w http.ResponseWriter, r *http.Request) {
	if err := renderPage(w, "home.jet", nil); err != nil {
		log.Println("error in Home while renderpage")
	}

}

type WebSocketConneciton struct {
	*websocket.Conn
}

// defines the response sent back from ws
type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConneciton `json:"-"`
}

// hanldes ws requests
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("error in WsEndpoint can not upgrade request to web socket", err)
	}
	log.Println("client connected to endpoint")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`

	conn := WebSocketConneciton{Conn: ws}
	clients[conn] = ""
	err = ws.WriteJSON(response)
	if err != nil {
		log.Fatalln("error in WsEndpoint can not write json", err)
	}
	go ListenForWs(&conn)

}

// listens for ws and send payload to a channel
func ListenForWs(conn *WebSocketConneciton) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("error listeningForWs ", fmt.Sprintf("%v", err))
		}
	}()
	var payload WsPayload
	for {
		if err := conn.ReadJSON(&payload); err != nil {

		} else {
			payload.Conn = *conn
			wsChan <- payload

		}
	}
}

// lisetns to channel and broadcasts response to all clients
func ListenToWsChannel() {
	var response WsJsonResponse

	for {
		e := <-wsChan
		switch e.Action {
		case "username":
			//get a list of user and send it back via broadcast
			clients[e.Conn] = e.Username
			response.Action = "list_users"
			response.ConnectedUsers = getUserList()
			broadcastToAll(response)

		}
		// response.Action = "Got it here!"
		// response.Message = fmt.Sprintf("How u doing? the action was %s", e.Action)
		// broadcastToAll(response)
	}
}
func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		if err := client.WriteJSON(response); err != nil {
			log.Println("error in broadcasting while writing json", err)
			_ = client.Close()
			delete(clients, client)
		}
	}

}

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		userList = append(userList, x)
	}
	sort.Strings(userList)
	return userList

}

//  renders jet template
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println("error in renderpage while getTemplate", err)
		return err
	}
	if err = view.Execute(w, data, nil); err != nil {
		log.Println("error in renderpage while execute", err)
		return err
	}
	return nil

}
