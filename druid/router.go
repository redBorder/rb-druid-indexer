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
	"log"
	"net/http"
	"rb-druid-indexer/logger"
	"strings"
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

func GetSupervisors(host string, port int) ([]string, error) {
	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor", host, port)

	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Errorf("Failed to send GET request: %v", err)
		return nil, fmt.Errorf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Log.Warnf("Failed to fetch supervisors, status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("failed to fetch supervisors, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var supervisors []string
	err = json.Unmarshal(body, &supervisors)
	if err != nil {
		logger.Log.Errorf("Failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	logger.Log.Infof("Successfully fetched supervisors: %v", supervisors)
	return supervisors, nil
}

func SubmitTask(host string, port int, task string) {
	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor", host, port)
	resp, err := http.Post(url, "application/json", strings.NewReader(task))
	if err != nil {
		logger.Log.Errorf("Error submitting task: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Error reading response: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		logger.Log.Warnf("Unexpected status code %d, response: %s", resp.StatusCode, string(body))
		return
	}

	logger.Log.Infof("Task submitted successfully: %s", string(body))
}

func DeleteTask(host string, port int, task string) {
	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor/%s/terminate", host, port, task)
	resp, err := http.Post(url, "application/json", strings.NewReader(task))
	if err != nil {
		logger.Log.Errorf("Error submitting task: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Error reading response: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		logger.Log.Warnf("Unexpected status code %d, response: %s", resp.StatusCode, string(body))
		return
	}

	logger.Log.Infof("DeleteTask submitted successfully: %s", string(body))
}

func CheckStats(host string, port int, task string) (map[string]map[string]SupervisorStats, error) {
	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor/%s/stats", host, port, task)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error sending GET request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code %d, response: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var supervisors map[string]map[string]SupervisorStats
	if err := json.Unmarshal(body, &supervisors); err != nil {
		log.Printf("Error unmarshalling JSON response: %v", err)
		return nil, err
	}

	return supervisors, nil
}

func ResetSupervisorOffset(host string, port int, task string) (bool, error) {
	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor/%s/resetOffsets", host, port, task)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error sending GET request: %v", err)
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code %d, response: %s", resp.StatusCode, string(body))
		return false, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	} else {
		return true, nil
	}
}
