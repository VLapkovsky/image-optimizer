package main

import (
	"fmt"
	"net/http"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func handler(writer http.ResponseWriter, request *http.Request) {
	// image, error := ioutil.ReadFile("/Users/vitaly.lapkovsky/Documents/test.png")
	// check(error)
	query := request.URL.Query()
	imageUrl := query.Get("imageUrl")

	if len(imageUrl) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "Bad url request: %s", request.URL)
		return
	}

	fmt.Fprint(writer, imageUrl)

	// writer.Header().Set("Content-Lenght", strconv.Itoa(len(image)))
	// writer.WriteHeader(http.StatusOK)
	// writer.Write(image)
	// fmt.Fprintf(writer, "Hi there, I love %s!", request.URL.Path[1:])
}

func main() {
	http.HandleFunc("/resize", handler)
	http.ListenAndServe(":8080", nil)
}
