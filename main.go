package main

import (
	"fmt"
	"flag"
	"time"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/websocket"
)

//globals that could be removed maybe (but are globals the real enemy?!)
var data_path string
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

var response_function_map map[string](func(*User, []string)string)

var world_map Map
var users []*User
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
	//get the data directory from the cli
	data_path_ptr := flag.String("data_dir", ".", "The directory where the data is")
	flag.Parse()
	if data_path_ptr != nil{
		data_path = *data_path_ptr
	}

	if load_map() == false{
		return
	}

	if load_help_text() == false{
		return
	}

	load_users() // try to load users

	response_function_map = map[string](func(*User, []string)string){
		"look"	: look_function,
		"l"		: look_function,
		"profile" : profile_function,
		"go" : go_function,
		"inv" : inv_function,
		"pickup" : pickup_function,
		"quest" : quest_function,
		"say" : say_function,
	}

	//save the game every 10 seconds or so
	ticker 	:= time.NewTicker(10 * time.Second)
	go func() {
	    for {
	       select {
	        case <- ticker.C:
	            save_users()
	        }
	    }
	 }()


	http.HandleFunc("/ws", handle_socket)
	http.Handle("/", http.FileServer(http.Dir(data_path+"html")))
	http.ListenAndServe(":8080", nil)
}
