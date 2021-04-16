package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	file, err := os.OpenFile("/home/senpai/docker-elk/logs/ololo.log",os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}

	log.SetOutput(file)

	type LogMessage struct {
		Time time.Time `json:"time"`
		Message string `json:"message"`
	}

	i:=0
	for range time.Tick(time.Second) {
		i++
		msg := LogMessage{
			Time:    time.Now(),
			Message: "tick number : " + strconv.Itoa(i),
		}
		bytes, _ :=json.Marshal(msg)

		fmt.Println(string(bytes))
		log.Println(string(bytes))
	}
}
