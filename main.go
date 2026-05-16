package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
}

func main(){


	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Print("error while running the server on port 3000", err)
	}
}