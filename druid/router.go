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

package druidrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rb-druid-indexer/logger"
	"rb-druid-indexer/zkclient"
	"strings"
	"time"
)

func GetSupervisors(routers []zkclient.DruidRouter) ([]string, error) {
	var allSupervisors []string

	randomIndex := int(time.Now().UnixNano() % int64(len(routers)))
	router := routers[randomIndex]

	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor", router.Address, router.Port)

	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Errorf("Failed to send GET request to %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Log.Warnf("Failed to fetch supervisors from %s, status code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Failed to read response body from %s: %v", url, err)
	}

	var supervisors []string
	err = json.Unmarshal(body, &supervisors)
	if err != nil {
		logger.Log.Errorf("Failed to unmarshal response from %s: %v", url, err)
	}

	logger.Log.Infof("Successfully fetched supervisors from %s: %v", url, supervisors)
	allSupervisors = append(allSupervisors, supervisors...)

	return allSupervisors, nil
}

func SubmitTask(routers []zkclient.DruidRouter, task string) {
	if len(routers) == 0 {
		logger.Log.Errorf("No available routers to submit the task")
		return
	}

	randomIndex := int(time.Now().UnixNano() % int64(len(routers)))
	router := routers[randomIndex]

	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor", router.Address, router.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(task))
	if err != nil {
		logger.Log.Errorf("Error submitting task to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Error reading response from %s: %v", url, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		logger.Log.Warnf("Unexpected status code %d from %s, response: %s", resp.StatusCode, url, string(body))
		return
	}

	logger.Log.Infof("Task submitted successfully to %s: %s", url, string(body))
}

func DeleteTask(routers []zkclient.DruidRouter, task string) {

	randomIndex := int(time.Now().UnixNano() % int64(len(routers)))
	router := routers[randomIndex]

	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor/%s/terminate", router.Address, router.Port, task)
	resp, err := http.Post(url, "application/json", strings.NewReader(task))
	if err != nil {
		logger.Log.Errorf("Error deleting task %s from %s: %v", task, url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Error reading response for task %s from %s: %v", task, url, err)
	}

	if resp.StatusCode != http.StatusOK {
		logger.Log.Warnf("Unexpected status code %d for task %s from %s, response: %s", resp.StatusCode, task, url, string(body))
	}

	logger.Log.Infof("Task %s deleted successfully from %s: %s", task, url, string(body))
}

type ActiveTask struct {
	Id              string           `json:"id"`
	StartingOffsets map[string]int64 `json:"startingOffsets"`
	Type            string           `json:"type"`
}

type SupervisorPayload struct {
	DataSource    string           `json:"dataSource"`
	Stream        string           `json:"stream"`
	Healthy       bool             `json:"healthy"`
	DetailedState string           `json:"detailedState"`
	LatestOffsets map[string]int64 `json:"latestOffsets"`
	ActiveTasks   []ActiveTask     `json:"activeTasks"`
}

type SupervisorStatus struct {
	Id             string            `json:"id"`
	GenerationTime string            `json:"generationTime"`
	Payload        SupervisorPayload `json:"payload"`
}

func GetSupervisorStatus(routers []zkclient.DruidRouter, supervisor string) (*SupervisorStatus, string, error) {
	if len(routers) == 0 {
		return nil, "", fmt.Errorf("no available routers")
	}

	randomIndex := int(time.Now().UnixNano() % int64(len(routers)))
	router := routers[randomIndex]

	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor/%s/status", router.Address, router.Port, supervisor)
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch supervisor status from %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code %d from %s, response: %s", resp.StatusCode, url, string(body))
	}

	var status SupervisorStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal response from %s: %w", url, err)
	}

	return &status, string(body), nil
}

func ResetSupervisor(routers []zkclient.DruidRouter, supervisor string) error {
	if len(routers) == 0 {
		return fmt.Errorf("no available routers")
	}

	randomIndex := int(time.Now().UnixNano() % int64(len(routers)))
	router := routers[randomIndex]

	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor/%s/reset", router.Address, router.Port, supervisor)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to reset supervisor from %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d from %s, response: %s", resp.StatusCode, url, string(body))
	}

	logger.Log.Infof("Supervisor %s reset successfully on %s: %s", supervisor, url, string(body))
	return nil
}

func ResetOffsetsSupervisor(routers []zkclient.DruidRouter, supervisor string) error {
	if len(routers) == 0 {
		return fmt.Errorf("no available routers")
	}

	randomIndex := int(time.Now().UnixNano() % int64(len(routers)))
	router := routers[randomIndex]

	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor/%s/resetOffsets", router.Address, router.Port, supervisor)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to resetOffsets supervisor from %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d from %s, response: %s", resp.StatusCode, url, string(body))
	}

	logger.Log.Infof("Supervisor %s resetOffsets successfully on %s: %s", supervisor, url, string(body))
	return nil
}

func SubmitCompaction(routers []zkclient.DruidRouter, dataSource string, frequency string, skipOffset string) error {
	if len(routers) == 0 {
		return fmt.Errorf("no available routers")
	}

	randomIndex := int(time.Now().UnixNano() % int64(len(routers)))
	router := routers[randomIndex]

	url := fmt.Sprintf("http://%s:%d/druid/coordinator/v1/config/compaction", router.Address, router.Port)
	
	payload := map[string]interface{}{
		"dataSource":           dataSource,
		"segmentGranularity":   frequency,
		"skipOffsetFromLatest": skipOffset,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal compaction config: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to submit compaction config to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d from %s, response: %s", resp.StatusCode, url, string(body))
	}

	logger.Log.Infof("Compaction config submitted successfully for %s on %s: %s", dataSource, url, string(body))
	return nil
}

func DeleteCompaction(routers []zkclient.DruidRouter, dataSource string) error {
	if len(routers) == 0 {
		return fmt.Errorf("no available routers")
	}

	randomIndex := int(time.Now().UnixNano() % int64(len(routers)))
	router := routers[randomIndex]

	url := fmt.Sprintf("http://%s:%d/druid/coordinator/v1/config/compaction/%s", router.Address, router.Port, dataSource)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete compaction config from %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status code %d from %s, response: %s", resp.StatusCode, url, string(body))
	}

	logger.Log.Infof("Compaction config deleted successfully for %s on %s: %s", dataSource, url, string(body))
	return nil
}



