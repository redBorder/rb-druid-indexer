// Copyright (C) 2025 Eneo Tecnologia S.L.
// Miguel Álvarez <malvarez@redborder.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package druidrouter

import (
	"encoding/json"
	"strings"
)

type KafkaConfiguration struct {
	Type string    `json:"type"`
	Spec KafkaSpec `json:"spec"`
}

type KafkaSpec struct {
	DataSchema   DataSchema   `json:"dataSchema"`
	IOConfig     IOConfig     `json:"ioConfig"`
	TuningConfig TuningConfig `json:"tuningConfig"`
}

type DataSchema struct {
	DataSource      string          `json:"dataSource"`
	MetricsSpec     []Metrics       `json:"metricsSpec"`
	GranularitySpec GranularitySpec `json:"granularitySpec"`
	TimestampSpec   TimestampSpec   `json:"timestampSpec"`
	DimensionsSpec  DimensionsSpec  `json:"dimensionsSpec"`
}

type Metrics struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	FieldName string `json:"fieldName" yaml:"fieldName"`
}

type GranularitySpec struct {
	Type               string `json:"type"`
	SegmentGranularity string `json:"segmentGranularity"`
	QueryGranularity   string `json:"queryGranularity"`
}

type TimestampSpec struct {
	Column string `json:"column"`
	Format string `json:"format"`
}

type DimensionsSpec struct {
	Dimensions          []string `json:"dimensions"`
	DimensionExclusions []string `json:"dimensionExclusions"`
}

type IOConfig struct {
	Type               string            `json:"type"`
	ConsumerProperties map[string]string `json:"consumerProperties"`
	Topic              string            `json:"topic"`
	InputFormat        InputFormat       `json:"inputFormat"`
}

type InputFormat struct {
	Type string `json:"type"`
}

type TuningConfig struct {
	Type                     string `json:"type"`
	ResetOffsetAutomatically bool   `json:"resetOffsetAutomatically"`
}

func GenerateConfig(dataSource string, KafkaBrokers []string, kafkaTopic, timestampColumn, timestampFormat string, dimensions []string, dimensionsExclusions []string, metrics []Metrics) (string, error) {
	config := KafkaConfiguration{
		Type: "kafka",
		Spec: KafkaSpec{
			DataSchema: DataSchema{
				DataSource:  dataSource,
				MetricsSpec: metrics,
				GranularitySpec: GranularitySpec{
					Type:               "uniform",
					SegmentGranularity: "HOUR",
					QueryGranularity:   "MINUTE",
				},
				TimestampSpec: TimestampSpec{
					Column: timestampColumn,
					Format: timestampFormat,
				},
				DimensionsSpec: DimensionsSpec{
					Dimensions:          dimensions,
					DimensionExclusions: dimensionsExclusions,
				},
			},
			IOConfig: IOConfig{
				Type: "kafka",
				ConsumerProperties: map[string]string{
					"bootstrap.servers": strings.Join(KafkaBrokers, ","),
					"auto.offset.reset": "latest",
				},
				Topic: kafkaTopic,
				InputFormat: InputFormat{
					Type: "json",
				},
			},
			TuningConfig: TuningConfig{
				Type:                     "kafka",
				ResetOffsetAutomatically: true,
			},
		},
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
