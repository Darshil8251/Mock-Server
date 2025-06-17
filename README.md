# Mock Server

A flexible mock server written in Go that supports multiple endpoints, pagination, headers, query parameters, and dynamic response modification.

## Features

- Multiple endpoint support
- YAML-based configuration
- Endpoint pagination
- Custom headers
- Query parameters
- Request body validation
- Expected response configuration
- Dynamic response modification based on pagination

## Installation

1. Make sure you have Go 1.21 or later installed
2. Clone this repository
3. Run `go mod download` to install dependencies

## Usage

1. Create a YAML configuration file (see `example-config.yaml` for reference)
2. Run the server:
   ```bash
   go run main.go -config your-config.yaml -port 8080
   ```

## Configuration Format

The configuration file should be in YAML format with the following structure:

```yaml
endpoints:
  - path: "/your/endpoint"
    method: "GET"  # HTTP method
    headers:       # Request headers
      Content-Type: "application/json"
    queryParams:   # Query parameters
      param1: "string"
    response:      # Response configuration
      status: 200
      headers:     # Response headers
        X-Custom: "value"
      body:        # Response body
        data: "value"
    pagination:    # Pagination configuration
      enabled: true
      pageSize: 10
      totalItems: 100
      pageParam: "page"
      sizeParam: "size"
    responseParams:  # Dynamic response parameters
      dynamicFields:  # Fields to modify per page
        - "id"
        - "value"
      incrementBy: 1  # Increment value for dynamic fields
```

## Example

See `example-config.yaml` for a complete example of the configuration format.

## API Endpoints

The server will create endpoints based on your configuration. For example, with the example configuration:

- GET `/api/users` - Returns a paginated list of users
- GET `/api/products` - Returns a paginated list of products

## Pagination

When pagination is enabled, the response will include pagination metadata:

```json
{
  "data": [...],
  "pagination": {
    "currentPage": 1,
    "pageSize": 10,
    "totalPages": 10,
    "totalItems": 100
  }
}
```

## License

MIT 