# Dida365 MCP Server

[English](README.md) | 简体中文

这是一个基于[滴答清单官方 API](https://developer.dida365.com/api#/openapi) 开发的 MCP Server 实现，用于连接 [滴答清单](https://dida365.com) 服务。

## ✨ 功能特点

- 🔐 基于官方 OAuth 2.0 认证机制
- 🎯 提供标准的 MCP Server 接口
- 💻 支持 STDIO 方式的 MCP Server 实现

## 📖 使用指南

### 1. 获取访问令牌（Access Token）

滴答清单采用 OAuth 2.0 授权机制来管理 API 访问权限。你可以通过以下两种方式获取访问令牌：

#### 方式一：使用在线服务（推荐）

1. 打开浏览器访问 <https://dida365.dcjanus.com/oauth/login>
2. 系统会自动重定向到滴答清单的授权页面
3. 在滴答清单页面确认授权
4. 授权成功后，系统会自动重定向回 `/oauth/callback` 端点
5. 你将收到一个包含 `access_token` 的 JSON 响应，请妥善保存该令牌

#### 方式二：本地部署

1. 访问滴答清单的[开发者中心](https://developer.dida365.com/manage)，创建新应用
2. 在应用设置中，将 `Redirect URL` 配置为 `http://localhost:8080/oauth/callback`
3. 启动本地服务器（默认端口：8080）
4. 访问 <http://localhost:8080/oauth/login>
5. 按照页面提示完成授权流程
6. 获取访问令牌

### 2. 配置 MCP Server

在配置文件中添加以下内容：

```json
{
    "mcpServers": {
        "dida365": {
            "command": "docker",
            "args": [
                "run",
                "-i",
                "--rm",
                "--init",
                "--pull",
                "always",
                "ghcr.io/dcjanus/dida365-mcp-server:latest",
                "dida365-mcp-server",
                "-access_token",
                "<YOUR_ACCESS_TOKEN>"
            ]
        }
    }
}
```

请将 `<YOUR_ACCESS_TOKEN>` 替换为你获取的访问令牌。

## ⚠️ 注意事项

- 请妥善保管访问令牌，避免泄露给他人
- 如遇授权失败，请检查网络连接或重新进行授权
