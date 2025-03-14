#!/bin/bash

DRUID_HOST="localhost:8888"

echo "Sending test message to Kafka topic: $KAFKA_TOPIC"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPT_PATH="$SCRIPT_DIR/rb_produce_syn_data.sh"
bash "$SCRIPT_PATH"

if [ $? -eq 0 ]; then
  echo "Message successfully sent to Kafka!"
else
  echo "Failed to send message to Kafka!"
  exit 1
fi

echo "Checking if Druid tasks are running..."

DRUID_TASK_STATUS=$(curl -s -X POST -H "Content-Type: application/json" -d '{"query":"SELECT * FROM druid_task WHERE status = \'RUNNING\'"}' "$DRUID_HOST/druid/v2/sql")

if echo "$DRUID_TASK_STATUS" | grep -q "RUNNING"; then
  echo "Druid task is running!"
else
  echo "No running tasks found in Druid!"
  exit 1
fi
