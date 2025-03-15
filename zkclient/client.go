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

package zkclient

import (
	"encoding/json"
	"fmt"
	"rb-druid-indexer/logger"

	"github.com/samuel/go-zookeeper/zk"
)

type DruidRouter struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	SSLPort     int    `json:"sslPort"`
	ServiceType string `json:"serviceType"`
}

type ZkClient interface {
	Children(path string) ([]string, *zk.Stat, error)
	Get(path string) ([]byte, *zk.Stat, error)
}

func GetDruidRouterInfo(conn ZkClient, RouterDiscoveryPath string) ([]DruidRouter, error) {
	if conn == nil {
		return nil, fmt.Errorf("zookeeper connection is nil")
	}

	children, _, err := conn.Children(RouterDiscoveryPath)
	if err != nil {
		return nil, err
	}

	if len(children) == 0 {
		logger.Log.Errorf("no routers found in Zookeeper")
		return nil, fmt.Errorf("no routers found in Zookeeper")
	}

	var routers []DruidRouter
	for _, child := range children {
		routerNode := RouterDiscoveryPath + "/" + child
		data, _, err := conn.Get(routerNode)
		if err != nil {
			logger.Log.Errorf("failed to get data for node %s: %v", routerNode, err)
			continue
		}

		var routerInfo DruidRouter
		err = json.Unmarshal(data, &routerInfo)
		if err != nil {
			logger.Log.Errorf("failed to unmarshal data for node %s: %v", routerNode, err)
			continue
		}

		routers = append(routers, routerInfo)
	}

	if len(routers) == 0 {
		return nil, fmt.Errorf("failed to retrieve any valid router information")
	}

	return routers, nil
}
