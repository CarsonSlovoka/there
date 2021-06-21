package main

import (
	. "github.com/Gebes/there"
)

type User struct {
	Id   int    `json:"id,omitempty" validate:"required"`
	Name string `json:"name,omitempty" validate:"required"`
}

func main() {
	router := Router{
		Port:              8080,
		LogRouteCalls:     true,
		LogResponseBodies: true,
		AlwaysLogErrors:   true,
	}

	router.HandleGet("/user", func(request Request) Response {
		return ResponseData(StatusOK, User{
			1, "John",
		})
	})
	router.HandlePost("/user", func(request Request) Response {

		var user User
		err := request.ReadBody(&user)
		if err != nil {
			return ResponseData(StatusBadRequest, err)
		}

		// code

		return ResponseData(StatusOK, "Saved the user to the database")

	})

	err := router.Listen()
	if err != nil {
		panic(err)
	}
}
