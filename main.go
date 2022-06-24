package main

import (
	"fmt"

	"creategoapp/models"
)

func main() {
	fmt.Println("CreatGoApp starting...")

	t := models.NewTemplate("First template")
	fmt.Println(t)

	fmt.Println("CreateGoApp done!")
}
