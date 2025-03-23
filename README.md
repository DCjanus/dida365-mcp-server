# Dida365 MCP Server

English | [简体中文](README.zh.md)

This is a MCP Server implementation for [Dida365](https://dida365.com) using the [Dida365 Official API](https://developer.dida365.com/api#/openapi).

## Features

- Official OAuth 2.0 authentication
- Standard MCP Server interface
- Support for both local and remote deployment

## Usage Guide

### 1. Obtain Access Token

Dida365 uses OAuth 2.0 authorization mechanism to manage API access. You can obtain an access token through either of the following methods:

#### Option 1: Using Online Service

Access our online OAuth service:
1. Visit `https://oauth.dida365.com/oauth/login` in your browser
2. The system will automatically redirect you to the Dida365 authorization page
3. Confirm the authorization on the Dida365 page
4. After successful authorization, you'll be redirected back to the `/oauth/callback` endpoint
5. You will receive an access token - store it securely for future use

#### Option 2: Local Deployment

Alternatively, you can run the OAuth service locally:
1. Start the local server (default port: 8080)
2. Visit `http://localhost:8080/oauth/login`
3. Follow the on-screen instructions to complete the authorization process
4. Upon success, you will receive an access token

### 2. Configure MCP Server

Here's a standard configuration file example:

```json
{
    "mcpServers": {
        "dida": {
            "command": "/path/to/dida365-mcp-server",
            "args": [
                "-access_token",
                "<YOUR_ACCESS_TOKEN>"
            ]
        }
    }
}
```

Replace `<YOUR_ACCESS_TOKEN>` with the access token you obtained in the previous step.

## Important Notes

- Keep your access token secure and never share it with others
- If authorization fails, check your network connection or try reauthorizing
- It's recommended to update your access token periodically for security 