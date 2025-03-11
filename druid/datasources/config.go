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

import druidrouter "rb-druid-indexer/druid"

type DataSourceConfig struct {
	DataSource           string
	Metrics              []druidrouter.Metrics
	Dimensions           []string
	DimensionsExclusions []string
}

var Configs = map[string]DataSourceConfig{
	"rb_flow": {
		DataSource:           FlowDataSource,
		Metrics:              FlowMetrics,
		Dimensions:           FlowDimensions,
		DimensionsExclusions: FlowDimensionsExclusions,
	},
	"rb_monitor": {
		DataSource:           MonitorDataSource,
		Metrics:              MonitorMetrics,
		Dimensions:           MonitorDimensions,
		DimensionsExclusions: MonitorDimensionsExclusions,
	},
	"rb_wireless": {
		DataSource:           WirelessDataSource,
		Metrics:              WirelessMetrics,
		Dimensions:           WirelessDimensions,
		DimensionsExclusions: WirelessDimensionsExclusions,
	},
	"rb_location": {
		DataSource:           LocationDataSource,
		Metrics:              LocationMetrics,
		Dimensions:           LocationDimensions,
		DimensionsExclusions: LocationDimensionsExclusions,
	},
	"rb_event": {
		DataSource:           EventDataSource,
		Metrics:              EventMetrics,
		Dimensions:           EventDimensions,
		DimensionsExclusions: EventDimensionsExclusions,
	},
	"rb_state": {
		DataSource:           StateDataSource,
		Metrics:              StateMetrics,
		Dimensions:           StateDimensions,
		DimensionsExclusions: StateDimensionsExclusions,
	},
	"rb_vault": {
		DataSource:           VaultDataSource,
		Metrics:              VaultMetrics,
		Dimensions:           VaultDimensions,
		DimensionsExclusions: VaultDimensionsExclusions,
	},
	"rb_scanner": {
		DataSource:           ScannerDataSource,
		Metrics:              ScannerMetrics,
		Dimensions:           ScannerDimensions,
		DimensionsExclusions: ScannerDimensionsExclusions,
	},
}

func GetDataSourceConfig(taskName string) (DataSourceConfig, bool) {
	config, exists := Configs[taskName]
	return config, exists
}
