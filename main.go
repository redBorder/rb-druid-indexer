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
	"os"
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

	zk, err := zkclient.NewZKClient(cfg.ZookeeperServers)
	if err != nil {
		logger.Log.Fatalf("Error connecting to ZooKeeper: %v", err)
	}

	nodePath, err := zk.CreateLeaderNode()
	if err != nil {
		logger.Log.Fatalf("Error creating leader node: %v", err)
	}

	cleaned := false

	for {
		if !zkclient.IsZKAlive(zk.GetConn()) {
			logger.Log.Warn("Zookeeper connection lost, reconnecting...")
		
			newZK, err := zkclient.NewZKClient(cfg.ZookeeperServers)
			if err != nil {
				logger.Log.Errorf("Failed to reconnect to Zookeeper: %v", err)
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}
			zk = newZK
		}

		if !zk.IsLeader(nodePath) {
			logger.Log.Info("I am not the leader. Waiting...")
			time.Sleep(60 * time.Second)
			continue
		}

		routers, err := zkclient.GetDruidRouterInfo(zk.GetConn())
		if err != nil {
			logger.Log.Fatalf("Error retrieving Druid Router info from ZooKeeper: %v", err)
		}

		supervisorTasks, err := druidrouter.GetSupervisors(routers)
		if err != nil {
			logger.Log.Fatalf("Failed to get supervisor tasks: %v", err)
		}

		if !cleaned {
			logger.Log.Info("Cleaning Supervisors (first time)...")
			for _, task := range supervisorTasks {
				druidrouter.DeleteTask(routers, task)
			}
			cleaned = true
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

			config, found := druiddatasources.GetDataSourceConfig(taskConfig.Spec)
			if !found {
				logger.Log.Fatalf("No configuration found for data source: %s", taskConfig.TaskName)
			}

			mergedDimensions := append(config.Dimensions, taskConfig.CustomDimensions...)

			jsonStr, err := druidrouter.GenerateConfig(
				taskConfig.TaskName,
				taskConfig.KafkaBrokers,
				taskConfig.Feed,
				"timestamp",
				"ruby",
				mergedDimensions,
				config.DimensionsExclusions,
				config.Metrics,
			)
			if err != nil {
				logger.Log.Fatalf("Error generating config for task %s: %v", taskConfig.TaskName, err)
			}

			druidrouter.SubmitTask(routers, jsonStr)
		}

		time.Sleep(60 * time.Second)
	}
}
