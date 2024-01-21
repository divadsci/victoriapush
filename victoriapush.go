/*
victoriapush pushes metrics to VictoriaMetrics using the prometheus exposition format

send metrics with custom labels in a map

assign global labels
*/
package victoriapush

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Vic struct {
	URL          string            //URL of victoriametrics server
	globalLabels map[string]string //labels to be applied to every metric pushed
	maxWait      float64           //maximum time to wait before sending whatever is in buffer
	maxQueueSize int               //maximum size for metrics queue before sending
	msgQ         chan string       //queue of messages waiting to be pushed#
	stopPushing  chan bool         //stops push loop
	pushTicker   *time.Ticker      //fires regular pushes
	instantPush  chan bool         //forces a push

}

type DataPoint struct {
	Metric string
	Value  float64
	Labels map[string]string
}

// add URL and globalLabels
func NewVictoriaPush(URL string, globalLabels map[string]string) *Vic {
	v := new(Vic)
	v.URL = URL
	v.globalLabels = globalLabels
	v.maxWait = 1
	v.maxQueueSize = 10

	//make msgQ 10x larger than the max queue size to prevent overflow data loss
	v.msgQ = make(chan string, v.maxQueueSize*10)

	//set the ticker time for regular pushes
	v.pushTicker = time.NewTicker(time.Duration(int64(v.maxWait * float64(time.Second))))
	v.stopPushing = make(chan bool)
	v.instantPush = make(chan bool)

	//start victoriapush loop that sends data based on time limit or queue size
	go v.pushLoop()
	return v
}

// Set VictoriaMetrics server URL
func (v *Vic) SetURL(URL string) {
	//todo: do some checking here?
	v.URL = URL
}

// Set the limits of the metric queue, maxWait for maximum wait for more data before sending, maxQueueSize for maximum number of metrics in push queue before sending
func (v *Vic) SetQueueLimits(maxWait float64, maxQueueSize int) {
	v.maxWait = maxWait
	v.maxQueueSize = maxQueueSize
}

// Replace globalLabels map with a new one
func (v *Vic) ReplaceGlobalLabels(globalLabels map[string]string) {
	v.globalLabels = globalLabels
}

// Add a map of labels to globalLabels
func (v *Vic) AddGlobalLabels(globalLabels map[string]string) {

	//iterate over new globalLabels map and add each to Vic data
	for label, val := range globalLabels {
		v.globalLabels[label] = val
	}
}

// Remove the list of labels from the globalLabels map IF it matches
func (v *Vic) RemGlobalLabels(labelsForRem []string) {

	//delete all labels in labelsForRem from globalLabels
	for _, label := range labelsForRem {
		delete(v.globalLabels, label)
	}
}

// Add dataPoint to queue ready to be pushed
func (v *Vic) EnqueueDataPoint(dataPoint DataPoint) {
	v.msgQ <- v.dataPointToExpo(dataPoint)
}

// convert DataPoint to prometheus exposition format ready for pushing to server
func (v *Vic) dataPointToExpo(dataPoint DataPoint) string {
	expoString := fmt.Sprintf("%s{", dataPoint.Metric)

	//go through all labels and convert into labe="value", format
	for label, data := range dataPoint.Labels {
		expoString = expoString + fmt.Sprintf("%s=\"%s\",", label, data)
	}

	//there will be an extra comma at the end. remove it pushing channel
	expoString = strings.TrimSuffix(expoString, ",")

	expoString = expoString + "} " + fmt.Sprintf("%f", dataPoint.Value) + " " + fmt.Sprintf("%d", time.Now().UnixMilli())
	return expoString
}

// loop continually, push data to server if maxWait elapsed or maxQueueSize reached
func (v *Vic) pushLoop() {
	println("starting push loop")
	for {
		select {
		case <-v.pushTicker.C:
			v.pushMetrics()
		case <-v.instantPush:
			v.pushMetrics()
		case <-v.stopPushing:
			println("push stopping")
			return
		}
	}
}

// combine metrics in msgQ into one command to send to server
func (v *Vic) pushMetrics() {
	totMetrics := ""
	chanLen := len(v.msgQ)
	if chanLen > 0 {
		for i := 0; i < chanLen; i++ {
			metric := <-v.msgQ
			totMetrics = totMetrics + "\n" + metric
		}
		println(totMetrics)
		println(v.URL)

		// Create a HTTP post request
		r, err := http.NewRequest("POST", v.URL, bytes.NewBuffer([]byte(totMetrics)))
		if err != nil {
			//panic(err)
			log.Println(err)
			return
		}

		client := &http.Client{}
		res, err := client.Do(r)
		if err != nil {
			log.Println(err)
			return
		}
		println(res.StatusCode)
	}
}

// stop pushing
func (v *Vic) StopPushing() {
	v.stopPushing <- true
}
