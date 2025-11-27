# SAI-IM - Secure & Advanced Instant Messaging System

<p align="center">
  <strong>基于微服务架构的高性能即时通讯系统</strong>
</p>

<p align="center">
  <a href="#核心特性">核心特性</a> •
  <a href="#技术架构">技术架构</a> •
  <a href="#项目结构">项目结构</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#详细设计">详细设计</a>
</p>

---

## 📖 项目概述

SAI-IM 是一个基于微服务架构的企业级即时通讯系统，采用 Go 语言开发，使用 WebSocket 实现实时双向通信，通过 gRPC 进行服务间调用。系统具备高性能、高可用、高安全性等特点，适用于企业内部通讯、社交应用等场景。

### ✨核心特性

- **🚀 高性能架构**
  - 支持百万级并发长连接
  - 采用 bitmap 算法优化消息已读未读状态管理
  - 基于 MongoDB 的高效消息存储方案
  - 消息队列异步处理，降低系统耦合度

- **🔒 安全机制**
  - JWT Token 认证，保证身份合法性
  - bcrypt 密码hash加密，防止密码泄露
  - 连接独占性，防止重复登录
  - MongoDB 持久化存储聊天记录
  - Bitmap 高效管理消息已读未读状态

- **💪 可靠性保障**
  - 参考 TCP 三次握手的 ACK 确认机制
  - 完善的心跳检测机制，实时监控连接状态
  - 消息序列号管理，确保消息顺序与可靠传输
  - 离线消息推送，保证消息不丢失

- **⚡ 微服务架构**
  - 基于 go-zero 框架的微服务设计
  - 使用 etcd 实现服务注册与发现
  - Apisix 网关统一路由与鉴权
  - 支持服务水平扩展与负载均衡

- **📊 完善的监控体系**
  - Jaeger 分布式链路追踪
  - ELK 日志收集与分析 (Elasticsearch + Logstash + Kibana)
  - Sail 配置中心支持热重载
  - 优雅重启机制

### 🔧技术栈

| 技术选型  | 说明                                               |
| --------- | -------------------------------------------------- |
| go-zero   | 微服务框架                                         |
| gRPC      | RPC 通信协议                                       |
| WebSocket | 双向通信协议                                       |
| MongoDB   | 消息记录存储                                       |
| MySQL     | 关系型数据存储                                     |
| Redis     | 缓存与会话管理                                     |
| Kafka     | 消息队列                                           |
| etcd      | 服务注册与配置中心                                 |
| Apisix    | API 网关                                           |
| SLF4J     | 日志框架                                           |
| Jaeger    | 分布式链路追踪                                     |
| ELK Stack | 日志收集与分析 (Elasticsearch + Logstash + Kibana) |

---

## 🏗️ 系统架构

### 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                          客户端层                                  │
│              Web Client / Mobile App / Desktop App              │
└─────────────────────────────────────────────────────────────────┘
                              ↓ HTTP/WebSocket
┌─────────────────────────────────────────────────────────────────┐
│                        API 网关层 (Apisix)                        │
│          路由 | 鉴权 | 限流 | 负载均衡 | 熔断降级                 │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                          微服务层                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ IM服务   │  │ User服务 │  │Social服务│  │ Task服务 │        │
│  ├──────────┤  ├──────────┤  ├──────────┤  ├──────────┤        │
│  │ • API    │  │ • API    │  │ • API    │  │ • Worker │        │
│  │ • RPC    │  │ • RPC    │  │ • RPC    │  │          │        │
│  │ • WS     │  │          │  │          │  │          │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────────────────┘
           ↓ gRPC              ↓ Kafka           ↓
