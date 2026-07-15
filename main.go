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
	"fmt"
	"os"
	"rb-druid-indexer/config"
	druidrouter "rb-druid-indexer/druid"
	"rb-druid-indexer/kafkaclient"
	"rb-druid-indexer/logger"
	zkclient "rb-druid-indexer/zkclient"
	"strings"
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

			if zk.GetConn() != nil {
				zk.GetConn().Close()
			}

			newZK, err := zkclient.NewZKClient(cfg.ZookeeperServers)
			if err != nil {
				logger.Log.Errorf("Failed to reconnect to Zookeeper: %v", err)
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}
			zk = newZK

			newNodePath, err := zk.CreateLeaderNode()
			if err != nil {
				logger.Log.Errorf("Failed to recreate leader node after reconnect: %v", err)
			} else {
				nodePath = newNodePath
			}
		}

		if !zk.IsLeader(nodePath) {
			logger.Log.Info("I am not the leader. Waiting...")
			time.Sleep(60 * time.Second)
			continue
		}

		routers, err := zkclient.GetDruidRouterInfo(zk.GetConn(), cfg.RouterDiscoveryPath)
		if err != nil {
			logger.Log.Errorf("Error retrieving Druid Router info from ZooKeeper: %v. Retrying in 60s...", err)
			time.Sleep(60 * time.Second)
			continue
		}

		supervisorTasks, err := druidrouter.GetSupervisors(routers)
		if err != nil {
			logger.Log.Errorf("Failed to get supervisor tasks: %v. Retrying in 60s...", err)
			time.Sleep(60 * time.Second)
			continue
		}

		if !cleaned {
			logger.Log.Info("Cleaning Supervisors (first time)...")
			for _, task := range supervisorTasks {
				druidrouter.DeleteTask(routers, task)
			}
			supervisorTasks = []string{}
			cleaned = true
		}

		for _, taskConfig := range cfg.Tasks {
			isEmpty, err := kafkaclient.IsTopicEmpty(taskConfig.KafkaBrokers, taskConfig.Feed)
			if err != nil {
				logger.Log.Warnf("Failed to check if Kafka topic %s is empty: %v. Assuming it has messages.", taskConfig.Feed, err)
				isEmpty = false
			}

			isSupervisorRunning := false
			for _, st := range supervisorTasks {
				if st == taskConfig.TaskName {
					isSupervisorRunning = true
					break
				}
			}

			if isEmpty {
				if isSupervisorRunning {
					logger.Log.Infof("Topic %s is empty. Terminating supervisor task %s to free resources.", taskConfig.Feed, taskConfig.TaskName)
					druidrouter.DeleteTask(routers, taskConfig.TaskName)
				} else {
					logger.Log.Debugf("Topic %s is empty and supervisor task %s is not running. Skipping.", taskConfig.Feed, taskConfig.TaskName)
				}
			} else {
				if !isSupervisorRunning {
					logger.Log.Infof("Topic %s has messages. Submitting supervisor task %s.", taskConfig.Feed, taskConfig.TaskName)
					dimensions := []string{}
					if taskConfig.Dimensions != nil {
						dimensions = taskConfig.Dimensions
					}

					dimensionsExclusions := []string{}
					if taskConfig.DimensionsExclusions != nil {
						dimensionsExclusions = taskConfig.DimensionsExclusions
					}

					jsonStr, err := druidrouter.GenerateConfig(
						taskConfig.TaskName,
						taskConfig.KafkaBrokers,
						taskConfig.Feed,
						"timestamp",
						"ruby",
						dimensions,
						dimensionsExclusions,
						taskConfig.Metrics,
					)
					if err != nil {
						logger.Log.Fatalf("Error generating config for task %s: %v", taskConfig.TaskName, err)
					}

					druidrouter.SubmitTask(routers, jsonStr)

					// Submit/Configure auto-compaction if enabled
					if taskConfig.Compaction != nil && *taskConfig.Compaction {
						err = druidrouter.SubmitCompaction(routers, taskConfig.TaskName, taskConfig.CompactionFrequency, taskConfig.SkipOffsetFromLatest)
						if err != nil {
							logger.Log.Errorf("Failed to configure compaction for %s: %v", taskConfig.TaskName, err)
						}
					} else {
						// Explicitly delete compaction config if disabled
						err = druidrouter.DeleteCompaction(routers, taskConfig.TaskName)
						if err != nil {
							logger.Log.Errorf("Failed to remove compaction for %s: %v", taskConfig.TaskName, err)
						}
					}
				} else {
					logger.Log.Debugf("Topic %s has messages and supervisor task %s is already running.", taskConfig.Feed, taskConfig.TaskName)

					// Check if Kafka offsets were reset/recreated
					status, rawStatus, err := druidrouter.GetSupervisorStatus(routers, taskConfig.TaskName)
					if err != nil {
						logger.Log.Warnf("Failed to get supervisor status for %s: %v", taskConfig.TaskName, err)
					} else {
						logger.Log.Infof("Supervisor %s - Healthy: %v, DetailedState: %s", taskConfig.TaskName, status.Payload.Healthy, status.Payload.DetailedState)
						kafkaOffsets, err := kafkaclient.GetLatestOffsets(taskConfig.KafkaBrokers, taskConfig.Feed)
						if err != nil {
							logger.Log.Warnf("Failed to get Kafka offsets for %s: %v", taskConfig.Feed, err)
						} else {
							wasReset := false

							// Check active tasks starting offsets vs Kafka latest offsets
							for _, task := range status.Payload.ActiveTasks {
								for partitionStr, startingOffset := range task.StartingOffsets {
									var partitionID int
									_, err := fmt.Sscan(partitionStr, &partitionID)
									if err == nil {
										if latestKafkaOffset, ok := kafkaOffsets[partitionID]; ok {
											if startingOffset > latestKafkaOffset {
												logger.Log.Warnf("Detected Kafka offset reset/recreation: task %s starting offset (%d) > Kafka latest offset (%d) on partition %d.", task.Id, startingOffset, latestKafkaOffset, partitionID)
												wasReset = true
												break
											}
										}
									}
								}
								if wasReset {
									break
								}
							}

							// Also compare supervisor latestOffsets vs Kafka latest offsets
							if !wasReset {
								for partitionStr, latestOffset := range status.Payload.LatestOffsets {
									var partitionID int
									_, err := fmt.Sscan(partitionStr, &partitionID)
									if err == nil {
										if latestKafkaOffset, ok := kafkaOffsets[partitionID]; ok {
											if latestOffset > latestKafkaOffset {
												logger.Log.Warnf("Detected Kafka offset reset/recreation: supervisor latestOffset (%d) > Kafka latest offset (%d) on partition %d.", latestOffset, latestKafkaOffset, partitionID)
												wasReset = true
												break
											}
										}
									}
								}
							}

							// Also check if supervisor is unhealthy due to OffsetOutOfRangeException
							hasOutOfRangeErr := strings.Contains(rawStatus, "OffsetOutOfRangeException") || strings.Contains(rawStatus, "offset out of range")
							if status.Payload.Healthy == false && hasOutOfRangeErr {
								wasReset = true
							}

							if wasReset {
								logger.Log.Warnf("Detected Kafka offset reset/recreation or out-of-range error for topic %s. Resetting supervisor %s.", taskConfig.Feed, taskConfig.TaskName)
								err := druidrouter.ResetSupervisor(routers, taskConfig.TaskName)
								if err != nil {
									logger.Log.Errorf("Failed to reset supervisor %s: %v", taskConfig.TaskName, err)
								}
								err = druidrouter.ResetOffsetsSupervisor(routers, taskConfig.TaskName)
								if err != nil {
									logger.Log.Errorf("Failed to resetOffsets supervisor %s: %v", taskConfig.TaskName, err)
								}
							}
						}
					}
				}
			}
		}

		time.Sleep(60 * time.Second)
	}
}
