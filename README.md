<p align="center">
    <img src="./assets/image.png" align="center" width="30%">
</p>
<p align="center"><h1 align="center">RB-DRUID-INDEXER</h1></p>
<p align="center">
	<em><code>Simple distributed druid-indexer task manager for kafka ingestion </code></em>
</p>
<p align="center">
	<img src="https://img.shields.io/github/license/redBorder/rb-druid-indexer?style=default&logo=opensourceinitiative&logoColor=white&color=ff2400" alt="license">
	<img src="https://img.shields.io/github/last-commit/redBorder/rb-druid-indexer?style=default&logo=git&logoColor=white&color=ff2400" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/redBorder/rb-druid-indexer?style=default&color=ff2400" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/redBorder/rb-druid-indexer?style=default&color=ff2400" alt="repo-language-count">
	<img src="https://img.shields.io/github/actions/workflow/status/redBorder/rb-druid-indexer/go.yml?branch=master&label=Go%20Build" alt="repo-language-count">
	<img src="https://img.shields.io/github/actions/workflow/status/redBorder/rb-druid-indexer/integration.yml?branch=master&label=Unit%20Tests%20and%20Integration%20Tests" alt="repo-language-count">
</p>
<p align="center"><!-- default option, no dependency badges. -->
</p>
<p align="center">
	<!-- default option, no dependency badges. -->
</p>
<br>

##  Table of Contents

