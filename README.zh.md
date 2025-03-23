# Dida365 MCP Server

[English](README.md) | 简体中文

这是一个[滴答清单](https://dida365.com)的 MCP Server 实现。使用的是[滴答清单官方 API](https://developer.dida365.com/api#/openapi)。

## 使用流程

1. 启动服务器后，在浏览器中访问 `/oauth/login` 接口
2. 系统会自动重定向到滴答清单的授权页面
3. 在滴答清单授权页面上确认授权
4. 授权完成后，系统会自动重定向回 `/oauth/callback` 接口
5. 成功授权后，您将获得一个 token，请妥善保存此 token 用于后续接口调用

注意：
- 请确保您已在滴答清单开发者平台注册了应用并获取了相应的 Client ID 和 Client Secret
- 回调地址需要与您在滴答清单开发者平台配置的回调地址一致