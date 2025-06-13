#!/bin/bash

KAFKA_CONTAINER="kafka"

TOPICS=("rb_monitor" "rb_flow_post")

PARTITIONS=1
REPLICATION_FACTOR=1

# Wait for Kafka to be ready
echo "Waiting for Kafka to be ready..."
until docker exec $KAFKA_CONTAINER kafka-topics --bootstrap-server localhost:9092 --list &>/dev/null; do
  echo "Kafka not ready, waiting..."
  sleep 5
done

for TOPIC in "${TOPICS[@]}"; do
  echo "Creating topic: $TOPIC"

  docker exec $KAFKA_CONTAINER kafka-topics --create --topic "$TOPIC" --bootstrap-server localhost:9092 --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR
  
  if [ $? -eq 0 ]; then
    echo "Topic '$TOPIC' created successfully!"
  else
    echo "Failed to create topic '$TOPIC'."
  fi
done
