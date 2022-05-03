package main

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"github.com/andersfylling/snowflake"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"github.com/vibin18/doods_client/webhooks"
	"io/ioutil"
	"net/http"
	"os"
)

type opts struct {
	File            string  `short:"f"  long:"file"      env:"FILE"  description:"Filename for detecting" required:"true"`
	DoodsServer     string  `           long:"server"      env:"DOODS_SERVER"  description:"Server name or IP of doods2 server and port number" default:"192.168.178.81:8099" required:"true"`
	DiscordToken    string  `           long:"token"      env:"DISCORD_TOKEN"  description:"Discord Webhook token" required:"true"`
	WebhookId       uint64  `           long:"webhook"      env:"DISCORD_WEBHOOK_ID"  description:"Discord Webhook ID" required:"true"`
	DetectorName    string  `long:"detector" env:"DETECTOR_NAME" description:"doods2 supports tflite, tensorflow, pytorch. If not specified,'default' will be used if it exists"`
	MinConfidence   float64 `           long:"mincon"      env:"MINIMUM_CONFIDENCE"  description:"Minimum confidence level and Max is 100" default:"50"`
	CameraId        string  `            long:"camera"      env:"CAMERA_NAME"  description:"Name of the camera" required:"false"`
	ShinobiExporter string  `       long:"exporter"      env:"SHINOBI_EXPORTER"  description:"Server name or IP of shinobi_exporter and port number" required:"false"`
}

var (
	argparser *flags.Parser
	arg       opts
)

func initArgparser() {
	argparser = flags.NewParser(&arg, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

func main() {

	initArgparser()
	webhook := snowflake.Snowflake(arg.WebhookId)
	minConfidence := arg.MinConfidence
	if !(minConfidence >= 0 && minConfidence <= 100) {
		log.Panicf("Minimum confidence should between 0-100, got %f", minConfidence)
	}

	fmt.Println(arg.ShinobiExporter)
	if len(arg.ShinobiExporter) > 0 && len(arg.CameraId) > 0 {
		shinobiExporterUrl := fmt.Sprintf("http://%s/hit", arg.ShinobiExporter)
		ReqBodyStr := fmt.Sprintf(`{"Name": "%s"}`, arg.CameraId)
		jsonStr := []byte(ReqBodyStr)
		req, err := http.NewRequest("GET", shinobiExporterUrl, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		log.Infof("Sending hit to exporter %s", shinobiExporterUrl)
		log.Infof("Sending data %s", ReqBodyStr)
		resp, err := client.Do(req)
		if err != nil {
			log.Panicf("Error sending exporter request %s ", err)
		}
		if resp.StatusCode == 200 {
			log.Infof("Hit succesfully sent")
		} else {
			log.Warnf("Failed! exporter returned %v", resp.StatusCode)
		}
		defer resp.Body.Close()

	}

	var ConfidenceMapList []map[string]float64

	if _, err := os.Stat(arg.File); os.IsNotExist(err) {
		log.Panicf("File %s NOT found!", arg.File)
	}

	byteImage, err := ioutil.ReadFile(arg.File)
	if err != nil {
		log.Panicf("Byte conversion failed!")
	}

	err, result := DetectImage(byteImage, minConfidence, arg.DetectorName)
	if err != nil {
		log.Panicf("Failed to detect image")
	}

	if len(result.Errors) > 0 {
		log.Panicf(result.Errors)
	}

	fmt.Print(result.Detections)
	fmt.Println(result.ImageData)
	//byteImageReturned := b64.StdEncoding.EncodeToString([]byte(result.ImageData))
	byteImageReturned, err := b64.StdEncoding.DecodeString(result.ImageData)
	if err != nil {
		log.Panicf("Failed to decode image data")
	}

	for _, v := range result.Detections {
		if v.Confidence >= minConfidence {
			itemMap := map[string]float64{
				v.Label: v.Confidence,
			}
			ConfidenceMapList = append(ConfidenceMapList, itemMap)
		}
	}

	webhooks.NotifyDiscord(webhook, arg.DiscordToken, byteImageReturned, "alert.jpg", minConfidence, ConfidenceMapList)
}
