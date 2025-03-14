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
	druidrouter "rb-druid-indexer/druid"
	datasources "rb-druid-indexer/druid/datasources"
)

type DataSourceConfig struct {
	DataSource           string
	Metrics              []druidrouter.Metrics
	Dimensions           []string
	DimensionsExclusions []string
}

var Configs = map[string]DataSourceConfig{
	"rb_flow": {
		DataSource:           datasources.FlowDataSource,
		Metrics:              datasources.FlowMetrics,
		Dimensions:           datasources.FlowDimensions,
		DimensionsExclusions: datasources.FlowDimensionsExclusions,
	},
	"rb_monitor": {
		DataSource:           datasources.MonitorDataSource,
		Metrics:              datasources.MonitorMetrics,
		Dimensions:           datasources.MonitorDimensions,
		DimensionsExclusions: datasources.MonitorDimensionsExclusions,
	},
	"rb_wireless": {
		DataSource:           datasources.WirelessDataSource,
		Metrics:              datasources.WirelessMetrics,
		Dimensions:           datasources.WirelessDimensions,
		DimensionsExclusions: datasources.WirelessDimensionsExclusions,
	},
	"rb_location": {
		DataSource:           datasources.LocationDataSource,
		Metrics:              datasources.LocationMetrics,
		Dimensions:           datasources.LocationDimensions,
		DimensionsExclusions: datasources.LocationDimensionsExclusions,
	},
	"rb_event": {
		DataSource:           datasources.EventDataSource,
		Metrics:              datasources.EventMetrics,
		Dimensions:           datasources.EventDimensions,
		DimensionsExclusions: datasources.EventDimensionsExclusions,
	},
	"rb_state": {
		DataSource:           datasources.StateDataSource,
		Metrics:              datasources.StateMetrics,
		Dimensions:           datasources.StateDimensions,
		DimensionsExclusions: datasources.StateDimensionsExclusions,
	},
	"rb_vault": {
		DataSource:           datasources.VaultDataSource,
		Metrics:              datasources.VaultMetrics,
		Dimensions:           datasources.VaultDimensions,
		DimensionsExclusions: datasources.VaultDimensionsExclusions,
	},
	"rb_scanner": {
		DataSource:           datasources.ScannerDataSource,
		Metrics:              datasources.ScannerMetrics,
		Dimensions:           datasources.ScannerDimensions,
		DimensionsExclusions: datasources.ScannerDimensionsExclusions,
	},
}

func GetDataSourceConfig(taskName string) (DataSourceConfig, bool) {
	config, exists := Configs[taskName]
	return config, exists
}
