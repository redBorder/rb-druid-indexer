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

func TestMonitorMetrics(t *testing.T) {
	expected := []druidrouter.Metrics{
		{Type: "count", Name: "events"},
		{Type: "doubleSum", Name: "sum_value", FieldName: "value"},
		{Type: "doubleMax", Name: "max_value", FieldName: "value"},
		{Type: "doubleMin", Name: "min_value", FieldName: "value"},
	}

	if !reflect.DeepEqual(MonitorMetrics, expected) {
		t.Errorf("MonitorMetrics does not match expected value.\nGot: %+v\nExpected: %+v", MonitorMetrics, expected)
	}
}

func TestMonitorDimensionsExclusions(t *testing.T) {
	expected := []string{"unit", "type", "value"}
	if !reflect.DeepEqual(MonitorDimensionsExclusions, expected) {
		t.Errorf("MonitorDimensionsExclusions does not match expected value.\nGot: %+v\nExpected: %+v", MonitorDimensionsExclusions, expected)
	}
}

func TestMonitorDimensions(t *testing.T) {
	if len(MonitorDimensions) != 0 {
		t.Errorf("Expected MonitorDimensions length to be 0, got %d", len(MonitorDimensions))
	}
}

func TestMonitorDataSource(t *testing.T) {
	expected := "rb_monitor"
	if MonitorDataSource != expected {
		t.Errorf("Expected MonitorDataSource to be %q, got %q", expected, MonitorDataSource)
	}
}
