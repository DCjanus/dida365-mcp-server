# Dida365 MCP Server

[English](README.md) | 简体中文

这是一个[滴答清单](https://dida365.com)的 MCP Server 实现。使用的是[滴答清单官方 API](https://developer.dida365.com/api#/openapi)。

## 功能特点

- 基于官方 OAuth 2.0 认证
- 提供标准的 MCP Server 接口
- 暂时只提供 STDIO 的 MCP Server 实现

## 使用说明

### 1. 获取访问令牌（Access Token）

滴答清单使用 OAuth 2.0 授权机制来管理 API 访问权限。你可以通过以下两种方式获取访问令牌：

#### 方式一：使用在线服务

访问我们提供的在线 OAuth 服务：
1. 打开浏览器访问 `https://oauth.dida365.dcjanus.com/oauth/login`
2. 系统会自动重定向到滴答清单的授权页面
3. 在滴答清单页面确认授权
4. 授权成功后，系统会自动重定向回 `/oauth/callback` 端点，我们的服务会返回一个 JSON 格式的响应，其中包含 `access_token` 字段，请妥善保存以供后续使用

#### 方式二：本地部署

你也可以选择在本地运行 OAuth 服务：
1. 前往滴答清单的[开发者中心](https://developer.dida365.com/manage)，创建一个新应用
2. 在上一步创建的应用中，设置 `Redirect URL` 为 `http://localhost:8080/oauth/callback`
3. 启动本地服务器，默认监听 8080 端口
4. 访问 `http://localhost:8080/oauth/login`
5. 按照页面提示完成授权流程
6. 成功后将获得访问令牌

### 2. 配置 MCP Server

以下是一个标准的配置文件示例：

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

请将 `<YOUR_ACCESS_TOKEN>` 替换为你在上一步获取的访问令牌。

## 注意事项

- 访问令牌（Access Token）请妥善保管，不要泄露给他人
- 如果遇到授权失败，请检查网络连接或重新进行授权

