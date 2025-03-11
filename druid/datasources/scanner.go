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

var ScannerMetrics = []druidrouter.Metrics{
	{Type: "count", Name: "events"},
}

var ScannerDimensionsExclusions = []string{}

var ScannerDimensions = []string{
	"pri", "pri_text", "syslogfacility", "syslogfacility_text", "syslogseverity", "syslogseverity_text",
	"hostname", "fromhost_ip", "app_name", "sensor_name", "proxy_uuid", "message", "status",
	"category", "source", "target", "sensor_uuid", "service_provider", "service_provider_uuid",
	"namespace", "namespace_uuid", "deployment", "deployment_uuid", "market", "market_uuid",
	"organization", "organization_uuid", "campus", "campus_uuid", "building", "building_uuid",
	"floor", "floor_uuid", "ipaddress", "scan_id", "scan_subtype", "scan_type", "result_data",
	"result", "cve_info", "vendor", "product", "version", "servicename", "protocol", "cpe",
	"cve", "port", "metric", "severity", "score", "mac", "subnet", "path", "layer", "ipv4",
	"port_state",
}

const ScannerDataSource = "rb_scanner"
