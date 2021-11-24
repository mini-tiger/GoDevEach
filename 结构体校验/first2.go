package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type User struct {
	FirstName      string `json:firstname`
	LastName       string `json:lastname`
	Age            uint8  `validate:"gte=0,lte=130"`
	Email          string `validate:"required,email"`
	FavouriteColor string `validate:"hexcolor|rgb|rgba"`
}

var validate *validator.Validate

func main() {
	validate = validator.New()

	validate.RegisterStructValidation(UserStructLevelValidation, User{})

	user := &User{
		FirstName:      "",
		LastName:       "",
		Age:            30,
		Email:          "TestFunc@126.com",
		FavouriteColor: "#000",
	}

	err := validate.Struct(user)
	if err != nil {
		fmt.Println(err)
	}
}

func UserStructLevelValidation(sl validator.StructLevel) {
	user := sl.Current().Interface().(User)

	if len(user.FirstName) == 0 && len(user.LastName) == 0 {
		sl.ReportError(user.FirstName, "FirstName", "firstname", "firstname", "")
		sl.ReportError(user.LastName, "LastName", "lastname", "lastname", "")
	}
}
