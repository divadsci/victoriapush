/*
victoriapush pushes metrics to VictoriaMetrics using the prometheus exposition format

send metrics with custom labels in a map

assign global labels
*/
package victoriapush

type Vic struct {
	URL          string
	globalLabels map[string]string
}

// add URL and globalLabels
func (v *Vic) InitVictoria(URL string, globalLabels map[string]string) {
	v.URL = URL
	v.globalLabels = globalLabels
}

// Replace globalLabels map with a new one
func (v *Vic) ReplaceGlobalLabels(globalLabels map[string]string) {
	v.globalLabels = globalLabels
}

// Add a map of labels to globalLabels
func (v *Vic) AddGlobalLabels(globalLabels map[string]string) {

	//get globalLabels from context, iterate over new map and add each to globalLabels from ctx

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
