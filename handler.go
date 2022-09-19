package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	url      = "https://would-you-rather-api.abaanshanid.repl.co/"
	fontFile = "UbuntuMono-R.ttf"
)

var (
	fontBytes []byte
	f         *truetype.Font
)

func init() {
	fontBytes, err := ioutil.ReadFile(fontFile)

	if err != nil {
		log.Printf("Error during init: %s", err)
	}

	f, err = freetype.ParseFont(fontBytes)

	if err != nil {
		log.Printf("Error during init: %s", err)
	}
}

type WyrResponse struct {
	id   string
	Data string
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("received request...")
	blankResponse := &WyrResponse{}
	grabContent(blankResponse)
	log.Printf("got response from grabContent: %s", blankResponse.Data)
	log.Printf("converting to PNG...")

	imgBytes, err := contentToImage(blankResponse.Data)
	//err := png.Encode(w, img)
	if err != nil {
		log.Printf("ERROR: while encoding image to PNG: %s", err)
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "image/png")
	w.Write(imgBytes)
	return
}

func genErrorResponse(target *WyrResponse, err error) {
	target.Data = fmt.Sprintf("Would you rather that this worked, or that it didn't (cause there seems to be something wrong this is an error message haha: \n %s)?", err)
}

//func grabContent(target *WyrResponse) {
//	log.Printf("received request")
//	resp, err := http.Get(url)
//
//	if err != nil {
//		log.Printf("Encountered error while trying to get question: %s", err)
//		genErrorResponse(target, err)
//	}
//
//	defer resp.Body.Close()
//
//	err = json.NewDecoder(resp.Body).Decode(target)
//
//	if err != nil {
//		genErrorResponse(target, err)
//	}
//}

func grabContent(target *WyrResponse) {
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(int64(len(questions))))

	if err != nil {
		log.Printf("encountered error while trying to generate question index: %s", err)
		genErrorResponse(target, err)
	} else {
		target.Data = questions[randomNumber.Int64()]
	}
}

func contentToImage(content string) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, 9*len(content), 50))
	col := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	point := fixed.Point26_6{X: fixed.I(20), Y: fixed.I(30)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(content)

	b := new(bytes.Buffer)
	if err := png.Encode(b, img); err != nil {
		log.Println("unable to encode image.")
		return nil, err
	}
	return b.Bytes(), nil
}

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/api/WouldYouRather", helloHandler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
