package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"os"
	"github.com/gorilla/websocket"
)

type Map struct
{
	Welcome_text string `json:"welcome_text"`
	Width int `json:"width"`
	Height int `json:"height"`
	Tiles []Tile `json:"worldmap"`
}

type Tile struct
{
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Here string `json:"here"`
	North string `json:"n"`
	South string `json:"s"`
	East string `json:"e"`
	West string `json:"w"`
	Requireditems []int `json:"required_items"`
	Passable bool `json:"passable"`
	Quests []int `json:"quests"`
	Items []int `json:"items"`
}

var data_path string = "/Users/williammeaton/go/src/github.com/vexparadox/gogame/"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var world_map Map

func load_map() bool {
	//load up the map data
	json_map, err := os.Open(data_path + "map.json")

	if err != nil{
		fmt.Println("Failed to open map file")
		return false
	} else {
		fmt.Printf("Opened map file %s\n", json_map.Name())
	}
	defer json_map.Close()

	decode_err := json.NewDecoder(json_map).Decode(&world_map)

	if decode_err != nil{
		fmt.Println("Failed to read map file")
		return false
	}

	fmt.Printf("Welcome text: \"%s\"\n", world_map.Welcome_text)
	fmt.Printf("World Size: %d x %d\n", world_map.Width, world_map.Height)

	return true
}



func main() {

	if load_map() == false{
		return
	}

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var reply string
			var msg_string = string(msg)

			if msg_string == "new_user"{
				reply = world_map.Welcome_text
				fmt.Printf("New user from: %s\n", conn.RemoteAddr())
			} else {
				// Print the message to the console
				fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), msg_string)
			}

			// Write message back to browser if there
			if reply != ""{
				if err = conn.WriteMessage(msgType, []byte(reply)); err != nil {
					return
				}
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, data_path+"websockets.html")
	})

	http.ListenAndServe(":8080", nil)
}
