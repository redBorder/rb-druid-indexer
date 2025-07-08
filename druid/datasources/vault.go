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

var VaultMetrics = []druidrouter.Metrics{
	{Type: "count", Name: "events"},
}

var VaultDimensionsExclusions = []string{
	"unit", "type", "valur",
}

var VaultDimensions = []string{
	"pri", "pri_text", "syslogfacility", "syslogfacility_text", "syslogseverity", "syslogseverity_text",
	"hostname", "fromhost_ip", "app_name", "sensor_name", "proxy_uuid", "message", "status",
	"ti_category", "ti_average_score", "ti_policy_name", "ti_policy_id", "ti_indicators",
	"category", "source", "target", "sensor_uuid", "service_provider", "service_provider_uuid",
	"namespace", "namespace_uuid", "deployment", "deployment_uuid", "market", "market_uuid",
	"organization", "organization_uuid", "campus", "campus_uuid", "building", "building_uuid",
	"floor", "floor_uuid", "action", "incident_uuid", "alarm_id", "alarm_name", "alarm_product_type",
	"alarm_condition", "alarm_user", "alarm_severity", "lan_ip", "wan_ip",
	"wireless_station", "asset_ip_address", "asset_mac_address",
}

// VaultDataSource is the name of the Vault data source in Druid.
const VaultDataSource = "rb_vault"
