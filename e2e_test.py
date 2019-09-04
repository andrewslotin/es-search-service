#!/usr/bin/env python3

import http.client, json, time

ELASTICSEARCH_URL = "localhost:9200"
SERVICE_URL = "localhost:8080"
DATA = ('{"index": {"_index": "products", "_type": "product"}}\n'
        '{"title": "AirMax", "brand": "Nike", "price": 1000, "stock": 10}\n'
        '{"index": {"_index": "products", "_type": "product"}}\n'
        '{"title": "SuperStar", "brand": "Adidas", "price": 999, "stock": 5}\n')

def seed(data):
    print("Uploading seed data")
    
    conn = http.client.HTTPConnection(ELASTICSEARCH_URL)
    conn.request("POST", "/_bulk", data, {"Content-Type": "application/x-ndjson"})

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
    path = "/v1/products?q=" + query
    conn = http.client.HTTPConnection(SERVICE_URL)
    conn.request("GET", path)

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
seed(DATA)
time.sleep(2) # give ES a chance to index

try:
    if not test(
        "Find 'Nike'",
        "Nike",
        '{"status":"success","results":[{"title":"AirMax","brand":"Nike","price":1000,"stock":10}]}'
    ):
        failed += 1

    if not test(
        "Find 'SuperStar'",
        "SuperStar",
        '{"status":"success","results":[{"title":"SuperStar","brand":"Adidas","price":999,"stock":5}]}'
    ):
        failed += 1

    if not test(
        "Don't find 'Puma'",
        "Puma",
        '{"status":"success","results":[]}'
    ):
        failed += 1
finally:
    clear()

exit(failed == 0 and 0 or 1)
