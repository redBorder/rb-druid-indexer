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

package logger

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/natefinch/lumberjack.v2"
)

var LogFiles = map[string][]logrus.Level{
	"/var/log/rb-druid-indexer/info.log":  {logrus.InfoLevel, logrus.WarnLevel},
	"/var/log/rb-druid-indexer/error.log": {logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	"/var/log/rb-druid-indexer/debug.log": {logrus.DebugLevel},
}

func TestInitLogger(t *testing.T) {
	err := os.MkdirAll("/var/log/rb-druid-indexer", 0755)
	assert.NoError(t, err)

	InitLogger()

	assert.NotNil(t, Log)
	assert.Equal(t, logrus.DebugLevel, Log.GetLevel())

	for file := range LogFiles {
		_, err := os.Stat(file)
		assert.NoError(t, err)
	}
}

func TestLogFileHook(t *testing.T) {
	tempFile := "/var/log/rb-druid-indexer/test.log"
	ensureLogFileExists(tempFile)

	hook := NewLogFileHook(&lumberjack.Logger{
		Filename:   tempFile,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}, logrus.InfoLevel)

	entry := &logrus.Entry{
		Logger:  Log,
		Level:   logrus.InfoLevel,
		Message: "Test log message",
	}

	err := hook.Fire(entry)
	assert.NoError(t, err)

	content, err := os.ReadFile(tempFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "Test log message")

	os.Remove(tempFile)
}

func TestLogFileHookLevels(t *testing.T) {
	hook := NewLogFileHook(&lumberjack.Logger{
		Filename:   "/var/log/rb-druid-indexer/test.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}, logrus.InfoLevel, logrus.WarnLevel)

	levels := hook.Levels()
	assert.Contains(t, levels, logrus.InfoLevel)
	assert.Contains(t, levels, logrus.WarnLevel)
	assert.NotContains(t, levels, logrus.ErrorLevel)
}
