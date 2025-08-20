# Shop - 电商微服务系统

## 项目概述

Shop 是一个基于 **Go** 和 **go-zero** 框架开发的微服务电商系统，支持商品管理、购物车、订单创建和秒杀功能。系统使用 MySQL 数据库（支持软删除）、Redis 缓存、gRPC 通信和 JWT 认证，适合高并发场景。

### 主要功能
- **商品管理**:
   - 查询商品列表（支持分页、分类、价格筛选）。
   - 查询商品详情（名称、描述、价格、库存等）。
   - 获取 Banner 和推荐商品。
   - 添加和更新商品（管理员功能，未暴露到网关）。
- **购物车**:
   - 添加、更新、删除商品，清空购物车。
   - 使用 UUID 作为 `cart_id`，支持软删除。
- **订单**:
   - 普通订单：支持从购物车下单或直接下单（传入商品列表）。
   - 秒杀订单：高并发场景，使用 Redis 分布式锁防止超卖。
   - 订单详情查询：限制用户只能查询自己的订单。
- **权限控制**: 通过 `JwtMiddleware` 验证用户身份（`UserID`）。
- **缓存**: Redis 缓存商品（`product:{pid}`）、购物车（`cart:{user_id}`）和订单详情（`order:detail:{order_id}`）。
- **数据库**: MySQL，支持软删除（`deleted_at`），`cart_id` 和 `order_id` 使用 UUID。

### 技术栈
| 组件         | 技术              |
|--------------|-------------------|
| 框架         | go-zero (gRPC, RESTful API) |
| 语言         | Go (>= 1.16)      |
| 数据库       | MySQL (>= 5.7, sqlx 模型) |
| 缓存         | Redis (>= 6.0)    |
| 通信         | gRPC (服务间), RESTful API (客户端) |
| 认证         | JWT (JwtMiddleware) |
| 其他         | ETCD (服务发现), UUID |

---

## 项目架构
- **product 服务**: 管理商品信息（列表、详情、Banner、推荐），支持分页和筛选，字段统一使用 `pid`。
- **cart 服务**: 管理购物车操作（添加、更新、删除、清空），`cart_id` 为 UUID，支持软删除。
- **order 服务**: 处理普通订单（购物车或直接下单）、秒杀订单、订单查询，使用 Redis 锁防止超卖。
- **gateway 服务**: 提供 RESTful API 入口，通过 `JwtMiddleware` 验证用户身份，调用 gRPC 接口。
- **payment 服务**: 支付系统
---

## 安装与部署

### 环境要求
- Go >= 1.16
- MySQL >= 5.7
- Redis >= 6.0
- ETCD >= 3.5
- `goctl`（go-zero 命令行工具）

### 安装步骤
1. **克隆项目**
   ```bash
   git clone <repository_url>
   cd shop

安装依赖bash
   ```bash
   go mod tidy
   ```

初始化数据库
创建 MySQL 数据库 mall 并执行以下建表语句：
   ```sql

