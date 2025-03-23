# Dida365 MCP Server

[English](README.md) | 简体中文

这是一个[滴答清单](https://dida365.com)的 MCP Server 实现。使用的是[滴答清单官方 API](https://developer.dida365.com/api#/openapi)。

## 使用流程

### 步骤 1: 获取滴答清单的 Access Token

滴答清单提供了 OAuth 授权机制，通过 OAuth 授权可以获取到用户的 Access Token。你可以本地运行一个 oauth 服务器，然后通过浏览器访问滴答清单的授权页面，授权后会重定向到你的 oauth 服务器，然后你就可以获取到用户的 Access Token，因为你可以设置 Redirect URL 为 localhost 的一个端口，所以可以本地运行。

也可以选择使用我已经搭建好的 oauth 服务器，只需要在浏览器中访问 `https://oauth.dida365.dcjanus.com/oauth/login`，遵循页面上的指引即可。

### 步骤 2：配置 MCP Server

以下是示例配置文件：

```json
{
    "mcpServers": {
        "dida": {
            "command": "/Users/dcjanus/Code/dida365-mcp-server/bin/dida365-mcp-server",
            "args": [
                "-access_token",
                "<ACCESS_TOKEN>"
            ]
        }
    }
}
```

