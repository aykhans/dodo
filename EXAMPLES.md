# Dodo Usage Examples

This document provides comprehensive examples of using Dodo with various configuration combinations. Each example includes three methods: CLI usage, YAML configuration, and JSON configuration.

## Table of Contents

1. [Basic HTTP Stress Testing](#1-basic-http-stress-testing)
2. [POST Request with Form Data](#2-post-request-with-form-data)
3. [API Testing with Authentication](#3-api-testing-with-authentication)
4. [Testing with Custom Headers and Cookies](#4-testing-with-custom-headers-and-cookies)
5. [Load Testing with Proxy Rotation](#5-load-testing-with-proxy-rotation)
6. [JSON API Testing with Dynamic Data](#6-json-api-testing-with-dynamic-data)
7. [File Upload Testing](#7-file-upload-testing)
8. [E-commerce Cart Testing](#8-e-commerce-cart-testing)
9. [GraphQL API Testing](#9-graphql-api-testing)
10. [WebSocket-style HTTP Testing](#10-websocket-style-http-testing)
11. [Multi-tenant Application Testing](#11-multi-tenant-application-testing)
12. [Rate Limiting Testing](#12-rate-limiting-testing)

---

## 1. Basic HTTP Stress Testing

Test a simple website with basic GET requests to measure performance under load.

### CLI Usage

```bash
dodo -u https://httpbin.org/get \
     -m GET \
     -d 5 \
     -r 100 \
     -t 5s \
     -o 30s \
     --skip-verify=false \
     -y
```

### YAML Configuration

```yaml
method: "GET"
url: "https://httpbin.org/get"
yes: true
timeout: "5s"
dodos: 5
requests: 100
duration: "30s"
skip_verify: false
```

### JSON Configuration

```json
{
    "method": "GET",
    "url": "https://httpbin.org/get",
    "yes": true,
    "timeout": "5s",
    "dodos": 5,
    "requests": 100,
    "duration": "30s",
    "skip_verify": false
}
```

---

## 2. POST Request with Form Data

Test form submission endpoints with randomized form data.

### CLI Usage

```bash
dodo -u https://httpbin.org/post \
     -m POST \
     -d 3 \
     -r 50 \
     -t 10s \
     --skip-verify=false \
     -H "Content-Type:application/x-www-form-urlencoded" \
     -b "username={{ fakeit_Username }}&password={{ fakeit_Password true true true true true 12 }}&email={{ fakeit_Email }}" \
     -b "username={{ fakeit_Username }}&password={{ fakeit_Password true true true true true 8 }}&email={{ fakeit_Email }}" \
     -y
```

### YAML Configuration

```yaml
method: "POST"
url: "https://httpbin.org/post"
yes: true
timeout: "10s"
dodos: 3
requests: 50
skip_verify: false

headers:
    - Content-Type: "application/x-www-form-urlencoded"

body:
    - "username={{ fakeit_Username }}&password={{ fakeit_Password true true true true true 12 }}&email={{ fakeit_Email }}"
    - "username={{ fakeit_Username }}&password={{ fakeit_Password true true true true true 8 }}&email={{ fakeit_Email }}"
```

### JSON Configuration

```json
{
    "method": "POST",
    "url": "https://httpbin.org/post",
    "yes": true,
    "timeout": "10s",
    "dodos": 3,
    "requests": 50,
    "skip_verify": false,
    "headers": [{ "Content-Type": "application/x-www-form-urlencoded" }],
    "body": [
        "username={{ fakeit_Username }}&password={{ fakeit_Password true true true true true 12 }}&email={{ fakeit_Email }}",
        "username={{ fakeit_Username }}&password={{ fakeit_Password true true true true true 8 }}&email={{ fakeit_Email }}"
    ]
}
```

---

## 3. API Testing with Authentication

Test protected API endpoints with various authentication methods.

### CLI Usage

```bash
dodo -u https://httpbin.org/bearer \
     -m GET \
     -d 4 \
     -r 200 \
     -t 8s \
     --skip-verify=false \
     -H "Authorization:Bearer {{ fakeit_LetterN 32 }}" \
     -H "User-Agent:{{ fakeit_UserAgent }}" \
     -H "X-Request-ID:{{ fakeit_Int }}" \
     -H "Accept:application/json" \
     -p "api_version=v1" \
     -p "format=json" \
     -p "client_id=mobile" -p "client_id=web" -p "client_id=desktop" \
     -y
```

### YAML Configuration

```yaml
method: "GET"
url: "https://httpbin.org/bearer"
yes: true
timeout: "8s"
dodos: 4
requests: 200
skip_verify: false

params:
    - api_version: "v1"
    - format: "json"
    - client_id: ["mobile", "web", "desktop"]

headers:
    - Authorization: "Bearer {{ fakeit_LetterN 32 }}"
    - User-Agent: "{{ fakeit_UserAgent }}"
    - X-Request-ID: "{{ fakeit_Int }}"
    - Accept: "application/json"
```

### JSON Configuration

```json
{
    "method": "GET",
    "url": "https://httpbin.org/bearer",
    "yes": true,
    "timeout": "8s",
    "dodos": 4,
    "requests": 200,
    "skip_verify": false,
    "params": [
        { "api_version": "v1" },
        { "format": "json" },
        { "client_id": ["mobile", "web", "desktop"] }
    ],
    "headers": [
        { "Authorization": "Bearer {{ fakeit_LetterN 32 }}" },
        { "User-Agent": "{{ fakeit_UserAgent }}" },
        { "X-Request-ID": "{{ fakeit_Int }}" },
        { "Accept": "application/json" }
    ]
}
```

---

## 4. Testing with Custom Headers and Cookies

Test applications that require specific headers and session cookies.

### CLI Usage

```bash
dodo -u https://httpbin.org/cookies \
     -m GET \
     -d 6 \
     -r 75 \
     -t 5s \
     --skip-verify=false \
     -H 'Accept-Language:{{ strings_Join "," (fakeit_LanguageAbbreviation) (fakeit_LanguageAbbreviation) (fakeit_LanguageAbbreviation) }}' \
     -H "X-Forwarded-For:{{ fakeit_IPv4Address }}" \
     -H "X-Real-IP:{{ fakeit_IPv4Address }}" \
     -H "Accept-Encoding:gzip" -H "Accept-Encoding:deflate" -H "Accept-Encoding:br" \
     -c "session_id={{ fakeit_UUID }}" \
     -c 'user_pref={{ fakeit_RandomString "a1" "b2" "c3" }}' \
     -c "theme=dark" -c "theme=light" -c "theme=auto" \
     -c "lang=en" -c "lang=es" -c "lang=fr" -c "lang=de" \
     -y
```

### YAML Configuration

```yaml
method: "GET"
url: "https://httpbin.org/cookies"
yes: true
timeout: "5s"
dodos: 6
requests: 75
skip_verify: false

headers:
    - Accept-Language: '{{ strings_Join "," (fakeit_LanguageAbbreviation) (fakeit_LanguageAbbreviation) (fakeit_LanguageAbbreviation) }}'
    - X-Forwarded-For: "{{ fakeit_IPv4Address }}"
    - X-Real-IP: "{{ fakeit_IPv4Address }}"
    - Accept-Encoding: ["gzip", "deflate", "br"]

cookies:
    - session_id: "{{ fakeit_UUID }}"
    - user_pref: '{{ fakeit_RandomString "a1" "b2" "c3" }}'
    - theme: ["dark", "light", "auto"]
    - lang: ["en", "es", "fr", "de"]
```

### JSON Configuration

```json
{
    "method": "GET",
    "url": "https://httpbin.org/cookies",
    "yes": true,
    "timeout": "5s",
    "dodos": 6,
    "requests": 75,
    "skip_verify": false,
    "headers": [
        {
            "Accept-Language": "{{ strings_Join \",\" (fakeit_LanguageAbbreviation) (fakeit_LanguageAbbreviation) (fakeit_LanguageAbbreviation) }}"
        },
        { "X-Forwarded-For": "{{ fakeit_IPv4Address }}" },
        { "X-Real-IP": "{{ fakeit_IPv4Address }}" },
        { "Accept-Encoding": ["gzip", "deflate", "br"] }
    ],
    "cookies": [
        { "session_id": "{{ fakeit_UUID }}" },
        { "user_pref": "{{ fakeit_RandomString \"a1\" \"b2\" \"c3\" }}" },
        { "theme": ["dark", "light", "auto"] },
        { "lang": ["en", "es", "fr", "de"] }
    ]
}
```

---

## 5. Load Testing with Proxy Rotation

Test through multiple proxies for distributed load testing.

### CLI Usage

```bash
dodo -u https://httpbin.org/ip \
     -m GET \
     -d 8 \
     -r 300 \
     -t 15s \
     --skip-verify=false \
     -x "http://proxy1.example.com:8080" \
     -x "http://proxy2.example.com:8080" \
     -x "socks5://proxy3.example.com:1080" \
     -x "http://username:password@proxy4.example.com:8080" \
     -H "User-Agent:{{ fakeit_UserAgent }}" \
     -H "Accept:application/json" \
     -y
```

### YAML Configuration

```yaml
method: "GET"
url: "https://httpbin.org/ip"
yes: true
timeout: "15s"
dodos: 8
requests: 300
skip_verify: false

proxy:
    - "http://proxy1.example.com:8080"
    - "http://proxy2.example.com:8080"
    - "socks5://proxy3.example.com:1080"
    - "http://username:password@proxy4.example.com:8080"

headers:
    - User-Agent: "{{ fakeit_UserAgent }}"
    - Accept: "application/json"
```

### JSON Configuration

```json
{
    "method": "GET",
    "url": "https://httpbin.org/ip",
    "yes": true,
    "timeout": "15s",
    "dodos": 8,
    "requests": 300,
    "skip_verify": false,
    "proxy": [
        "http://proxy1.example.com:8080",
        "http://proxy2.example.com:8080",
        "socks5://proxy3.example.com:1080",
        "http://username:password@proxy4.example.com:8080"
    ],
    "headers": [
        { "User-Agent": "{{ fakeit_UserAgent }}" },
        { "Accept": "application/json" }
    ]
}
```

---

## 6. JSON API Testing with Dynamic Data

Test REST APIs with realistic JSON payloads.

### CLI Usage

```bash
dodo -u https://httpbin.org/post \
     -m POST \
     -d 5 \
     -r 150 \
     -t 12s \
     --skip-verify=false \
     -H "Content-Type:application/json" \
     -H "Accept:application/json" \
     -H "X-API-Version:2023-10-01" \
     -b '{"user_id":{{ fakeit_Uint }},"name":"{{ fakeit_Name }}","email":"{{ fakeit_Email }}","created_at":"{{ fakeit_Date }}"}' \
     -b '{"product_id":{{ fakeit_Uint }},"name":"{{ fakeit_ProductName }}","price":{{ fakeit_Price 10 1000 }},"category":"{{ fakeit_ProductCategory }}"}' \
     -b '{"order_id":"{{ fakeit_UUID }}","items":[{"id":{{ fakeit_Uint }},"quantity":{{ fakeit_IntRange 1 10 }}}],"total":{{ fakeit_Price 50 500 }}}' \
     -y
```

### YAML Configuration

```yaml
method: "POST"
url: "https://httpbin.org/post"
yes: true
timeout: "12s"
dodos: 5
requests: 150
skip_verify: false

headers:
    - Content-Type: "application/json"
    - Accept: "application/json"
    - X-API-Version: "2023-10-01"

body:
    - '{"user_id":{{ fakeit_Uint }},"name":"{{ fakeit_Name }}","email":"{{ fakeit_Email }}","created_at":"{{ fakeit_Date }}"}'
    - '{"product_id":{{ fakeit_Uint }},"name":"{{ fakeit_ProductName }}","price":{{ fakeit_Price 10 1000 }},"category":"{{ fakeit_ProductCategory }}"}'
    - '{"order_id":"{{ fakeit_UUID }}","items":[{"id":{{ fakeit_Uint }},"quantity":{{ fakeit_IntRange 1 10 }}}],"total":{{ fakeit_Price 50 500 }}}'
```

### JSON Configuration

```json
{
    "method": "POST",
    "url": "https://httpbin.org/post",
    "yes": true,
    "timeout": "12s",
    "dodos": 5,
    "requests": 150,
    "skip_verify": false,
    "headers": [
        { "Content-Type": "application/json" },
        { "Accept": "application/json" },
        { "X-API-Version": "2023-10-01" }
    ],
    "body": [
        "{\"user_id\":{{ fakeit_Uint }},\"name\":\"{{ fakeit_Name }}\",\"email\":\"{{ fakeit_Email }}\",\"created_at\":\"{{ fakeit_Date }}\"}",
        "{\"product_id\":{{ fakeit_Uint }},\"name\":\"{{ fakeit_ProductName }}\",\"price\":{{ fakeit_Price 10 1000 }},\"category\":\"{{ fakeit_ProductCategory }}\"}",
        "{\"order_id\":\"{{ fakeit_UUID }}\",\"items\":[{\"id\":{{ fakeit_Uint }},\"quantity\":{{ fakeit_IntRange 1 10 }}}],\"total\":{{ fakeit_Price 50 500 }}}"
    ]
}
```

---

## 7. File Upload Testing

Test file upload endpoints with multipart form data.

### CLI Usage

```bash
dodo -u https://httpbin.org/post \
     -m POST \
     -d 3 \
     -r 25 \
     -t 30s \
     --skip-verify=false \
     -H "X-Upload-Source:dodo-test" \
     -H "User-Agent:{{ fakeit_UserAgent }}" \
     -b '{{ body_FormData (dict_Str "filename" (fakeit_UUID) "content" (fakeit_Paragraph 3 5 10 " ")) }}' \
     -b '{{ body_FormData (dict_Str "file" (fakeit_UUID) "description" (fakeit_Sentence 10) "category" "image") }}' \
     -y
```

### YAML Configuration

```yaml
method: "POST"
url: "https://httpbin.org/post"
yes: true
timeout: "30s"
dodos: 3
requests: 25
skip_verify: false

headers:
    - X-Upload-Source: "dodo-test"
    - User-Agent: "{{ fakeit_UserAgent }}"

body:
    - '{{ body_FormData (dict_Str "filename" (fakeit_UUID) "content" (fakeit_Paragraph 3 5 10 " ")) }}'
    - '{{ body_FormData (dict_Str "file" (fakeit_UUID) "description" (fakeit_Sentence 10) "category" "image") }}'
```

### JSON Configuration

```json
{
    "method": "POST",
    "url": "https://httpbin.org/post",
    "yes": true,
    "timeout": "30s",
    "dodos": 3,
    "requests": 25,
    "skip_verify": false,
    "headers": [
        { "X-Upload-Source": "dodo-test" },
        { "User-Agent": "{{ fakeit_UserAgent }}" }
    ],
    "body": [
        "{{ body_FormData (dict_Str \"filename\" (fakeit_UUID) \"content\" (fakeit_Paragraph 3 5 10 \" \")) }}",
        "{{ body_FormData (dict_Str \"file\" (fakeit_UUID) \"description\" (fakeit_Sentence 10) \"category\" \"image\") }}"
    ]
}
```

---

## 8. E-commerce Cart Testing

Test shopping cart operations with realistic product data.

### CLI Usage

```bash
dodo -u https://api.example-shop.com/cart \
     -m POST \
     -d 10 \
     -r 500 \
     -t 8s \
     --skip-verify=false \
     -H "Content-Type:application/json" \
     -H "Authorization:Bearer {{ fakeit_LetterN 32 }}" \
     -H "X-Client-Version:1.2.3" \
     -H "User-Agent:{{ fakeit_UserAgent }}" \
     -c "cart_session={{ fakeit_UUID }}" \
     -c "user_pref=guest" -c "user_pref=member" -c "user_pref=premium" \
     -c "region=US" -c "region=EU" -c "region=ASIA" \
     -p "currency=USD" -p "currency=EUR" -p "currency=GBP" \
     -p "locale=en-US" -p "locale=en-GB" -p "locale=de-DE" -p "locale=fr-FR" \
     -b '{"action":"add","product_id":"{{ fakeit_UUID }}","quantity":{{ fakeit_IntRange 1 5 }},"user_id":"{{ fakeit_UUID }}"}' \
     -b '{"action":"remove","product_id":"{{ fakeit_UUID }}","user_id":"{{ fakeit_UUID }}"}' \
     -b '{"action":"update","product_id":"{{ fakeit_UUID }}","quantity":{{ fakeit_IntRange 1 10 }},"user_id":"{{ fakeit_UUID }}"}' \
     -y
```

### YAML Configuration

```yaml
method: "POST"
url: "https://api.example-shop.com/cart"
yes: true
timeout: "8s"
dodos: 10
requests: 500
skip_verify: false

headers:
    - Content-Type: "application/json"
    - Authorization: "Bearer {{ fakeit_LetterN 32 }}"
    - X-Client-Version: "1.2.3"
    - User-Agent: "{{ fakeit_UserAgent }}"

cookies:
    - cart_session: "{{ fakeit_UUID }}"
    - user_pref: ["guest", "member", "premium"]
    - region: ["US", "EU", "ASIA"]

params:
    - currency: ["USD", "EUR", "GBP"]
    - locale: ["en-US", "en-GB", "de-DE", "fr-FR"]

body:
    - '{"action":"add","product_id":"{{ fakeit_UUID }}","quantity":{{ fakeit_IntRange 1 5 }},"user_id":"{{ fakeit_UUID }}"}'
    - '{"action":"remove","product_id":"{{ fakeit_UUID }}","user_id":"{{ fakeit_UUID }}"}'
    - '{"action":"update","product_id":"{{ fakeit_UUID }}","quantity":{{ fakeit_IntRange 1 10 }},"user_id":"{{ fakeit_UUID }}"}'
```

### JSON Configuration

```json
{
    "method": "POST",
    "url": "https://api.example-shop.com/cart",
    "yes": true,
    "timeout": "8s",
    "dodos": 10,
    "requests": 500,
    "skip_verify": false,
    "headers": [
        { "Content-Type": "application/json" },
        { "Authorization": "Bearer {{ fakeit_LetterN 32 }}" },
        { "X-Client-Version": "1.2.3" },
        { "User-Agent": "{{ fakeit_UserAgent }}" }
    ],
    "cookies": [
        { "cart_session": "{{ fakeit_UUID }}" },
        { "user_pref": ["guest", "member", "premium"] },
        { "region": ["US", "EU", "ASIA"] }
    ],
    "params": [
        { "currency": ["USD", "EUR", "GBP"] },
        { "locale": ["en-US", "en-GB", "de-DE", "fr-FR"] }
    ],
    "body": [
        "{\"action\":\"add\",\"product_id\":\"{{ fakeit_UUID }}\",\"quantity\":{{ fakeit_IntRange 1 5 }},\"user_id\":\"{{ fakeit_UUID }}\"}",
        "{\"action\":\"remove\",\"product_id\":\"{{ fakeit_UUID }}\",\"user_id\":\"{{ fakeit_UUID }}\"}",
        "{\"action\":\"update\",\"product_id\":\"{{ fakeit_UUID }}\",\"quantity\":{{ fakeit_IntRange 1 10 }},\"user_id\":\"{{ fakeit_UUID }}\"}"
    ]
}
```

---

## 9. GraphQL API Testing

Test GraphQL endpoints with various queries and mutations.

### CLI Usage

```bash
dodo -u https://api.example.com/graphql \
     -m POST \
     -d 4 \
     -r 100 \
     -t 10s \
     --skip-verify=false \
     -H "Content-Type:application/json" \
     -H "Authorization:Bearer {{ fakeit_UUID }}" \
     -H "X-GraphQL-Client:dodo-test" \
     -b '{"query":"query GetUser($id: ID!) { user(id: $id) { id name email } }","variables":{"id":"{{ fakeit_UUID }}"}}' \
     -b '{"query":"query GetPosts($limit: Int) { posts(limit: $limit) { id title content } }","variables":{"limit":{{ fakeit_IntRange 5 20 }}}}' \
     -b '{"query":"mutation CreatePost($input: PostInput!) { createPost(input: $input) { id title } }","variables":{"input":{"title":"{{ fakeit_Sentence 5 }}","content":"{{ fakeit_Paragraph 2 3 5 " "}}","authorId":"{{ fakeit_UUID }}"}}}' \
     -y
```

### YAML Configuration

```yaml
method: "POST"
url: "https://api.example.com/graphql"
yes: true
timeout: "10s"
dodos: 4
requests: 100
skip_verify: false

headers:
    - Content-Type: "application/json"
    - Authorization: "Bearer {{ fakeit_UUID }}"
    - X-GraphQL-Client: "dodo-test"

body:
    - '{"query":"query GetUser($id: ID!) { user(id: $id) { id name email } }","variables":{"id":"{{ fakeit_UUID }}"}}'
    - '{"query":"query GetPosts($limit: Int) { posts(limit: $limit) { id title content } }","variables":{"limit":{{ fakeit_IntRange 5 20 }}}}'
    - '{"query":"mutation CreatePost($input: PostInput!) { createPost(input: $input) { id title } }","variables":{"input":{"title":"{{ fakeit_Sentence 5 }}","content":"{{ fakeit_Paragraph 2 3 5 " "}}","authorId":"{{ fakeit_UUID }}"}}}'
```

### JSON Configuration

```json
{
    "method": "POST",
    "url": "https://api.example.com/graphql",
    "yes": true,
    "timeout": "10s",
    "dodos": 4,
    "requests": 100,
    "skip_verify": false,
    "headers": [
        { "Content-Type": "application/json" },
        { "Authorization": "Bearer {{ fakeit_UUID }}" },
        { "X-GraphQL-Client": "dodo-test" }
    ],
    "body": [
        "{\"query\":\"query GetUser($id: ID!) { user(id: $id) { id name email } }\",\"variables\":{\"id\":\"{{ fakeit_UUID }}\"}}",
        "{\"query\":\"query GetPosts($limit: Int) { posts(limit: $limit) { id title content } }\",\"variables\":{\"limit\":{{ fakeit_IntRange 5 20 }}}}",
        "{\"query\":\"mutation CreatePost($input: PostInput!) { createPost(input: $input) { id title } }\",\"variables\":{\"input\":{\"title\":\"{{ fakeit_Sentence 5 }}\",\"content\":\"{{ fakeit_Paragraph 2 3 5 \\\" \\\"}}\",\"authorId\":\"{{ fakeit_UUID }}\"}}}"
    ]
}
```

---

## 10. WebSocket-style HTTP Testing

Test real-time applications with WebSocket-like HTTP endpoints.

### CLI Usage

```bash
dodo -u https://api.realtime-app.com/events \
     -m POST \
     -d 15 \
     -r 1000 \
     -t 5s \
     -o 60s \
     --skip-verify=false \
     -H "Content-Type:application/json" \
     -H "X-Event-Type:{{ fakeit_LetterNN 4 12 }}" \
     -H "Connection:keep-alive" \
     -H "Cache-Control:no-cache" \
     -c "connection_id={{ fakeit_UUID }}" \
     -c "session_token={{ fakeit_UUID }}" \
     -p "channel=general" -p "channel=notifications" -p "channel=alerts" -p "channel=updates" \
     -p "version=v1" -p "version=v2" \
     -b '{"event":"{{ fakeit_Word }}","data":{"timestamp":"{{ fakeit_Date }}","user_id":"{{ fakeit_UUID }}","message":"{{ fakeit_Sentence 8 }}"}}' \
     -b '{"event":"ping","data":{"timestamp":"{{ fakeit_Date }}","client_id":"{{ fakeit_UUID }}"}}' \
     -b '{"event":"status_update","data":{"status":"{{ fakeit_Word }}","user_id":"{{ fakeit_UUID }}","timestamp":"{{ fakeit_Date }}"}}' \
     -y
```

### YAML Configuration

```yaml
method: "POST"
url: "https://api.realtime-app.com/events"
yes: true
timeout: "5s"
dodos: 15
requests: 1000
duration: "60s"
skip_verify: false

headers:
    - Content-Type: "application/json"
    - X-Event-Type: "{{ fakeit_LetterNN 4 12 }}"
    - Connection: "keep-alive"
    - Cache-Control: "no-cache"

cookies:
    - connection_id: "{{ fakeit_UUID }}"
    - session_token: "{{ fakeit_UUID }}"

params:
    - channel: ["general", "notifications", "alerts", "updates"]
    - version: ["v1", "v2"]

body:
    - '{"event":"{{ fakeit_Word }}","data":{"timestamp":"{{ fakeit_Date }}","user_id":"{{ fakeit_UUID }}","message":"{{ fakeit_Sentence 8 }}"}}'
    - '{"event":"ping","data":{"timestamp":"{{ fakeit_Date }}","client_id":"{{ fakeit_UUID }}"}}'
    - '{"event":"status_update","data":{"status":"{{ fakeit_Word }}","user_id":"{{ fakeit_UUID }}","timestamp":"{{ fakeit_Date }}"}}'
```

### JSON Configuration

```json
{
    "method": "POST",
    "url": "https://api.realtime-app.com/events",
    "yes": true,
    "timeout": "5s",
    "dodos": 15,
    "requests": 1000,
    "duration": "60s",
    "skip_verify": false,
    "headers": [
        { "Content-Type": "application/json" },
        { "X-Event-Type": "{{ fakeit_LetterNN 4 12 }}" },
        { "Connection": "keep-alive" },
        { "Cache-Control": "no-cache" }
    ],
    "cookies": [
        { "connection_id": "{{ fakeit_UUID }}" },
        { "session_token": "{{ fakeit_UUID }}" }
    ],
    "params": [
        { "channel": ["general", "notifications", "alerts", "updates"] },
        { "version": ["v1", "v2"] }
    ],
    "body": [
        "{\"event\":\"{{ fakeit_Word }}\",\"data\":{\"timestamp\":\"{{ fakeit_Date }}\",\"user_id\":\"{{ fakeit_UUID }}\",\"message\":\"{{ fakeit_Sentence 8 }}\"}}",
        "{\"event\":\"ping\",\"data\":{\"timestamp\":\"{{ fakeit_Date }}\",\"client_id\":\"{{ fakeit_UUID }}\"}}",
        "{\"event\":\"status_update\",\"data\":{\"status\":\"{{ fakeit_Word }}\",\"user_id\":\"{{ fakeit_UUID }}\",\"timestamp\":\"{{ fakeit_Date }}\"}}"
    ]
}
```

---

## 11. Multi-tenant Application Testing

Test SaaS applications with tenant-specific configurations.

### CLI Usage

```bash
dodo -u https://app.saas-platform.com/api/data \
     -m GET \
     -d 12 \
     -r 600 \
     -t 15s \
     --skip-verify=false \
     -H "X-Tenant-ID:{{ fakeit_UUID }}" \
     -H "Authorization:Bearer {{ fakeit_LetterN 64 }}" \
     -H "X-Client-Type:web" -H "X-Client-Type:mobile" -H "X-Client-Type:api" \
     -H "Accept:application/json" \
     -c "tenant_session={{ fakeit_UUID }}" \
     -c "user_role=admin" -c "user_role=user" -c "user_role=viewer" \
     -c "subscription_tier=free" -c "subscription_tier=pro" -c "subscription_tier=enterprise" \
     -p "page={{ fakeit_IntRange 1 10 }}" \
     -p "limit={{ fakeit_IntRange 10 100 }}" \
     -p "sort=created_at" -p "sort=updated_at" -p "sort=name" \
     -p "order=asc" -p "order=desc" \
     -p "filter_by=active" -p "filter_by=inactive" -p "filter_by=pending" \
     -y
```

### YAML Configuration

```yaml
method: "GET"
url: "https://app.saas-platform.com/api/data"
yes: true
timeout: "15s"
dodos: 12
requests: 600
skip_verify: false

headers:
    - X-Tenant-ID: "{{ fakeit_UUID }}"
    - Authorization: "Bearer {{ fakeit_LetterN 64 }}"
    - X-Client-Type: ["web", "mobile", "api"]
    - Accept: "application/json"

cookies:
    - tenant_session: "{{ fakeit_UUID }}"
    - user_role: ["admin", "user", "viewer"]
    - subscription_tier: ["free", "pro", "enterprise"]

params:
    - page: "{{ fakeit_IntRange 1 10 }}"
    - limit: "{{ fakeit_IntRange 10 100 }}"
    - sort: ["created_at", "updated_at", "name"]
    - order: ["asc", "desc"]
    - filter_by: ["active", "inactive", "pending"]
```

### JSON Configuration

```json
{
    "method": "GET",
    "url": "https://app.saas-platform.com/api/data",
    "yes": true,
    "timeout": "15s",
    "dodos": 12,
    "requests": 600,
    "skip_verify": false,
    "headers": [
        { "X-Tenant-ID": "{{ fakeit_UUID }}" },
        { "Authorization": "Bearer {{ fakeit_LetterN 64 }}" },
        { "X-Client-Type": ["web", "mobile", "api"] },
        { "Accept": "application/json" }
    ],
    "cookies": [
        { "tenant_session": "{{ fakeit_UUID }}" },
        { "user_role": ["admin", "user", "viewer"] },
        { "subscription_tier": ["free", "pro", "enterprise"] }
    ],
    "params": [
        { "page": "{{ fakeit_IntRange 1 10 }}" },
        { "limit": "{{ fakeit_IntRange 10 100 }}" },
        { "sort": ["created_at", "updated_at", "name"] },
        { "order": ["asc", "desc"] },
        { "filter_by": ["active", "inactive", "pending"] }
    ]
}
```

---

## 12. Rate Limiting Testing

Test API rate limits and throttling mechanisms.

### CLI Usage

```bash
dodo -u https://api.rate-limited.com/endpoint \
     -m GET \
     -d 20 \
     -r 2000 \
     -t 3s \
     -o 120s \
     --skip-verify=false \
     -H "X-API-Key:{{ fakeit_UUID }}" \
     -H "X-Client-ID:{{ fakeit_UUID }}" \
     -H "X-Rate-Limit-Test:true" \
     -H "User-Agent:{{ fakeit_UserAgent }}" \
     -c "rate_limit_bucket={{ fakeit_UUID }}" \
     -c "client_tier=tier1" -c "client_tier=tier2" -c "client_tier=tier3" \
     -p "burst_test=true" \
     -p "client_type=premium" -p "client_type=standard" -p "client_type=free" \
     -p "request_id={{ fakeit_UUID }}" \
     -y
```

### YAML Configuration

```yaml
method: "GET"
url: "https://api.rate-limited.com/endpoint"
yes: true
timeout: "3s"
dodos: 20
requests: 2000
duration: "120s"
skip_verify: false

headers:
    - X-API-Key: "{{ fakeit_UUID }}"
    - X-Client-ID: "{{ fakeit_UUID }}"
    - X-Rate-Limit-Test: "true"
    - User-Agent: "{{ fakeit_UserAgent }}"

params:
    - burst_test: "true"
    - client_type: ["premium", "standard", "free"]
    - request_id: "{{ fakeit_UUID }}"

cookies:
    - rate_limit_bucket: "{{ fakeit_UUID }}"
    - client_tier: ["tier1", "tier2", "tier3"]
```

### JSON Configuration

```json
{
    "method": "GET",
    "url": "https://api.rate-limited.com/endpoint",
    "yes": true,
    "timeout": "3s",
    "dodos": 20,
    "requests": 2000,
    "duration": "120s",
    "skip_verify": false,
    "headers": [
        { "X-API-Key": "{{ fakeit_UUID }}" },
        { "X-Client-ID": "{{ fakeit_UUID }}" },
        { "X-Rate-Limit-Test": "true" },
        { "User-Agent": "{{ fakeit_UserAgent }}" }
    ],
    "params": [
        { "burst_test": "true" },
        { "client_type": ["premium", "standard", "free"] },
        { "request_id": "{{ fakeit_UUID }}" }
    ],
    "cookies": [
        { "rate_limit_bucket": "{{ fakeit_UUID }}" },
        { "client_tier": ["tier1", "tier2", "tier3"] }
    ]
}
```

---

## Notes

- All examples use template functions for dynamic data generation
- Adjust `dodos`, `requests`, `duration`, and `timeout` values based on your testing requirements
- Use `skip_verify: true` for testing with self-signed certificates
- Set `yes: true` to skip confirmation prompts in automated testing
- Template functions like `{{ fakeit_* }}` generate random realistic data for each request
- Multiple values in arrays (e.g., `["value1", "value2"]`) will be randomly selected per request
- Use the `body_FormData` function for multipart form uploads
- Proxy configurations support HTTP, SOCKS5, and SOCKS5H protocols

For more template functions and advanced configuration options, refer to the main documentation and `utils/templates.go`.
