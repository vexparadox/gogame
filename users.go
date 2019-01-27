package main

import(
	"fmt"
	"os"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
    "io/ioutil"
)

type User struct
{
	Username string `json:"username"`
	Id [16]byte `json:"id"`
	Password []byte `json:"password"`
	Location_id int `json:"location_id"`
	Inventory Inventory `json:"inventory"`
}


//type fo writing out and loading users
type UserBundle struct
{
	Users []*User `json:"users"`
}

func save_users() bool {
	user_json, err := json.Marshal(UserBundle{users})
	if err != nil{
		return false
	}
	ioutil.WriteFile(data_path+"users.json", user_json, 0644)
	return true
}

func load_users() bool{
	json_users, err := os.Open(data_path + "users.json")

	if err != nil{
		fmt.Println("Failed to open users file.")
		return false
	} else {
		fmt.Printf("Opened users file %s.\n", json_users.Name())
	}
	defer json_users.Close()

	var bundle UserBundle

	decode_err := json.NewDecoder(json_users).Decode(&bundle)

	if decode_err != nil{
		fmt.Println("Failed to read users file")
		return false
	}

	users = bundle.Users
	fmt.Printf("Number of users: \"%d\"\n", len(users))
	return true
}

func is_user_id_valid(user_id string) bool{
	unique_id, err := uuid.ParseHex(user_id)
	if unique_id == nil || err != nil{
		return false
	}

	for _, user := range users{
		if *unique_id == user.Id{
			return true
		}
	}
	return false
}

func get_user_data(user_id string) *User{
	unique_id, err := uuid.ParseHex(user_id)
	if unique_id == nil || err != nil{
		return nil
	}

	for _, user := range users{
		if *unique_id == user.Id{
			return user
		}
	}
	return nil
}

func username_exists(username string) bool{
	for _, user := range users{
		if username == user.Username{
			return true
		}
	}
	return false
}

func create_new_user(parameters []string) (*User, string){
	if check_parameters(2, parameters) == false{
		return nil, incorrect_parameters
	}

	username := parameters[0]
	password := parameters[1]

	if username_exists(username){
		return nil, "The username '" + username + "' already exists."
	}

	unique_id, _ 		:= uuid.NewV4()
	hashed_password, _ 	:= bcrypt.GenerateFromPassword([]byte(password), 10)
	new_user := new(User)
	*new_user = User{username, *unique_id, hashed_password, 0, Inventory{}}
	users = append(users, new_user)
	save_users()
	return new_user, ""
}