CREATE TABLE products
(
   id          BIGINT AUTO_INCREMENT PRIMARY KEY,
   pid         VARCHAR(50)  NOT NULL UNIQUE,
   name        VARCHAR(100) NOT NULL,
   description TEXT,
   detail      TEXT,
   main_image  VARCHAR(255),
   thumbnail   VARCHAR(255),
   price DOUBLE NOT NULL,
   stock       INT          NOT NULL,
   category_id VARCHAR(50),
   is_banner   BOOLEAN               DEFAULT FALSE,
   created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE carts
(
   cart_id    VARCHAR(50) NOT NULL UNIQUE,
   user_id    VARCHAR(50) NOT NULL UNIQUE,
   created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NULL,
   PRIMARY KEY (cart_id)
);

CREATE TABLE cart_items
(
   id         BIGINT AUTO_INCREMENT PRIMARY KEY,
   cart_id    VARCHAR(50) NOT NULL,
   product_id VARCHAR(50) NOT NULL,
   quantity   INT         NOT NULL,
   created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NULL,
   UNIQUE KEY uk_cart_product (cart_id, product_id)
);

CREATE TABLE orders
(
   order_id   VARCHAR(50) NOT NULL UNIQUE,
   user_id    VARCHAR(50) NOT NULL,
   total_price DOUBLE NOT NULL,
   status     VARCHAR(20) NOT NULL,
   created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (order_id)
);

CREATE TABLE order_items
(
   id         BIGINT AUTO_INCREMENT PRIMARY KEY,
   order_id   VARCHAR(50) NOT NULL,
   product_id VARCHAR(50) NOT NULL,
   quantity   INT         NOT NULL,
   created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_pid ON products (pid);
CREATE INDEX idx_carts_user_id ON carts (user_id);
CREATE INDEX idx_cart_items_cart_id ON cart_items (cart_id);
CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_order_items_order_id ON order_items (order_id);

```
配置服务shop/product/etc/product.yaml:
```yaml

Name: product
Host: 0.0.0.0
Port: 8086
Mysql:
DataSource: user:password@tcp(localhost:3306)/mall?charset=utf8mb4&parseTime=true
Redis:
   Host: localhost:6379
   Pass: ""
   Type: node
```
shop/cart/etc/cart.yaml:
```yaml

Name: cart
Host: 0.0.0.0
Port: 8084
Mysql:
DataSource: user:password@tcp(localhost:3306)/mall?charset=utf8mb4&parseTime=true
Redis:
   Host: localhost:6379
   Pass: ""
   Type: node
ProductRpc:
   Etcd:
      Hosts:
      - localhost:2379
      Key: product.rpc
```
shop/order/etc/order.yaml:
```yaml

Name: order
Host: 0.0.0.0
Port: 8085
Mysql:
DataSource: user:password@tcp(localhost:3306)/mall?charset=utf8mb4&parseTime=true
Redis:
   Host: localhost:6379
   Pass: ""
   Type: node
ProductRpc:
   Etcd:
      Hosts:
      - localhost:2379
      Key: product.rpc
CartRpc:
   Etcd:
      Hosts:
      - localhost:2379
      Key: cart.rpc
```
shop/gateway/etc/gatewayapi.yaml:
```yaml

Name: gatewayapi
Host: 0.0.0.0
Port: 8080
Jwt:
Secret: your_jwt_secret
ProductRpc:
   Etcd:
      Hosts:
      - localhost:2379
      Key: product.rpc
CartRpc:
   Etcd:
      Hosts:
      - localhost:2379
      Key: cart.rpc
OrderRpc:
   Etcd:
      Hosts:
      - localhost:2379
      Key: order.rpc
```

启动服务
```bash

# 启动 product 服务
cd shop/product
go run product.go

# 启动 cart 服务
cd shop/cart
go run cart.go

# 启动 order 服务
cd shop/order
go run order.go

# 启动 gateway 服务
cd shop/gateway
go run gatewayapi.go
```

错误码细化:
```go
const (
ErrCodeOK             = 0
ErrCodeInvalidParam   = 1000
ErrCodeEmptyCart      = 1001
ErrCodeEmptyItems     = 1002
ErrCodeStockOut       = 1003
ErrCodeProductNotFound = 1004
ErrCodeOrderNotFound  = 1005
ErrCodeUnauthorized    = 1006
)
```

## 后续优化
未登录用户支持: 使用 session_id 存储临时购物车。
性能优化:批量查询商品信息，减少 ProductRpc 调用。
缓存商品数据（product:{pid}）。

管理员功能: 暴露 AddProduct 和 UpdateProduct API，支持商品管理。
监控与日志: 集成 Prometheus 和 Loki。

FAQ如何生成 JWT Token?
配置 gatewayapi.yaml 中的 Jwt.Secret，使用 JWT 库生成 token，包含 UserID.
为何收到 "无权访问该订单"?
确保 user_id 与订单的 user_id 匹配，且 JWT Token 有效。
如何测试高并发?
使用 wrk 或 ab 模拟多用户请求 SeckillOrder 接口，验证 Redis 锁效果。
如何添加管理员功能?
扩展 product.api，添加 AddProduct 和 UpdateProduct 接口，结合权限验证。

Docker 部署（建议）编写 Dockerfile:
```dockerfile

FROM golang:1.16
WORKDIR /app
COPY . .
RUN go mod tidy
CMD ["go", "run", "<service>.go"]
```
使用 Docker Compose:
```yaml
version: '3'
services:
   product:
      build: ./product
      ports:
         - "8086:8086"
      depends_on:
         - mysql
         - redis
   cart:
      build: ./cart
      ports:
         - "8084:8084"
      depends_on:
         - mysql
         - redis
         - product
   order:
      build: ./order
      ports:
         - "8085:8085"
      depends_on:
         - mysql
         - redis
         - product
         - cart
   gateway:
      build: ./gateway
      ports:
         - "8080:8080"
      depends_on:
         - product
         - cart
         - order
   mysql:
      image: mysql:5.7
      environment:
         MYSQL_ROOT_PASSWORD: password
         MYSQL_DATABASE: mall
      volumes:
         - ./init.sql:/docker-entrypoint-initdb.d/init.sql
   redis:
      image: redis:6.0
   etcd:
      image: quay.io/coreos/etcd:v3.5
```
运行:
```bash

docker-compose up -d
```
注意事项数据库索引: 确保 products, carts, cart_items, orders, order_items 表已创建索引。
JWT 认证: 配置 gatewayapi.yaml 中的 Jwt.Secret.
Redis 缓存: 确保 Redis 连接正常，缓存键有效期为 3600 秒。
并发安全: 秒杀订单已使用 Redis 锁，普通订单建议添加锁。
字段统一: 所有服务使用 pid，数据库仍为 product_id.

联系方式如有问题，请联系项目维护者或提交 Issue.完成日期: 2025-08-18