┌─────────────────────────────────────────────────────────────────┐
│                          中间件层                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │  etcd    │  │  Kafka   │  │  Redis   │  │   ...    │        │
│  │ (注册发现)│  │ (消息队列)│  │  (缓存)  │  │          │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                          存储层                                    │
│  ┌──────────┐  ┌──────────┐                                     │
│  │  MySQL   │  │ MongoDB  │                                     │
│  │(业务数据) │  │(消息记录) │                                     │
│  └──────────┘  └──────────┘                                     │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                       监控与运维层                                 │
│  Prometheus | Grafana | ELK | Jaeger                           │
└─────────────────────────────────────────────────────────────────┘
```

### 🎁服务模块说明

#### 1. IM 服务 (即时通讯核心)
- **API 层**: 提供 HTTP REST API，处理消息查询、会话管理等请求
- **RPC 层**: 提供 gRPC 服务，供其他微服务调用推送消息
- **WebSocket 层**: 维护客户端长连接，实现实时消息推送

#### 2. User 服务 (用户管理)
- **API 层**: 用户注册、登录、信息管理等接口
- **RPC 层**: 提供用户信息查询、认证等服务

#### 3. Social 服务 (社交关系)
- **API 层**: 好友关系、群组管理等接口
- **RPC 层**: 社交关系查询、验证等服务

#### 4. Task 服务 (异步任务)
- 处理离线消息推送
- 消息持久化
- 统计分析等异步任务

---

## 📁 项目结构

```
SAI-IM/
├── apps/                          # 微服务应用目录
│   ├── im/                        # IM 即时通讯服务
│   │   ├── api/                   # HTTP API 服务
│   │   │   ├── internal/          # 内部实现
│   │   │   │   ├── config/        # 配置定义
│   │   │   │   ├── handler/       # 请求处理器
│   │   │   │   ├── logic/         # 业务逻辑
│   │   │   │   ├── svc/           # 服务上下文
│   │   │   │   └── types/         # 类型定义
│   │   │   └── im.api             # API 定义文件
│   │   ├── rpc/                   # gRPC 服务
│   │   │   ├── internal/          # 内部实现
│   │   │   │   ├── config/        # 配置定义
│   │   │   │   ├── logic/         # 业务逻辑
│   │   │   │   ├── server/        # gRPC 服务实现
│   │   │   │   └── svc/           # 服务上下文
│   │   │   ├── pb/                # Protobuf 生成文件
│   │   │   └── im.proto           # Protobuf 定义
│   │   ├── ws/                    # WebSocket 服务
│   │   │   ├── internal/          # 内部实现
│   │   │   │   ├── config/        # 配置定义
│   │   │   │   ├── handler/       # WebSocket 处理器
│   │   │   │   └── svc/           # 服务上下文
│   │   │   └── ws.go              # WebSocket 服务入口
│   │   ├── immodels/              # MongoDB 模型
│   │   └── exec.sh                # 启动脚本
│   │
│   ├── user/                      # 用户服务
│   │   ├── api/                   # HTTP API 服务
│   │   ├── rpc/                   # gRPC 服务
│   │   ├── models/                # MySQL 数据模型
│   │   └── exec.sh                # 启动脚本
│   │
│   ├── social/                    # 社交服务
│   │   ├── api/                   # HTTP API 服务
│   │   ├── rpc/                   # gRPC 服务
│   │   ├── socialmodels/          # 数据模型
│   │   └── exec.sh                # 启动脚本
│   │
│   └── task/                      # 异步任务服务
│       └── mq/                    # 消息队列消费者
│
├── pkg/                           # 公共包/工具库
│   ├── bitmap/                    # Bitmap 已读未读算法
│   ├── configserver/              # 配置中心客户端
│   ├── constants/                 # 常量定义
│   ├── ctxdata/                   # 上下文数据工具
│   ├── encrypt/                   # 加密工具(抗量子算法)
│   ├── interceptor/               # gRPC 拦截器
│   ├── middleware/                # HTTP 中间件
│   ├── resultx/                   # 统一响应格式
│   ├── retry/                     # 重试机制
│   ├── wuid/                      # 分布式唯一 ID 生成
│   └── xerr/                      # 错误码定义
│
├── components/                    # 中间件组件配置
│   ├── apisix/                    # API 网关配置
│   ├── apisix-dashboard/          # 网关管理面板
│   ├── prometheus/                # 监控配置
│   ├── grafana/                   # 可视化配置
│   ├── filebeat/                  # 日志采集
│   ├── kibana/                    # 日志查询
│   ├── logstash/                  # 日志处理
│   └── sail/                      # 配置中心
│
├── deploy/                        # 部署相关
│   ├── dockerfile/                # 各服务 Dockerfile
│   ├── mk/                        # Makefile 片段
│   ├── script/                    # 部署脚本
│   └── sql/                       # 数据库初始化脚本
│
├── docker-compose.yaml            # Docker Compose 编排
├── Makefile                       # 构建脚本
├── go.mod                         # Go 模块定义
└── README.md                      # 项目说明文档
```

---

## 🚀 快速开始

### 环境要求

- Go 1.24+
- Docker & Docker Compose
- MySQL 8.0+
- MongoDB 4.4+
- Redis 6.0+
- Kafka 2.8+
- etcd 3.5+

### 快速部署

#### 1. 克隆项目

```bash
git clone https://github.com/your-repo/SAI-IM.git
cd SAI-IM
```

#### 2. 启动基础设施

```bash
# 启动所有中间件服务 (MySQL, MongoDB, Redis, Kafka, etcd, Apisix 等)
docker-compose up -d
```

#### 3. 初始化数据库

```bash
# 执行 SQL 初始化脚本
mysql -h 127.0.0.1 -u root -p < deploy/sql/user.sql
mysql -h 127.0.0.1 -u root -p < deploy/sql/social.sql
mysql -h 127.0.0.1 -u root -p < deploy/sql/im.sql
```

#### 4. 编译服务

```bash
# 使用 Makefile 编译所有服务
make build

