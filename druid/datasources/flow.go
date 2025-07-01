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

var FlowMetrics = []druidrouter.Metrics{
	{Type: "count", Name: "events"},
	{Type: "longSum", Name: "sum_bytes", FieldName: "bytes"},
	{Type: "longSum", Name: "sum_pkts", FieldName: "pkts"},
	{Type: "longSum", Name: "sum_rssi", FieldName: "client_rssi_num"},
	{Type: "hyperUnique", Name: "clients", FieldName: "client_mac"},
	{Type: "hyperUnique", Name: "wireless_stations", FieldName: "wireless_station"},
}

// var FlowDimensionsExclusions = []string{
// 	"bytes", "pkts", "flow_end_reason", "first_switched", "wan_ip_name",
// }

// var FlowDimensions = []string{
// 	"application_id_name", "building", "building_uuid", "campus", "campus_uuid",
// 	"client_accounting_type", "client_auth_type", "client_fullname", "client_gender",
// 	"client_id", "client_latlong", "client_loyality", "client_mac", "client_mac_vendor",
// 	"client_rssi", "client_vip", "conversation", "coordinates_map", "deployment",
// 	"deployment_uuid", "direction", "dot11_protocol", "dot11_status", "dst_map", "duration",
// 	"engine_id_name", "floor", "floor_uuid", "host", "host_l2_domain", "http_social_media",
// 	"http_user_agent", "https_common_name", "interface_name", "ip_as_name", "ip_country_code",
// 	"ip_protocol_version", "l4_proto", "lan_interface_description", "lan_interface_name",
// 	"lan_ip", "lan_ip_as_name", "lan_ip_country_code", "lan_ip_name",
// 	"lan_ip_net_name", "lan_l4_port", "lan_name", "lan_vlan", "market", "market_uuid",
// 	"namespace", "namespace_uuid", "organization", "organization_uuid", "product_name",
// 	"public_ip", "public_ip_mac", "referer", "referer_l2",
// 	"scatterplot", "selector_name", "sensor_ip", "sensor_name", "sensor_uuid", "service_provider",
// 	"service_provider_uuid", "src_map", "tcp_flags", "tos", "type", "url", "wan_interface_description",
// 	"wan_interface_name", "wan_ip",	"wan_ip_as_name", "wan_ip_country_code",
// 	"wan_ip_map", "wan_ip_net_name", "wan_l4_port", "wan_name", "wan_vlan", "wireless_id",
// 	"ti_category", "ti_average_score", "ti_policy_name", "ti_policy_id", "ti_indicators",
// 	"wireless_operator", "wireless_station", "zone", "zone_uuid",
// }

	

const FlowDataSource = "rb_flow"
