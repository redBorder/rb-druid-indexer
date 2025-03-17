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

func TestScannerMetrics(t *testing.T) {
	expected := []druidrouter.Metrics{
		{Type: "count", Name: "events"},
	}

	if !reflect.DeepEqual(ScannerMetrics, expected) {
		t.Errorf("ScannerMetrics does not match expected value.\nGot: %+v\nExpected: %+v", ScannerMetrics, expected)
	}
}

func TestScannerDimensionsExclusions(t *testing.T) {
	expected := []string{}
	if !reflect.DeepEqual(ScannerDimensionsExclusions, expected) {
		t.Errorf("ScannerDimensionsExclusions does not match expected value.\nGot: %+v\nExpected: %+v", ScannerDimensionsExclusions, expected)
	}
}

func TestScannerDimensions(t *testing.T) {
	const expectedLength = 57
	if len(ScannerDimensions) != expectedLength {
		t.Errorf("Expected ScannerDimensions length to be %d, got %d", expectedLength, len(ScannerDimensions))
	}

	if ScannerDimensions[0] != "pri" {
		t.Errorf("Expected first dimension to be 'pri', got %q", ScannerDimensions[0])
	}
	if ScannerDimensions[1] != "pri_text" {
		t.Errorf("Expected second dimension to be 'pri_text', got %q", ScannerDimensions[1])
	}
	if ScannerDimensions[56] != "port_state" {
		t.Errorf("Expected last dimension to be 'port_state', got %q", ScannerDimensions[56])
	}
}

func TestScannerDataSource(t *testing.T) {
	expected := "rb_scanner"
	if ScannerDataSource != expected {
		t.Errorf("Expected ScannerDataSource to be %q, got %q", expected, ScannerDataSource)
	}
}
