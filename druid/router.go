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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rb-druid-indexer/logger"
	"rb-druid-indexer/zkclient"
	"strings"
	"time"
)

type BuildSegmentsStats struct {
	Processed          float64 `json:"processed"`
	ProcessedBytes     float64 `json:"processedBytes"`
	Unparseable        float64 `json:"unparseable"`
	ThrownAway         float64 `json:"thrownAway"`
	ProcessedWithError float64 `json:"processedWithError"`
}

type BuildSegments struct {
	FiveM    BuildSegmentsStats `json:"5m"`
	FifteenM BuildSegmentsStats `json:"15m"`
	OneM     BuildSegmentsStats `json:"1m"`
}

type MovingAverages struct {
	BuildSegments BuildSegments `json:"buildSegments"`
}

type TotalsBuildSegments struct {
	Processed          int     `json:"processed"`
	ProcessedBytes     float64 `json:"processedBytes"`
	ProcessedWithError int     `json:"processedWithError"`
	ThrownAway         int     `json:"thrownAway"`
	Unparseable        int     `json:"unparseable"`
}

type Totals struct {
	BuildSegments TotalsBuildSegments `json:"buildSegments"`
}

type SupervisorStats struct {
	MovingAverages MovingAverages `json:"movingAverages"`
	Totals         Totals         `json:"totals"`
}

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
