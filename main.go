package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/ajstarks/svgo"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	output = flag.String("o", "output_1.svg", "output filename")

	firstLayer  = flag.String("pics_1", "", "pictures")
	firstRepeat = flag.Int("repeat_1", 0, "svg repeat. default indefinite")
	firstDur    = flag.String("dur_1", "0", "svg duration. default 1")
	firstBegin  = flag.String("begin_1", "0", "svg begin. default 0")

	secondLayer  = flag.String("pics_2", "", "pictures")
	secondRepeat = flag.Int("repeat_2", 0, "svg repeat. default indefinite")
	secondDur    = flag.String("dur_2", "0", "svg duration. default pics count")
	secondBegin  = flag.String("begin_2", "0", "svg begin. default 0")

	thirdLayer  = flag.String("pics_3", "", "pictures")
	thirdRepeat = flag.Int("repeat_3", 0, "svg repeat. default indefinite")
	thirdDur    = flag.String("dur_3", "0", "svg duration. default pics count")
	thirdBegin  = flag.String("begin_3", "0", "svg begin. default 0")

	fourthLayer  = flag.String("pics_4", "", "pictures")
	fourthRepeat = flag.Int("repeat_4", 0, "svg repeat. default indefinite")
	fourthDur    = flag.String("dur_4", "0", "svg duration. default pics count")
	fourthBegin  = flag.String("begin_4", "0", "svg begin. default 0")

	fifthLayer  = flag.String("pics_5", "", "pictures")
	fifthRepeat = flag.Int("repeat_5", 0, "svg repeat. default indefinite")
	fifthDur    = flag.String("dur_5", "0", "svg duration. default pics count")
	fifthBegin  = flag.String("begin_5", "0", "svg begin. default 0")
)

type Capturer struct {
	saved         *os.File
	bufferChannel chan string
	out           *os.File
	in            *os.File
}

// stdout capture start
func (c *Capturer) StartCapturingStdout() {
	c.saved = os.Stdout
	var err error
	c.in, c.out, err = os.Pipe()
	if err != nil {
		panic(err)
	}

	os.Stdout = c.out
	c.bufferChannel = make(chan string)
	go func() {
		var b bytes.Buffer
		io.Copy(&b, c.in)
		c.bufferChannel <- b.String()
	}()
}

// stop capture
func (c *Capturer) StopCapturingStdout() string {
	c.out.Close()
	os.Stdout = c.saved
	return <-c.bufferChannel
}

func base64Encode(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("A file read error occured: " + filePath)
		log.Fatal(err)
	}
	defer file.Close()
	fi, _ := file.Stat()
	size := fi.Size()
	data := make([]byte, size)
	file.Read(data)
	return base64.StdEncoding.EncodeToString(data)
}

func getImageSize(filePath string) (int, int) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("A file read error occured: " + filePath)
		log.Fatal(err)
	}
	defer file.Close()

	conf, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("image recognition error: " + filePath)
		log.Fatal(err)
	}
	return conf.Width, conf.Height
}

func makeSvgDef(id int, layer string, repeat int, duration string, begin string, s *svg.SVG) {
	layerList := strings.Split(layer, ",")

	var animateValues []string

	s.Def()
	for i := 0; i < len(layerList); i++ {
		layerId := "layer_" + strconv.Itoa(id+1) + "_" + strconv.Itoa(i+1)
		imageFilePath := layerList[i]
		ext := strings.ToLower(path.Ext(imageFilePath))

		width, height := getImageSize(imageFilePath)

		animateValues = append(animateValues, "#"+layerId)

		imageMimeType := ""
		if ext == ".png" {
			imageMimeType = "png"
		} else if ext == ".jpeg" {
			imageMimeType = "jpeg"
		}
		fmt.Println("<image id=\"" + layerId + "\" x=\"0\" y=\"0\" width=\"" + strconv.Itoa(width) + "\" height=\"" + strconv.Itoa(height) + "\" xlink:href=\"data:image/" + imageMimeType + ";base64," + base64Encode(imageFilePath) + "\" shape-rendering=\"crispEdges\" image-rendering=\"optimizeQuality\" />")
	}
	s.DefEnd()
	animateValuesStr := strings.Join(animateValues, ";")

	repeatCount := "indefinite"

	if repeat != 0 {
		repeatCount = strconv.Itoa(repeat)
	}

	durationNum := strconv.Itoa(1)

	if duration != "0" {
		durationNum = duration
	}

	beginNum := strconv.Itoa(0)

	if begin != "0" {
		beginNum = begin
	}

	fmt.Println("<use x=\"0\" y=\"0\">")
	fmt.Println("<animate attributeName=\"xlink:href\" begin=\"" + beginNum + "s\" dur=\"" + durationNum + "s\" repeatCount=\"" + repeatCount + "\" values=\"" + animateValuesStr + "\" />")
	fmt.Println("</use>")
}

func main() {
	flag.Parse()

	firstLayerList := strings.Split(*firstLayer, ",")
	svgStartViewWidth, svgStartViewHeight := getImageSize(firstLayerList[0])

	allLayer := []string{*firstLayer, *secondLayer, *thirdLayer, *fourthLayer, *fifthLayer}
	allLayerRepeat := []int{*firstRepeat, *secondRepeat, *thirdRepeat, *fourthRepeat, *fifthRepeat}
	allLayerDur := []string{*firstDur, *secondDur, *thirdDur, *fourthDur, *fifthDur}
	allLayerBegin := []string{*firstBegin, *secondBegin, *thirdBegin, *fourthBegin, *fifthBegin}

	c := &Capturer{}
	c.StartCapturingStdout()

	s := svg.New(os.Stdout)
	s.Startview(svgStartViewWidth, svgStartViewHeight, 0, 0, svgStartViewWidth, svgStartViewHeight)

	for i := 0; i < len(allLayer); i++ {
		if allLayer[i] != "" {
			makeSvgDef(i, allLayer[i], allLayerRepeat[i], allLayerDur[i], allLayerBegin[i], s)
		}
	}

	s.End()

	captured := c.StopCapturingStdout()

	output, _ := os.Create(*output)

	output.WriteString(captured)

	fmt.Println("done!")
}
