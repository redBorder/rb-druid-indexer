#!/bin/bash

KAFKA_CONTAINER="kafka"

TOPICS=("rb_monitor_post" "rb_flow_post")

PARTITIONS=1
REPLICATION_FACTOR=1

for TOPIC in "${TOPICS[@]}"; do
    echo "Creating topic: $TOPIC"
    
    docker exec $KAFKA_CONTAINER kafka-topics --create --topic "$TOPIC" --bootstrap-server localhost:9092 --partitions $PARTITIONS --replication-factor $REPLICATION_FACTOR

    if [ $? -eq 0 ]; then
      echo "Topic '$TOPIC' created successfully!"
    else
      echo "Failed to create topic '$TOPIC'."
    fi
done
