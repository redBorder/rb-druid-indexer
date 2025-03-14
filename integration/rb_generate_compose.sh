#!/bin/bash

# Function to generate a random string
generate_random_string() {
  echo $(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 8)
}

# Generate random strings for each volume
METADATA_VOLUME=$(generate_random_string)
INDEXER_VOLUME=$(generate_random_string)
HISTORICAL_VOLUME=$(generate_random_string)
BROKER_VOLUME=$(generate_random_string)
COORDINATOR_VOLUME=$(generate_random_string)
ROUTER_VOLUME=$(generate_random_string)
DRUID_SHARED_VOLUME=$(generate_random_string)
KAFKA_DATA_VOLUME=$(generate_random_string)
ZK_DATA_VOLUME=$(generate_random_string)

# Print the docker-compose file with updated random volume names
cat << EOF
version: "2.2"

volumes:
  ${METADATA_VOLUME}: {}
  ${INDEXER_VOLUME}: {}
  ${HISTORICAL_VOLUME}: {}
  ${BROKER_VOLUME}: {}
  ${COORDINATOR_VOLUME}: {}
  ${ROUTER_VOLUME}: {}
  ${DRUID_SHARED_VOLUME}: {}
  ${KAFKA_DATA_VOLUME}: {}
  ${ZK_DATA_VOLUME}: {}

services:
  postgres:
    container_name: postgres
    image: postgres:latest
    ports:
      - "5432:5432"
    volumes:
      - ${METADATA_VOLUME}:/var/lib/postgresql/data
    environment:
      - POSTGRES_PASSWORD=FoolishPassword
      - POSTGRES_USER=druid
      - POSTGRES_DB=druid

  zookeeper:
    container_name: zookeeper
    image: zookeeper:3.8.4
    ports:
      - "2181:2181"
    environment:
      - ZOO_MY_ID=1
    volumes:
      - ${ZK_DATA_VOLUME}:/data
      - ${ZK_DATA_VOLUME}:/datalog

  coordinator:
    image: apache/druid:32.0.0
    container_name: coordinator
    volumes:
      - ${DRUID_SHARED_VOLUME}:/opt/shared
      - ${COORDINATOR_VOLUME}:/opt/druid/var
    depends_on:
      - zookeeper
      - postgres
    ports:
      - "8081:8081"
    command:
      - coordinator
    env_file:
      - environment

  broker:
    image: apache/druid:32.0.0
    container_name: broker
    volumes:
      - ${BROKER_VOLUME}:/opt/druid/var
    depends_on:
      - zookeeper
      - postgres
      - coordinator
    ports:
      - "8082:8082"
    command:
      - broker
    env_file:
      - environment

  historical:
    image: apache/druid:32.0.0
    container_name: historical
    volumes:
      - ${DRUID_SHARED_VOLUME}:/opt/shared
      - ${HISTORICAL_VOLUME}:/opt/druid/var
    depends_on: 
      - zookeeper
      - postgres
      - coordinator
    ports:
      - "8083:8083"
    command:
      - historical
    env_file:
      - environment

  indexer:
    image: apache/druid:32.0.0
    container_name: indexer
    volumes:
      - ${DRUID_SHARED_VOLUME}:/opt/shared
      - ${INDEXER_VOLUME}:/opt/druid/var
    depends_on: 
      - zookeeper
      - postgres
      - coordinator
    ports:
      - "8091:8091"
      - "8100-8105:8100-8105"
    command:
      - indexer
    env_file:
      - environment

  router:
    image: apache/druid:32.0.0
    container_name: router
    volumes:
      - ${ROUTER_VOLUME}:/opt/druid/var
    depends_on:
      - zookeeper
      - postgres
      - coordinator
    ports:
      - "8888:8888"
    command:
      - router
    env_file:
      - environment

  kafka:
    image: confluentinc/cp-kafka:latest
    container_name: kafka
    ports:
      - "9092:9092"
    environment:
      - KAFKA_BROKER_ID=1
      - KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181
      - KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092
      - KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092
      - KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1
    volumes:
      - ${KAFKA_DATA_VOLUME}:/var/lib/kafka/data
    depends_on:
      - zookeeper
      - postgres
      - coordinator
      - broker
      - historical
      - indexer
      - router
EOF