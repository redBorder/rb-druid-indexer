package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	KAFKA_HOST     = "kafka.service:9092"
	ZOOKEEPER_HOST = "zookeeper.service:2181"
)

type TaskConfig struct {
	TaskName  string `yaml:"task_name"`
	Namespace string `yaml:"namespace"`
	Feed      string `yaml:"feed"`
	KafkaHost string `yaml:"kafka_host"`
}

type Config struct {
	ZookeeperServers []string     `yaml:"zookeeper_servers"`
	Tasks            []TaskConfig `yaml:"tasks"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	var config Config

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	if len(config.ZookeeperServers) == 0 {
		config.ZookeeperServers = []string{ZOOKEEPER_HOST}
	}

	if len(config.ZookeeperServers) == 0 {
		return nil, fmt.Errorf("zookeeper_servers must be provided in the config file")
	}

	for i, task := range config.Tasks {
		if task.KafkaHost == "" {
			config.Tasks[i].KafkaHost = KAFKA_HOST
		}
	}

	return &config, nil
}