- [Table of Contents](#table-of-contents)
- [Overview](#overview)
- [How rb-druid-indexer fits in our new indexing system or yours](#how-rb-druid-indexer-fits-in-our-new-indexing-system-or-yours)
- [Features](#features)
- [Configuration](#configuration)
- [zookeeper\_servers](#zookeeper_servers)
- [discovery\_path](#discovery_path)
- [tasks](#tasks)
  - [task\_name](#task_name)
  - [spec](#spec)
  - [feed](#feed)
  - [kafka\_brokers](#kafka_brokers)
  - [dimensions](#dimensions)
  - [dimensions\_exclusions](#dimensions_exclusions)
  - [metrics](#metrics)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Author](#author)

---

##  Overview

`rb-druid-indexer` is a cluster-compatible service designed to manage the indexing of Kafka data streams into Druid. It handles task announcements, generates configuration specification files, and submits tasks to the Druid Supervisor.

---

##  How rb-druid-indexer fits in our new indexing system or yours
<p align="center">
    <img src="./assets/arch_img_new.png" align="center" width="120%">
</p>
In the old system, Druid indexing relied on ShardSpec with druid-realtime, where tasks were split into multiple shards across nodes for parallel processing. This approach, defined in static realtime spec files & hard-to-deploy nodes introduced complexity in shard management and scalability. In contrast, the new system uses the rb-druid-indexer, which simplifies the process by submitting single tasks without shard splitting to druid router wich automatically distribute task in druid indexer nodes and we leave overlord to manage balancing.

You can notice this fast with this diagram
<p align="center">
    <img src="./assets/old_vs_new.png" align="center" width="120%">
</p>


##  Features

- Multi Druid Router compatible
- Auto Finder for Druid Routers
- Cluster compatible & FailOver support using ZooKeeper 
- Automatic task managment and load balancing when submiting / deleting tasks

---

## Configuration

The configuration for `rb-druid-indexer` is defined in a YAML file that is generated with the [druid-indexer cookbook](https://github.com/redBorder/cookbook-druid-indexer) in **/etc/rb-druid-indexer/config.yml**. It includes settings for both Zookeeper, the tasks that should be executed, the dimensions and their metrics. Below is an example configuration file:

```yaml
zookeeper_servers:
  - "rb-malvarez1.node:2181"
  - "rb-malvarez3.node:2181"
  - "rb-malvarez2.node:2181"
discovery_path: "/druid/discovery/druid:router"

tasks:
  - task_name: "rb_monitor"
    feed: "rb_monitor"
    spec: "rb_monitor"
    kafka_brokers:
      - "rb-malvarez1.node:9092"
    dimensions_exclusions:
        - "unit"
        - "type"
        - "value"
    metrics:
        - type: count
          name: events
        - type: doubleSum
          name: sum_value
          fieldName: value
        - type: doubleMax
          name: max_value
          fieldName: value
        - type: doubleMin
          name: min_value
          fieldName: value
  - task_name: "rb_flow"
    feed: "rb_flow_post"
    spec: "rb_flow"
    kafka_brokers:
      - "rb-malvarez1.node:9092"
      - "rb-malvarez3.node:9092"
      - "rb-malvarez2.node:9092"
    dimensions:
      - "application_id_name"
      - "building"
      - "building_uuid"
      - "campus"
      - "campus_uuid"
      - "client_accounting_type"
      - "client_auth_type"
      - "client_fullname"
      - "client_gender"
      - "client_id"
      - "client_latlong"
      - "client_loyality"
      - "client_mac"
      - "client_mac_vendor"
      - "client_rssi"
      - "client_vip"
      - "conversation"
      - "coordinates_map"
      - "deployment"
      - "deployment_uuid"
      - "direction"
    dimensions_exclusions:
      - "bytes"
      - "pkts"
      - "flow_end_reason"
      - "first_switched"
      - "wan_ip_name"
    metrics:
      - type: count
        name: events
      - type: longSum
        name: sum_bytes
        fieldName: bytes
      - type: longSum
        name: sum_pkts
        fieldName: pkts
      - type: longSum
        name: sum_rssi
        fieldName: client_rssi_num
      - type: hyperUnique
        name: clients
        fieldName: client_mac
      - type: hyperUnique
        name: wireless_stations
        fieldName: wireless_station
...
```


## zookeeper_servers
- **Description**: A list of Zookeeper servers used for leadership checks and coordination.
- **Type**: Array of strings.
- **Example**: 
    - `"127.0.0.1:2181"`
    - `"127.0.0.2:2181"`
## discovery_path
- **Description**: (optional field) ZooKeeper path where Druid routers are announced
- **Type**: String.
- **Example**: 
    - `"/druid/discovery/druid:router"`
## tasks
- **Description**: A list of tasks to be managed by the indexer. Each task contains the following attributes:

### task_name
- **Description**: The name of the task. This is used to identify the task in the system.
- **Type**: String.
- **Example**: 
    - `"rb_monitor"`
    - `"rb_flow"`

### spec
- **Description**: The spec file name associated with the task (for realtime configuration)
- **Type**: String.
- **Example**: 
    - `"rb_flow"`

### feed
- **Description**: The name of the Kafka feed associated with the task. This specifies which feed to listen to.
- **Type**: String.
- **Example**: 
    - `"rb_monitor"`
    - `"rb_flow_post"`

### kafka_brokers
- **Description**: The list of kafka brokers for supervisor
- **Type**: Array.
- **Example**: 
    kafka_brokers:
      - `"kafka.service:9092"`

### dimensions
- **Description**: The list of dimensions for the datasource
- **Type**: Array.
- **Example**: 
    dimensions:
      - `"lan_ip"`

### dimensions_exclusions
- **Description**: The list of dimensions that will be excluded for the datasource
- **Type**: Array.
- **Example**: 
    dimensions_exclusions:
      - `"wan_ip"`

### metrics
- **Description**: The list of metrics that will be used for the datasource
- **Type**: Array of objects.
- **Example**: 
    metrics:
      - type: count
        name: events
      - type: longSum
        name: sum_bytes
        fieldName: bytes

---

##  Project Structure

```sh
rb-druid-indexer
├── LICENSE
├── Makefile
├── README.md
├── assets
│   ├── arch_img_new.png
│   ├── image.png
│   └── old_vs_new.png
├── config
│   ├── config.go
│   └── config_test.go
├── druid
│   ├── config
│   │   ├── config.go
│   │   └── config_test.go
│   ├── realtime.go
│   ├── realtime_test.go
│   ├── router.go
│   └── router_test.go
├── example_config.yml
├── go.mod
├── go.sum
├── integration
│   ├── config.yml
│   ├── docker-compose.yml
│   ├── environment
│   ├── rb_create_topics.sh
│   ├── rb_generate_compose.sh
│   ├── rb_produce_syn_data.sh
│   └── rb_run_integration_tests.sh
├── logger
│   ├── logger.go
│   └── logger_test.go
├── main.go
├── main_test.go
├── packaging
│   └── rpm
│       ├── Makefile
│       ├── rb-druid-indexer.service
│       └── rb-druid-indexer.spec
├── rb-druid-indexer
└── zkclient
    ├── client.go
    ├── client_test.go
    ├── election.go
    ├── election_test.go
    ├── task_announcer.go
    └── task_announcer_test.go
```

---
##  Getting Started

###  Prerequisites

Before getting started with rb-druid-indexer, ensure your runtime environment meets the following requirements:

- **Programming Language:** Go
- **Package Manager:** Go modules


###  Installation

Install rb-druid-indexer using one of the following methods:

**Build from source:**

1. Clone the rb-druid-indexer repository:
```sh
❯ git clone https://github.com/redBorder/rb-druid-indexer
```

2. Navigate to the project directory:
```sh
❯ cd rb-druid-indexer
```

3. Install the project dependencies:


**Using `go modules`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Go-00ADD8.svg?style={badge_style}&logo=go&logoColor=white" />](https://golang.org/)

```sh
❯ go build
```

###  Usage
Run rb-druid-indexer using the following command:
**Using `go modules`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Go-00ADD8.svg?style={badge_style}&logo=go&logoColor=white" />](https://golang.org/)

```sh
❯ ./rb-druid-indexer --config example_config.yml
```


---

##  Contributing

- **💬 [Join the Discussions](https://github.com/redBorder/rb-druid-indexer/discussions)**: Share your insights, provide feedback, or ask questions.
- **🐛 [Report Issues](https://github.com/redBorder/rb-druid-indexer/issues)**: Submit bugs found or log feature requests for the `rb-druid-indexer` project.
- **💡 [Submit Pull Requests](https://github.com/redBorder/rb-druid-indexer/blob/main/CONTRIBUTING.md)**: Review open PRs, and submit your own PRs.

<details closed>
<summary>Contributing Guidelines</summary>

1. **Fork the Repository**: Start by forking the project repository to your github account.
2. **Clone Locally**: Clone the forked repository to your local machine using a git client.
   ```sh
   git clone https://github.com/redBorder/rb-druid-indexer
   ```
3. **Create a New Branch**: Always work on a new branch, giving it a descriptive name.
   ```sh
   git checkout -b new-feature-x
   ```
4. **Make Your Changes**: Develop and test your changes locally.
5. **Commit Your Changes**: Commit with a clear message describing your updates.
   ```sh
   git commit -m 'Implemented new feature x.'
   ```
6. **Push to github**: Push the changes to your forked repository.
   ```sh
   git push origin new-feature-x
   ```
7. **Submit a Pull Request**: Create a PR against the original project repository. Clearly describe the changes and their motivations.
8. **Review**: Once your PR is reviewed and approved, it will be merged into the main branch. Congratulations on your contribution!
</details>

<details closed>
<summary>Contributor Graph</summary>
<br>
<p align="left">
   <a href="https://github.com{/redBorder/rb-druid-indexer/}graphs/contributors">
      <img src="https://contrib.rocks/image?repo=redBorder/rb-druid-indexer">
   </a>
</p>
</details>

---

##  License

This project is protected under the [AGPL-3.0](https://www.gnu.org/licenses/agpl-3.0.txt) License. For more details, refer to the [LICENSE](https://www.gnu.org/licenses/agpl-3.0.txt) file.

---

##  Author

This project is developed for redBorder and the OS community by Miguel Álvarez <malvarez@redborder.com>

---
