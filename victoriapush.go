/*
victoriapush pushes metrics to VictoriaMetrics using the prometheus exposition format

send metrics with custom labels in a map

assign global labels
*/
package victoriapush

type Vic struct {
	URL          string            //URL of victoriametrics server
	globalLabels map[string]string //labels to be applied to every metric pushed
	maxWait      float64           //maximum time to wait before sending whatever is in buffer
	maxQueueSize int               //maximum size for metrics queue before sending
	msgQ         chan string       //queue of messages waiting to be pushed
}

type DataPoint struct {
	metric string
	value  float64
	labels map[string]string
}

// add URL and globalLabels
func NewVictoriaPush(URL string, globalLabels map[string]string) *Vic {
	v := new(Vic)
	v.URL = URL
	v.globalLabels = globalLabels
	v.maxWait = 1
	v.maxQueueSize = 10

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
	dataPointExpo := v.dataPointToExpo(dataPoint)
	println(dataPointExpo)

}

// convert DataPoint to prometheus exposition format ready for pushing to server
func (v *Vic) dataPointToExpo(dataPoint DataPoint) string {
	return dataPoint.metric
}

// loop continually, push data to server if maxWait elapsed or maxQueueSize reached
func (v *Vic) pushLoop() {

}
