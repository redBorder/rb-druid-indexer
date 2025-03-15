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
	"log"
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

func NewZKClient(zookeeperServers []string) (*ZKClient, error) {
	conn, _, err := zk.Connect(zookeeperServers, 5*time.Second)
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
		if err != nil {
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
	sort.Slice(children, func(i, j int) bool {
		seqI, errI := extractSeq(children[i])
		seqJ, errJ := extractSeq(children[j])
		if errI != nil || errJ != nil {
			return false
		}
		return seqI < seqJ
	})
	if len(children) > 0 {
		return children[0], nil
	}
	logger.Log.Errorf("no zookeeper leader found")
	return "", fmt.Errorf("no leader found")
}

func (zkClient *ZKClient) IsLeader(nodePath string) bool {
	leader, err := zkClient.GetLeader()
	if err != nil {
		logger.Log.Errorf("Error getting leader")
		log.Fatalf("Error getting leader: %v", err)
		return false
	}
	return nodePath == LEADER_ELECTION_PATH+"/"+leader
}