# 或单独编译某个服务
make build-im-api
make build-im-rpc
make build-im-ws
```

#### 5. 启动服务

```bash
# 启动 User 服务
cd apps/user
./exec.sh

# 启动 Social 服务
cd apps/social
./exec.sh

# 启动 IM 服务
cd apps/im
./exec.sh
```

#### 6. 访问服务

- **API 网关 (Apisix)**: http://localhost:9080
- **Apisix Dashboard**: http://localhost:9000
- **Jaeger UI**: http://localhost:16686
- **Kibana**: http://localhost:5601
- **Sail 配置中心**: http://localhost:8108
- **Elasticsearch**: http://localhost:9200

---

## 📚 核心设计与实现

### 1. 业务功能设计

#### IM 服务核心功能

**对外服务**
- HTTP REST API: 消息记录查询、会话管理、用户状态查询
- WebSocket 长连接: 实时消息推送、在线状态维护

**对内服务**
- gRPC 接口: 供其他微服务调用推送消息
- 消息队列消费: 处理异步消息推送任务

**核心业务场景**
- ✅ 单聊 (Private Chat): 一对一实时通讯
- ✅ 群聊 (Group Chat): 多人实时通讯
- ✅ 消息已读未读: 基于 Bitmap 算法的高效状态管理
- ✅ 用户在线离线: 实时状态检测与更新
- ✅ 历史消息: 基于时序号的分页查询
- ✅ 离线消息: Kafka 异步推送保证消息必达

### 2. 数据库设计

#### 存储架构

采用 **MySQL + MongoDB + Redis** 混合存储架构:

| 存储层 | 用途 | 说明 |
|--------|------|------|
| **MySQL** | 业务数据 | 用户信息、社交关系、群组信息等结构化数据 |
| **MongoDB** | 消息记录 | 聊天记录、离线消息等海量非结构化数据 |
| **Redis** | 缓存层 | 在线状态、会话列表、消息已读未读 Bitmap |

#### 核心数据表

**MySQL 表结构**

```sql
-- 用户表
user (
    id, username, password, nickname, avatar, 
    email, phone, status, created_at, updated_at
)

