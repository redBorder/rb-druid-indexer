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
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"rb-druid-indexer/logger"
	"rb-druid-indexer/zkclient"

	"github.com/sirupsen/logrus"
)

func init() {
	if logger.Log == nil {
		l := logrus.New()
		l.Out = io.Discard
		logger.Log = l
	}
}

func TestGetSupervisors(t *testing.T) {
	expectedSupervisors := []string{"sup1", "sup2"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		if r.URL.Path != "/druid/indexer/v1/supervisor" {
			t.Errorf("Unexpected URL path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedSupervisors); err != nil {
			t.Errorf("Error encoding JSON response: %v", err)
		}
	}))
	defer ts.Close()

	parts := strings.Split(strings.TrimPrefix(ts.URL, "http://"), ":")
	if len(parts) != 2 {
		t.Fatalf("Unexpected test server URL format: %s", ts.URL)
	}
	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		t.Fatalf("Failed to convert port: %v", err)
	}

	routers := []zkclient.DruidRouter{
		{Address: host, Port: port},
	}

	supervisors, err := GetSupervisors(routers)
	if err != nil {
		t.Fatalf("GetSupervisors returned error: %v", err)
	}

	if !reflect.DeepEqual(supervisors, expectedSupervisors) {
		t.Errorf("Expected supervisors %v, got %v", expectedSupervisors, supervisors)
	}
}

func TestSubmitTask(t *testing.T) {
	task := `{"id": "task1", "type": "index"}`

	var received string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/druid/indexer/v1/supervisor" {
			t.Errorf("Unexpected URL path: %s", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading request body: %v", err)
		}
		received = string(body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("submitted"))
	}))
	defer ts.Close()

	parts := strings.Split(strings.TrimPrefix(ts.URL, "http://"), ":")
	if len(parts) != 2 {
		t.Fatalf("Unexpected test server URL format: %s", ts.URL)
	}
	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		t.Fatalf("Failed to convert port: %v", err)
	}

	routers := []zkclient.DruidRouter{
		{Address: host, Port: port},
	}

	SubmitTask(routers, task)

	if received != task {
		t.Errorf("Expected task body %q, got %q", task, received)
	}
}

func TestDeleteTask(t *testing.T) {
	task := "task1"
	expectedPath := "/druid/indexer/v1/supervisor/" + task + "/terminate"

	var received string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("Unexpected URL path: got %s, expected %s", r.URL.Path, expectedPath)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading request body: %v", err)
		}
		received = string(body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("deleted"))
	}))
	defer ts.Close()

	parts := strings.Split(strings.TrimPrefix(ts.URL, "http://"), ":")
	if len(parts) != 2 {
		t.Fatalf("Unexpected test server URL format: %s", ts.URL)
	}
	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		t.Fatalf("Failed to convert port: %v", err)
	}

	routers := []zkclient.DruidRouter{
		{Address: host, Port: port},
	}

	DeleteTask(routers, task)

	if received != task {
		t.Errorf("Expected task body %q, got %q", task, received)
	}
}
