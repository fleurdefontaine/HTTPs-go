# HTTPS-Golang for GTPS

This project is a Go-based HTTP server designed to handle GTPS server requests securely using TLS encryption. The server processes incoming HTTP requests, enforces rate-limiting, checks for valid user agents, blocks specific IP addresses, and handles static file requests with a fallback to a CDN if the file is not found locally.

## Features

- **Rate Limiting**: Enforces a rate limit for each IP address to prevent abuse.
- **IP Blocking**: Blocks requests from certain IP address prefixes.
- **User Agent Validation**: Only allows specific user agents to access certain endpoints.
- **Static File Handling**: Serves files from the local server or redirects to a CDN if the file is not found locally.
- **TLS Encryption**: Uses TLS 1.2 or higher for secure communication.

## Requirements

- Go 1.18 or higher

## Setup

1. Clone the repository:

    ```bash
    git clone https://github.com/fleurdefontaine/HTTPs-go.git
    cd HTTPs-go
    ```

2. Install Go dependencies:

    ```bash
    go mod tidy
    ```

3. Create the necessary directories for configuration files and SSL certificates:

    - `config/config.json` for the server configuration.
    - `config/SSL/` for the TLS certificates (`server.crt` and `server.key`).

    Example `config.json`:

    ```json
    {
      "ip": "127.0.0.1",
      "port": 17091,
      "loginurl": "your-login-url.com",
      "ratelimit": 100,
      "cdn": "cdn-path"
    }
    ```

4. Ensure you have valid TLS certificates in the `config/SSL/` folder.

## Running the Server

To run the server, use the following command:

```bash
go run main.go
```

You can also specify a custom port by passing it as an argument:

```bash
go run main.go 8080
```

By default, the server will listen on port `443`.

## Endpoints

- **/growtopia/server_data.php**
  - Method: `POST`
  - Requires a valid `User-Agent` header from the allowed list.

- **/cache/{file}**
  - Method: `GET`
  - Requires a valid `User-Agent` header from the allowed list.

## License

This project is licensed under the AGPL-3.0 - see the [LICENSE](LICENSE) file for details.