-- 社交关系表  
social_relation (
    id, user_id, friend_id, group_id, 
    relation_type, status, created_at
)

-- 群组表
group_info (
    id, name, avatar, owner_id, member_count,
    created_at, updated_at
)
```

**MongoDB 集合**

```go
// 聊天记录集合 - ChatLog
type ChatLog struct {
    ID             primitive.ObjectID  // MongoDB 文档ID
    ConversationId string              // 会话ID (基于 sendId+recvId 计算)
    SendId         string              // 发送者ID
    RecvId         string              // 接收者ID
    MsgFrom        int                 // 消息来源
    ChatType       constants.ChatType  // 聊天类型 (单聊/群聊)
    MsgType        constants.MType     // 消息类型 (文本/图片/文件等)
    MsgContent     string              // 消息内容
    SendTime       int64               // 发送时间戳
    Status         int                 // 消息状态
    ReadRecords    []byte              // 已读未读记录 (Bitmap 字节数组)
    CreateAt       time.Time
    UpdateAt       time.Time
}

// 会话列表集合 - Conversations
type Conversations struct {
    ID               primitive.ObjectID       // MongoDB 文档ID
    UserId           string                   // 用户ID
    ConversationList map[string]*Conversation // 会话列表 (key: conversationId)
    CreateAt         time.Time
    UpdateAt         time.Time
}
```

#### 关键设计

**1. 会话ID (ConversationId)**
- 算法: `conversationId = hash(min(sendId, recvId) + max(sendId, recvId))`
- 作用: 将双向聊天记录统一到同一会话下，简化查询逻辑

**2. 消息时间戳 (SendTime)**
- 使用 Unix 时间戳记录消息发送时间，用于:
  - 消息排序（按时间倒序）
  - 分页查询（基于时间范围）
  - 消息去重（基于 MongoDB ObjectID）

**3. Bitmap 已读未读**

```go
// pkg/bitmap/bitmap.go
type Bitmap struct {
    bits []byte  // 字节数组存储
    size int     // 总bit数 = len(bits) * 8
}
```

- 将已读未读状态存储在 MongoDB `ChatLog.ReadRecords` 字段（`[]byte`）
- 使用 Bitmap 算法：每个 bit 表示一个用户的已读状态
- 通过 hash(userId) 定位到具体的 bit 位置
- 优势: 
  - 极大节省存储空间（1个用户只占1bit）
  - O(1) 查询和更新复杂度
  - 支持导出/导入，方便持久化

### 3. 消息存储策略

**客户端缓存 + 服务端持久化**

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Client    │────────>│  WebSocket   │────────>│   MongoDB   │
│   Cache     │         │    Server    │         │   (持久化)   │
└─────────────┘         └──────────────┘         └─────────────┘
      │                        │                         │
      │                        ↓                         │
      │                  ┌──────────┐                    │
      └─────────────────>│  Kafka   │───────────────────┘
                         │  (队列)   │
                         └──────────┘
```

**优点**
- ✅ 客户端缓存提升用户体验
- ✅ 服务端持久化保证消息不丢失
- ✅ 支持跨设备消息同步
- ✅ 基于时间戳分页拉取历史消息（倒序查询）

**MongoDB vs MySQL 选型**
- **MongoDB**: 高写入性能、水平扩展、灵活 Schema，适合海量消息存储
- **MySQL**: ACID 事务、关系查询，适合业务数据存储

### 4. 通讯机制与消息流转

#### 消息传输流程

