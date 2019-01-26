package main
import(
	"encoding/json"
	"os"
	"fmt"
)

type Map struct
{
	Welcome_text string `json:"welcome_text"`
	Width int `json:"width"`
	Height int `json:"height"`
	Tiles []Tile `json:"worldmap"`
}

type Tile struct
{
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Here string `json:"here"`
	North string `json:"n"`
	South string `json:"s"`
	East string `json:"e"`
	West string `json:"w"`
	Requireditems []int `json:"required_items"`
	Passable bool `json:"passable"`
	Quests []int `json:"quests"`
	Items []int `json:"items"`
}

func load_map() bool {
	//load up the map data
	json_map, err := os.Open(data_path + "map.json")

	if err != nil{
		fmt.Println("Failed to open map file")
		return false
	} else {
		fmt.Printf("Opened map file %s\n", json_map.Name())
	}
	defer json_map.Close()

	decode_err := json.NewDecoder(json_map).Decode(&world_map)

	if decode_err != nil{
		fmt.Println("Failed to read map file")
		return false
	}

	fmt.Printf("Welcome text: \"%s\"\n", world_map.Welcome_text)
	fmt.Printf("World Size: %d x %d\n", world_map.Width, world_map.Height)

	return true
}
