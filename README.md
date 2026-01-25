# Microservices + Pub/Sub Architecture
This project simulates an e-commerce platform. It does the following:
* PUBLIC ENDPOINTS
    * Account Sign Up
    * Account Login
    * Get All Products
    * Get Products By IDs
    * Get Product By ID
    * Search Product
* BUYER ENDPOINTs
    * Account Create Order
    * Account View Orders
* ADMIN ENDPOINTs
    * Create Orders

This project uses `Microservices + Pub/Sub architecture`, completely written in Go.

# Tech Stack
* Go
* gRPC
* Elastic Search
* RabbitMQ
* Docker

# Topology
Below is an overview of the whole program.

# API Endpoints
Below is documentation for the API Endpoints.

## Sign Up
Registers a new user and returns a jwt token.

URL: `POST /api/signup`

Body:
```bash
{
  "email": "user@example.com", // email must be valid because email service will send message to this email
  "password": "securepassword123"
}
```

## Login
Authenticates a user and returns a jwt token.

URL: `POST /api/login`

Body:
```bash
{
  "email": "user@example.com", // email must be valid because email service will send message to this email
  "password": "securepassword123"
}
```

## Get Products
Fetches a list of products from Elastic Search. Supports filtering by specific IDs.

URL: `GET /api/products`

Query Params: ids=id1,id2 (optional)

## Get Product By ID
Get Product by specific ID.

URL: `GET /api/products/{id}`

## Search Products
Searches the Elasticsearch index for products matching a query string.

URL: `GET /api/products/search`

Body:
```bash
{ 
    "query": "laptop" 
}
```

## Create Order
Create a new order. **Upon success**, a message is sent to *RabbitMQ* to notify the Email Service.

URL: `POST /api/order`

Body:
```bash
{
  "products": [
    { 
        "id": "prod_123", 
        "quantity": 2 
    },
    { 
        "id": "prod_345", 
        "quantity": 3
    },
  ]
}
```

## List Orders
Retrieves all orders associated with the authenticated account.

URL: `GET /api/order`

## Create Product
Adds a new product to the Elastic Search index.

URL: `POST /admin/products`

Body:
```bash
{
  "name": "Gaming Mouse",
  "description": "RGB high-dpi mouse",
  "price": 59.99
}
```