# 秒杀系统（Seckill System）

一个基于 Go 构建的简易秒杀系统，支持商品初始化、库存管理和 Redis 缓存预热，适合学习和入门高并发场景下的系统设计。

---

## 📁 项目结构说明（internal/）

```

internal/
├── config          # 配置相关（如数据库、Redis配置加载）
├── handler         # HTTP 接口层（处理请求、参数校验）
├── model           # 数据模型定义（数据库结构 & Redis缓存结构）
├── repository      # 数据访问层（数据库 & Redis 的操作）
├── router          # 路由注册
├── service         # 核心业务逻辑层（商品初始化、缓存写入等）
├── util            # 工具类（时间处理、错误处理等）

````

---

## 🔁 请求处理流程图（商品初始化）

```plaintext
前端发送 POST /initialize-product 请求
           ↓
Router 路由匹配到 handler
           ↓
Handler 解析 JSON 请求体 + 参数校验
           ↓
调用 Service 进行商品初始化
           ↓
 ┌──────────────┬─────────────────────┐
 ↓              ↓                     ↓
写入数据库   写入 Redis 缓存       返回创建成功
(repository)  (repository)          (HTTP 201)
````

---

## 🧱 数据库表设计（商品表）

```sql
CREATE TABLE products (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  price NUMERIC(10,2) NOT NULL,
  stock BIGINT NOT NULL,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

>  数据库字段说明：

* `stock`：用于秒杀库存
* `start_time` / `end_time`：控制秒杀活动时间段
* `price`：商品秒杀价格

---

## Redis 缓存结构

商品缓存结构以 `product_info:<id>` 为 key，value 为 JSON 序列化后的结构：

```json
{
  "name": "秒杀商品A",
  "description": "限时抢购商品",
  "price": 99.9,
  "stock": 100,
  "start_time": "2025-06-02T12:00:00Z",
  "end_time": "2025-06-02T12:10:00Z"
}
```

---

## 🛠 技术栈

| 技术           | 说明                      |
| ------------ | ----------------------- |
| Go (Golang)  | 核心开发语言                  |
| net/http     | 原生 HTTP 路由与服务           |
| GORM         | 操作数据库（PostgreSQL/MySQL） |
| Redis        | 高性能缓存，用于库存信息等热点数据       |
| JSON         | 数据交互格式                  |
| time.RFC3339 | 秒杀时间格式解析                |

---

## 已实现功能

* [x] 商品初始化（数据库 + Redis）
* [x] 请求参数校验
* [x] 秒杀时间设置
* [x] 基础架构分层（Handler / Service / Repository）

---

## 下一步可扩展

* 秒杀下单接口（高并发控制、幂等）
* 限流防刷（令牌桶/滑动窗口等）
* 消息队列异步扣库存
* Redis 缓存库存控制（预减）

---

