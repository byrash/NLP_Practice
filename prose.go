package main

import (
	"fmt"

	"gopkg.in/jdkato/prose.v2"
)

func prose_practice() {
	// document, err :=
	// 	prose.NewDocument("Go is an open-source programming language created at Google.",
	// 		prose.WithExtraction(false))
	// document, err :=
	// 	prose.NewDocument("@jdkato, go to https://foo.com thanks :)",
	// 		prose.WithExtraction(false))
	document, err :=
		prose.NewDocument("Lebron James plays basketball in Los Angeles.")
	if err != nil {
		fmt.Println(err)
		panic("Error Creating Document")
	}
	fmt.Println("---------TOKENS--------")
	for _, token := range document.Tokens() {
		fmt.Println(token)
	}
	fmt.Println("---------ENTITYS--------")
	for _, enity := range document.Entities() {
		fmt.Println(enity.Text, enity.Label)
	}
}
