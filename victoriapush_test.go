/*
victoriapush pushes metrics to VictoriaMetrics using the prometheus exposition format

send metrics with custom labels in a map

assign global labels
*/
package victoriapush

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewVictoriaPush(t *testing.T) {

	globalLabelsTest := make(map[string]string)
	globalLabelsTest["label_1"] = "a value"
	globalLabelsTest["label_2"] = "another value"

	type args struct {
		URL          string
		globalLabels map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Push returned",
			args: args{
				URL:          "victoriametrics:8289",
				globalLabels: globalLabelsTest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewVictoriaPush(tt.args.URL, tt.args.globalLabels)
			if got.URL != tt.args.URL {
				t.Errorf("NewVictoriaPush() = %v, want %v", got, tt.args.URL)
			}
		})
	}
}

func TestVic_SetQueueLimits(t *testing.T) {
	type fields struct {
		URL          string
		globalLabels map[string]string
		maxWait      float64
		maxQueueSize int
		msgQ         chan string
	}
	type args struct {
		maxWait      float64
		maxQueueSize int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "queue lims",
			fields: fields{},
			args: args{
				maxWait:      17.4,
				maxQueueSize: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vic{
				URL:          tt.fields.URL,
				globalLabels: tt.fields.globalLabels,
				maxWait:      tt.fields.maxWait,
				maxQueueSize: tt.fields.maxQueueSize,
				msgQ:         tt.fields.msgQ,
			}
			v.SetQueueLimits(tt.args.maxWait, tt.args.maxQueueSize)
			got := args{v.maxWait, v.maxQueueSize}
			if got != tt.args {
				t.Errorf("SetQueueLimits() got: %v, wanted: %v", got, tt.args)
			}
		})
	}
}

func TestVic_AddGlobalLabels(t *testing.T) {

	globalLabelsTest := map[string]string{
		"label 1": "a value",
		"label 2": "another value",
	}

	globalLabelsAdditional := map[string]string{
		"label_3": "3rd value",
	}

	globalLabelsCombined := map[string]string{
		"label 1": "a value",
		"label 2": "another value",
		"label_3": "3rd value",
	}

	type fields struct {
		URL          string
		globalLabels map[string]string
		maxWait      float64
		maxQueueSize int
		msgQ         chan string
	}
	type args struct {
		globalLabels map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "adding labels",
			fields: fields{globalLabels: globalLabelsTest},
			args: args{
				globalLabelsAdditional,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vic{
				URL:          tt.fields.URL,
				globalLabels: tt.fields.globalLabels,
				maxWait:      tt.fields.maxWait,
				maxQueueSize: tt.fields.maxQueueSize,
				msgQ:         tt.fields.msgQ,
			}
			v.AddGlobalLabels(tt.args.globalLabels)

			if !reflect.DeepEqual(v.globalLabels, globalLabelsCombined) {
				t.Errorf("SetQueueLimits() got: %v, wanted: %v", v.globalLabels, globalLabelsCombined)
			}
		})
	}
}

func TestVic_dataPointToExpo(t *testing.T) {
	type fields struct {
		URL          string
		globalLabels map[string]string
		maxWait      float64
		maxQueueSize int
		msgQ         chan string
	}
	type args struct {
		dataPoint DataPoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "expotest",
			args: args{DataPoint{
				Metric: "metrictest",
				Value:  21.9,
				Labels: map[string]string{
					"label1": "a value",
					"label2": "another value",
				},
			}},
			want: "metrictest{label1=\"a value\",label2=\"another value\"} 21.900000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vic{
				URL:          tt.fields.URL,
				globalLabels: tt.fields.globalLabels,
				maxWait:      tt.fields.maxWait,
				maxQueueSize: tt.fields.maxQueueSize,
				msgQ:         tt.fields.msgQ,
			}
			got := v.dataPointToExpo(tt.args.dataPoint)
			println(got)
			got = got[0 : len(got)-14]
			fmt.Println(got)
			if got != tt.want {
				t.Errorf("Vic.dataPointToExpo() = %v, want %v", got, tt.want)
			}
		})
	}
}
