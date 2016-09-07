package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func handler(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	inputImageurl := query.Get("imageUrl")

	if len(inputImageurl) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "Bad url request: %s", request.URL)
		return
	}

	var webClient = &http.Client{Timeout: time.Second * 10}

	responce, error := webClient.Get(inputImageurl)
	check(error)
	defer responce.Body.Close()

	imageBody, error := ioutil.ReadAll(responce.Body)
	check(error)

	for key, value := range responce.Header {
		writer.Header().Set(key, strings.Join(value, ""))
	}

	// writer.Header().Set("Content-Length", responce.Header.Get("Content-Length"))
	writer.WriteHeader(http.StatusOK)

	writer.Write(imageBody)

	// fmt.Fprint(writer, imageUrl)

	// writer.Header().Set("Content-Lenght", strconv.Itoa(len(image)))
	// writer.WriteHeader(http.StatusOK)
	// writer.Write(image)
	// fmt.Fprintf(writer, "Hi there, I love %s!", request.URL.Path[1:])
}

func main() {
	http.HandleFunc("/resize", handler)
	http.ListenAndServe(":8080", nil)
}

// https://www.socketloop.com/tutorials/golang-resize-image
// https://godoc.org/?q=resize+image
// https://github.com/h2non/imaginary
// https://github.com/willnorris/imageproxy
