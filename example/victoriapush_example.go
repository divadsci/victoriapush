package main

import (
	"time"

	"github.com/divadsci/victoriapush"
)

func main() {
	globalLabels := make(map[string]string)
	v := victoriapush.NewVictoriaPush("http://localhost:8428", globalLabels)
	dataPoint := victoriapush.DataPoint{
		Metric: "example_metric",
		Value:  21.9,
		Labels: map[string]string{
			"label1": "a value",
			"label2": "another value",
		},
	}

	v.EnqueueDataPoint(dataPoint)
	time.Sleep(time.Second)
	v.EnqueueDataPoint(dataPoint)
	time.Sleep(time.Second)
	v.EnqueueDataPoint(dataPoint)
	time.Sleep(time.Second)
	v.EnqueueDataPoint(dataPoint)
	time.Sleep(time.Second)
	v.EnqueueDataPoint(dataPoint)
	time.Sleep(time.Second)
	for {
		time.Sleep(time.Second)
	}
}
