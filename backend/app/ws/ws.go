package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	auth "social-network/app/Auth"
	"social-network/app/Auth/jwt"
	"social-network/app/modules"
	"social-network/app/structs"

	"social-network/app/logs"

	"github.com/gorilla/websocket"
)

// Define a WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type update struct {
	Sender  string `json:"sender"`
	Uid     int    `json:"uid"`
	Command string `json:"command"`
	Value   any    `json:"value"`
}

type message struct {
	Sender   int    `json:"sender"`
	SName    string `json:"author_name"`
	Receiver int    `json:"receiver"`
	Message  string `json:"content"`
}

var (
	sockets = structs.Sockets
	mutex   sync.Mutex
)

type group struct {
	sync.Mutex
	subs map[int]struct{}
	id   int
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.ErrorLog.Println(err)
		return
	}
	defer conn.Close()

	uId, uName := getData(r)
	if uId == 0 || uName == "" {
		logs.ErrorLog.Println("Invalid user ID or username", uId, uName)
		err = conn.WriteMessage(websocket.TextMessage, []byte(`{"sender":"system","content":"invalid user"}`))
		if err != nil {
			logs.ErrorLog.Println("Error sending invalid user message:", err)
		}
		return
	}
	// fmt.Printf("New connection: %s\n", uName)

	if !addConnToMap(uId, conn) {
		return
	}
	defer deleteConnFromMap(uId)

	for {
		// Read message from the WebSocket connection
		azer, msg, err := conn.ReadMessage()
		if err != nil || azer != websocket.TextMessage {
			log.Println(err)
			return
		}

		var request message
		err = json.Unmarshal(msg, &request)
		if err != nil {
			logs.ErrorLog.Println("Error parsing JSON:", err)
			continue
		}
		request.Sender = uId
		request.SName = uName
		if len(request.Message) > 100 {
			request.Message = request.Message[:100]
		}
		// Respond back with a JSON message
		err = request.send()
		var status_response string
		if err != nil {
			logs.ErrorLog.Printf("Error handling request: %q", err.Error())
			status_response = `{"author_name":"system","content":"failed to send message"}`
			err = conn.WriteMessage(websocket.TextMessage, []byte(status_response))
			if err != nil {
				logs.ErrorLog.Println(err)
				return
			}
		}
	}
}

func addConnToMap(uID int, connection *websocket.Conn) bool {
	mutex.Lock()
	defer mutex.Unlock()
	var u update
	if conn, exists := sockets[uID]; exists {
		if c, is := conn.(*websocket.Conn); is {
			c.Close()
			sockets[uID] = connection
		} else {
			return false
		}
	} else {
		sockets[uID] = connection
	}

	u = update{"<system>", uID, "online", true}
	err := u.send()
	if err != nil {
		logs.ErrorLog.Printf("Error sending update message: %v", err)
		return false
	}
	logs.InfoLog.Printf("User %d connected\n", uID)

	groups, err := modules.GetGroupImIn(uID)
	if err != nil {
		logs.ErrorLog.Printf("add conn to map %q", err)
	}
	for _, v := range groups {
		id := int(v.ID)
		if ws, ok := sockets[id]; ok {
			if g, ok := ws.(*group); ok {
				g.subs[id] = struct{}{}
			} else {
				logs.ErrorLog.Fatalln("user with group id", sockets, uID, groups, id)
			}
		} else {
			g := group{}
			g.subs = map[int]struct{}{uID: {}}
			g.id = id
			sockets[id] = &g
		}
	}
	return true
}

func deleteConnFromMap(uID int) {
	mutex.Lock()
	delete(sockets, uID)

	u := update{"<system>", uID, "online", false}
	err := u.send()
	if err != nil {
		logs.ErrorLog.Printf("Error sending update message: %v", err)
	}
	logs.InfoLog.Printf("User %d disconnected\n", uID)
	mutex.Unlock()
}

func getData(r *http.Request) (int, string) {
	payload := r.Context().Value(auth.UserContextKey)
	data, ok := payload.(*jwt.JwtPayload)
	if ok {
		return data.Sub, data.Username
	} else {
		return 0, ""
	}
}

func (m *message) send() error {
	err := modules.AddDm(m.Sender, m.Receiver, m.Message)
	if err != nil {
		err = errors.New("failed to store message in db with error: " + err.Error())
		logs.ErrorLog.Printf("Error storing message in database: %v", err)
		return err
	}
	responseData, err := json.Marshal(m)
	if err != nil {
		logs.ErrorLog.Println("Error marshaling response:", err)
		return err
	}
	if profile, exist := sockets[m.Sender]; !exist {
		logs.ErrorLog.Printf("User %d not found or not connected\n", m.Receiver)
		return fmt.Errorf("user not found or not connected")
	} else {
		err = profile.WriteMessage(websocket.TextMessage, responseData)
		if err != nil {
			logs.ErrorLog.Println(err)
			return errors.New("failed to send message to receiver with error: " + err.Error())
		}
	}

	if profile, exist := sockets[m.Receiver]; !exist {
		return nil
	} else {
		profile.WriteMessage(websocket.TextMessage, responseData)
	}
	return nil
}

func (g *group) WriteMessage(messageType int, data []byte) error {
	g.Lock()
	defer g.Unlock()

	for id := range g.subs {
		if profile, exist := sockets[id]; !exist {
			if ws, is := profile.(*websocket.Conn); is {
				ws.WriteMessage(websocket.TextMessage, data)
			} else {
				logs.ErrorLog.Printf("how did we get here ??%v %v %v\n", sockets, id, data)
			}
		} else {
			delete(g.subs, id)
		}
	}
	if len(g.subs) == 0 {
		mutex.Lock()
		delete(sockets, g.id)
		mutex.Unlock()
	}
	return nil
}

func (u *update) send() error {
	responseData, err := json.Marshal(u)
	if err != nil {
		logs.ErrorLog.Println("Error marshaling response:", err)
	}
	for _, profile := range sockets {
		if ws, is := profile.(*websocket.Conn); is {
			err = ws.WriteMessage(websocket.TextMessage, responseData)
			if err != nil {
				logs.ErrorLog.Println("Error sending update message:", err)
				return err
			}
		}
	}
	return nil
}

// send to user

func (u *update) sendToUser() error {
	responseData, err := json.Marshal(u)
	if err != nil {
		logs.ErrorLog.Println("Error marshaling response:", err)
		return err
	}

	if err := modules.AddDm(u.Uid, u.Uid, string(responseData)); err != nil {
		err = errors.New("failed to store message in db with error: " + err.Error())
		logs.ErrorLog.Printf("Error storing message in database: %v", err)
		return err
	}

	if profile, exist := sockets[u.Uid]; !exist {
		logs.ErrorLog.Printf("User %d not found or not connected\n", u.Uid)
		return fmt.Errorf("user not found or not connected")
	} else {
		err = profile.WriteMessage(websocket.TextMessage, responseData)
		if err != nil {
			logs.ErrorLog.Println("Error sending update message:", err)
			return err
		}
	}
	return nil
}

func notifyUser(uId int, command string, value any) error {
	u := update{"<system>", uId, command, value}
	err := u.sendToUser()
	if err != nil {
		logs.ErrorLog.Printf("Error sending update message: %v", err)
	}
	return err
}

func NotifyAll(command string, value any) error {
	u := update{"<system>", 0, command, value}
	err := u.send()
	if err != nil {
		logs.ErrorLog.Printf("Error sending update message: %v", err)
	}
	return err
}

func init() {
	structs.NotifyUser = notifyUser
}
