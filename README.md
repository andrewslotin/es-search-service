Task
====

Create a REST API Search Service
Let’s say our company has e-commerce platform with selling fashion products. Each product has title, brand, price and stock. The platform should have search functionality.

Your task is to implement microservice to serve search requests.

The Service has to support:
* GET method to perform search queries. E.g. https://example.com/products?q=black shoes 
* Authentication
* Api should have versioning
* Pagination and sorting
* Filtering. E.g. https://example.com/products?q=black shoes&filter=brand:brand_name

Requirements:
* Language: Go
* Storage: ElasticSearch
* Service should be dockerized
* Automated tests

Nice to have:
* Lightweight UI to show that the service works

Please push the service to GitHub (your public/private repo). If you have any questions about the test please contact us.

