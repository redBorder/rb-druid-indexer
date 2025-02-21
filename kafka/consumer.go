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

package rbkafka

import (
	"context"
	"fmt"
	"rb-druid-indexer/logger"
	"sync"

	"github.com/segmentio/kafka-go"
)

var topicFlags = make(map[string]bool)
var mu sync.Mutex

func NewKafkaConsumer(brokers, topic, group string) (*kafka.Reader, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokers},
		Topic:       topic,
		GroupID:     group,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: kafka.LastOffset,
	})

	if reader == nil {
		logger.Log.Errorf("Failed to create consumer for topic %s", topic)
		return nil, fmt.Errorf("failed to create consumer for topic %s", topic)
	}

	return reader, nil
}

func ReadMessagesForTopic(reader *kafka.Reader, topic string) {
	mu.Lock()
	topicFlags[topic] = false
	mu.Unlock()

	for {
		_, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}

		mu.Lock()
		topicFlags[topic] = true
		mu.Unlock()
	}

	mu.Lock()
	topicFlags[topic] = false
	mu.Unlock()

	logger.Log.Infof("No more messages for topic %s", topic)
}

func StartConsumer(topicBrokers map[string]string, group string) {
	for topic, broker := range topicBrokers {
		reader, err := NewKafkaConsumer(broker, topic, group)
		if err != nil {
			logger.Log.Errorf("Error creating consumer for topic %s on broker %s: %v", topic, broker, err)
			continue
		}

		go ReadMessagesForTopic(reader, topic)
	}
}

func CheckFlag(topic string) bool {
	mu.Lock()
	defer mu.Unlock()

	flag, exists := topicFlags[topic]
	if !exists {
		logger.Log.Infof("Topic %s not found", topic)
		return false
	}
	return flag
}

func SetFalseFlag(topic string) {
	mu.Lock()
	defer mu.Unlock()
	_, exists := topicFlags[topic]
	if exists {
		logger.Log.Infof("Topic %s found locking", topic)
		topicFlags[topic] = false
	}
}
