package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

//globals

var data_path string = "/Users/williammeaton/go/src/github.com/vexparadox/gogame/"
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}
var world_map Map
var users []User
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
	func_map := map[string](func(*User, []string)string){
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

			if msg_string == "new_conn"{
				reply = world_map.Welcome_text
				fmt.Printf("New connection from: %s\n", conn.RemoteAddr())
			} else if msg_string != "" {
				// Print the message to the console
				fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), msg_string)

				//split to an array of parameters given
				parameters = strings.Split(msg_string, " ")

				//allow login before validation check
				if parameters[1] == "login"{
					if is_user_id_valid(parameters[0]){
						reply = "You're already logged in!"
					} else {
						user_id := login_function(parameters[2:])
						//no user id, so create new one!
						if user_id == ""{

							unique_id, _ := uuid.NewV4()
							hashed_password, _ := bcrypt.GenerateFromPassword([]byte(parameters[2:][1]), 10)

							new_user := User{parameters[2:][0], *unique_id, hashed_password, 0}
							users = append(users, new_user)
							fmt.Printf("New user has registered! %s : %s\n", parameters[2:][0], unique_id.String())
							reply = "userid:"+unique_id.String()
						}else {
							//otherwise send back userid
							reply = "userid:"+user_id
						}
					}
				} else {
					if is_user_id_valid(parameters[0]) {
						user_data := get_user_data(parameters[0])
						//call the matching function with the said parameters
						func_to_call := func_map[parameters[1]]
						if func_to_call != nil{
							//remove first 2 as they are the login id and function name
							reply = func_to_call(user_data, parameters[2:])
						}	
					} else {
						reply = "Invalid login token, try refreshing or logging in again!"
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
