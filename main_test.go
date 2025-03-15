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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"rb-druid-indexer/config"
	"rb-druid-indexer/logger"
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockZKClient struct {
	mock.Mock
}

type Metrics struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	FieldName string `json:"fieldName"`
}

func (m *MockZKClient) CreateLeaderNode() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockZKClient) IsLeader(nodePath string) bool {
	args := m.Called(nodePath)
	return args.Bool(0)
}

func (m *MockZKClient) GetConn() *zk.Conn {
	args := m.Called()
	return args.Get(0).(*zk.Conn)
}

type MockDruidRouter struct {
	mock.Mock
}

func (m *MockDruidRouter) GetSupervisors(routers []string) ([]string, error) {
	args := m.Called(routers)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockDruidRouter) DeleteTask(routers []string, taskName string) error {
	args := m.Called(routers, taskName)
	return args.Error(0)
}

func (m *MockDruidRouter) GenerateConfig(taskName string, kafkaBrokers []string, feed, timestampField, parseSpec string, dimensions, dimensionsExclusions []string, metrics []Metrics) (string, error) {
	args := m.Called(taskName, kafkaBrokers, feed, timestampField, parseSpec, dimensions, dimensionsExclusions, metrics)
	return args.String(0), args.Error(1)
}

func (m *MockDruidRouter) SubmitTask(routers []string, jsonStr string) error {
	args := m.Called(routers, jsonStr)
	return args.Error(0)
}

func createTestConfig(t *testing.T) string {
	content := `
zookeeper_servers:
  - "zookeeper1.service:2181"
  - "zookeeper2.service:2181"

tasks:
  - task_name: "rb_monitor"
    feed: "rb_monitor"
    spec: "rb_monitor"
    kafka_brokers:
      - "kafka1.service:9092"
      - "kafka2.service:9093"
  - task_name: "rb_state"
    feed: "rb_state_post"
    spec: "rb_state"
    kafka_brokers:
      - "kafka1.service:9092"
      - "kafka2.service:9093"
  - task_name: "rb_flow"
    feed: "rb_flow_post"
    spec: "rb_flow"
    kafka_brokers:
      - "kafka1.service:9092"
      - "kafka2.service:9093"
  - task_name: "rb_event"
    feed: "rb_event_post"
    spec: "rb_event"
    kafka_brokers:
      - "kafka1.service:9092"
      - "kafka2.service:9093"
  - task_name: "rb_vault"
    feed: "rb_vault_post"
    spec: "rb_vault"
    kafka_brokers:
      - "kafka1.service:9092"
      - "kafka2.service:9093"
  - task_name: "rb_scanner"
    feed: "rb_scanner_post"
    spec: "rb_scanner"
    kafka_brokers:
      - "kafka1.service:9092"
      - "kafka2.service:9093"
  - task_name: "rb_location"
    feed: "rb_loc_post"
    spec: "rb_location"
    kafka_brokers:
      - "kafka1.service:9092"
      - "kafka2.service:9093"
  - task_name: "rb_wireless"
    feed: "rb_wireless"
    spec: "rb_wireless"
    kafka_brokers:
      - "kafka1.service:9092"
      - "kafka2.service:9093"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	assert.NoError(t, err)
	return configPath
}

func TestMainFlow(t *testing.T) {
	logger.InitLogger()

	configPath := createTestConfig(t)

	cfg, err := config.LoadConfig(configPath)
	assert.NoError(t, err)

	fmt.Println(cfg)
	assert.NotNil(t, cfg.Tasks)
	assert.Equal(t, "rb_monitor", cfg.Tasks[0].TaskName)

	mockZK := new(MockZKClient)
	mockDruid := new(MockDruidRouter)

	mockZK.On("CreateLeaderNode").Return("/leader/node1", nil)
	mockZK.On("IsLeader", "/leader/node1").Return(true)

	routers := []string{"http://druid-router:8888"}
	supervisors := []string{"existing-task"}
	mockDruid.On("GetSupervisors", routers).Return(supervisors, nil)
	mockDruid.On("DeleteTask", routers, "existing-task").Return(nil)

	nodePath, err := mockZK.CreateLeaderNode()

	assert.NoError(t, err)
	assert.Equal(t, "/leader/node1", nodePath)

	isLeader := mockZK.IsLeader(nodePath)
	assert.True(t, isLeader)

	supervisorTasks, err := mockDruid.GetSupervisors(routers)
	assert.NoError(t, err)
	assert.Equal(t, supervisors, supervisorTasks)

	err = mockDruid.DeleteTask(routers, supervisorTasks[0])
	assert.NoError(t, err)

	taskNames := []string{"rb_monitor"}
	tasksToAnnounce := taskNames

	for _, taskName := range tasksToAnnounce {
		var taskConfig *config.TaskConfig
		for _, t := range cfg.Tasks {
			if t.TaskName == taskName {
				taskConfig = &t
				break
			}
		}

		assert.NotNil(t, taskConfig)

		dimensions := []string{"dim1", "dim2"}
		dimensionsExclusions := []string{}
		metrics := []Metrics{
			{Name: "metric1", Type: "count"},
			{Name: "metric2", Type: "doubleSum"},
		}

		mergedDimensions := append(dimensions, taskConfig.CustomDimensions...)

		mockDruid.On("GenerateConfig",
			taskConfig.TaskName,
			taskConfig.KafkaBrokers,
			taskConfig.Feed,
			"timestamp",
			"ruby",
			mergedDimensions,
			dimensionsExclusions,
			metrics,
		).Return(`{"test": "config"}`, nil)

		mockDruid.On("SubmitTask", routers, `{"test": "config"}`).Return(nil)

		jsonStr, err := mockDruid.GenerateConfig(
			taskConfig.TaskName,
			taskConfig.KafkaBrokers,
			taskConfig.Feed,
			"timestamp",
			"ruby",
			mergedDimensions,
			dimensionsExclusions,
			metrics,
		)
		assert.NoError(t, err)
		assert.Equal(t, `{"test": "config"}`, jsonStr)

		err = mockDruid.SubmitTask(routers, jsonStr)
		assert.NoError(t, err)
	}

	mockZK.AssertExpectations(t)
	mockDruid.AssertExpectations(t)
}

func TestMainFlowNotLeader(t *testing.T) {
	logger.InitLogger()

	configPath := createTestConfig(t)

	_, err := config.LoadConfig(configPath)
	assert.NoError(t, err)

	mockZK := new(MockZKClient)
	mockZK.On("CreateLeaderNode").Return("/leader/node1", nil)
	mockZK.On("IsLeader", "/leader/node1").Return(false)

	nodePath, err := mockZK.CreateLeaderNode()
	assert.NoError(t, err)
	assert.Equal(t, "/leader/node1", nodePath)

	isLeader := mockZK.IsLeader(nodePath)
	assert.False(t, isLeader)

	mockZK.AssertExpectations(t)
}

func TestZKReconnection(t *testing.T) {
	logger.InitLogger()

	configPath := createTestConfig(t)

	_, err := config.LoadConfig(configPath)
	assert.NoError(t, err)

	mockZK := new(MockZKClient)
	mockZK.On("CreateLeaderNode").Return("/leader/node1", nil)
	mockZK.On("GetConn").Return(&zk.Conn{})

	nodePath, err := mockZK.CreateLeaderNode()
	assert.NoError(t, err)
	assert.Equal(t, "/leader/node1", nodePath)

	conn := mockZK.GetConn()
	assert.NotNil(t, conn)

	mockZK.AssertExpectations(t)
}