```
┌──────────┐                                              ┌──────────┐
│ Client A │                                              │ Client B │
└─────┬────┘                                              └────┬─────┘
      │ 1. Send Message (WebSocket)                           │
      │────────────────────────────────>                      │
      │                                  ┌──────────────┐     │
      │                                  │ WebSocket    │     │
      │                                  │   Server     │     │
      │                                  └──────┬───────┘     │
      │                                         │             │
      │                                  2. Save to MongoDB   │
      │                                         │             │
      │                                         ↓             │
      │                                  ┌──────────────┐     │
      │                                  │   MongoDB    │     │
      │                                  └──────────────┘     │
      │                                         │             │
      │                                  3. Push via Kafka    │
      │                                         │             │
      │                                         ↓             │
      │                                  ┌──────────────┐     │
      │                                  │    Kafka     │     │
      │                                  └──────┬───────┘     │
      │                                         │             │
      │                                  4. Deliver Message   │
      │                                         │             │
      │                                         └────────────>│
      │                                                       │
      │ 5. ACK Confirmation (Seq)                            │
      │<──────────────────────────────────────────────────────│
```

#### 消息队列设计

**Kafka Topic 设计**

根据 docker-compose.yaml 配置，系统使用以下 Topic：

```yaml
KAFKA_CREATE_TOPICS: "ws2ms_chat:8:1,ms2ps_chat:8:1,msg_to_mongo:8:1"
```

- `ws2ms_chat`: WebSocket 到微服务的消息队列
- `ms2ps_chat`: 微服务到推送服务的消息队列
- `msg_to_mongo`: 消息持久化到 MongoDB 的队列

**消息可靠性保障**
1. **At Least Once**: Kafka 消息至少投递一次
2. **Idempotent**: 基于消息序列号去重
3. **ACK 机制**: 客户端确认接收后更新状态
4. **重试机制**: 推送失败自动重试 (指数退避)



### 5. 鉴权设计

#### 认证流程

```
┌─────────┐   1. Login         ┌──────────┐
│ Client  │──────────────────> │   User   │
│         │                    │  Service │
│         │   2. JWT Token     │          │
│         │ <──────────────────└──────────┘
│         │
│         │   3. Connect (Token in Header)
│         │──────────────────>  ┌──────────┐
│         │                     │    WS    │
│         │   4. Verify Token   │  Server  │
│         │                     └─────┬────┘
│         │                           │
│         │   5. WebSocket Upgrade    │
│         │<──────────────────────────┘
└─────────┘
```

#### JWT 认证实现

**Authentication 接口定义**

```go
type Authentication interface {
    // 验证 Token 有效性
    Auth(w http.ResponseWriter, r *http.Request) bool
    // 获取用户 ID
    UserId(r *http.Request) string
}
```

**WebSocket 握手认证**

WebSocket 连接建立时通过 `sec-websocket-protocol` header 传递 JWT Token:

```go
func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {
    // 1. 从 WebSocket 握手头中提取 Token
    if tok := r.Header.Get("sec-websocket-protocol"); tok != "" {
        r.Header.Set("Authorization", tok)
    }
    
    // 2. 解析并验证 JWT Token
    tok, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
    if err != nil || !tok.Valid {
        return false
    }
    
    // 3. 提取用户 Claims 并存入 Context
    claims, ok := tok.Claims.(jwt.MapClaims)
    if !ok {
        return false
    }
    
    *r = *r.WithContext(
        context.WithValue(r.Context(), ctxdata.IdentityKey, claims[ctxdata.IdentityKey])
    )
    return true
}
```

**安全特性**
- ✅ Token 仅在建立连接时验证，不重复鉴权
- ✅ 基于 go-zero 的 JWT 中间件
- ✅ 用户身份信息存储在 Context 中
- ✅ 支持 Token 过期自动断开连接

### 6. 心跳检测机制

#### 为什么需要心跳检测

**问题场景**
1. WebSocket 长时间无数据传输，网关或中间层可能主动断开连接
2. 网络异常导致连接断开，但客户端/服务端无法感知
3. 客户端进程崩溃，服务端无法及时清理连接资源

**解决方案**
- ✅ **保活 (Keep-Alive)**: 定期发送心跳包维持连接活跃
- ✅ **死链检测 (Dead Connection Detection)**: 及时发现并清理断开的连接

