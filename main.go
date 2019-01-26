package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/websocket"
)

//globals

var data_path string = "/Users/williammeaton/go/src/github.com/vexparadox/gogame/"
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var world_map Map
var help_text string

func load_help_text() bool{
	help_file_bytes, err := ioutil.ReadFile(data_path+"help.txt")
	if err != nil{
		fmt.Println("Failed to load help file")
		return false
	}
	help_text = string(help_file_bytes)
	return true
}

func main() {

	if load_map() == false{
		return
	}

	if load_help_text() == false{
		return
	}

	//map of functions in functions.go to string commands
	func_map := map[string](func([]string)string){
		"help"  : help_function,
		"?"		: help_function,
		"look"	: look_function,
		"l"		: look_function,
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
			var parameters []string

			if msg_string == "new_user"{
				reply = world_map.Welcome_text
				fmt.Printf("New connection from: %s\n", conn.RemoteAddr())
			} else if msg_string != "" {
				// Print the message to the console
				fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), msg_string)

				//split to an array of parameters given
				parameters = strings.Split(msg_string, " ")

				//check for userid
				if parameters[0] == ""{
					reply = "Invalid login token, try using the login command again"
				} else {

					//now check if user id is valid

					//call the matching function with the said parameters
					func_to_call := func_map[parameters[1]]
					if func_to_call != nil{
						reply = func_to_call(parameters)
					}	
				}
			}

			//if we have no reply, it's because we didn't parse correctly
			if reply == ""{
				reply = "\""+parameters[1]+"\" was an unrecognised command, try using help to see what you can do."
			}
			// Write message back to browser if there
			if err = conn.WriteMessage(msgType, []byte(reply)); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, data_path+"websockets.html")
	})

	http.ListenAndServe(":8080", nil)
}
