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

var LocationMetrics = []druidrouter.Metrics{
	{Type: "count", Name: "events"},
	{Type: "doubleSum", Name: "sum_latitude", FieldName: "latitude"},
	{Type: "doubleSum", Name: "sum_longitude", FieldName: "longitude"},
	{Type: "hyperUnique", Name: "unique_locations", FieldName: "location_id"},
}

var LocationDimensionsExclusions = []string{}

var LocationDimensions = []string{
	"location_id", "latitude", "longitude", "address", "city", "region",
	"country", "postal_code", "sensor_name", "sensor_uuid", "deployment",
	"deployment_uuid", "namespace", "namespace_uuid", "organization",
	"organization_uuid", "market", "market_uuid", "floor", "floor_uuid",
	"zone", "zone_uuid", "building", "building_uuid", "campus",
	"campus_uuid", "service_provider", "service_provider_uuid",
}

const LocationDataSource = "rb_location"
