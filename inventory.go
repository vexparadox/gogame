package main

type Inventory struct{
	Items []Item `json:"items"`
}


type Item struct{
	Id int `json:"id"`
}