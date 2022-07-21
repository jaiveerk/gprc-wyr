package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	blankResponse := &WyrResponse{}
	grabContent(blankResponse)
	log.Printf("got response from grabContent: %s", blankResponse.Data)

	log.Printf("converting to PNG...")

	//imgBytes, err := contentToImage(blankResponse.Data)
	img := contentToImage(blankResponse.Data)
	err := png.Encode(w, img)
	if err != nil {
		log.Printf("ERROR: while encoding image to PNG: %s", err)
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/octet-stream")
	//w.Write(imgBytes)
	return
}

func genErrorResponse(target *WyrResponse, err error) {
	target.Data = fmt.Sprintf("Would you rather that this worked, or that it didn't (cause there seems to be something wrong this is an error message haha: \n %s)?", err)
}

func grabContent(target *WyrResponse) {
	log.Printf("received request")
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Encountered error while trying to get question: %s", err)
		genErrorResponse(target, err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(target)

	if err != nil {
		genErrorResponse(target, err)
	}
}

func contentToImage(content string) image.Image {
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

	return img
	//text := strings.Split(content, "\n")
	//
	//var fontSize float64 = 12
	//
	//fgColor := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	//bgColor := color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
	//
	//fg := image.NewUniform(fgColor)
	//bg := image.NewUniform(bgColor)
	//
	//imgLength := 7 * len(content)
	//
	//rgba := image.NewRGBA(image.Rect(0, 0, imgLength, 50))
	//draw.Draw(rgba, rgba.Bounds(), bg, image.Pt(0, 0), draw.Src)
	//c := freetype.NewContext()
	//c.SetDPI(72)
	//c.SetFont(f)
	//c.SetFontSize(fontSize)
	//c.SetClip(rgba.Bounds())
	//c.SetDst(rgba)
	//c.SetSrc(fg)
	//c.SetHinting(font.HintingNone)
	//
	//textXOffset := len(content) / 2
	//textYOffset := 10 + int(c.PointToFixed(fontSize)>>6) // Note shift/truncate 6 bits first
	//
	//pt := freetype.Pt(textXOffset, textYOffset)
	//for _, s := range text {
	//	_, err := c.DrawString(strings.Replace(s, "\r", "", -1), pt)
	//	if err != nil {
	//		return nil, err
	//	}
	//	pt.Y += c.PointToFixed(fontSize * 1.5)
	//}
	//
	//b := new(bytes.Buffer)
	//if err := png.Encode(b, rgba); err != nil {
	//	log.Printf("ERROR: unable to encode image: %s", err)
	//	return nil, err
	//}
	//return b.Bytes(), nil
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