#### 心跳实现设计

参考 gRPC KeepAlive 机制，实现 **Idle Timer (空闲检测)**:

```go
// 参照 grpc 源码心跳检测定时器实现
// connection.go 中的 keepalive() 方法
func (c *Conn) keepalive() {
    idleTimer := time.NewTimer(c.maxConnectionIdle)
    defer idleTimer.Stop()
    
    for {
        select {
        case <-idleTimer.C:
            c.idleMu.Lock()
            idle := c.idle
            if idle.IsZero() { // 连接非空闲
                c.idleMu.Unlock()
                idleTimer.Reset(c.maxConnectionIdle)
                continue
            }
            val := c.maxConnectionIdle - time.Since(idle)
            c.idleMu.Unlock()
            if val <= 0 {
                // 连接空闲时间超过最大空闲时间，优雅关闭连接
                c.s.Close(c)
                return
            }
            idleTimer.Reset(val)
        case <-c.done:
            return
        }
    }
}
```

**Idle Timer (空闲检测) 工作原理**
- **作用**: 检测连接是否长时间无数据传输
- **触发**: 超过 `maxConnectionIdle` 无活动
- **机制**: 
  1. 每次 `ReadMessage`/`WriteMessage` 时更新 `idle` 时间
  2. 定时器检查 `idle.IsZero()`（Zero表示有活动）
  3. 如果非Zero且超过最大空闲时间，则关闭连接
  4. 否则重置定时器为剩余时间
- **动作**: 优雅关闭连接（调用 `c.s.Close(c)`）

**客户端心跳**
- 客户端可发送 `FramePing` 类型消息保持连接活跃
- 服务端收到后回复 `FramePing` 响应

#### 工作流程

```
Client                          Server
  │                               │
  │   1. WebSocket Connect        │
  │──────────────────────────────>│
  │                               │ Start Timers
  │                               │
  │   2. Send Message             │
  │──────────────────────────────>│ Reset lastRead
  │                               │
  │        ... idle ...           │
  │                               │
  │   3. PING                     │
  │<──────────────────────────────│ KeepAlive Timer
  │                               │
  │   4. PONG                     │
  │──────────────────────────────>│ Connection OK
  │                               │
  │        ... timeout ...        │
  │                               │
  │   5. PING (no response)       │
  │<──────────────────────────────│
  │                               │ Wait PongTimeout
  │   X (connection lost)         │
  │                               │ Close Connection
```

#### 关键参数配置

```go
// 启动 WebSocket 服务器时配置
srv := websocket.NewServer(c.ListenOn,
    websocket.WithAuthentication(handler.NewJwtAuth(ctx)),
    // 设置最大空闲时间为10秒
    websocket.WithServerMaxConnectionIdle(10*time.Second),
    // 设置ACK确认模式（NoAck/OnlyAck/RigorAck）
    websocket.WithServerAck(websocket.NoAck),
)
```

**⚠️ 注意事项**
- 服务端检测间隔应 **大于** 客户端发送间隔
- 避免频繁的连接建立和断开
- 客户端断线自动重连机制

### 7. ACK 确认与消息序列

#### 消息序列机制

参考 **TCP 三次握手** 实现可靠消息传输:

```
Client                                Server
  │                                     │
  │  1. Send Msg (Seq=100)             │
  │───────────────────────────────────>│
  │                                     │ Save to DB
  │                                     │
  │  2. ACK (Seq=100, Status=Received) │
  │<───────────────────────────────────│
  │                                     │
  │  3. Read Msg (Seq=100)             │
  │───────────────────────────────────>│
  │                                     │ Update Read Status
  │                                     │
  │  4. ACK (Seq=100, Status=Read)     │
  │<───────────────────────────────────│
```

#### 消息状态流转

```
Sending ──> Sent ──> Received ──> Read
  (0)      (1)       (2)         (3)
```

