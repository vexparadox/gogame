package main

import(
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var incorrect_parameters string = "Incorrect number of parameters for that command!"

func check_parameters(required int, parameters []string) bool {
	return len(parameters) == required
}

//login function is special okay!
//if valid, it returns the logged in user
//otherwise user will be nil and a string reason why
func login_function(parameters []string) (*User, string) {
	if check_parameters(2, parameters) == false{
		return nil, incorrect_parameters
	}
	for _, user := range users{
		if user.username == parameters[0]{
			err := bcrypt.CompareHashAndPassword(user.password, []byte(parameters[1]))
			if err != nil{
				return nil, "Invalid password provided!"
			} else {
				fmt.Printf("User has logged in: %s\n", user.username)
				return &user, ""
			}
		}
	}
	return nil, "Username '" + parameters[0] + "' doesn't exist! Use the 'register' command to create a new account."
}


func help_function(_ *User, _ []string) string{
	return help_text
}

func profile_function(user_data *User, _[]string) string{
	return "Username: " + user_data.username
}

func go_function(_ *User, _[]string) string{
	return "Not done yet!"
}

func inv_function(_ *User, _[]string) string{
	return "Not done yet!"
}

func pickup_function(_ *User, _[]string) string{
	return "Not done yet!"
}

func say_function(_ *User, _[]string) string{
	return "Not done yet!"
}

func quest_function(_ *User, _[]string) string{
	return "Not done yet!"
}

func look_function(user_data *User, parameters []string) string{
	current_tile := world_map.get_tile_for_user(user_data)

	if current_tile != nil{
		if check_parameters(1, parameters){
			if parameters[0] == "e"{
				return current_tile.East
			} else if parameters[0] == "w"{
				return current_tile.West
			} else if parameters[0] == "n"{
				return current_tile.North
			} else if parameters[0] == "s"{
				return current_tile.South
			}else{
				return "Invalid direction given."
			}
		} else {
			return current_tile.Here
		}
	}
	fmt.Printf("Failed to find valid tile for user %s", user_data.username)
	return "Internal server error!"
}