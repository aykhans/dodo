<h1 align="center">Dodo - A Fast and Easy-to-Use HTTP Benchmarking Tool</h1>
<p align="center">
<img width="30%" height="30%" src="https://ftp.aykhans.me/web/client/pubshares/VzPtSHS7yPQT7ngoZzZSNU/browse?path=%2Fdodo.png">
</p>

## Installation

### Using Docker (Recommended)

Pull the Dodo image from Docker Hub:

```sh
docker pull aykhans/dodo:latest
```

When using Dodo with Docker and a local config file, you must provide the config.json file as a volume to the Docker run command (not as the "-f config.json" argument):

```sh
docker run -v /path/to/config.json:/config.json aykhans/dodo
```

If you're using Dodo with Docker and providing a config file via URL, you don't need to set a volume:

```sh
docker run aykhans/dodo -f https://raw.githubusercontent.com/aykhans/dodo/main/config.json
```

### Using Binary Files

You can download pre-built binaries from the [releases](https://github.com/aykhans/dodo/releases) section.

### Building from Source

To build Dodo from source, you need to have [Go 1.24+](https://golang.org/dl/) installed.
Follow these steps:

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

You can use Dodo with CLI arguments, a JSON config file, or both. When using both, CLI arguments will override JSON config values if there's a conflict.

### 1. CLI

Send 1000 GET requests to https://example.com with 10 parallel dodos (threads) and a timeout of 2 seconds:

```sh
dodo -u https://example.com -m GET -d 10 -r 1000 -t 2s
```

With Docker:

```sh
docker run --rm -i aykhans/dodo -u https://example.com -m GET -d 10 -r 1000 -t 2s
```

### 2. JSON Config File

Send 1000 GET requests to https://example.com with 10 parallel dodos (threads) and a timeout of 800 milliseconds:

```jsonc
{
    "method": "GET",
    "url": "https://example.com",
    "yes": false,
    "timeout": "800ms",
    "dodos": 10,
    "requests": 1000,

    "params": [
        // A random value will be selected from the list for first "key1" param on each request
        // And always "value" for second "key1" param on each request
        // e.g. "?key1=value2&key1=value"
        { "key1": ["value1", "value2", "value3", "value4"] },
        { "key1": "value" },

        // A random value will be selected from the list for param "key2" on each request
        // e.g. "?key2=value2"
        { "key2": ["value1", "value2"] },
    ],

    "headers": [
        // A random value will be selected from the list for first "key1" header on each request
        // And always "value" for second "key1" header on each request
        // e.g. "key1: value3", "key1: value"
        { "key1": ["value1", "value2", "value3", "value4"] },
        { "key1": "value" },

        // A random value will be selected from the list for header "key2" on each request
        // e.g. "key2: value2"
        { "key2": ["value1", "value2"] },
    ],

    "cookies": [
        // A random value will be selected from the list for first "key1" cookie on each request
        // And always "value" for second "key1" cookie on each request
        // e.g. "key1=value4; key1=value"
        { "key1": ["value1", "value2", "value3", "value4"] },
        { "key1": "value" },

        // A random value will be selected from the list for cookie "key2" on each request
        // e.g. "key2=value1"
        { "key2": ["value1", "value2"] },
    ],

    "body": "body-text",
    // OR
    // A random body value will be selected from the list for each request
    "body": ["body-text1", "body-text2", "body-text3"],

    "proxy": "http://example.com:8080",
    // OR
    // A random proxy will be selected from the list for each request
    "proxy": [
        "http://example.com:8080",
        "http://username:password@example.com:8080",
        "socks5://example.com:8080",
        "socks5h://example.com:8080",
    ],
}
```

```sh
dodo -f /path/config.json
# OR
dodo -f https://example.com/config.json
```

With Docker:

```sh
docker run --rm -i -v /path/to/config.json:/config.json aykhans/dodo
# OR
docker run --rm -i aykhans/dodo -f https://example.com/config.json
```

### 3. Combined (CLI & JSON)

Override the config file arguments with CLI arguments:

```sh
dodo -f /path/to/config.json -u https://example.com -m GET -d 10 -r 1000 -t 5s
```

With Docker:

```sh
docker run --rm -i -v /path/to/config.json:/config.json aykhans/dodo -u https://example.com -m GET -d 10 -r 1000 -t 5s
```

## CLI and JSON Config Parameters

If `Headers`, `Params`, `Cookies`, `Body`, or `Proxy` fields have multiple values, each request will choose a random value from the list.

| Parameter       | JSON config file | CLI Flag     | CLI Short Flag | Type                           | Description                                                     | Default |
| --------------- | ---------------- | ------------ | -------------- | ------------------------------ | --------------------------------------------------------------- | ------- |
| Config file     | -                | -config-file | -f             | String                         | Path to local config file or http(s) URL of the config file     | -       |
| Yes             | yes              | -yes         | -y             | Boolean                        | Answer yes to all questions                                     | false   |
| URL             | url              | -url         | -u             | String                         | URL to send the request to                                      | -       |
| Method          | method           | -method      | -m             | String                         | HTTP method                                                     | GET     |
| Requests        | requests         | -requests    | -r             | UnsignedInteger                | Total number of requests to send                                | 1000    |
| Dodos (Threads) | dodos            | -dodos       | -d             | UnsignedInteger                | Number of dodos (threads) to send requests in parallel          | 1       |
| Timeout         | timeout          | -timeout     | -t             | Duration                       | Timeout for canceling each request                              | 10s     |
| Params          | params           | -param       | -p             | [{String: String OR [String]}] | Request parameters                                              | -       |
| Headers         | headers          | -header      | -H             | [{String: String OR [String]}] | Request headers                                                 | -       |
| Cookies         | cookies          | -cookie      | -c             | [{String: String OR [String]}] | Request cookies                                                 | -       |
| Body            | body             | -body        | -b             | String OR [String]             | Request body or list of request bodies                          | -       |
| Proxy           | proxies          | -proxy       | -x             | String OR [String]             | Proxy URL or list of proxy URLs                                 | -       |
