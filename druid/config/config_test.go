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
)

func TestGetDataSourceConfig_Exists(t *testing.T) {
	testCases := []struct {
		taskName       string
		expectedConfig DataSourceConfig
	}{
		{
			taskName:       "rb_flow",
			expectedConfig: Configs["rb_flow"],
		},
		{
			taskName:       "rb_monitor",
			expectedConfig: Configs["rb_monitor"],
		},
		{
			taskName:       "rb_wireless",
			expectedConfig: Configs["rb_wireless"],
		},
		{
			taskName:       "rb_location",
			expectedConfig: Configs["rb_location"],
		},
		{
			taskName:       "rb_event",
			expectedConfig: Configs["rb_event"],
		},
		{
			taskName:       "rb_state",
			expectedConfig: Configs["rb_state"],
		},
		{
			taskName:       "rb_vault",
			expectedConfig: Configs["rb_vault"],
		},
		{
			taskName:       "rb_scanner",
			expectedConfig: Configs["rb_scanner"],
		},
	}

	for _, tc := range testCases {
		config, exists := GetDataSourceConfig(tc.taskName)
		if !exists {
			t.Errorf("Expected config to exist for taskName %q, but it does not", tc.taskName)
		}
		if !reflect.DeepEqual(config, tc.expectedConfig) {
			t.Errorf("For taskName %q, expected config %+v, got %+v", tc.taskName, tc.expectedConfig, config)
		}
	}
}

func TestGetDataSourceConfig_NotExists(t *testing.T) {
	taskName := "non_existing"
	_, exists := GetDataSourceConfig(taskName)
	if exists {
		t.Errorf("Expected config not to exist for taskName %q, but it does", taskName)
	}
}
