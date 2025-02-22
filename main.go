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
	"flag"
	"rb-druid-indexer/config"
	druidrouter "rb-druid-indexer/druid"
	druiddatasources "rb-druid-indexer/druid/datasources"
	"rb-druid-indexer/logger"
	zkclient "rb-druid-indexer/zkclient"
	"time"
)

func main() {

	logger.InitLogger()

	configFilePath := flag.String("config", "config.yml", "Path to the configuration file (YAML)")

	flag.Parse()

	cfg, err := config.LoadConfig(*configFilePath)
	if err != nil {
		logger.Log.Fatalf("Error loading configuration: %v", err)
	}

	router, err := zkclient.GetDruidRouterInfo(cfg.ZookeeperServers)
	if err != nil {
		logger.Log.Fatalf("Error retrieving Druid Router info from ZooKeeper: %v", err)
	}

	zk, err := zkclient.NewZKClient(cfg.ZookeeperServers)
	if err != nil {
		logger.Log.Fatalf("Error connecting to ZooKeeper: %v", err)
	}

	nodePath, err := zk.CreateLeaderNode()
	if err != nil {
		logger.Log.Fatalf("Error creating leader node: %v", err)
	}

	for {
		if !zk.IsLeader(nodePath) {
			logger.Log.Info("I am not the leader. Waiting...")
			time.Sleep(60 * time.Second)
			continue
		}

		supervisorTasks, err := druidrouter.GetSupervisors(router.Address, router.Port)
		if err != nil {
			logger.Log.Fatalf("Failed to get supervisor tasks: %v", err)
		}

		var taskNames []string
		for _, taskConfig := range cfg.Tasks {
			taskNames = append(taskNames, taskConfig.TaskName)
		}

		tasksToAnnounce := zkclient.TaskAnnouncer(supervisorTasks, taskNames)

		for _, taskName := range tasksToAnnounce {
			var taskConfig *config.TaskConfig
			for _, t := range cfg.Tasks {
				if t.TaskName == taskName {
					taskConfig = &t
					break
				}
			}
			if taskConfig == nil {
				logger.Log.Fatalf("No configuration found for task: %s", taskName)
			}

			config, found := druiddatasources.GetDataSourceConfig(taskConfig.TaskName)
			if !found {
				logger.Log.Fatalf("No configuration found for data source: %s", taskConfig.TaskName)
			}

			jsonStr, err := druidrouter.GenerateConfig(
				taskConfig.Namespace,
				config.DataSource,
				taskConfig.KafkaHost,
				taskConfig.Feed,
				"timestamp",
				"ruby",
				config.Dimensions,
				config.Metrics,
			)
			if err != nil {
				logger.Log.Fatalf("Error generating config for task %s: %v", taskConfig.TaskName, err)
			}

			druidrouter.SubmitTask(router.Address, router.Port, jsonStr)
		}

		time.Sleep(60 * time.Second)
	}
}
