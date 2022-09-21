package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
)

const (
	url        = "https://would-you-rather-api.abaanshanid.repl.co/"
	fontFile   = "UbuntuMono-R.ttf"
	lineLength = 150
)

var (
	fontBytes []byte
	f         *truetype.Font
)

func init() {
	fontBytes, err := os.ReadFile(fontFile)

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
	contentLength := len(content)
	width := 8 * lineLength
	height := 25 * ((contentLength / lineLength) + 1)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	blk := color.RGBA{R: 50, G: 49, B: 48, A: 255}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.Set(i, j, blk)
		}
	}

	whiteText := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	lines := splitText(content)

	for i, line := range lines {
		point := fixed.Point26_6{X: fixed.I(20), Y: fixed.I(20 * (i + 1))}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(whiteText),
			Face: basicfont.Face7x13,
			Dot:  point,
		}
		d.DrawString(line)

	}

	b := new(bytes.Buffer)
	if err := png.Encode(b, img); err != nil {
		log.Println("unable to encode image.")
		return nil, err
	}
	return b.Bytes(), nil
}

func splitText(content string) []string {
	contentLength := len(content)
	if contentLength < lineLength {
		return []string{content}
	}
	var lines []string

	words := strings.Split(content, " ")

	for i := 0; i < len(words); {
		currentBuilder := strings.Builder{}
		currentLength := 0

		for currentLength < lineLength && i < len(words) {
			currentBuilder.WriteString(words[i] + " ")
			currentLength = currentLength + len(words[i]) + 1
			i++
		}
		lines = append(lines, currentBuilder.String())
	}

	log.Printf("lines are %s", lines)
	return lines
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
