# Traefik Require Auth Headers Plugin

A Traefik middleware plugin that validates the Authorization header prefix against a list of required headers. This plugin is useful for enforcing specific authentication schemes in your API gateway.

## Features

- Validates the presence of the `Authorization` header
- Checks if the authorization prefix matches the configured required headers
- Customizable error responses with support for:
  - Custom HTTP status code
  - Custom headers
  - Custom response body
  - Multiple content types (JSON/Plain text)

## Configuration

### Static Configuration

To enable the plugin in your Traefik instance:

```yaml
experimental:
  plugins:
    require-auth-headers:
      moduleName: "github.com/unsoon/traefik-require-auth-headers"
      version: "v1.0.0"
```

### Dynamic Configuration

```yaml
http:
  middlewares:
    my-require-auth-headers:
      plugin:
        require-auth-headers:
          requiredHeaders:
            - "Bearer"
            - "Basic"
          errorResponse:
            statusCode: 401
            contentType: "application/json"
            headers:
              X-Custom-Header: "custom value"
            body:
              statusCode: 401
              message: "Unauthorized"
              description: "You are not authorized to access this resource"
```

### Configuration Options

| Option                      | Type                | Required | Default      | Description                                                         |
| --------------------------- | ------------------- | -------- | ------------ | ------------------------------------------------------------------- |
| `requiredHeaders`           | `[]string`          | Yes      | `[]`         | List of accepted Authorization header prefixes                      |
| `errorResponse.statusCode`  | `int`               | No       | `401`        | HTTP status code for error responses                                |
| `errorResponse.contentType` | `string`            | No       | `text/plain` | Content type of error response (`text/plain` or `application/json`) |
| `errorResponse.headers`     | `map[string]string` | No       | `{}`         | Additional headers to include in error response                     |
| `errorResponse.body`        | `interface{}`       | No       | `nil`        | Custom response body                                                |

## How It Works

1. The plugin intercepts incoming HTTP requests
2. Checks for the presence of the `Authorization` header
3. Extracts the authorization scheme prefix (e.g., "Bearer" from "Bearer xyz123")
4. Validates the prefix against the list of `requiredHeaders`
5. If validation fails:
   - Returns configured error response
   - Applies custom headers
   - Sets specified status code
   - Returns configured body with appropriate content type
6. If validation succeeds:
   - Passes request to next middleware/handler

## Example Usage

### Basic Authentication Check

```yaml
http:
  middlewares:
    basic-auth-check:
      plugin:
        require-auth-headers:
          requiredHeaders: ["Basic"]
          errorResponse:
            statusCode: 401
            contentType: "application/json"
            body:
              error: "Basic authentication required"
```

### Multiple Auth Schemes

```yaml
http:
  middlewares:
    multi-auth-check:
      plugin:
        require-auth-headers:
          requiredHeaders: ["Bearer", "Basic"]
          errorResponse:
            statusCode: 401
            contentType: "application/json"
            headers:
              WWW-Authenticate: 'Bearer realm="example"'
            body:
              message: "Require either Bearer or Basic authentication"
```

## Error Response Examples

### JSON Response

```yaml
errorResponse:
  statusCode: 401
  contentType: "application/json"
  body:
    error: "Unauthorized"
    message: "Invalid authentication scheme"
    code: 401
```

### Plain Text Response

```yaml
errorResponse:
  statusCode: 401
  contentType: "text/plain"
  body: "Authentication required"
```

## License

This plugin is distributed under the [MIT License](LICENSE).
