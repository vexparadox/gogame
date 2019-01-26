package main

import(
	"github.com/nu7hatch/gouuid"
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