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

func TestStateMetrics(t *testing.T) {
	expected := []druidrouter.Metrics{
		{Type: "count", Name: "events"},
		{Type: "doubleSum", FieldName: "value", Name: "sum_value"},
		{Type: "hyperUnique", FieldName: "wireless_station", Name: "wireless_stations"},
		{Type: "hyperUnique", FieldName: "wireless_channel", Name: "wireless_channels"},
		{Type: "longSum", FieldName: "wireless_tx_power", Name: "sum_wireless_tx_power"},
	}

	if !reflect.DeepEqual(StateMetrics, expected) {
		t.Errorf("StateMetrics does not match expected value.\nGot: %+v\nExpected: %+v", StateMetrics, expected)
	}
}

func TestStateDimensionsExclusions(t *testing.T) {
	expected := []string{}
	if !reflect.DeepEqual(StateDimensionsExclusions, expected) {
		t.Errorf("StateDimensionsExclusions does not match expected value.\nGot: %+v\nExpected: %+v", StateDimensionsExclusions, expected)
	}
}

func TestStateDimensions(t *testing.T) {
	const expectedLength = 32
	if len(StateDimensions) != expectedLength {
		t.Errorf("Expected StateDimensions length to be %d, got %d", expectedLength, len(StateDimensions))
	}

	if StateDimensions[0] != "wireless_station" {
		t.Errorf("Expected first dimension to be 'wireless_station', got %q", StateDimensions[0])
	}
	if StateDimensions[1] != "type" {
		t.Errorf("Expected second dimension to be 'type', got %q", StateDimensions[1])
	}
	if StateDimensions[31] != "client_count" {
		t.Errorf("Expected last dimension to be 'client_count', got %q", StateDimensions[31])
	}
}

func TestStateDataSource(t *testing.T) {
	expected := "rb_state"
	if StateDataSource != expected {
		t.Errorf("Expected StateDataSource to be %q, got %q", expected, StateDataSource)
	}
}
