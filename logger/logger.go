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

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

func ensureLogFileExists(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logrus.Fatalf("Failed to create log file %s: %v", filename, err)
		}
		f.Close()
	}
}

func InitLogger() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	Log.SetLevel(logrus.DebugLevel)

	ensureLogFileExists("/var/log/rb-druid-indexer/info.log")
	ensureLogFileExists("/var/log/rb-druid-indexer/error.log")
	ensureLogFileExists("/var/log/rb-druid-indexer/debug.log")

	infoLog := &lumberjack.Logger{
		Filename:   "/var/log/rb-druid-indexer/info.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	errorLog := &lumberjack.Logger{
		Filename:   "/var/log/rb-druid-indexer/error.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	debugLog := &lumberjack.Logger{
		Filename:   "/var/log/rb-druid-indexer/debug.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	Log.AddHook(NewLogFileHook(infoLog, logrus.InfoLevel, logrus.WarnLevel))
	Log.AddHook(NewLogFileHook(errorLog, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel))
	Log.AddHook(NewLogFileHook(debugLog, logrus.DebugLevel))

	Log.SetOutput(os.Stdout)
}

type LogFileHook struct {
	Writer    *lumberjack.Logger
	LogLevels []logrus.Level
}

func NewLogFileHook(writer *lumberjack.Logger, levels ...logrus.Level) *LogFileHook {
	return &LogFileHook{Writer: writer, LogLevels: levels}
}

func (hook *LogFileHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func (hook *LogFileHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}
