<h1 align="center">Dodo is a simple and easy-to-use HTTP benchmarking tool.</h1>
<p align="center">
<img width="30%" height="30%" src="https://raw.githubusercontent.com/aykhans/dodo/main/assets/dodo.png">
</p>

## Installation
### With Docker (Recommended)
Pull the Dodo image from Docker Hub:
```sh
docker pull aykhans/dodo:latest
```
If you use Dodo with Docker and a config file, you must provide the config.json file as a volume to the Docker run command (not as the "-c config.json" argument), as shown in the examples in the [usage](#usage) section.

### With Binary File
You can grab binaries in the [releases](https://github.com/aykhans/dodo/releases) section.

### Build from Source
To build Dodo from source, you need to have [Go1.22+](https://golang.org/dl/) installed. <br>
Follow the steps below to build dodo:

1. **Clone the repository:**

    ```sh
    git clone https://github.com/aykhans/dodo.git
    ```

2. **Navigate to the project directory:**

    ```sh
    cd dodo
    ```

3. **Build the project:**

    ```sh
    go build -o dodo
    ```

This will generate an executable named `dodo` in the project directory.

## Usage
You can use Dodo with CLI arguments, a JSON config file, or both. If you use both, CLI arguments will always override JSON config arguments if there is a conflict.

### 1. CLI
Send 1000 GET requests to https://example.com with 10 parallel dodos (threads) and a timeout of 2000 milliseconds:

```sh
dodo -u https://example.com -m GET -d 10 -r 1000 -t 2000
```
With Docker:
```sh
docker run --rm aykhans/dodo -u https://example.com -m GET -d 10 -r 1000 -t 2000
```

### 2. JSON config file
You can find an example config structure in the [config.json](https://github.com/aykhans/dodo/blob/main/config.json) file:
```json
{
    "method": "GET",
    "url": "https://example.com",
    "timeout": 10000,
    "dodos_count": 50,
    "request_count": 1000,
    "params": {},
    "headers": {},
    "cookies": {},
    "body": "",
    "proxies": [
        {
            "url": "http://example.com:8080",
            "username": "username",
            "password": "password"
        },
        {
            "url": "http://example.com:8080"
        }
    ]
}
```
Send 1000 GET requests to https://example.com with 5 parallel dodos (threads) and a timeout of 10000 milliseconds:

```sh
dodo -c /path/config.json
```
With Docker:
```sh
docker run --rm -v ./path/config.json:/dodo/config.json -i aykhans/dodo
```

### 3. Both (CLI & JSON)
Override the config file arguments with CLI arguments:

```sh
dodo -c /path/config.json -u https://example.com -m GET -d 10 -r 1000 -t 2000
```
With Docker:
```sh
docker run --rm -v ./path/config.json:/dodo/config.json -i aykhans/dodo -u https://example.com -m GET -d 10 -r 1000 -t 2000
```

## CLI and JSON Config Parameters
| Parameter | JSON config file | CLI Flag | CLI Short Flag | Type | Description | Default |
| ----------- | ----------- | ----------- | ----------- | ----------- | ----------- | ----------- |
| Config file | - | --config-file | -c | String | Path to the JSON config file | - |
| URL | url | --url | -u | String | URL to send the request to | - |
| Method | method | --method | -m | String | HTTP method | GET |
| Request count | request_count | --request-count | -r | Integer | Total number of requests to send | 1000 |
| Dodos count (Threads) | dodos_count | --dodos-count | -d | Integer | Number of dodos (threads) to send requests in parallel | 1 |
| Timeout | timeout | --timeout | -t | Integer | Timeout for canceling each request (milliseconds) | 10000 |
| Params | params | - | - | Key-Value {string: string} | Request parameters | - |
| Headers | headers | - | - | Key-Value {string: string} | Request headers | - |
| Cookies | cookies | - | - | Key-Value {string: string} | Request cookies | - |
| Body | body | - | - | String | Request body | - |
| Proxy | proxies | - | - | List[Key-Value {string: string}] | List of proxies (will check active proxies before sending requests) | - |

## Examples
![dodo_example](https://raw.githubusercontent.com/aykhans/dodo/main/assets/dodo-example.gif)