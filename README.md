# Api Gatekeeper

A simple API gateway and user management gateway.

## Use Cases

The API Gatekeeper is a application that sits between the public web requests and your backends (monoliths, microsservices, etc.). It focus on the following use cases:

- **User Management:** api-gatekeeper comes with a full featured user resource system, with scope based authentication and authorization.
- **API Gateway:** api-gatekeeper allows you to expose multiple HTTP backends on the same host using the included auth.

![](https://raw.githubusercontent.com/gustapinto/api-gatekeeper/main/docs/diagram-dark.drawio.png#gh-dark-mode-only)
![](https://raw.githubusercontent.com/gustapinto/api-gatekeeper/main/docs/diagram-light.drawio.png#gh-light-mode-only)

## Configuration

The configuration is done by a yaml file. This file path must be provided by the `-config=<path to yaml>` when running the application. The configuration properties are:
```yaml
api:                        # The application api configuration
  address: "localhost:3000" # The application listen address

database:                                                                                        # The application database configurations
  provider: "postgres"                                                                           # The database provider
  dsn: "postgresql://api-gatekeeper:api-gatekeeper@postgres:5432/api-gatekeeper?sslmode=disable" # The database connection dsn

backends:                         # The backends configurations
  - name: "ping-backend"          # The backend name
    host: "http://localhost:8080" # The backend host
    scopes:                       # (Optional) The authentication scopes required for every route in this backend
      - "ping-backend-scope"
    headers: # (Optional) The headers to be included in the request for every route in this backend
      Authorization: "Bearer foobar"
      X-Example-Header-backend: "example backend header"
    routes:                     # The backend routes
      - method: "GET"           # The route method
        backendPath: "/ping"    # The absolute path on the backend backend (your application)
        gatekeeperPath: "/ping" # The absolute path that will be exposed by the api-gatekeeper
        timeoutSeconds: 30      # (Optional) The timeout in seconds, set to 0 to dont timeout
        isPublic: false         # (Optional) If this route is public
        scopes:                 # (Optional) The authentication scopes required for this route, they will be added with the backend scopes
          - "ping-backend.get-ping-scope"
        headers: # (Optional) The headers to be included in the request for this route
          X-Example-Header-Route: "example route header"
```

## FAQ

### Is this application production ready?

In the moment **No**.