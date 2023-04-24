package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

// User contains user information
type User struct {
	FirstName      string     `validate:"required"`
	LastName       string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=130"`
	Email          string     `validate:"required,email"`
	FavouriteColor string     `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
	String         string     `validate:"is-awesome"`             // 自定义
}

// Address houses a users address information
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func ValidateMyVal(fl validator.FieldLevel) bool {

	return strings.Contains(fl.Field().String(), "zidingyi")
}

func main() {

	validate = validator.New()
	validate.RegisterValidation("is-awesome", ValidateMyVal)
	validateStruct()
	validateVariable()
}

func validateStruct() {

	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
	}

	user := &User{
		FirstName:      "Badger",
		LastName:       "Smith",
		Age:            135, // 不在  0 - 130
		Email:          "Badger.Smith@gmail.com",
		FavouriteColor: "#d3dbdc", // 正确 hexcolor|rgb|rgba|hsl|hsla
		Addresses:      []*Address{address},
		String:         "zidingyi",
	}

	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(user)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			//log.Fatalln(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("%#v\n", err)
			fmt.Println("Namespace:", err.Namespace())
			fmt.Println("Field:", err.Field())
			fmt.Println("StructNamespace:", err.StructNamespace())
			fmt.Println("StructField:", err.StructField())
			fmt.Println("Tag:", err.Tag(), err)
			fmt.Println("ActualTag:", err.ActualTag())
			fmt.Println("Kind:", err.Kind())
			fmt.Println("Type:", err.Type())
			fmt.Println("Value:", err.Value())
			fmt.Println("Param:", err.Param())
			fmt.Println()
		}
		// from here you can create your own error messages in whatever language you wish
		return
	}

	// save user to database
}

func validateVariable() {

	myEmail := "joeybloggs.gmail.com"
	myEmail2 := "joeybloggs@gmail.com"

	errs := validate.Var(myEmail, "required,email")

	if errs != nil {
		fmt.Println("myEmail", errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
		return
	}
	errs = validate.Var(myEmail2, "required,email")
	fmt.Println(errs)
	// email ok, move on
}