| 状态 | 说明 | 触发条件 |
|------|------|----------|
| **Sending** | 客户端发送中 | 消息正在发送 |
| **Sent** | 已发送到服务器 | 服务器收到消息 |
| **Received** | 对方已接收 | 对方客户端收到消息 |
| **Read** | 对方已读 | 对方打开会话查看消息 |

#### 核心设计

**1. ACK 模式选择**

```go
type AckType int

const (
    NoAck    AckType = iota  // 不进行ACK应答
    OnlyAck                  // 服务端响应一次应答
    RigorAck                 // 严格ACK，服务端应答后客户端再进行一次应答
)
```

**2. 服务端 ACK 确认队列**

```go
type Conn struct {
    // 读消息队列
    readMessage []*Message
    // 记录消息序列化 (key: 消息id, value: 具体消息)
    readMessageSeq map[string]*Message
    // ACK确认后将消息发送给任务处理
    message chan *Message
}
```

- 服务端维护待确认消息队列 `readMessage`
- 超时未收到客户端 ACK 自动重发（RigorAck 模式）
- 基于消息 ID 和 `AckSeq` 序列号去重

**3. 消息 ID 生成**
- 使用 MongoDB ObjectID 作为消息唯一标识
- 或自定义 UUID/雪花ID

### 8. 安全设计

#### 密码加密

项目使用 `bcrypt` 进行密码hash加密：

```go
// pkg/encrypt/hash.go

// hash加密
func GenPasswordHash(password []byte) ([]byte, error) {
    return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// hash校验
func ValidatePasswordHash(password string, hashed string) bool {
    if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
        return false
    }
    return true
}
```

#### 会话ID计算

使用 MD5 hash 生成会话ID：

```go
func Md5(str []byte) string {
    h := md5.New()
    h.Write(str)
    return hex.EncodeToString(h.Sum(nil))
}

// 会话 ID = md5(min(userId1, userId2) + max(userId1, userId2))
```

#### 安全特性

**1. JWT Token 认证**
- 使用 JWT Token 进行用户身份验证
- Token 通过 `sec-websocket-protocol` 头传递

**2. 连接独占性**
- 同一用户只能保持一个 WebSocket 连接
- 新连接建立时会关闭旧连接

**3. 消息存储**
- MongoDB 存储聊天记录
- Bitmap 存储消息已读未读状态（`readRecords` 字段）

### 9. 高可用与扩展设计

#### 服务水平扩展

**当前实现**

```go
type Server struct {
    sync.RWMutex
    
    // 连接管理：内存Map存储
    connToUser map[*Conn]string  // 连接 -> 用户ID
    userToConn map[string]*Conn  // 用户ID -> 连接
}
```

- 单机部署，使用内存 Map 管理连接
- 通过 `sync.RWMutex` 保证并发安全

**扩展方案**

如需多实例部署，可考虑：

1. **Redis 存储连接路由**
   - Key: `ws:user:{userId}`, Value: `{serverIp}:{port}`
   - 消息推送时查询 Redis 定位目标服务器

2. **Kafka 跨服务器消息转发**
   - 每个 WS 服务器订阅 Kafka Topic
   - 根据 userId 路由到对应分区

#### 配置中心 (Sail)

**动态配置管理**
```go
// 使用 Sail 配置中心
import "github.com/HYY-yu/sail-client"

// 监听配置变更
sail.Watch("app.config", func(config Config) {
    // 热更新配置
    updateConfig(config)
})
```

**优势**
- 配置集中管理
- 动态更新无需重启
- 多环境配置隔离
- 配置版本管理与回滚

#### 优雅重启 (Graceful Restart)

**问题**
- 直接 `kill -9` 会导致连接中断、数据丢失
- 服务重启期间无法处理新请求

**解决方案: go-zero Shutdown 机制**

