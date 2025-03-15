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

func TestVaultMetrics(t *testing.T) {
	expected := []druidrouter.Metrics{
		{Type: "count", Name: "events"},
	}

	if !reflect.DeepEqual(VaultMetrics, expected) {
		t.Errorf("VaultMetrics does not match expected value.\nGot: %+v\nExpected: %+v", VaultMetrics, expected)
	}
}

func TestVaultDimensionsExclusions(t *testing.T) {
	expected := []string{
		"unit", "type", "valur",
	}
	if !reflect.DeepEqual(VaultDimensionsExclusions, expected) {
		t.Errorf("VaultDimensionsExclusions does not match expected value.\nGot: %+v\nExpected: %+v", VaultDimensionsExclusions, expected)
	}
}

func TestVaultDimensions(t *testing.T) {
	const expectedLength = 41
	if len(VaultDimensions) != expectedLength {
		t.Errorf("Expected VaultDimensions length to be %d, got %d", expectedLength, len(VaultDimensions))
	}

	if VaultDimensions[0] != "pri" {
		t.Errorf("Expected first dimension to be 'pri', got %q", VaultDimensions[0])
	}
	if VaultDimensions[1] != "pri_text" {
		t.Errorf("Expected second dimension to be 'pri_text', got %q", VaultDimensions[1])
	}
	if VaultDimensions[40] != "alarm_severity" {
		t.Errorf("Expected last dimension to be 'alarm_severity', got %q", VaultDimensions[33])
	}
}

func TestVaultDataSource(t *testing.T) {
	expected := "rb_vault"
	if VaultDataSource != expected {
		t.Errorf("Expected VaultDataSource to be %q, got %q", expected, VaultDataSource)
	}
}
