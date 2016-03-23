package main

import (
	"github.com/dzlab/random.go"
	"log"
)

var ()

func main() {
	p := random.NewParser()
	object, err := p.Parse("main.yml")
	if err != nil {
		panic(err)
	}
	// generate some data
	for i := 0; i < 10; i++ {
		log.Println(i, ">", object.GetKV("=", "&"))
		log.Println(i, ">", object.GetJSON())

	}
}
