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

var EventMetrics = []druidrouter.Metrics{
	{Type: "count", Name: "events"},
	{Type: "hyperUnique", Name: "signatures", FieldName: "msg"},
}

var EventDimensionsExclusions = []string{
	"payload",
}

var EventDimensions = []string{
        "src", "src_is_malicious", "dst", "dst_is_malicious", "sensor_uuid", "src_port",
	"dst_port", "src_as_name", "src_country_code", "dst_map", "src_map", "service_provider",
	"sha256", "sha256_is_malicious", "file_uri", "file_uri_is_malicious", "file_size",
	"file_hostname", "file_hostname_is_malicious", "action", "ethlength_range", "icmptype",
	"ethsrc", "ethsrc_vendor", "ethdst", "ethdst_vendor", "ttl", "vlan", "classification",
	"domain_name", "group_name", "sig_generator", "rev", "priority", "msg", "sig_id",
	"dst_country_code", "dst_as_name", "namespace", "deployment", "market", "organization",
	"campus", "building", "floor", "floor_uuid", "conversation", "iplen_range",
	"l4_proto", "sensor_name", "scatterplot", "src_net_name", "dst_net_name", "tos",
	"service_provider_uuid", "namespace_uuid", "market_uuid", "organization_uuid",
	"campus_uuid", "building_uuid", "deployment_uuid", "incident_uuid", "event_uuid",
}

const EventDataSource = "rb_event"
