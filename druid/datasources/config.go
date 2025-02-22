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
	DataSource string
	Metrics    []druidrouter.Metrics
	Dimensions []string
}

var Configs = map[string]DataSourceConfig{
	"rb_flow": {
		DataSource: FlowDataSource,
		Metrics:    FlowMetrics,
		Dimensions: FlowDimensions,
	},
	"rb_monitor": {
		DataSource: MonitorDataSource,
		Metrics:    MonitorMetrics,
		Dimensions: MonitorDimensions,
	},
	"rb_wireless": {
		DataSource: WirelessDataSource,
		Metrics:    WirelessMetrics,
		Dimensions: WirelessDimensions,
	},
	"rb_loc": {
		DataSource: LocationDataSource,
		Metrics:    LocationMetrics,
		Dimensions: LocationDimensions,
	},
	"rb_event": {
		DataSource: EventDataSource,
		Metrics:    EventMetrics,
		Dimensions: EventDimensions,
	},
	"rb_state": {
		DataSource: StateDataSource,
		Metrics:    StateMetrics,
		Dimensions: StateDimensions,
	},
	"rb_vault": {
		DataSource: VaultDataSource,
		Metrics:    VaultMetrics,
		Dimensions: VaultDimensions,
	},
}

func GetDataSourceConfig(taskName string) (DataSourceConfig, bool) {
	config, exists := Configs[taskName]
	return config, exists
}
