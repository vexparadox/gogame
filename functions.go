package main

import(
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"github.com/nu7hatch/gouuid"
)

var incorrect_parameters string = "Incorrect number of parameters for that command!"

func check_parameters(required int, parameters []string) bool {
	return len(parameters) == required
}

//login function is special okay!
//if valid, it will return the user's id
func login_function(parameters []string) string{
	if check_parameters(2, parameters) == false{
		return incorrect_parameters
	}
	for _, user := range users{
		if user.username == parameters[0]{
			err := bcrypt.CompareHashAndPassword(user.password, []byte(parameters[1]))
			if err == nil{
				return "Invalid password provided!"
			} else {
				fmt.Printf("User has logged in: %s", user.username)
				return (*uuid.UUID)(&(user.id)).String()
			}
		}
	}
	return ""
}


func help_function(_ []string) string{
	return help_text
}


func look_function(parameters []string) string{
	if check_parameters(1, parameters){

	} else {

	}
	return help_text
}