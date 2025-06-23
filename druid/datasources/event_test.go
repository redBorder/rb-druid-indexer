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

func TestEventMetrics(t *testing.T) {
	expected := []druidrouter.Metrics{
		{Type: "count", Name: "events"},
		{Type: "hyperUnique", Name: "signatures", FieldName: "msg"},
	}

	if !reflect.DeepEqual(EventMetrics, expected) {
		t.Errorf("EventMetrics does not match expected value.\nGot: %+v\nExpected: %+v", EventMetrics, expected)
	}
}

func TestEventDimensionsExclusions(t *testing.T) {
	expected := []string{"payload"}
	if !reflect.DeepEqual(EventDimensionsExclusions, expected) {
		t.Errorf("EventDimensionsExclusions does not match expected value.\nGot: %+v\nExpected: %+v", EventDimensionsExclusions, expected)
	}
}

func TestEventDimensions(t *testing.T) {
	const expectedLength = 63
	if len(EventDimensions) != expectedLength {
		t.Errorf("Expected EventDimensions length to be %d, got %d", expectedLength, len(EventDimensions))
	}

	if EventDimensions[0] != "src" {
		t.Errorf("Expected first dimension to be 'src', got %q", EventDimensions[0])
	}
	if EventDimensions[1] != "dst" {
		t.Errorf("Expected second dimension to be 'dst', got %q", EventDimensions[1])
	}
	if EventDimensions[62] != "event_uuid" {
		t.Errorf("Expected last dimension to be 'event_uuid', got %q", EventDimensions[66])
	}
}

func TestEventDataSource(t *testing.T) {
	expected := "rb_event"
	if EventDataSource != expected {
		t.Errorf("Expected EventDataSource to be %q, got %q", expected, EventDataSource)
	}
}
