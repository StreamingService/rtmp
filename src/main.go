package main

import (
	"log"

	"server"
)

func main() {
	log.Println("rtmp serv")

	s, err := server.New(":1935")
	if (err != nil) {
		log.Fatalln(err)
		return
	}
	defer s.Close()

	err = s.Serv()
	if (err != nil) {
		log.Fatalln(err)
	}
	
}