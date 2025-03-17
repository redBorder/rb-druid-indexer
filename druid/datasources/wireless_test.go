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

func TestWirelessMetrics(t *testing.T) {
	expected := []druidrouter.Metrics{
		{Type: "count", Name: "events"},
		{Type: "hyperUnique", Name: "wireless_stations", FieldName: "wireless_station"},
		{Type: "hyperUnique", Name: "wireless_channels", FieldName: "wireless_channel"},
		{Type: "longSum", Name: "sum_wireless_tx_power", FieldName: "wireless_tx_power"},
	}

	if !reflect.DeepEqual(WirelessMetrics, expected) {
		t.Errorf("WirelessMetrics does not match expected value.\nGot: %+v\nExpected: %+v", WirelessMetrics, expected)
	}
}

func TestWirelessDimensionsExclusions(t *testing.T) {
	expected := []string{}
	if !reflect.DeepEqual(WirelessDimensionsExclusions, expected) {
		t.Errorf("WirelessDimensionsExclusions does not match expected value.\nGot: %+v\nExpected: %+v", WirelessDimensionsExclusions, expected)
	}
}

func TestWirelessDimensions(t *testing.T) {
	const expectedLength = 32
	if len(WirelessDimensions) != expectedLength {
		t.Errorf("Expected WirelessDimensions length to be %d, got %d", expectedLength, len(WirelessDimensions))
	}

	if WirelessDimensions[0] != "wireless_station" {
		t.Errorf("Expected first dimension to be 'wireless_station', got %q", WirelessDimensions[0])
	}
	if WirelessDimensions[1] != "type" {
		t.Errorf("Expected second dimension to be 'type', got %q", WirelessDimensions[1])
	}
	if WirelessDimensions[31] != "client_count" {
		t.Errorf("Expected last dimension to be 'client_count', got %q", WirelessDimensions[31])
	}
}

func TestWirelessDataSource(t *testing.T) {
	expected := "rb_wireless"
	if WirelessDataSource != expected {
		t.Errorf("Expected WirelessDataSource to be %q, got %q", expected, WirelessDataSource)
	}
}
