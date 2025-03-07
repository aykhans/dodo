<h1 align="center">Dodo is a simple and easy-to-use HTTP benchmarking tool.</h1>
<p align="center">
<img width="30%" height="30%" src="https://ftp.aykhans.me/web/client/pubshares/hB6VSdCnBCr8gFPeiMuCji/browse?path=%2Fdodo.png">
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
    go build -ldflags "-s -w" -o dodo
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
docker run --rm -i aykhans/dodo -u https://example.com -m GET -d 10 -r 1000 -t 2000
```

### 2. JSON config file
You can find an example config structure in the [config.json](https://github.com/aykhans/dodo/blob/main/config.json) file:
```jsonc
{
    "method": "GET",
    "url": "https://example.com",
    "no_proxy_check": false,
    "timeout": 10000,
    "dodos": 1,
    "requests": 1,
    "params": {
        // Random param value will be selected from the param-key1 and param-key2 list for each request
        "param-key1": ["param-value1", "param-value2", "param-value3"],
        "param-key2": ["param-value1", "param-value2", "param-value3"]
    },
    "headers": {
        // Random header value will be selected from the header-key1 and header-key2 list for each request
        "header-key1": ["header-value1", "header-value2", "header-value3"],
        "header-key2": ["header-value2", "header-value2", "header-value3"]
    },
    "cookies": {
        // Random cookie value will be selected from the cookie-key1 and cookie-key2 list for each request
        "cookie-key1": ["cookie-value1", "cookie-value2", "cookie-value3"],
        "cookie-key2": ["cookie-value2", "cookie-value2", "cookie-value3"]
    },
    // Random body value will be selected from the body list for each request
    "body": ["body1", "body2", "body3"],
    // Random proxy will be selected from the proxy list for each request
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
Send 1000 GET requests to https://example.com with 10 parallel dodos (threads) and a timeout of 2000 milliseconds:

```sh
dodo -c /path/config.json
```
With Docker:
```sh
docker run --rm -i -v ./path/config.json:/dodo/config.json aykhans/dodo
```

### 3. Both (CLI & JSON)
Override the config file arguments with CLI arguments:

```sh
dodo -c /path/config.json -u https://example.com -m GET -d 10 -r 1000 -t 2000
```
With Docker:
```sh
docker run --rm -i -v ./path/config.json:/dodo/config.json aykhans/dodo -u https://example.com -m GET -d 10 -r 1000 -t 2000
```

## CLI and JSON Config Parameters
If the Headers, Params, Cookies and Body fields have multiple values, each request will choose a random value from the list.

| Parameter             | JSON config file | CLI Flag        | CLI Short Flag | Type                             | Description                                                         | Default     |
| --------------------- | ---------------- | --------------- | -------------- | -------------------------------- | ------------------------------------------------------------------- | ----------- |
| Config file           | -                | --config-file   | -c             | String                           | Path to the JSON config file                                        | -           |
| Yes                   | -                | --yes           | -y             | Boolean                          | Answer yes to all questions                                         | false       |
| URL                   | url              | --url           | -u             | String                           | URL to send the request to                                          | -           |
| Method                | method           | --method        | -m             | String                           | HTTP method                                                         | GET         |
| Requests              | requests         | --requests      | -r             | Integer                          | Total number of requests to send                                    | 1000        |
| Dodos (Threads)       | dodos            | --dodos         | -d             | Integer                          | Number of dodos (threads) to send requests in parallel              | 1           |
| Timeout               | timeout          | --timeout       | -t             | Integer                          | Timeout for canceling each request (milliseconds)                   | 10000       |
| No Proxy Check        | no_proxy_check   | --no-proxy-check| -              | Boolean                          | Disable proxy check                                                 | false       |
| Params                | params           | -               | -              | Key-Value {String: [String]}     | Request parameters                                                  | -           |
| Headers               | headers          | -               | -              | Key-Value {String: [String]}     | Request headers                                                     | -           |
| Cookies               | cookies          | -               | -              | Key-Value {String: [String]}     | Request cookies                                                     | -           |
| Body                  | body             | -               | -              | [String]                         | Request body                                                        | -           |
| Proxy                 | proxies          | -               | -              | List[Key-Value {string: string}] | List of proxies (will check active proxies before sending requests) | -           |
