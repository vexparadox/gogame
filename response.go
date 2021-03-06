/*
TODO:
	- fix the way commands are called, the odd if login etc is dumb, get rid of it
	- add in remaining commands
*/

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

				user_data := get_user_data(client_message.Token)

				if command_string == "login"{
					if user_data != nil{
						server_reply.Message = "You're already logged in."
					}else{
						user_data, reason_str := login_function(parameters)
						//if we get a nil back, it means it was invalid, so send the reason back to the client
						if user_data == nil{
							server_reply.Message = reason_str
						}else{
							server_reply.Token = (*uuid.UUID)(&(user_data.Id)).String()
							server_reply.Message = "Welcome back, " + user_data.Username + "."
						}
					}
				} else if command_string == "register" {
					if user_data != nil{
						server_reply.Message = "You can't register when you're logged in."
					}else{
						new_user, reason_str := create_new_user(parameters)
						//if we get a nil back, it means it was invalid, so send the reason back to the client
						if new_user == nil{
							server_reply.Message = reason_str
						}else{
							server_reply.Token = (*uuid.UUID)(&(new_user.Id)).String()
							server_reply.Message = "You've registered with the username '" + parameters[0] + "'."
						}
					}
				} else if command_string == "help" || command_string == "?"{
					server_reply.Message = help_text
				} else if user_data != nil {

					//call the matching function with the said parameters
					func_to_call := response_function_map[command_string]
					if func_to_call != nil{
						//send back the response of the function
						server_reply.Message = func_to_call(user_data, parameters)
						fmt.Printf("%s: %s used command '%s'\n", conn.RemoteAddr(), user_data.Username, command_string)
					}
				}else{
					server_reply.Message = "You're not logged in. Try using 'login' or 'register' if you're new."
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