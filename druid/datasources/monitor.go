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

var MonitorMetrics = []druidrouter.Metrics{
	{Type: "count", Name: "events"},
	{Type: "doubleSum", Name: "sum_value", FieldName: "value"},
	{Type: "doubleMax", Name: "max_value", FieldName: "value"},
	{Type: "doubleMin", Name: "min_value", FieldName: "value"},
}

var MonitorDimensionsExclusions = []string{
	"unit", "type", "value",
}

var MonitorDimensions = []string{}

const MonitorDataSource = "rb_monitor"
