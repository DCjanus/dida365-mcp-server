# Dida365 MCP Server

English | [简体中文](README.zh.md)

This is a MCP Server implementation for [Dida365](https://dida365.com) using the [Dida365 Official API](https://developer.dida365.com/api#/openapi).

## ✨ Features

- 🔐 Official OAuth 2.0 authentication mechanism
- 🎯 Standard MCP Server interface
- 💻 STDIO-based MCP Server implementation

## 📖 Usage Guide

### 1. Obtain Access Token

Dida365 uses OAuth 2.0 authorization mechanism to manage API access. You can obtain an access token through either of the following methods:

#### Option 1: Using Online Service (Recommended)

1. Visit <https://dida365.dcjanus.com/oauth/login> in your browser
2. The system will automatically redirect you to the Dida365 authorization page
3. Confirm the authorization on the Dida365 page
4. After successful authorization, you'll be redirected back to the `/oauth/callback` endpoint
5. You will receive a JSON response containing the `access_token` - store it securely

#### Option 2: Local Deployment

1. Visit Dida365's [Developer Center](https://developer.dida365.com/manage) to create a new application
2. In the application settings, configure the `Redirect URL` as `http://localhost:8080/oauth/callback`
3. Start the local server (default port: 8080)
4. Visit <http://localhost:8080/oauth/login>
5. Follow the on-screen instructions to complete the authorization process
6. Obtain the access token

### 2. Configure MCP Server

Add the following to your configuration file:

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

Replace `<YOUR_ACCESS_TOKEN>` with the access token you obtained.

## ⚠️ Important Notes

- Keep your access token secure and never share it with others
- If authorization fails, check your network connection or try reauthorizing
