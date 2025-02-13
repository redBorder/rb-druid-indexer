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
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetSupervisors(host string, port int) ([]string, error) {
	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor", host, port)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch supervisors, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var supervisors []string
	err = json.Unmarshal(body, &supervisors)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return supervisors, nil
}

func SubmitTask(host string, port int, task string) {
	url := fmt.Sprintf("http://%s:%d/druid/indexer/v1/supervisor", host, port)
	resp, err := http.Post(url, "application/json", strings.NewReader(task))
	if err != nil {
		log.Printf("Error submitting task: %v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code %d, response: %s", resp.StatusCode, string(body))
		return
	}
	log.Printf("Task submitted successfully: %s", string(body))
}
