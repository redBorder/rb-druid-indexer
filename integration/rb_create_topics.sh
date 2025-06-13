#!/bin/bash

KAFKA_CONTAINER="kafka"

TOPICS=("rb_monitor" "rb_flow_post")

PARTITIONS=1
REPLICATION_FACTOR=1

if ! docker ps --filter "name=$KAFKA_CONTAINER" --filter "status=running" | grep -q $KAFKA_CONTAINER; then
    echo "Error: Container '$KAFKA_CONTAINER' is not running"
    echo "Available containers:"
    docker ps
    exit 1
fi

echo "Kafka container is running. Proceeding with topic creation..."

for TOPIC in "${TOPICS[@]}"; do
    # Wait for Kafka to be ready
    echo "Waiting for Kafka to be ready..."
    RETRY_COUNT=0
    MAX_RETRIES=30
    
    until docker exec $KAFKA_CONTAINER kafka-topics --bootstrap-server localhost:9092 --list &>/dev/null; do
        echo "Kafka not ready, waiting... (attempt $((++RETRY_COUNT))/$MAX_RETRIES)"
        if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
            echo "Timeout waiting for Kafka to be ready"
            exit 1
        fi
        sleep 5
    done

    echo "Creating topic: $TOPIC"
    
    if docker exec $KAFKA_CONTAINER kafka-topics --create --topic "$TOPIC" --bootstrap-server localhost:9092 --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR; then
        echo "Topic '$TOPIC' created successfully!"
    else
        echo "Failed to create topic '$TOPIC'."
        exit 1
    fi
done
