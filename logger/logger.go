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
		if f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			f.Close()
		} else {
			logrus.Fatalf("Failed to create log file %s: %v", filename, err)
		}
	}
}

func InitLogger() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	Log.SetLevel(logrus.DebugLevel)

	logFiles := map[string][]logrus.Level{
		"/var/log/rb-druid-indexer/info.log":  {logrus.InfoLevel, logrus.WarnLevel},
		"/var/log/rb-druid-indexer/error.log": {logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
		"/var/log/rb-druid-indexer/debug.log": {logrus.DebugLevel},
	}

	for file, levels := range logFiles {
		ensureLogFileExists(file)
		Log.AddHook(NewLogFileHook(&lumberjack.Logger{
			Filename:   file,
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		}, levels...))
	}

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