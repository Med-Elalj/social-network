package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"sync"

	"social-network/server/logs"
	"social-network/sn/db"
	"social-network/sn/security"
	"social-network/sn/security/jwt"

	"github.com/gorilla/websocket"
)

// Define a WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type update struct {
	Sender string `json:"sender"`
	Type   string `json:"type"`
	Uname  string `json:"username"`
	Online bool   `json:"online"`
}

type message struct {
	Sender   int    `json:"sender"`
	Receiver int    `json:"receiver"`
	Message  string `json:"message"`
}

var (
	sockets = make(map[int]*websocket.Conn)
	mutex   sync.Mutex
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	uName := getUid(r)
	// fmt.Printf("New connection: %s\n", uName)

	addConnToMap(uName, conn)
	defer deleteConnFromMap(uName)

	for {
		// Read message from the WebSocket connection
		azer, msg, err := conn.ReadMessage()
		if err != nil || azer != websocket.TextMessage {
			log.Println(err)
			return
		}

		var request message
		// if len(msg) > 7 && string(msg)[:7] == "typing:" {
		// 	m := string(msg)[7:]
		// 	conn, exist := sockets[m]
		// 	if !exist || conn == nil {
		// 		continue
		// 	}

		// 	err := conn.WriteMessage(websocket.TextMessage, []byte(`{"sender":"internal","type":"typing","username":"`+uName+`"}`))
		// 	if err != nil {
		// 		logs.Errorf("Error sending typing notification:", err)
		// 	}
		// 	continue
		// } else if len(msg) > 11 && string(msg)[:11] == "stoptyping:" {
		// 	m := string(msg)[11:]
		// 	conn, exist := sockets[m]
		// 	if !exist || conn == nil {
		// 		continue
		// 	}

		// 	err := conn.WriteMessage(websocket.TextMessage, []byte(`{"sender":"internal","type":"stoptyping","username":"`+uName+`"}`))
		// 	if err != nil {
		// 		logs.Errorf("Error sending stop typing notification:", err)
		// 	}
		// 	continue
		// }
		err = json.Unmarshal(msg, &request)
		if err != nil {
			log.Println("Error parsing JSON:", err)
			continue
		}
		request.Sender = uName
		request.Message = html.EscapeString(request.Message)
		// Respond back with a JSON message
		err = request.send()
		var status_response string
		if err != nil {
			logs.Errorf("Error handling request: %q", err.Error())
			status_response = `{"sender":"system","message":"failed to send message"}`
			err = conn.WriteMessage(websocket.TextMessage, []byte(status_response))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func addConnToMap(uID int, conn *websocket.Conn) {
	mutex.Lock()
	if conn, exists := sockets[uID]; exists {
		log.Printf("User %d already connected\n", uID)
		conn.Close()
	} else {
		for _, v := range sockets {
			if err := v.WriteJSON(update{"internal", "toggle", fmt.Sprint(uID), true}); err != nil {
				logs.Errorf("azer %v", err)
			}
		}
	}
	sockets[uID] = conn
	mutex.Unlock()
}

func deleteConnFromMap(uID int) {
	mutex.Lock()
	delete(sockets, uID)
	for _, v := range sockets {
		if err := v.WriteJSON(update{"internal", "toggle", fmt.Sprint(uID), false}); err != nil {
			logs.Errorf("qsdf %v", err)
		}
	}
	mutex.Unlock()
}

func getUid(r *http.Request) int {
	payload := r.Context().Value(security.UserContextKey)
	data, ok := payload.(*jwt.JwtPayload)
	if ok {
		return data.Sub
	} else {
		return 0
	}
}

func (m *message) send() error {
	err := db.AddDm(m.Sender, m.Receiver, m.Message)
	if err != nil {
		err = errors.New("failed to store message in db with error: " + err.Error())
		logs.Errorf("Error storing message in database: %v", err)
		return err
	}
	responseData, err := json.Marshal(m)
	if err != nil {
		log.Println("Error marshaling response:", err)
		return err
	}
	conn, exist := sockets[m.Sender]
	if !exist || conn == nil {
		log.Printf("User %d not found or not connected\n", m.Receiver)
		return fmt.Errorf("user not found or not connected")
	}

	err = conn.WriteMessage(websocket.TextMessage, responseData)
	if err != nil {
		log.Println(err)
		return errors.New("failed to send message to receiver with error: " + err.Error())
	}

	conn, exist = sockets[m.Receiver]
	if !exist || conn == nil {
		return nil
	}

	conn.WriteMessage(websocket.TextMessage, responseData)
	return nil
}
