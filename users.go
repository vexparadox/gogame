package main

import(
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)
type User struct
{
	username string
	id [16]byte
	password []byte
	location_id int
}

func is_user_id_valid(user_id string) bool{
	unique_id, err := uuid.ParseHex(user_id)
	if unique_id == nil || err != nil{
		return false
	}

	for _, user := range users{
		if *unique_id == user.id{
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
		if *unique_id == user.id{
			return &user
		}
	}
	return nil
}

func username_exists(username string) bool{
	for _, user := range users{
		if username == user.username{
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
	new_user := User{username, *unique_id, hashed_password, 0}
	users = append(users, new_user)
	return &new_user, ""
}