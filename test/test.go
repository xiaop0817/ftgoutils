package main

import (
	"github.com/xiaop0817/ftgoutils/http"
	"log"
)

func main() {
	param := map[string]string{
		"accNum": "P10326383ABFEEEC0D1F9726E24752C0B12A6",
		"pwd":    "312222222aaaaaaaa22",
		"ip":     "192.168.0.7",
	}
	m := map[string]string{}
	log.Println(param)
	//testPost(param, m)

	http.Get("http://localhost:7090/3a/flow/available", map[string]interface{}{"accNum": "abc"}, &m)
	log.Println(m)
}

func testPost(param map[string]string, m map[string]string) {
	http.PostJson("http://localhost:7090/3a/account", param, &m, nil)
	log.Println(m)
}
