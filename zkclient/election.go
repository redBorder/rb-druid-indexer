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
	"fmt"
	"rb-druid-indexer/logger"
	"sort"
	"strconv"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	LEADER_ELECTION_PATH = "/rb-druid-indexer"
)

type ZKClient struct {
	conn *zk.Conn
}

type SilentLogger struct{}

func (s SilentLogger) Printf(format string, a ...interface{}) {
	// Do nothing to suppress the logs
}

func NewZKClient(zookeeperServers []string) (*ZKClient, error) {
	conn, _, err := zk.Connect(zookeeperServers, 5*time.Second, zk.WithLogger(SilentLogger{}))
	if err != nil {
		return nil, err
	}
	return &ZKClient{conn: conn}, nil
}

func (z *ZKClient) GetConn() *zk.Conn {
	return z.conn
}

func IsZKAlive(conn *zk.Conn) bool {
	if conn == nil {
		return false
	}

	for i := 0; i < 3; i++ {
		state := conn.State()

		if state == zk.StateConnected || state == zk.StateHasSession {
			return true
		}

		time.Sleep(1 * time.Second)
	}
	return false
}

func (zkClient *ZKClient) CreateLeaderNode() (string, error) {
	exists, _, err := zkClient.conn.Exists(LEADER_ELECTION_PATH)
	if err != nil {
		return "", err
	}
	if !exists {
		_, err := zkClient.conn.Create(LEADER_ELECTION_PATH, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return "", err
		}
	}
	return zkClient.conn.Create(LEADER_ELECTION_PATH+"/node-", []byte{}, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
}

func extractSeq(nodeName string) (int64, error) {
	if len(nodeName) < 10 {
		return 0, fmt.Errorf("nodeName %q is too short to extract sequence", nodeName)
	}
	seq := nodeName[len(nodeName)-10:]
	return strconv.ParseInt(seq, 10, 64)
}

func (zkClient *ZKClient) GetLeader() (string, error) {
	children, _, err := zkClient.conn.Children(LEADER_ELECTION_PATH)
	if err != nil {
		return "", err
	}

	var validChildren []string
	for _, child := range children {
		if _, err := extractSeq(child); err == nil {
			validChildren = append(validChildren, child)
		}
	}

	if len(validChildren) == 0 {
		logger.Log.Warnf("no zookeeper leader found under %s", LEADER_ELECTION_PATH)
		return "", fmt.Errorf("no leader found")
	}

	sort.Slice(validChildren, func(i, j int) bool {
		seqI, _ := extractSeq(validChildren[i])
		seqJ, _ := extractSeq(validChildren[j])
		return seqI < seqJ
	})

	return validChildren[0], nil
}

func (zkClient *ZKClient) IsLeader(nodePath string) bool {
	if nodePath == "" {
		return false
	}
	leader, err := zkClient.GetLeader()
	if err != nil {
		logger.Log.Errorf("Error getting leader: %v", err)
		return false
	}
	return nodePath == LEADER_ELECTION_PATH+"/"+leader
}
