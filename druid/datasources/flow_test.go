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
package datasources

import (
	"reflect"
	"testing"

	druidrouter "rb-druid-indexer/druid"
)

func TestFlowMetrics(t *testing.T) {
	expected := []druidrouter.Metrics{
		{Type: "count", Name: "events"},
		{Type: "longSum", Name: "sum_bytes", FieldName: "bytes"},
		{Type: "longSum", Name: "sum_pkts", FieldName: "pkts"},
		{Type: "longSum", Name: "sum_rssi", FieldName: "client_rssi_num"},
		{Type: "hyperUnique", Name: "clients", FieldName: "client_mac"},
		{Type: "hyperUnique", Name: "wireless_stations", FieldName: "wireless_station"},
	}

	if !reflect.DeepEqual(FlowMetrics, expected) {
		t.Errorf("FlowMetrics does not match expected value.\nGot: %+v\nExpected: %+v", FlowMetrics, expected)
	}
}

func TestFlowDimensionsExclusions(t *testing.T) {
	expected := []string{"bytes", "pkts", "flow_end_reason", "first_switched", "wan_ip_name"}
	if !reflect.DeepEqual(FlowDimensionsExclusions, expected) {
		t.Errorf("FlowDimensionsExclusions does not match expected value.\nGot: %+v\nExpected: %+v", FlowDimensionsExclusions, expected)
	}
}

func TestFlowDimensions(t *testing.T) {
	const expectedLength = 93
	if len(FlowDimensions) != expectedLength {
		t.Errorf("Expected FlowDimensions length to be %d, got %d", expectedLength, len(FlowDimensions))
	}

	if FlowDimensions[0] != "application_id_name" {
		t.Errorf("Expected first dimension to be 'application_id_name', got %q", FlowDimensions[0])
	}
	if FlowDimensions[1] != "building" {
		t.Errorf("Expected second dimension to be 'building', got %q", FlowDimensions[1])
	}
	if FlowDimensions[92] != "zone_uuid" {
		t.Errorf("Expected last dimension to be 'zone_uuid', got %q", FlowDimensions[92])
	}
}

func TestFlowDataSource(t *testing.T) {
	expected := "rb_flow"
	if FlowDataSource != expected {
		t.Errorf("Expected FlowDataSource to be %q, got %q", expected, FlowDataSource)
	}
}
