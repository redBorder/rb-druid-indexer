package druidrouter

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	dataSource := "testDataSource"
	kafkaBrokers := []string{"broker1:9092", "broker2:9092"}
	kafkaTopic := "testTopic"
	timestampColumn := "ts"
	timestampFormat := "auto"
	dimensions := []string{"dim1", "dim2"}
	dimensionsExclusions := []string{"exclusion1"}
	metrics := []Metrics{
		{
			Type:      "count",
			Name:      "count",
			FieldName: "",
		},
		{
			Type:      "sum",
			Name:      "sum_metric",
			FieldName: "value",
		},
	}

	jsonStr, err := GenerateConfig(dataSource, kafkaBrokers, kafkaTopic, timestampColumn, timestampFormat, dimensions, dimensionsExclusions, metrics)
	if err != nil {
		t.Fatalf("GenerateConfig returned an error: %v", err)
	}

	var config KafkaConfiguration
	err = json.Unmarshal([]byte(jsonStr), &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if config.Type != "kafka" {
		t.Errorf("Expected Type to be 'kafka', got '%s'", config.Type)
	}

	if config.Spec.DataSchema.DataSource != dataSource {
		t.Errorf("Expected DataSource to be '%s', got '%s'", dataSource, config.Spec.DataSchema.DataSource)
	}
	if config.Spec.DataSchema.TimestampSpec.Column != timestampColumn {
		t.Errorf("Expected TimestampSpec.Column to be '%s', got '%s'", timestampColumn, config.Spec.DataSchema.TimestampSpec.Column)
	}
	if config.Spec.DataSchema.TimestampSpec.Format != timestampFormat {
		t.Errorf("Expected TimestampSpec.Format to be '%s', got '%s'", timestampFormat, config.Spec.DataSchema.TimestampSpec.Format)
	}
	if !reflect.DeepEqual(config.Spec.DataSchema.DimensionsSpec.Dimensions, dimensions) {
		t.Errorf("Expected Dimensions to be %v, got %v", dimensions, config.Spec.DataSchema.DimensionsSpec.Dimensions)
	}
	if !reflect.DeepEqual(config.Spec.DataSchema.DimensionsSpec.DimensionExclusions, dimensionsExclusions) {
		t.Errorf("Expected DimensionExclusions to be %v, got %v", dimensionsExclusions, config.Spec.DataSchema.DimensionsSpec.DimensionExclusions)
	}
	if !reflect.DeepEqual(config.Spec.DataSchema.MetricsSpec, metrics) {
		t.Errorf("Expected MetricsSpec to be %v, got %v", metrics, config.Spec.DataSchema.MetricsSpec)
	}

	if config.Spec.IOConfig.Topic != kafkaTopic {
		t.Errorf("Expected Topic to be '%s', got '%s'", kafkaTopic, config.Spec.IOConfig.Topic)
	}
	expectedBootstrap := "broker1:9092,broker2:9092"
	if config.Spec.IOConfig.ConsumerProperties["bootstrap.servers"] != expectedBootstrap {
		t.Errorf("Expected bootstrap.servers to be '%s', got '%s'", expectedBootstrap, config.Spec.IOConfig.ConsumerProperties["bootstrap.servers"])
	}

	if config.Spec.TuningConfig.Type != "kafka" {
		t.Errorf("Expected TuningConfig.Type to be 'kafka', got '%s'", config.Spec.TuningConfig.Type)
	}
	if config.Spec.TuningConfig.ResetOffsetAutomatically != true {
		t.Errorf("Expected TuningConfig.ResetOffsetAutomatically to be 'true', got '%t'", config.Spec.TuningConfig.ResetOffsetAutomatically)
	}
}
