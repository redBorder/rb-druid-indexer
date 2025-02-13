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
</p>
<p align="center"><!-- default option, no dependency badges. -->
</p>
<p align="center">
	<!-- default option, no dependency badges. -->
</p>
<br>

##  Table of Contents

- [ Overview](#-overview)
- [ Features](#-features)
- [ Project Structure](#-project-structure)
  - [ Project Index](#-project-index)
- [ Getting Started](#-getting-started)
  - [ Prerequisites](#-prerequisites)
  - [ Installation](#-installation)
  - [ Usage](#-usage)
  - [ Testing](#-testing)
- [ Project Roadmap](#-project-roadmap)
- [ Contributing](#-contributing)
- [ License](#-license)
- [ Acknowledgments](#-acknowledgments)

---

##  Overview

<code>`rb-druid-indexer` is a cluster-compatible service designed to manage the indexing of Kafka data streams into Druid. It handles task announcements, generates configuration specification files, and submits tasks to the Druid Supervisor.
</code>

---

##  Features

<code>❯ Cluster compatible using ZooKeeper </code>
<code>❯ Automatic task ingestion & specfile configuration </code>
<code>❯ FailOver support for long-term ingestion </code>

---

## Configuration

<code>The configuration for `rb-druid-indexer` is defined in a YAML file and includes settings for both Zookeeper and the tasks that should be executed. Below is an example configuration file:</code>

```yaml
zookeeper_servers:
  - "127.0.0.1:2181"

tasks:
  - task_name: "rb_monitor"
    namespace: ""
    feed: "rb_monitor"
    kafka_host: "kafka.service:9092"
  - task_name: "rb_flow"
    namespace: ""
    feed: "rb_flow_post"
    kafka_host: "kafka.service:9092"
```


## zookeeper_servers
- **Description**: A list of Zookeeper servers used for leadership checks and coordination.
- **Type**: Array of strings.
- **Example**: 
    - `"127.0.0.1:2181"`

## tasks
- **Description**: A list of tasks to be managed by the indexer. Each task contains the following attributes:

### task_name
- **Description**: The name of the task. This is used to identify the task in the system.
- **Type**: String.
- **Example**: 
    - `"rb_monitor"`
    - `"rb_flow"`

### namespace (optional)
- **Description**: The namespace associated with the task. This can be left empty if not needed.
- **Type**: String.
- **Example**: 
    - `""` (empty)

### feed
- **Description**: The name of the Kafka feed associated with the task. This specifies which feed to listen to.
- **Type**: String.
- **Example**: 
    - `"rb_monitor"`
    - `"rb_flow_post"`

### kafka_host
- **Description**: The host and port for the Kafka service where the feed is being published.
- **Type**: String.
- **Example**: 
    - `"kafka.service:9092"`

<code>Every dataSource is managed in ```/druid/datasources/${datasource}.go for example</code>

```go
package datasources

import druidrouter "rb-druid-indexer/druid"

var FlowMetrics = []druidrouter.Metrics{
	{Type: "count", Name: "events"},
	{Type: "longSum", Name: "sum_bytes", FieldName: "bytes"},
	{Type: "longSum", Name: "sum_pkts", FieldName: "pkts"},
	{Type: "longSum", Name: "sum_rssi", FieldName: "client_rssi_num"},
	{Type: "hyperUnique", Name: "clients", FieldName: "client_mac"},
	{Type: "hyperUnique", Name: "wireless_stations", FieldName: "wireless_station"},
	{Type: "longSum", Name: "sum_dl_score", FieldName: "darklist_score"},
}

var FlowDimensions = []string{
	"application_id_name", "building", "building_uuid", "campus", "campus_uuid",
	"client_accounting_type", "client_auth_type", "client_fullname", "client_gender",
	"client_id", "client_latlong", "client_loyality", "client_mac", "client_mac_vendor",
	"client_rssi", "client_vip", "conversation", "coordinates_map", "darklist_category",
	"darklist_direction", "darklist_score_name", "darklist_score", "deployment",
	"deployment_uuid", "direction", "dot11_protocol", "dot11_status", "dst_map", "duration",
	"engine_id_name", "floor", "floor_uuid", "host", "host_l2_domain", "http_social_media",
	"http_user_agent", "https_common_name", "interface_name", "ip_as_name", "ip_country_code",
	"ip_protocol_version", "l4_proto", "lan_interface_description", "lan_interface_name",
	"lan_ip", "lan_ip_is_malicious", "lan_ip_as_name", "lan_ip_country_code", "lan_ip_name",
	"lan_ip_net_name", "lan_l4_port", "lan_name", "lan_vlan", "market", "market_uuid",
	"namespace", "namespace_uuid", "organization", "organization_uuid", "product_name",
	"public_ip", "public_ip_is_malicious", "public_ip_mac", "referer", "referer_l2",
	"scatterplot", "selector_name", "sensor_ip", "sensor_name", "sensor_uuid", "service_provider",
	"service_provider_uuid", "src_map", "tcp_flags", "tos", "type", "url", "wan_interface_description",
	"wan_interface_name", "wan_ip", "wan_ip_is_malicious", "wan_ip_as_name", "wan_ip_country_code",
	"wan_ip_map", "wan_ip_net_name", "wan_l4_port", "wan_name", "wan_vlan", "wireless_id",
	"wireless_operator", "wireless_station", "zone", "zone_uuid",
}

const FlowDataSource = "rb_flow"
```

<code>and later published in the `config.go` file in `/druid/datasources/config.go`</code>

```go
var Configs = map[string]DataSourceConfig{
	"rb_flow": {
		DataSource: FlowDataSource,
		Metrics:    FlowMetrics,
		Dimensions: FlowDimensions,
	},
	"rb_monitor": {
		DataSource: MonitorDataSource,
		Metrics:    MonitorMetrics,
		Dimensions: MonitorDimensions,
	},
}
```

<code>So if you want to add your own you have to make a copy of any datasource and include in the config.go of datasource for later call it with your `config.yml`</code>
---

##  Project Structure

```sh
└── rb-druid-indexer/
    ├── LICENSE
    ├── config
    │   └── config.go
    ├── druid
    │   ├── datasources
    │   ├── realtime.go
    │   └── router.go
    ├── example_config.yml
    ├── go.mod
    ├── go.sum
    ├── main.go
    └── zkclient
        ├── client.go
        ├── election.go
        └── task_announcer.go
```


###  Project Index
<details open>
	<summary><b><code>RB-DRUID-INDEXER/</code></b></summary>
	<details> <!-- __root__ Submodule -->
		<summary><b>__root__</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/main.go'>main.go</a></b></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/go.mod'>go.mod</a></b></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/go.sum'>go.sum</a></b></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/example_config.yml'>example_config.yml</a></b></td>
			</tr>
			</table>
		</blockquote>
	</details>
	<details> <!-- config Submodule -->
		<summary><b>config</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/config/config.go'>config.go</a></b></td>
			</tr>
			</table>
		</blockquote>
	</details>
	<details> <!-- zkclient Submodule -->
		<summary><b>zkclient</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/zkclient/election.go'>election.go</a></b></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/zkclient/client.go'>client.go</a></b></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/zkclient/task_announcer.go'>task_announcer.go</a></b></td>
			</tr>
			</table>
		</blockquote>
	</details>
	<details> <!-- druid Submodule -->
		<summary><b>druid</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/druid/realtime.go'>realtime.go</a></b></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/druid/router.go'>router.go</a></b></td>
			</tr>
			</table>
			<details>
				<summary><b>datasources</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/druid/datasources/location.go'>location.go</a></b></td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/druid/datasources/config.go'>config.go</a></b></td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/druid/datasources/event.go'>event.go</a></b></td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/druid/datasources/wireless.go'>wireless.go</a></b></td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/druid/datasources/monitor.go'>monitor.go</a></b></td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/druid/datasources/state.go'>state.go</a></b></td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/redBorder/rb-druid-indexer/blob/master/druid/datasources/flow.go'>flow.go</a></b></td>
					</tr>
					</table>
				</blockquote>
			</details>
		</blockquote>
	</details>
</details>

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
❯ ./rb-druid-indexer -c config.yml
```


---
##  Project Roadmap

- [X] **`Task 1`**: <strike>Implement feature one.</strike>
- [ ] **`Task 2`**: Implement feature two.
- [ ] **`Task 3`**: Implement feature three.

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

##  Acknowledgments

- List any resources, contributors, inspiration, etc. here.

---