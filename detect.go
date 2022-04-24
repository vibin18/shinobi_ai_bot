package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type RequestData struct {
	Id           string `json:"id"`
	DetectorName string `json:"detector_name"`
	Data         string `json:"data"`
	Image        string `json:"image"`
	Throttle     int    `json:"throttle"`
	Detect       struct {
		Person     float64 `json:"person"`
		Cat        float64 `json:"cat"`
		Bicycle    float64 `json:"bicycle"`
		Car        float64 `json:"car"`
		Motorcycle float64 `json:"motorcycle"`
		Truck      float64 `json:"truck"`
		Bird       float64 `json:"bird"`
		Dog        float64 `json:"dog"`
		Horse      float64 `json:"horse"`
		cow        float64 `json:"cow"`
		elephant   float64 `json:"elephant"`
		bear       float64 `json:"bear"`
		umbrella   float64 `json:"umbrella"`
		handbag    float64 `json:"handbag"`
	} `json:"detect"`
}

type ResponseData struct {
	Id         string       `json:"id"`
	ImageData  string       `json:"image"`
	Detections []Detections `json:"detections"`
	Errors     string       `json:"error"`
}

type Detections struct {
	Bottom     float64 `json:"bottom"`
	Confidence float64 `json:"confidence"`
	Label      string  `json:"label"`
	Left       float64 `json:"left"`
	Right      float64 `json:"right"`
	Top        float64 `json:"top"`
}

func NewRequestData() *RequestData {
	return &RequestData{}
}

func (c *RequestData) SetdetectOption(val float64) {
	c.Detect.Person = val
	c.Detect.Cat = val
	c.Detect.Bicycle = val
	c.Detect.Car = val
	c.Detect.Motorcycle = val
	c.Detect.Truck = val
	c.Detect.Bird = val
	c.Detect.Dog = val
	c.Detect.Horse = val
	c.Detect.cow = val
	c.Detect.elephant = val
	c.Detect.bear = val
	c.Detect.umbrella = val
	c.Detect.handbag = val
}

func (c *RequestData) Setdetector_name(val string) {
	c.DetectorName = val
}

func (c *RequestData) SetFiledata(val []byte) {
	encodedImage := b64.StdEncoding.EncodeToString(val)
	c.Data = encodedImage
}

func (c *RequestData) SetId(val string) {
	c.Id = val
}

func (c *RequestData) SetImageType(val string) {
	c.Image = val
}

func (c *RequestData) SetThrottle(val int) {
	c.Throttle = val
}

//func DetectImage(imageFile string, minProb int64, wg *sync.WaitGroup) (error, ResponseData ) {
func DetectImage(imageFile []byte, minProb float64) (error, ResponseData) {
	var ret ResponseData
	con := NewRequestData()
	con.SetFiledata(imageFile)
	con.Setdetector_name("default")
	con.SetdetectOption(minProb)
	con.SetId("123")
	con.SetImageType("jpeg")
	con.SetThrottle(1)

	prettyJSON, err := json.MarshalIndent(con, "", "    ")

	if err != nil {
		log.Panic(err)
	}

	requestBody := bytes.NewBuffer(prettyJSON)

	fmt.Print(requestBody)

	resp, err := http.Post("http://"+arg.DoodsServer+"/detect", "application/json", requestBody)
	if err != nil {
		log.Infof("Error requesting")
		log.Panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Infof("Error from response")
		log.Panic(err)
	}

	err = json.Unmarshal(body, &ret)
	if err != nil {
		log.Infof("Error unmarsheling response")
		log.Panic(err)
	}
	log.Infof("Detection finished")
	//wg.Done()
	return err, ret
}
