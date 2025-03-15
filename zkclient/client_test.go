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
	"rb-druid-indexer/logger"
	"testing"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockZkConn struct {
	mock.Mock
}

func (m *MockZkConn) Children(path string) ([]string, *zk.Stat, error) {
	args := m.Called(path)
	return args.Get(0).([]string), nil, args.Error(1)
}

func (m *MockZkConn) Get(path string) ([]byte, *zk.Stat, error) {
	args := m.Called(path)
	return args.Get(0).([]byte), nil, args.Error(1)
}

func TestMain(m *testing.M) {
	logger.Log = logrus.New()
	logger.Log.SetLevel(logrus.DebugLevel)

	m.Run()
}

func TestGetDruidRouterInfo_InvalidJson(t *testing.T) {
	mockConn := new(MockZkConn)

	routerPath := "/druid/router"
	mockChildren := []string{"router1"}

	mockConn.On("Children", routerPath).Return(mockChildren, nil)
	mockConn.On("Get", routerPath+"/router1").Return([]byte("invalid json"), nil)

	routers, err := GetDruidRouterInfo(mockConn, routerPath)

	assert.Error(t, err)
	assert.Nil(t, routers)
}
