# Dida365 MCP Server

This is a MCP Server implementation for [Dida365](https://dida365.com) using the [Dida365 Official API](https://developer.dida365.com/api#/openapi).

## Usage

1. After starting the server, visit the `/oauth/login` endpoint in your browser
2. The system will automatically redirect you to the Dida365 authorization page
3. Confirm the authorization on the Dida365 authorization page
4. After authorization, the system will automatically redirect back to the `/oauth/callback` endpoint
5. Upon successful authorization, you will receive a token - please store it securely for future API calls

Note:
- Make sure you have registered your application on the Dida365 developer platform and obtained the Client ID and Client Secret
- The callback URL must match the one configured in your Dida365 developer platform settings 