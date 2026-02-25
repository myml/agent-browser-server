# Agent Browser SSH Server

在 Docker 中运行 agent-browser，提供 SSH 和 MCP 服务。

## 项目简介

本项目是一个 MCP (Model Context Protocol) 服务器，通过 Docker 容器提供远程 agent-browser 访问能力。容器内同时运行 SSH 服务，方便进行调试和管理。

## 功能特性

- 🐳 **Docker 容器化部署**：一键启动完整的 agent-browser 环境
- 🔌 **MCP 协议支持**：通过 MCP 协议暴露 agent-browser 工具
- 🖥️ **SSH 访问**：内置 Dropbear SSH 服务器，支持远程登录
- ⚡ **HTTP 流式传输**：支持通过 HTTP 流式传输 MCP 消息
- 🛠️ **完整依赖**：预装 agent-browser 及其所有依赖

## 快速开始

### 运行容器

```bash
docker run -d -p 8080:8080 wrj97/agent-browser-server:main
```

### SSH 连接

```bash
docker run -d -p 2222:22 -p 8080:8080 -e SSH_PASSWORD=password wrj97/agent-browser-server:main
ssh root@localhost -p 22222
# 密码为环境变量 SSH_PASSWORD 设置的值（默认为 root）
```

## MCP 工具说明

### agent-browser

执行 agent-browser 命令的工具。

**参数：**

- `args` (必需): 传递给 agent-browser 的命令行参数数组

**使用示例：**

首先查看帮助信息：

```json
{
  "args": ["--help"]
}
```

执行浏览器操作：

```json
{
  "args": ["navigate", "https://example.com"]
}
```

**返回结果：**

```json
{
  "status": "success",
  "exit_code": 0,
  "stdout": "命令输出内容",
  "stderr": "错误输出内容",
  "command": "agent-browser navigate https://example.com",
  "execution_time": "1.234s"
}
```

## 配置说明

### 环境变量

| 变量名         | 默认值 | 说明         |
| -------------- | ------ | ------------ |
| `SSH_PASSWORD` | `root` | SSH 登录密码 |

### 端口映射

| 端口 | 说明              |
| ---- | ----------------- |
| 22   | SSH 服务端口      |
| 8080 | MCP HTTP 服务端口 |

## MCP 客户端配置

### Claude Desktop 配置

在 Claude Desktop 的配置文件中添加：

```json
{
  "mcpServers": {
    "agent-browser": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

### 其他 MCP 客户端

连接到 `http://localhost:8080/mcp` 即可使用 MCP 服务。

## 开发

### 本地运行

```bash
# 安装依赖
go mod download

# 运行服务器
go run main.go
```

服务器将在 `http://127.0.0.1:8080/mcp` 启动。

## 技术栈

- **Go**: MCP 服务器实现
- **Node.js**: agent-browser 运行环境
- **Dropbear**: 轻量级 SSH 服务器
- **Docker**: 容器化部署

## 许可证

MIT License
