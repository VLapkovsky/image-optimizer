package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

	responce, error := http.Get(inputImageurl)
	check(error)
	defer responce.Body.Close()

	imageBody, error := ioutil.ReadAll(responce.Body)
	check(error)

	contentLenght, _ := strconv.ParseUint(request.URL.Query().Get("Content-Lenght"), 10, 64)
	writer.Header().Set("Content-Lenght", strconv.FormatUint(contentLenght, 10))
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
