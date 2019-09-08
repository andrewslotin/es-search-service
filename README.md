es-search-service
=================

A service that abstracts the underlying Elasticsearch cluster by proxying search requests to it and
ensuring consistent JSON API.

Installation
------------

```bash
go get github.com/andrewslotin/es-search-service
```

Usage
-----

```bash
$GOPATH/bin/es-search-service --nodes=<comma-separated list of ES nodes>
```

This will launch an instance of search service accepting connections on `:8080` and using provided
nodes list as a storage. You can override the default listed address by either setting `LISTEN_ADDR`
env variable or by passing it via `-l=` flag.

Search service also provides a simple UI to build and run search queries. You can access it at `http://<listen addr>/`.

On startup the search service ensures that provided Elasticsearch cluster is reachable and in green state.

To allow service to wait until the ES cluster boots up provide connection timeout either via
`ELASTICSEARCH_CONN_TIMEOUT` env variable or by passing a duration value with `--timeout=` flag.

To learn about all possible configuration options, run `es-search-service --help`.

### Using Docker

`es-search-service` provides a `Dockerfile` allowing to run it as a Docker container.

```bash
cd <project directory>
docker build -t es-search-service .
docker run --rm -p 8080:8080 -e "LISTEN_ADDR=:8080" -e "ELASTICSEARCH_NODES=<nodes>" es-search-service
```

This would build and run an image of the service connected to Elasticsearch cluster `<nodes>`
and listening on port `8080`

Search API
----------

```
GET /v1/products?q=<query>
Authorization: Basic <credentials>
```

The Search API requires Basic authentication. It does not perform any kind of authorization, so
any login/password pair will work.

### Example responses

**Success**
```javascript
{
    "status": "success",
    "results": [
        // ... an array of matching documents
    ]
}
```

**Error**
```javascript
{
    "status": "error",
    "code": 400,                         // HTTP status code
    "message": "missing query parameter" // error details
}
```

### Pagination

To enable pagination of search results, add `from` and `size` parameters to your query
with the number of documents to skip and the page size respectively.

```
GET /v1/products?q=<query>&from=10&size=5
Authorization: Basic <credentials>
```

### Sorting

The sorting order for results can be provided by passing the sort field name followed by a colon and
search direction (`asc`/`desc`). You can provide more that one sort field by sending multiple `sort`
parameters. In this case the sort fields are used in the same order as they are specified in query.

```
GET /v1/products?q=<query>&sort=title:asc&sort=price:desc
Authorization: Basic <credentials>
```

### Filtering

To filter the search results based on certain field values provide the filtering query in the `filter`
request parameter. The filter query has to be in [Lucene syntax](https://lucene.apache.org/core/2_9_4/queryparsersyntax.html)
and will be sent to the Elasticsearch as-is.

```
GET /v1/products?q=<query>&filter=price:1500
Authorization: Basic <credentials>
```

Testing
-------

```bash
go test github.com/andrewslotin/es-search-service/...
```

This will run all unit tests bundled with the service source code.

```bash
cd <project directory>
docker-compose up
python3 e2e_test.py
```

This will launch an Elasticsearch cluster container and a container running the search service
listening on `localhost:8080`, and the end-to-end test suite against it. To run the test suite
against an already running service instance provide `ELASTICSEARCH_NODES` list and service `LISTEN_ADDR`
via env variables before calling `python3 e2e_test.py`.
