# Mock Server

Mock Server help to test cloud-to-cloud integrations, connector, and addon pagination part.

## Features

- Multiple endpoint support
- JSON-based configuration
- Page, offset, and link-based pagination, are supported
- Custom headers and query parameters
- Request body validation
- Dynamic response configuration via external JSON files

## Installation

1. Make sure you have Go 1.21 or later installed
2. Clone this repository
3. Run `go mod download` to install dependencies

## Usage

1. Create a JSON configuration file (see below for an example).
2. Set the environment variables:
   - `CONFIG_FILE_PATH` to the path of your config file (e.g., `export CONFIG_FILE_PATH=./config.json`)
   - `PORT` to the port you want the server to run on (e.g., `export PORT=8080`)
3. Run the server:
   ```bash
   go run cmd/main.go
   ```

## Configuration Format

The configuration file should be in JSON format with the following structure:

```json
{
  "endpoints": [
    {
      "path": "/api/threats",
      "method": "GET",
      "headers": {},
      "queryParams": {},
      "requestBody": {},
      "rateLimit": 0,
      "pagination": {
        "type": "page",
        "location": "query",
        "options": {
          "pageKey": "page",
          "pageSizeKey": "pageSize",
          "pageSize": 10,
          "pageCount": 5,
          "totalRecordCount": 50
        }
      },
      "responseObjFilePath": "response/threatResponse.json",
      "responseField": "threats"
    }
  ]
}
```

- `path`: Provide the path/endpoint of API.
- `method`: Provide the type of API(e.g GET,DELETE,PUT,POST,etc).
- `header`: Enter the supported header parameter by API.
- `queryParams`: Enter API supported Query Parameters.
- `requestBody`: Provide the request body for the API.
- `pagination.type`: Pagination type. Supported: `page`, `offset`, `link`,`token`
- `pagination.location`: Enter the location of the pagination parameter in req.(e.g header,body,query).
- `options`: Provide the pagination parameters.

  - `pageKey`: Share the page key use by the vendor API, default will be page.e.g page,pageNo,pageCount,etc.
  - `pageSizeKey`: Provide the vendor supported page size key, default is pageSize. e.g limit,size,etc.
  - `pageSize`: Enter the no of records per page, default is 100.
  - `totalPage`: Enter the total no of page you want to fetch, default will be 2.
  - `totalRecord`: Provide the no of records you want to fetch, default will be 200.
  - `linkKey`: Provide the link field in present in response object, default will be link.
  - `tokenKet`: Share the token field name present in response object, applicable for only token base pagination, default will be token.

- `responseObjFilePath`: Path to a JSON file containing the response object template.
- `responseField`: The key in the response object that is an array and will be paginated.
