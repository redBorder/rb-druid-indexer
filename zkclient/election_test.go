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
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func TestExtractSeq(t *testing.T) {
	nodeName := "node-0000000010"
	seq, err := extractSeq(nodeName)
	if err != nil {
		t.Fatalf("extractSeq error: %v", err)
	}
	if seq != 10 {
		t.Errorf("expected sequence 10, got %d", seq)
	}

	invalidNodeName := "node-123"
	_, err = extractSeq(invalidNodeName)
	if err == nil {
		t.Errorf("expected error for invalid node name, got nil")
	}
}

func TestIsZKAlive(t *testing.T) {
	if IsZKAlive(nil) {
		t.Errorf("expected false for nil connection")
	}

	client, err := NewZKClient([]string{"127.0.0.1:2181"})
	if err != nil {
		t.Skipf("Skipping TestIsZKAlive: unable to connect to ZooKeeper: %v", err)
	}
	defer client.conn.Close()

	if !IsZKAlive(client.conn) {
		t.Errorf("expected true for a valid connection")
	}
}

func cleanupZK(t *testing.T, conn *zk.Conn) {
	children, _, err := conn.Children(LEADER_ELECTION_PATH)
	if err == nil {
		for _, child := range children {
			err = conn.Delete(LEADER_ELECTION_PATH+"/"+child, -1)
			if err != nil {
				t.Logf("cleanup: failed to delete child %s: %v", child, err)
			}
		}
		err = conn.Delete(LEADER_ELECTION_PATH, -1)
		if err != nil {
			t.Logf("cleanup: failed to delete leader election node: %v", err)
		}
	}
}

func TestLeaderElection(t *testing.T) {
	client1, err := NewZKClient([]string{"127.0.0.1:2181"})
	if err != nil {
		t.Skipf("Skipping TestLeaderElection: unable to connect to ZooKeeper: %v", err)
	}
	defer client1.conn.Close()

	client2, err := NewZKClient([]string{"127.0.0.1:2181"})
	if err != nil {
		t.Skipf("Skipping TestLeaderElection: unable to connect to ZooKeeper: %v", err)
	}
	defer client2.conn.Close()

	cleanupZK(t, client1.conn)

	nodePath1, err := client1.CreateLeaderNode()
	if err != nil {
		t.Fatalf("client1 CreateLeaderNode error: %v", err)
	}
	nodePath2, err := client2.CreateLeaderNode()
	if err != nil {
		t.Fatalf("client2 CreateLeaderNode error: %v", err)
	}

	time.Sleep(1 * time.Second)

	seq1, err := extractSeq(nodePath1)
	if err != nil {
		t.Fatalf("extractSeq error on nodePath1: %v", err)
	}
	seq2, err := extractSeq(nodePath2)
	if err != nil {
		t.Fatalf("extractSeq error on nodePath2: %v", err)
	}

	leader, err := client1.GetLeader()
	if err != nil {
		t.Fatalf("GetLeader error: %v", err)
	}

	var expectedLeader string
	if seq1 < seq2 {
		expectedLeader = nodePath1[len(LEADER_ELECTION_PATH)+1:]
	} else {
		expectedLeader = nodePath2[len(LEADER_ELECTION_PATH)+1:]
	}

	if leader != expectedLeader {
		t.Errorf("expected leader %s, got %s", expectedLeader, leader)
	}

	isLeader1 := client1.IsLeader(nodePath1)
	isLeader2 := client2.IsLeader(nodePath2)
	if seq1 < seq2 {
		if !isLeader1 {
			t.Errorf("client1 expected to be leader")
		}
		if isLeader2 {
			t.Errorf("client2 expected not to be leader")
		}
	} else {
		if !isLeader2 {
			t.Errorf("client2 expected to be leader")
		}
		if isLeader1 {
			t.Errorf("client1 expected not to be leader")
		}
	}

	cleanupZK(t, client1.conn)
}

func TestNewZKClientInvalid(t *testing.T) {
	_, err := NewZKClient([]string{"invalid:2181"})
	if err == nil {
		t.Errorf("expected error for invalid ZooKeeper server address")
	}
}
