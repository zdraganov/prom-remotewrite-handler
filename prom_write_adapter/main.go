// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"github.com/golang/protobuf/proto"

	"github.com/golang/snappy"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
)

// Samples ...
type Samples struct {
	Value     *float64 `json:"value"`
	Timestamp int64    `json:"timestamp"`
}

// TimeSeries ...
type TimeSeries struct {
	Name    string            `json:"name"`
	Labels  map[string]string `json:"labels"`
	Samples []*Samples        `json:"samples"`
}

func prometheusTsToJSON(ts *prompb.TimeSeries) ([]byte, error) {
	metric := &TimeSeries{}

	metric.Labels = make(map[string]string)
	metric.Samples = make([]*Samples, len(ts.Samples))

	for _, l := range ts.Labels {
		if model.LabelName(l.Name) == model.MetricNameLabel {
			metric.Name = string(model.LabelValue(l.Value))
		} else {
			metric.Labels[string(model.LabelName(l.Name))] = string(model.LabelValue(l.Value))
		}
	}

	for i, s := range ts.Samples {
		value := &s.Value

		if math.IsNaN(s.Value) || math.IsInf(s.Value, 0) {
			value = nil
		}

		metric.Samples[i] = &Samples{
			Value:     value,
			Timestamp: s.Timestamp,
		}
	}

	return json.Marshal(metric)
}

func main() {
	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		compressed, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reqBuf, err := snappy.Decode(nil, compressed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req prompb.WriteRequest
		if err := proto.Unmarshal(reqBuf, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, ts := range req.Timeseries {
			tsJSON, err := prometheusTsToJSON(ts)

			if err != nil {
				log.Fatal(err)
				return
			}

			fmt.Println(string(tsJSON))
		}
	})

	log.Fatal(http.ListenAndServe(":1234", nil))
}
