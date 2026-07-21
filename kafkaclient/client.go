// Copyright (C) 2026 Eneo Tecnologia S.L.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

package kafkaclient

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// IsTopicEmpty checks if a Kafka topic has no messages available to read.
// It returns true if the sum of (LastOffset - FirstOffset) for all partitions is 0.
// If the topic does not exist or has no partitions, it returns true, nil.
func IsTopicEmpty(brokers []string, topic string) (bool, error) {
	if len(brokers) == 0 {
		return false, fmt.Errorf("no kafka brokers provided")
	}

	client := &kafka.Client{
		Addr:    kafka.TCP(brokers...),
		Timeout: 5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch metadata for the topic to discover partitions
	req := &kafka.MetadataRequest{
		Addr:   kafka.TCP(brokers...),
		Topics: []string{topic},
	}

	resp, err := client.Metadata(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to get metadata for topic %s: %w", topic, err)
	}

	// Find the topic in the response
	var topicMetadata kafka.Topic
	var found bool
	for _, t := range resp.Topics {
		if t.Name == topic {
			topicMetadata = t
			found = true
			break
		}
	}

	if !found {
		// Topic not found in metadata response, treat as empty
		return true, nil
	}

	if topicMetadata.Error != nil {
		// If topic does not exist, kafka returns UnknownTopicOrPartition error.
		// We treat non-existent topics as empty.
		return true, nil
	}

	if len(topicMetadata.Partitions) == 0 {
		return true, nil
	}


	// Query offsets for all partitions of this topic
	var offsetRequests []kafka.OffsetRequest
	for _, p := range topicMetadata.Partitions {
		offsetRequests = append(offsetRequests, kafka.FirstOffsetOf(p.ID))
		offsetRequests = append(offsetRequests, kafka.LastOffsetOf(p.ID))
	}

	offsetsResp, err := client.ListOffsets(ctx, &kafka.ListOffsetsRequest{
		Addr: kafka.TCP(brokers...),
		Topics: map[string][]kafka.OffsetRequest{
			topic: offsetRequests,
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to list offsets for topic %s: %w", topic, err)
	}

	partitionOffsets, exists := offsetsResp.Topics[topic]
	if !exists {
		return true, nil
	}

	totalMessages := int64(0)
	for _, po := range partitionOffsets {
		if po.Error != nil {
			return false, fmt.Errorf("failed to list offset for partition %d: %w", po.Partition, po.Error)
		}
		lag := po.LastOffset - po.FirstOffset
		if lag > 0 {
			totalMessages += lag
		}
	}

	return totalMessages == 0, nil
}

// GetLatestOffsets returns a map of partition ID to its latest offset (LastOffset) in Kafka.
// If the topic does not exist, it returns an empty map.
func GetLatestOffsets(brokers []string, topic string) (map[int]int64, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no kafka brokers provided")
	}

	client := &kafka.Client{
		Addr:    kafka.TCP(brokers...),
		Timeout: 5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch metadata for the topic to discover partitions
	req := &kafka.MetadataRequest{
		Addr:   kafka.TCP(brokers...),
		Topics: []string{topic},
	}

	resp, err := client.Metadata(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for topic %s: %w", topic, err)
	}

	// Find the topic in the response
	var topicMetadata kafka.Topic
	var found bool
	for _, t := range resp.Topics {
		if t.Name == topic {
			topicMetadata = t
			found = true
			break
		}
	}

	if !found || topicMetadata.Error != nil || len(topicMetadata.Partitions) == 0 {
		return map[int]int64{}, nil
	}

	// Query latest offsets for all partitions of this topic
	var offsetRequests []kafka.OffsetRequest
	for _, p := range topicMetadata.Partitions {
		offsetRequests = append(offsetRequests, kafka.LastOffsetOf(p.ID))
	}

	offsetsResp, err := client.ListOffsets(ctx, &kafka.ListOffsetsRequest{
		Addr: kafka.TCP(brokers...),
		Topics: map[string][]kafka.OffsetRequest{
			topic: offsetRequests,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list offsets for topic %s: %w", topic, err)
	}

	partitionOffsets, exists := offsetsResp.Topics[topic]
	if !exists {
		return map[int]int64{}, nil
	}

	latestOffsets := make(map[int]int64)
	for _, po := range partitionOffsets {
		if po.Error != nil {
			return nil, fmt.Errorf("failed to list offset for partition %d: %w", po.Partition, po.Error)
		}
		latestOffsets[po.Partition] = po.LastOffset
	}

	return latestOffsets, nil
}

