#!/bin/bash

DRUID_HOST="localhost:8888"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPT_PATH="$SCRIPT_DIR/rb_produce_syn_data.sh"
bash "$SCRIPT_PATH"

echo "Checking for all Druid supervisors..."

SUPERVISOR_TASK_STATUS=$(curl -s -X GET "$DRUID_HOST/druid/indexer/v1/supervisor")

echo "Supervisor Task Status:"
echo "$SUPERVISOR_TASK_STATUS"

if [[ "$SUPERVISOR_TASK_STATUS" == "[]" ]]; then
  echo "No supervisors found. Exiting with failure."
  exit 1
else
  echo "Supervisor tasks found. Exiting with success."
  exit 0
fi

DRUID_TASK_DATA=$(curl -s -o response.json -w "%{http_code}" -X POST -H "Content-Type: application/json" \
  -d '{"query": "SELECT * FROM rb_flow"}' \
  "$DRUID_HOST/druid/v2/sql")

HTTP_STATUS="${DRUID_TASK_DATA: -3}"

echo "HTTP Status: $HTTP_STATUS"
echo "Druid Query Result:"
cat response.json

if [ "$HTTP_STATUS" -eq 200 ]; then
  echo "Druid query successful. Exiting with success."
  exit 0
else
  echo "Druid query failed with status code $HTTP_STATUS. Exiting with failure."
  exit 1
fi
