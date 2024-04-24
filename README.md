httpcheck
=========

httpcheck is a command-line tool for measuring HTTP performance. 

![screenshot](screenshot.png)

### Installation

```bash
go install github.com/ptrhng/httpcheck@latest
```

### Usage

Default:

```bash
$ httpcheck httpie.io/hello
```

Custom HTTP method, HTTP header, and JSON data:

```bash
$ httpcheck PUT pie.dev/put X-API-Token:123 name=John obj:='{"k": "v"}'
```

Sending form data:

```bash
$ httpcheck PUT pie.dev/put name=john --form 
```

Adding query parameters:

```bash
$ httpcheck PUT pie.dev/put q==search page==1
```

### Request Items

Request item can be used to specify HTTP header, query parameters, and data. Each item consists of a key, value, and separator.

- HTTP Header `key:value`
- Query Parameter `key==value`
- String Data Field `key=value`
- JSON Data field `key:=value` 

Example:

```bash
httpcheck PUT pie.dev/put \
    X-header:value \
    query_param==value \
    data_string=value \
    data_json:='{"date": "today"}'
```
