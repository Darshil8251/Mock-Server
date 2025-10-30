# Mock Server

Mock Server helps test cloud-to-cloud integrations, connectors, and addon pagination.

## Features

- Multiple endpoint support
- JSON-based configuration
- Support for page, offset, link, and token based pagination
- Dynamic response configuration via external JSON files

## Installation

1. Make sure you have Go 1.21 or later installed.
2. Clone this repository.
3. Run `go mod download` to install dependencies.

## Usage

1. Create a JSON configuration file (see below for an example).
2. Set the environment variables:
   - `CONFIG_FILE_PATH`: Path to your config file (e.g., `export CONFIG_FILE_PATH=./config.json`)
   - `PORT`: Port you want the server to run on (e.g., `export PORT=8080`)
3. Run the server:
   ```bash
   go run cmd/main.go
   ```

## Configuration Format

The configuration file should be in JSON format with the following structure:

```json
{
    "path": "link/cpc-alerts",
    "method": "GET",
    "pagination": {
        "type": "link",
        "options": {
            "totalRecord": 15,
            "pageSize": 5,
            "linkKey": "next_url"
        }
    },
    "response": {
        "filePath": "response/cpc_alerts.json",
        "fieldName": "cpcAlerts"
    }
}
```

- `path`: The endpoint path of the API.
- `method`: The HTTP method for the endpoint (e.g., GET, DELETE, PUT, POST).
- `pagination.type`: Pagination type. Supported values: `page`, `offset`, `link`, `token`, `none`.
- `options`: Pagination parameters.
  - `pageSize`: Number of records per page. Default is 100.
  - `totalRecord`: Total number of records expected from the server. Default is 200.
  - `linkKey`: The link field name in the response object for link-based pagination.
  - `tokenKey`: The token field name in the response object, applicable only for token-based pagination. Default is `token`.
  - `offsetKey`: The offset field name in the response object for offset-based pagination.
- `response.filePath`: Path to a JSON file containing the response object template.
- `response.fieldName`: The key in the response object that is an array and will be paginated.
