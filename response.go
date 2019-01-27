package main

import(
	"strings"
	"fmt"
	"net/http"
	"github.com/nu7hatch/gouuid"
)

func handle_socket(w http.ResponseWriter, r *http.Request){
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


			if is_user_id_valid(parameters[0]) {
				if parameters[1] == "login"{
					reply = "You're already logged in!"
				}else{
					user_data := get_user_data(parameters[0])
					//call the matching function with the said parameters
					func_to_call := response_function_map[parameters[1]]
					if func_to_call != nil{
						//remove first 2 as they are the login id and function name
						reply = func_to_call(user_data, parameters[2:])
					}	
				}
			}else{
				if parameters[1] == "login"{
					user_id := login_function(parameters[2:])
					//no user id, so create new one!
					if user_id == ""{

						var new_user *User = create_new_user(parameters[2:][0], parameters[2:][1])

						fmt.Printf("New user has registered! %s\n", new_user.username)
						reply = "userid:"+(*uuid.UUID)(&(new_user.id)).String()
					}else {
						//otherwise send back userid
						reply = "userid:"+user_id
					}
				}else if parameters[1] == "help" || parameters[1] == "?"{
					//this check is jank af
					reply = help_function(nil, parameters[2:])
				}else{
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
}