#!/usr/bin/env python3

import os, http.client, json, time
from base64 import b64encode

ELASTICSEARCH_URL = os.environ.get("ELASTICSEARCH_NODES", "localhost:9200")
SERVICE_URL = os.environ.get("LISTEN_ADDR", "localhost:8080")
DATA = ('{"index": {"_index": "products", "_type": "product"}}\n'
        '{"title": "AirMax", "brand": "Nike", "price": 1000, "stock": 10}\n'
        '{"index": {"_index": "products", "_type": "product"}}\n'
        '{"title": "Pegasus Shield", "brand": "Nike", "price": 1500, "stock": 12}\n'
        '{"index": {"_index": "products", "_type": "product"}}\n'
        '{"title": "Pegasus Shield", "brand": "Nike", "price": 2000, "stock": 20}\n'
        '{"index": {"_index": "products", "_type": "product"}}\n'
        '{"title": "Zoom", "brand": "Nike", "price": 2000, "stock": 1}\n'
        '{"index": {"_index": "products", "_type": "product"}}\n'
        '{"title": "SuperStar", "brand": "Adidas", "price": 999, "stock": 5}\n'
        
)
MAPPINGS = """
{
  "properties": {
    "title": {
      "type": "text",
      "fielddata": true
    },
    "brand": {
      "type": "text",
      "fielddata": true
    },
    "price": {
      "type": "long"
    },
    "stock": {
      "type": "long"
    }
  }
}
"""

def seed(data, mappings):
    print("Uploading seed data")
    
    conn = http.client.HTTPConnection(ELASTICSEARCH_URL)
    conn.request("POST", "/_bulk", data, headers={"Content-Type": "application/x-ndjson"})

    resp = conn.getresponse()
    if resp.status >= 400:
        raise BaseException("failed to seed cluster: " + str(resp.read()))

    conn.close()

    conn = http.client.HTTPConnection(ELASTICSEARCH_URL)
    conn.request("PUT", "/products/_mapping", mappings, headers={"Content-Type": "application/json"})

    resp = conn.getresponse()
    if resp.status >= 400:
        raise BaseException("failed to seed cluster: " + str(resp.read()))

    conn.close()

def clear():
    print("Removing seed data")
    
    conn = http.client.HTTPConnection(ELASTICSEARCH_URL)
    conn.request("DELETE", "/products")

    resp = conn.getresponse()
    if resp.status >= 400:
        raise BaseException("failed to clear cluster: " + str(resp.read()))

    conn.close()

def query_search_api(query):
    path = "/v1/products?" + query
    conn = http.client.HTTPConnection(SERVICE_URL)
    conn.request("GET", path, headers={
        "Authorization": "Basic " + b64encode(b"user1:password2").decode("ascii")
    })

    resp = conn.getresponse()
    result = resp.read()
    conn.close()

    return result

def compare_json(expected, actual):
    return normalize_json(expected) == normalize_json(actual)

def normalize_json(data):
    return json.dumps(json.loads(data), sort_keys = True)

def test(name, query, expected):
    print(name, "...", end=" ")
    result = query_search_api(query)

    if not compare_json(result, expected):
        print("failed")
        print("\t", "expected: ", expected)
        print("\t", "got:      ", result)
        return False

    print("succeeded")

    return True


failed = 0
seed(DATA, MAPPINGS)
time.sleep(2) # give ES a chance to index

try:
    if not test(
        "Find all 'Nike' with price 1500",
        "q=Nike&filter=price:1500",
        """
        {
          "status": "success",
          "results": [
            {
              "title": "Pegasus Shield",
              "brand": "Nike",
              "price": 1500,
              "stock": 12
            }
          ]
        }
        """
    ):
        failed += 1

    if not test(
        "Find all 'Nike Pegasus Shield' sorted by title and price (expencieve first), return 2nd page of results with one result per page",
        "q=Nike&sort=title:asc&sort=price:desc&filter=title:Pegasus+Shield&from=1&size=1",
        """
        {
          "status": "success",
          "results": [
            {
              "title": "Pegasus Shield",
              "brand": "Nike",
              "price": 1500,
              "stock": 12
            }
          ]
        }
        """
    ):
        failed += 1

    if not test(
        "Don't find 'Puma'",
        "q=Puma",
        '{"status":"success","results":[]}'
    ):
        failed += 1
finally:
    clear()

exit(failed == 0 and 0 or 1)
