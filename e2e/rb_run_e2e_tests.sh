#!/bin/bash

DRUID_HOST="localhost:8888"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPT_PATH="$SCRIPT_DIR/rb_produce_syn_data.sh"
bash "$SCRIPT_PATH"

echo "Checking if Druid tasks are running..."

DRUID_TASK_STATUS=$(curl -s -X POST -H "Content-Type: application/json" \
  -d '{"query": "SELECT * FROM druid_task"}' \
  "$DRUID_HOST/druid/v2/sql")

echo "Druid Query Result:"
echo "$DRUID_TASK_STATUS"

if echo "$DRUID_TASK_STATUS" | grep -q "RUNNING"; then
  echo "Druid task is running!"
else
  echo "No running tasks found in Druid!"
  exit 1
fi
