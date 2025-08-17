<h1 align="center">Dodo - A Fast and Easy-to-Use HTTP Benchmarking Tool</h1>

![Usage](https://ftp.aykhans.me/web/client/pubshares/VzPtSHS7yPQT7ngoZzZSNU/browse?path=/dodo_demonstrate.gif)

<div align="center">
  <h4>
      <a href="./EXAMPLES.md">
          Examples
      </a>
      <span> | </span>
      <a href="#installation">
          Install
      </a>
      <span> | </span>
      <a href="https://hub.docker.com/r/aykhans/dodo">
          Docker
      </a>
  </h4>
  <br>
    <a href="https://coff.ee/aykhan">
        <img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 40px !important;width: 150px !important;">
    </a>
</div>

## Table of Contents

- [Installation](#installation)
    - [Using Docker (Recommended)](#using-docker-recommended)
    - [Using Pre-built Binaries](#using-pre-built-binaries)
    - [Building from Source](#building-from-source)
- [Usage](#usage)
    - [1. CLI Usage](#1-cli-usage)
    - [2. Config File Usage](#2-config-file-usage)
        - [2.1 YAML/YML Example](#21-yamlyml-example)
        - [2.2 JSON Example](#22-json-example)
    - [3. CLI & Config File Combination](#3-cli--config-file-combination)
- [Config Parameters Reference](#config-parameters-reference)
- [Template Functions](#template-functions)

## Installation

### Using Docker (Recommended)

Pull the latest Dodo image from Docker Hub:

```sh
docker pull aykhans/dodo:latest
```

To use Dodo with Docker and a local config file, mount the config file as a volume and pass it as an argument:

```sh
docker run -v /path/to/config.json:/config.json aykhans/dodo -f /config.json
```

If you're using a remote config file via URL, you don't need to mount a volume:

```sh
docker run aykhans/dodo -f https://raw.githubusercontent.com/aykhans/dodo/main/config.yaml
```

### Using Pre-built Binaries

Download the latest binaries from the [releases](https://github.com/aykhans/dodo/releases) section.

### Building from Source

To build Dodo from source, ensure you have [Go 1.24+](https://golang.org/dl/) installed.

```sh
go install -ldflags "-s -w" github.com/aykhans/dodo@latest
```

## Usage

Dodo supports CLI arguments, configuration files (JSON/YAML), or a combination of both. If both are used, CLI arguments take precedence.

### 1. CLI Usage

Send 1000 GET requests to https://example.com with 10 parallel dodos (threads), each with a timeout of 2 seconds, within a maximum duration of 1 minute:

```sh
dodo -u https://example.com -m GET -d 10 -r 1000 -o 1m -t 2s
```

With Docker:

```sh
docker run --rm -i aykhans/dodo -u https://example.com -m GET -d 10 -r 1000 -o 1m -t 2s
```

### 2. Config File Usage

Send 1000 GET requests to https://example.com with 10 parallel dodos (threads), each with a timeout of 800 milliseconds, within a maximum duration of 250 seconds:

#### 2.1 YAML/YML Example

```yaml
method: "GET"
url: "https://example.com"
yes: false
timeout: "800ms"
dodos: 10
requests: 1000
duration: "250s"
skip_verify: false

params:
    # A random value will be selected from the list for first "key1" param on each request
    # And always "value" for second "key1" param on each request
    # e.g. "?key1=value2&key1=value"
    - key1: ["value1", "value2", "value3", "value4"]
    - key1: "value"

    # A random value will be selected from the list for param "key2" on each request
    # e.g. "?key2=value2"
    - key2: ["value1", "value2"]

headers:
    # A random value will be selected from the list for first "key1" header on each request
    # And always "value" for second "key1" header on each request
    # e.g. "key1: value3", "key1: value"
    - key1: ["value1", "value2", "value3", "value4"]
    - key1: "value"

    # A random value will be selected from the list for header "key2" on each request
    # e.g. "key2: value2"
    - key2: ["value1", "value2"]

cookies:
    # A random value will be selected from the list for first "key1" cookie on each request
    # And always "value" for second "key1" cookie on each request
    # e.g. "key1=value4; key1=value"
    - key1: ["value1", "value2", "value3", "value4"]
    - key1: "value"

    # A random value will be selected from the list for cookie "key2" on each request
    # e.g. "key2=value1"
    - key2: ["value1", "value2"]

body: "body-text"
# OR
# A random body value will be selected from the list for each request
body:
    - "body-text1"
    - "body-text2"
    - "body-text3"

proxy: "http://example.com:8080"
# OR
# A random proxy will be selected from the list for each request
proxy:
    - "http://example.com:8080"
    - "http://username:password@example.com:8080"
    - "socks5://example.com:8080"
    - "socks5h://example.com:8080"
```

```sh
dodo -f /path/config.yaml
# OR
dodo -f https://example.com/config.yaml
```

With Docker:

```sh
docker run --rm -i -v /path/to/config.yaml:/config.yaml aykhans/dodo -f /config.yaml
# OR
docker run --rm -i aykhans/dodo -f https://example.com/config.yaml
```

#### 2.2 JSON Example

```jsonc
{
    "method": "GET",
    "url": "https://example.com",
    "yes": false,
    "timeout": "800ms",
    "dodos": 10,
    "requests": 1000,
    "duration": "250s",
    "skip_verify": false,

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

### 3. CLI & Config File Combination

CLI arguments override config file values:

```sh
dodo -f /path/to/config.yaml -u https://example.com -m GET -d 10 -r 1000 -o 1m -t 5s
```

With Docker:

```sh
docker run --rm -i -v /path/to/config.json:/config.json aykhans/dodo -f /config.json -u https://example.com -m GET -d 10 -r 1000 -o 1m -t 5s
```

You can find more usage examples in the [EXAMPLES.md](./EXAMPLES.md) file.

## Config Parameters Reference

If `Headers`, `Params`, `Cookies`, `Body`, or `Proxy` fields have multiple values, each request will choose a random value from the list.

| Parameter       | config file | CLI Flag     | CLI Short Flag | Type                           | Description                                                 | Default |
| --------------- | ----------- | ------------ | -------------- | ------------------------------ | ----------------------------------------------------------- | ------- |
| Config file     |             | -config-file | -f             | String                         | Path to local config file or http(s) URL of the config file | -       |
| Yes             | yes         | -yes         | -y             | Boolean                        | Answer yes to all questions                                 | false   |
| URL             | url         | -url         | -u             | String                         | URL to send the request to                                  | -       |
| Method          | method      | -method      | -m             | String                         | HTTP method                                                 | GET     |
| Dodos (Threads) | dodos       | -dodos       | -d             | UnsignedInteger                | Number of dodos (threads) to send requests in parallel      | 1       |
| Requests        | requests    | -requests    | -r             | UnsignedInteger                | Total number of requests to send                            | -       |
| Duration        | duration    | -duration    | -o             | Time                           | Maximum duration for the test                               | -       |
| Timeout         | timeout     | -timeout     | -t             | Time                           | Timeout for canceling each request                          | 10s     |
| Params          | params      | -param       | -p             | [{String: String OR [String]}] | Request parameters                                          | -       |
| Headers         | headers     | -header      | -H             | [{String: String OR [String]}] | Request headers                                             | -       |
| Cookies         | cookies     | -cookie      | -c             | [{String: String OR [String]}] | Request cookies                                             | -       |
| Body            | body        | -body        | -b             | String OR [String]             | Request body or list of request bodies                      | -       |
| Proxy           | proxies     | -proxy       | -x             | String OR [String]             | Proxy URL or list of proxy URLs                             | -       |
| Skip Verify     | skip_verify | -skip-verify |                | Boolean                        | Skip SSL/TLS certificate verification                       | false   |

## Template Functions

Dodo supports template functions in `Headers`, `Params`, `Cookies`, and `Body` fields. These functions allow you to generate dynamic values for each request.

You can use Go template syntax to include dynamic values in your requests. Here's how to use template functions:

In CLI config:

```sh
dodo -u https://example.com -r 1 \
    -header "User-Agent:{{ fakeit_UserAgent }}" \ # e.g. "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)"
    -param "username={{ strings_ToUpper fakeit_Username }}" \ # e.g. "username=JOHN BOB"
    -cookie "token={{ fakeit_Password true true true true true 10 }}" \ # e.g. token=1234567890abcdef1234567890abcdef
    -body '{"email":"{{ fakeit_Email }}", "password":"{{ fakeit_Password true true true true true 10 }}"}' # e.g. {"email":"john.doe@example.com", "password":"12rw4d-78d"}
```

In YAML/YML config:

```yaml
headers:
    - User-Agent: "{{ fakeit_UserAgent }}" # e.g. "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)"
    - "Random-Header-{{fakeit_FirstName}}": "static_value" # e.g. "Random-Header-John: static_value"

cookies:
    - token: "Bearer {{ fakeit_UUID }}" # e.g. "token=Bearer 1234567890abcdef1234567890abcdef"

params:
    - id: "{{ fakeit_Uint }}" # e.g. "id=1234567890"
    - username: "{{ fakeit_Username }}" # e.g. "username=John Doe"

body:
    - '{ "username": "{{ fakeit_Username }}", "password": "{{ fakeit_Password }}" }' # e.g. { "username": "john.doe", "password": "password123" }
    - '{ "email": "{{ fakeit_Email }}", "phone": "{{ fakeit_Phone }}" }' # e.g. { "email": "john.doe@example.com", "phone": "1234567890" }
    - '{{ body_FormData (dict_Str "username" fakeit_Username "password" "secret123") }}' # Creates multipart form data for form submissions, automatically sets the appropriate Content-Type header.
```

In JSON config:

```jsonc
{
    "headers": [
        { "User-Agent": "{{ fakeit_UserAgent }}" }, // e.g. "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)"
    ],
    "body": [
        "{ \"username\": \"{{ strings_RemoveSpaces fakeit_Username }}\", \"password\": \"{{ fakeit_Password }}\" }", // e.g. { "username": "johndoe", "password": "password123" }
        "{{ body_FormData (dict_Str \"username\" fakeit_Username \"password\" \"12345\") }}", // Creates multipart form data for form submissions, automatically sets the appropriate Content-Type header.
    ],
}
```

For the full list of template functions over 200 functions, refer to the `NewFuncMap` function in `utils/templates.go`.
