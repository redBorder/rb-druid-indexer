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

package config

import (
	"fmt"
	"os"
	"rb-druid-indexer/logger"

	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_DRUID_ROUTER  = "/druid/discoveryPath/druid:router"
	DEFAULT_KAFKA_BROKERS = "kafka.service:9092"
	ZOOKEEPER_HOST        = "zookeeper.service:2181"
)

type TaskConfig struct {
	TaskName         string   `yaml:"task_name"`
	Feed             string   `yaml:"feed"`
	Spec             string   `yaml:"spec"`
	KafkaBrokers     []string `yaml:"kafka_brokers"`
	CustomDimensions []string `yaml:"custom_dimensions"`
}

type Config struct {
	RouterDiscoveryPath string       `yaml:"discovery_path"`
	ZookeeperServers    []string     `yaml:"zookeeper_servers"`
	Tasks               []TaskConfig `yaml:"tasks"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Log.Errorf("could not open config file")
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	var config Config

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		logger.Log.Errorf("could not decode file")
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	if config.RouterDiscoveryPath == "" {
		config.RouterDiscoveryPath = DEFAULT_DRUID_ROUTER
	}

	if len(config.ZookeeperServers) == 0 {
		config.ZookeeperServers = []string{ZOOKEEPER_HOST}
	}

	if len(config.ZookeeperServers) == 0 {
		logger.Log.Errorf("zookeeper_servers must be provided in the config file")
		return nil, fmt.Errorf("zookeeper_servers must be provided in the config file")
	}

	for i, task := range config.Tasks {
		if len(task.KafkaBrokers) == 0 {
			config.Tasks[i].KafkaBrokers = []string{DEFAULT_KAFKA_BROKERS}
		}
		if task.CustomDimensions == nil {
			config.Tasks[i].CustomDimensions = []string{}
		}
	}

	return &config, nil
}
