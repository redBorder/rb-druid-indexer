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

func TestLocationMetrics(t *testing.T) {
	expected := []druidrouter.Metrics{
		{Type: "count", Name: "events"},
		{Type: "doubleSum", Name: "sum_latitude", FieldName: "latitude"},
		{Type: "doubleSum", Name: "sum_longitude", FieldName: "longitude"},
		{Type: "hyperUnique", Name: "unique_locations", FieldName: "location_id"},
	}

	if !reflect.DeepEqual(LocationMetrics, expected) {
		t.Errorf("LocationMetrics does not match expected value.\nGot: %+v\nExpected: %+v", LocationMetrics, expected)
	}
}

func TestLocationDimensionsExclusions(t *testing.T) {
	expected := []string{}
	if !reflect.DeepEqual(LocationDimensionsExclusions, expected) {
		t.Errorf("LocationDimensionsExclusions does not match expected value.\nGot: %+v\nExpected: %+v", LocationDimensionsExclusions, expected)
	}
}

func TestLocationDimensions(t *testing.T) {
	const expectedLength = 28
	if len(LocationDimensions) != expectedLength {
		t.Errorf("Expected LocationDimensions length to be %d, got %d", expectedLength, len(LocationDimensions))
	}

	if LocationDimensions[0] != "location_id" {
		t.Errorf("Expected first dimension to be 'location_id', got %q", LocationDimensions[0])
	}
	if LocationDimensions[1] != "latitude" {
		t.Errorf("Expected second dimension to be 'latitude', got %q", LocationDimensions[1])
	}
	if LocationDimensions[27] != "service_provider_uuid" {
		t.Errorf("Expected last dimension to be 'campus_uuid', got %q", LocationDimensions[27])
	}
}

func TestLocationDataSource(t *testing.T) {
	expected := "rb_location"
	if LocationDataSource != expected {
		t.Errorf("Expected LocationDataSource to be %q, got %q", expected, LocationDataSource)
	}
}
