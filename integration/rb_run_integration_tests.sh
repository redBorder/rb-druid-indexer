#!/bin/bash

DRUID_HOST="localhost:8888"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPT_PATH="$SCRIPT_DIR/rb_produce_syn_data.sh"

echo "Producing test data to Kafka..."
bash "$SCRIPT_PATH"

echo "Waiting for Druid supervisor 'rb_flow' to be active..."

MAX_RETRIES=30
RETRY_COUNT=0
SUPERVISOR_FOUND=false

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  SUPERVISOR_TASK_STATUS=$(curl -s -X GET "$DRUID_HOST/druid/indexer/v1/supervisor")
  
  if echo "$SUPERVISOR_TASK_STATUS" | grep -q "rb_flow"; then
    echo "Supervisor task 'rb_flow' found!"
    SUPERVISOR_FOUND=true
    break
  fi

  echo "Waiting for supervisor 'rb_flow'... (Attempt $((RETRY_COUNT+1))/$MAX_RETRIES). Current response: $SUPERVISOR_TASK_STATUS"
  sleep 5
  RETRY_COUNT=$((RETRY_COUNT+1))
done

if [ "$SUPERVISOR_FOUND" = false ]; then
  echo "Supervisor task 'rb_flow' was not registered in time. Exiting with failure."
  exit 1
fi

echo "Waiting for Druid SQL query on 'rb_flow' to return data..."

MAX_QUERY_RETRIES=24
QUERY_RETRY_COUNT=0
QUERY_SUCCESS=false

while [ $QUERY_RETRY_COUNT -lt $MAX_QUERY_RETRIES ]; do
  DRUID_TASK_DATA=$(curl -s -o response.json -w "%{http_code}" -X POST -H "Content-Type: application/json" \
    -d '{"query": "SELECT * FROM rb_flow"}' \
    "$DRUID_HOST/druid/v2/sql")

  HTTP_STATUS="${DRUID_TASK_DATA: -3}"

  if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "Druid query successful (HTTP status 200)."
    echo "Druid Query Result:"
    cat response.json
    QUERY_SUCCESS=true
    break
  fi

  echo "Druid query pending (HTTP $HTTP_STATUS)... (Attempt $((QUERY_RETRY_COUNT+1))/$MAX_QUERY_RETRIES)"
  sleep 5
  QUERY_RETRY_COUNT=$((QUERY_RETRY_COUNT+1))
done

if [ "$QUERY_SUCCESS" = false ]; then
  echo "Druid query failed after retries. Last HTTP Status: $HTTP_STATUS"
  echo "Last Druid Query Response:"
  cat response.json
  exit 1
fi

echo "Integration test completed successfully!"
exit 0

