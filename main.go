package main

import (
	"fmt"
	"image"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const defaultImageWidth int = 400
const defaultImageHeight int = 400

func sendBadRequest(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(writer, "Bad url request: %s", request.URL)
}

func getImageFormat(response *http.Response) (imageFormat imaging.Format) {
	switch response.Header.Get("Content-Type") {
	case "image/jpeg", "image/jpg":
		imageFormat = imaging.JPEG

	case "image/gif":
		imageFormat = imaging.GIF

	case "image/png":
		imageFormat = imaging.PNG

	case "image/bmp":
		imageFormat = imaging.BMP

	case "image/tiff":
		imageFormat = imaging.TIFF

	default:
		imageFormat = imaging.JPEG
	}

	return imageFormat
}

func getImageSize(originalSize image.Point, inputWidth, inputHeight string) (imageWidth, imageHeight int, err error) {
	if len(inputWidth) == 0 &&
		len(inputHeight) == 0 {
		if originalSize.X > originalSize.Y {
			imageWidth = defaultImageWidth
			imageHeight = 0
		} else {
			imageWidth = 0
			imageHeight = defaultImageHeight
		}

	} else if len(inputWidth) > 0 {
		imageWidth, err = strconv.Atoi(inputWidth)
	} else if len(inputHeight) > 0 {
		imageHeight, err = strconv.Atoi(inputHeight)
	}

	return imageWidth, imageHeight, err
}

func handler(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	inputImageurl := query.Get("imageUrl")
	inputImageWidth := query.Get("imageWidth")
	inputImageHeight := query.Get("imageHeight")

	if len(inputImageurl) == 0 {
		sendBadRequest(writer, request)
		return
	}

	var webClient = &http.Client{Timeout: time.Second * 30}

	response, error := webClient.Get(inputImageurl)
	check(error)
	defer response.Body.Close()

	// image, error := imaging.Decode(response.Body)
	image, _, error := image.Decode(response.Body)
	check(error)

	imageFormat := getImageFormat(response)

	originImageSize := image.Bounds().Max

	if originImageSize.X <= defaultImageWidth &&
		originImageSize.Y <= defaultImageHeight {
		writer.WriteHeader(http.StatusOK)
		imaging.Encode(writer, image, imageFormat)
		return
	}

	imageWidth, imageHeight, error := getImageSize(originImageSize, inputImageWidth, inputImageHeight)
	check(error)

	destImage := imaging.Resize(image, imageWidth, imageHeight, imaging.Box)

	writer.WriteHeader(http.StatusOK)
	// fmt.Fprintln(writer, formatStr)
	imaging.Encode(writer, destImage, imageFormat)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.HandleFunc("/resize", handler)
	http.ListenAndServe(":8080", nil)
}

// https://www.socketloop.com/tutorials/golang-resize-image
// https://godoc.org/?q=resize+image
// https://github.com/h2non/imaginary
// https://github.com/willnorris/imageproxy
