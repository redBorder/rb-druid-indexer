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
	"os"
	"rb-druid-indexer/logger"
	"fmt"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DataSources map[string]DataSource `yaml:"DataSources"` // Mapa de nombre de fuente a su configuración
}

type DataSource struct {
	DataSource       string   `yaml:"DataSource"`
	Metrics          []druidrouter.Metrics `yaml:"metrics"`
	Dimensions           []string `yaml:"dimensions"`
	DimensionsExclusions []string `yaml:"dimensions_exclusions"`
}

var DataSourceConfig, err = LoadDimensions("/var/www/rb-rails/config/dimensions.yml")

func GetDataSourceConfig(taskName string) (DataSource, bool) {
	config, exists := DataSourceConfig.DataSources[taskName]
	return config, exists
}

func LoadDimensions(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Log.Errorf("could not open config file")
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()
	
	var config Config

	// Usamos un decodificador YAML
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	// logger.Log.Errorf("\n\nDECODED: %#v", config)
	// Devolvemos la estructura de configuración parseada
	return &config, nil
}