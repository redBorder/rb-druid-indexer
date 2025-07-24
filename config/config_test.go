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
	"os"
	"path/filepath"
	"testing"

	"rb-druid-indexer/logger"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	druidrouter "rb-druid-indexer/druid"

)

type Metrics struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	FieldName string `json:"fieldName"`
}

func TestLoadConfig(t *testing.T) {
	logger.Log = logrus.New()

	validDir := t.TempDir()
	missingDir := t.TempDir()

	tests := []struct {
		name        string
		filePath    string
		expectedErr bool
	}{
		{
			name:        "Valid config",
			filePath:    filepath.Join(validDir, "example_config.yaml"),
			expectedErr: false,
		},
		{
			name:        "Missing config file",
			filePath:    filepath.Join(missingDir, "non_existing_config.yaml"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if !tt.expectedErr {
				createTestConfigFile(t, tt.filePath)
			}

			cfg, err := LoadConfig(tt.filePath)

			if (err != nil) != tt.expectedErr {
				t.Errorf("LoadConfig() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			if !tt.expectedErr {
				if cfg.RouterDiscoveryPath != DEFAULT_DRUID_ROUTER {
					t.Errorf("expected RouterDiscoveryPath %q, got %q", DEFAULT_DRUID_ROUTER, cfg.RouterDiscoveryPath)
				}
				if len(cfg.ZookeeperServers) != 1 || cfg.ZookeeperServers[0] != ZOOKEEPER_HOST {
					t.Errorf("expected ZookeeperServers [%q], got %v", ZOOKEEPER_HOST, cfg.ZookeeperServers)
				}
				if len(cfg.Tasks) != 1 {
					t.Errorf("expected 1 task, got %d", len(cfg.Tasks))
				} else {
					task := cfg.Tasks[0]
					if task.TaskName != "task1" {
						t.Errorf("expected task name %q, got %q", "task1", task.TaskName)
					}
					if task.Feed != "feed1" {
						t.Errorf("expected feed %q, got %q", "feed1", task.Feed)
					}
					if task.Spec != "spec1" {
						t.Errorf("expected spec %q, got %q", "spec1", task.Spec)
					}
					if len(task.KafkaBrokers) != 2 ||
						task.KafkaBrokers[0] != "kafka1" ||
						task.KafkaBrokers[1] != "kafka2" {
						t.Errorf("unexpected KafkaBrokers: got %v, expected [\"kafka1\", \"kafka2\"]", task.KafkaBrokers)
					}
					if len(task.Dimensions) != 0 {
						t.Errorf("expected empty dimensions, got %v", task.Dimensions)
					}
					if len(task.DimensionsExclusions) != 0 {
						t.Errorf("expected empty dimensions exclusions, got %v", task.DimensionsExclusions)
					}
				}
			}
		})
	}
}

func createTestConfigFile(t *testing.T, filePath string) {
	config := &Config{
		RouterDiscoveryPath: DEFAULT_DRUID_ROUTER,
		ZookeeperServers:    []string{ZOOKEEPER_HOST},
		Tasks: []TaskConfig{
			{
				TaskName:     "task1",
				Feed:         "feed1",
				Spec:         "spec1",
				KafkaBrokers: []string{"kafka1", "kafka2"},
				Dimensions: []string{},
				DimensionsExclusions: []string {},
				Metrics: []druidrouter.Metrics{},
			},
		},
	}

	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("could not create test config file: %v", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(config); err != nil {
		t.Fatalf("could not encode test config file: %v", err)
	}
}
