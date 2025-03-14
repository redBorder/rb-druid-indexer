#!/bin/bash

KAFKA_CONTAINER="kafka"
KAFKA_BROKER="localhost:9092"

TOPIC="rb_flow_post"

FlowDimensions=("application_id_name" "building" "building_uuid" "campus" "campus_uuid"
"client_accounting_type" "client_auth_type" "client_fullname" "client_gender" "client_id"
"client_latlong" "client_loyality" "client_mac" "client_mac_vendor" "client_rssi"
"client_vip" "conversation" "coordinates_map" "darklist_category" "darklist_direction"
"darklist_score_name" "darklist_score" "deployment" "deployment_uuid" "direction" "dot11_protocol"
"dot11_status" "dst_map" "duration" "engine_id_name" "floor" "floor_uuid" "host" "host_l2_domain"
"http_social_media" "http_user_agent" "https_common_name" "interface_name" "ip_as_name" "ip_country_code"
"ip_protocol_version" "l4_proto" "lan_interface_description" "lan_interface_name" "lan_ip"
"lan_ip_is_malicious" "lan_ip_as_name" "lan_ip_country_code" "lan_ip_name" "lan_ip_net_name"
"lan_l4_port" "lan_name" "lan_vlan" "market" "market_uuid" "namespace" "namespace_uuid" "organization"
"organization_uuid" "product_name" "public_ip" "public_ip_is_malicious" "public_ip_mac" "referer"
"referer_l2" "scatterplot" "selector_name" "sensor_ip" "sensor_name" "sensor_uuid" "service_provider"
"service_provider_uuid" "src_map" "tcp_flags" "tos" "type" "url" "wan_interface_description"
"wan_interface_name" "wan_ip" "wan_ip_is_malicious" "wan_ip_as_name" "wan_ip_country_code" "wan_ip_map"
"wan_ip_net_name" "wan_l4_port" "wan_name" "wan_vlan" "wireless_id" "wireless_operator" "wireless_station"
"zone" "zone_uuid")

generate_random_string() {
  echo $(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 8)
}

generate_random_value() {
  local field=$1
  case $field in
    "client_gender" | "client_auth_type" | "client_accounting_type" | "client_vip" | "darklist_category" | "direction" | "dot11_status" | "ip_protocol_version" | "l4_proto" | "lan_ip_is_malicious" | "public_ip_is_malicious")
      echo $(shuf -e "male" "female" "other" "unknown" "public" "private" "secure" "open" "ipv4" "ipv6" "true" "false" -n 1)
      ;;
    "client_latlong" | "coordinates_map" | "dst_map" | "src_map" | "lan_ip" | "public_ip" | "wan_ip" | "ip_country_code" | "market_uuid" | "namespace_uuid" | "organization_uuid" | "zone_uuid")
      echo "$(shuf -i -1800-1800 -n 1),$(shuf -i -900-900 -n 1)" 
      ;;
    "duration" | "darklist_score" | "lan_l4_port" | "wan_l4_port" | "tcp_flags" | "tos")
      echo $(shuf -i 1-1000 -n 1)
      ;;
    "client_id" | "building_uuid" | "campus_uuid" | "deployment_uuid" | "floor_uuid" | "sensor_uuid" | "organization_uuid" | "market_uuid")
      echo "$(generate_random_string)-$(shuf -i 1-1000 -n 1)"
      ;;
    "url" | "http_user_agent" | "http_social_media" | "referer")
      echo "https://www.example.com/$(generate_random_string)"
      ;;
    "product_name" | "organization" | "building" | "campus" | "market" | "sensor_name" | "service_provider")
      echo "$(generate_random_string)-product"
      ;;
    *)
      echo "$(generate_random_string)"
      ;;
  esac
}

generate_random_json() {
  local timestamp=$(date +%s)
  local json="{\"timestamp\": \"$timestamp\""
  
  for field in "${FlowDimensions[@]}"; do
    local value=$(generate_random_value $field)
    json="$json, \"$field\": \"$value\""
  done
  
  json="$json}"
  echo $json
}

for i in {1..10}; do
  echo "Sending random JSON data to topic: $TOPIC"
  
  JSON_DATA=$(generate_random_json)
  
  echo $JSON_DATA | docker exec -i $KAFKA_CONTAINER kafka-console-producer --broker-list $KAFKA_BROKER --topic $TOPIC
  
  echo "Data sent to topic '$TOPIC': $JSON_DATA"
done