```go
// 并发同步管道：启动 N 个 goroutine，等它们全部完成之后再继续下一步
var wg sync.WaitGroup

proc.WrapUp() // 通知当前服务优雅停止

// 启动新的服务实例（使用新配置），并加入 WaitGroup
wg.Add(1) // 阻塞
go func(c config.Config) {
	defer wg.Done()
	Run(c) // 新服务开始运行，阻塞于 server.Start()
}(c)

// main goroutine 等待所有服务（当前和未来因配置更新启动的服务）优雅退出
wg.Wait()
```

**关键步骤**
1. 调用 `proc.WrapUp()` 发送停止信号
2. go-zero 内部 `inShutdown` 原子变量标记，拒绝新请求
3. 等待活跃连接/请求处理完成
4. 新服务实例启动并接管流量
5. 旧服务实例优雅退出

---

## 🔧 开发指南

### 代码生成

项目使用 `goctl` 工具生成代码:

```bash
# 生成 API 代码
goctl api go -api apps/im/api/im.api -dir apps/im/api/

# 生成 RPC 代码  
goctl rpc protoc apps/im/rpc/im.proto --go_out=. --go-grpc_out=. --zrpc_out=apps/im/rpc/

# 生成 MongoDB Model
goctl model mongo --type chatLog --dir apps/im/immodels/
```

### 项目规范

**目录结构**
- `internal/`: 内部实现，不对外暴露
- `api/`: HTTP API 定义
- `rpc/`: gRPC 服务定义
- `models/`: 数据模型

**命名规范**
- 文件名: 小写 + 下划线 (snake_case)
- 包名: 小写，简短
- 接口: I 开头或 able 结尾
- 错误码: 统一在 `pkg/xerr` 定义

---

## 📊 性能优化策略

**1. 连接复用**
- WebSocket 长连接复用
- gRPC 连接池
- go-zero TaskRunner 并发任务管理

**2. 内存优化**
- 内存 Map 管理 WebSocket 连接（`connToUser`, `userToConn`）
- Bitmap 算法节省已读未读状态存储空间（1用户=1bit）
- 消息通道缓冲控制（`chan *Message`，容量为1）

**3. 异步处理**
- Kafka 异步消息推送（3个 Topic，8个分区）
- MongoDB 批量写入优化
- go-zero threading.TaskRunner 并发处理

**4. 查询优化**
- MongoDB 索引优化（conversationId + sendTime）
- 按时间倒序分页查询（`SetSort(bson.M{"sendTime": -1})`）
- 默认每次加载100条消息

---

## 🤝 贡献指南

欢迎贡献代码、提交 Issue 和 Pull Request!

### 🎨提交规范

提交代码前，在根目录下执行命令

```
go fmt .
```

在 Go 语言中，最常用的代码规范命令是 go fmt，它会自动格式化你的 Go 代码，使其符合官方推荐的风格。

**Git Commit 类型**

| **Emoji** | **类型**    | **说明**                     |
| --------- | ----------- | ---------------------------- |
| ✨         | feat        | 新功能                       |
| 🐛         | fix         | 修复 bug                     |
| 📝         | docs        | 文档变更                     |
| 🎨         | style       | 代码格式（不影响逻辑）       |
| ♻️         | refactor    | 重构代码（非功能、非修复）   |
| ⚡         | perf        | 性能优化                     |
| 📦         | build       | 构建系统或依赖变更           |
| 🔧         | chore       | 日常杂项（非 src 内容）      |
| 🚀         | deploy      | 发布部署相关                 |
| 🔥         | remove      | 删除代码或文件               |
| ♿         | a11y        | 无障碍优化（例如 aria 标签） |
| 💄         | ui          | 改动 UI 样式                 |
| 🧪         | test        | 测试代码                     |
| 🚚         | mv / rename | 移动或重命名文件             |
| 🗃         | database    | 数据库相关                   |
| 🐎         | perf        | 性能优化                     |
| 🚑         | hotfix      | 紧急修复                     |
| 🧹         | clean       | 清理代码/不再使用的文件      |

---

