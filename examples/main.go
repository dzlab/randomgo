package main

import (
	"github.com/dzlab/randomgo"
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
	kv := random.NewKVEncoder("=", "&")
	json := random.NewJsonEncoder()
	for i := 0; i < 10; i++ {
		log.Println(i, ">", kv.Encode(object))
		log.Println(i, ">", json.Encode(object))

	}
}
