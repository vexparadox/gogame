package main

import(
	"strings"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/nu7hatch/gouuid"
)

type SocketMessage struct
{
	Token string `json:"token"`
	Message string `json:"message"`
}

func handle_socket(w http.ResponseWriter, r *http.Request){
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		//out reply struct, we will send this back as json
		var server_reply SocketMessage

		//get data out of json client sent
		var client_message SocketMessage
		json_err := json.Unmarshal(msg, &client_message)
		if json_err != nil{
			server_reply.Message = "Internal server error!"
			fmt.Printf("Invalid json from client: %s", string(msg))
		}else {
			var parameters []string = strings.Split(client_message.Message, " ")	//split to an array of parameters given
			//cut the first one out
			command_string := parameters[0]
			parameters = parameters[1:]

			//on first connection we send the welcome text of the world
			if client_message.Message == "new_conn"{
				server_reply.Message = world_map.Welcome_text
				fmt.Printf("%s: New connection\n", conn.RemoteAddr())
			} else if client_message.Message != "" {
				//otherwise we check if the token they've provided it valid
				if is_user_id_valid(client_message.Token) {
					if command_string == "login"{
						server_reply.Message = "You're already logged in!"
					}else{
						user_data := get_user_data(client_message.Token)
						//call the matching function with the said parameters
						func_to_call := response_function_map[command_string]
						if func_to_call != nil{
							//remove first 2 as they are the login id and function name
							server_reply.Message = func_to_call(user_data, parameters)
							fmt.Printf("%s: %s used command '%s'\n", conn.RemoteAddr(), user_data.username, command_string)
						}	
					}
				}else{
					if command_string == "login"{
						user_id := login_function(parameters)
						//no user id, so create new one!
						if user_id == ""{
							var new_user *User = create_new_user(parameters[0], parameters[1])
							fmt.Printf("%s: New user has registered! %s\n", conn.RemoteAddr(), new_user.username)
							server_reply.Token = (*uuid.UUID)(&(new_user.id)).String()
						}else {
							//otherwise send back userid of the exisiting user
							server_reply.Token = user_id
						}
					}else if command_string == "help" || command_string == "?"{
						//this check is jank af but we need the player to be able to ask for help before they login
						server_reply.Message = help_function(nil, parameters)
					}else{
						server_reply.Message = "Invalid login token, try refreshing or logging in again!"
					}
				}
			}

			//if we have no reply, it's because we didn't parse correctly
			if server_reply.Message == ""{
				server_reply.Message = "\""+command_string+"\" was an unrecognised command, try using help to see what you can do."
			}
		}

		//turn our reply into json
		json_bytes, json_marshal_err := json.Marshal(server_reply)

		if json_marshal_err != nil{
			fmt.Printf("Interal server error! Failed to marshal json data.\n")
		}

		// Write message back to browser if there
		if err = conn.WriteMessage(msgType, json_bytes); err != nil {
			return
		}
	}
}