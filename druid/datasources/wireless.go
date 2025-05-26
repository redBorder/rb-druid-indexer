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

var WirelessMetrics = []druidrouter.Metrics{
	{Type: "count", Name: "events"},
	{Type: "hyperUnique", Name: "wireless_stations", FieldName: "wireless_station"},
	{Type: "hyperUnique", Name: "wireless_channels", FieldName: "wireless_channel"},
	{Type: "longSum", Name: "sum_wireless_tx_power", FieldName: "wireless_tx_power"},
}

var WirelessDimensionsExclusions = []string{}

var WirelessDimensions = []string{
	"wireless_station", "type", "wireless_channel", "wireless_tx_power",
	"wireless_admin_state", "wireless_op_state", "wireless_mode", "wireless_slot",
	"ti_category", "ti_score", "ti_policy_name", "ti_policy_id", "ti_indicators",
	"sensor_name", "sensor_uuid", "deployment", "deployment_uuid", "namespace",
	"namespace_uuid", "organization", "organization_uuid", "market", "market_uuid",
	"floor", "floor_uuid", "zone", "zone_uuid", "building", "building_uuid",
	"campus", "campus_uuid", "service_provider", "service_provider_uuid",
	"wireless_station_ip", "status", "wireless_station_name", "client_count",
}

const WirelessDataSource = "rb_wireless"
