package main

import (
	"io/ioutil"
	"net/http"
)

func GetRequestBody(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	text, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	return text
}
