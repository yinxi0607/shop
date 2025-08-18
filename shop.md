---
title: shop
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# shop

Base URLs:

* <a href="http://127.0.0.1:8080/">Develop Env: http://127.0.0.1:8080/</a>

# Authentication

# Default

## GET health

GET /api/v1/health

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

# user

## POST register

POST /user/register

> Body Parameters

```json
{
  "username": "string",
  "password": "string",
  "email": "string",
  "avatar": "string",
  "bio": "string",
  "address": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|body|body|object| no |none|
|» username|body|string| yes |none|
|» password|body|string| yes |none|
|» email|body|string| yes |none|
|» avatar|body|string| yes |none|
|» bio|body|string| yes |none|
|» address|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "user_id": "2fc10a54-f1d8-40c5-b8c6-f5de69ba2baf"
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» user_id|string|true|none||none|
|» msg|string|true|none||none|

## POST login

POST /user/login

> Body Parameters

```json
{
  "username": "string",
  "password": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|body|body|object| no |none|
|» username|body|string| yes |none|
|» password|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTYxMDQxNDgsInJvbGUiOiJ1c2VyIiwidXNlcl9pZCI6IjdlMWU4MjNhLTU4ZWUtNGZiNi05NWU4LWUwNjIzMzgzZWFmMSJ9.9_p8K7RLEmTZ2jYmGTGAB3bnvofN8ZDgyxYoGKsLRwc"
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» token|string|true|none||none|
|» msg|string|true|none||none|

## PUT change-avatar

PUT /api/v1/user/change-avatar

> Body Parameters

```json
{
  "new_avatar": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|
|body|body|object| no |none|
|» new_avatar|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "success": true
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» success|boolean|true|none||none|
|» msg|string|true|none||none|

## PUT change-password

PUT /api/v1/user/change-password

> Body Parameters

```json
{
  "old_password": "string",
  "new_password": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|
|body|body|object| no |none|
|» old_password|body|string| yes |none|
|» new_password|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "success": true
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» success|boolean|true|none||none|
|» msg|string|true|none||none|

## PUT change-role

PUT /api/v1/user/change-role

> Body Parameters

```json
{
  "role": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|
|body|body|object| no |none|
|» role|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 1000,
  "data": null,
  "msg": "user_id is not admin"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|null|true|none||none|
|» msg|string|true|none||none|

## PUT change-username

PUT /api/v1/user/change-username

> Body Parameters

```json
{
  "new_username": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|
|body|body|object| no |none|
|» new_username|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "success": true
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» success|boolean|true|none||none|
|» msg|string|true|none||none|

## GET info

GET /api/v1/user/info

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "user_id": "7e1e823a-58ee-4fb6-95e8-e0623383eaf1",
    "username": "test888",
    "email": "test88@test88.com",
    "avatar": "8888",
    "bio": "888888",
    "address": "beijing",
    "created_at": "2025-08-18 06:39:58"
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» user_id|string|true|none||none|
|»» username|string|true|none||none|
|»» email|string|true|none||none|
|»» avatar|string|true|none||none|
|»» bio|string|true|none||none|
|»» address|string|true|none||none|
|»» created_at|string|true|none||none|
|» msg|string|true|none||none|

# product

## GET banner

GET /api/v1/product/banner

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|limit|query|integer| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "products": [
      {
        "pid": "87b9934b-b68d-43c3-a956-02c827b3b2ab",
        "name": "ipad",
        "description": "ipad",
        "detail": "apple ipad",
        "main_image": "https://picsum.photos/seed/DJBziq6OuT/3434/3704",
        "thumbnail": "https://loremflickr.com/3484/1417?lock=6562046881672518",
        "price": 554.49,
        "stock": 5,
        "category_id": "a21d8a9a-f17c-4521-96aa-a82a5216e6ee",
        "created_at": "2025-08-12T09:59:08Z",
        "updated_at": "2025-08-12T09:59:08Z"
      }
    ]
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» products|[object]|true|none||none|
|»»» pid|string|false|none||none|
|»»» name|string|false|none||none|
|»»» description|string|false|none||none|
|»»» detail|string|false|none||none|
|»»» main_image|string|false|none||none|
|»»» thumbnail|string|false|none||none|
|»»» price|number|false|none||none|
|»»» stock|integer|false|none||none|
|»»» category_id|string|false|none||none|
|»»» created_at|string|false|none||none|
|»»» updated_at|string|false|none||none|
|» msg|string|true|none||none|

## GET list

GET /api/v1/product/list

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|integer| no |none|
|page_size|query|integer| no |none|
|category_id|query|integer| no |none|
|min_price|query|integer| no |none|
|max_price|query|integer| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "products": [
      {
        "pid": "87b9934b-b68d-43c3-a956-02c827b3b2ab",
        "name": "ipad",
        "description": "ipad",
        "detail": "apple ipad",
        "main_image": "https://picsum.photos/seed/DJBziq6OuT/3434/3704",
        "thumbnail": "https://loremflickr.com/3484/1417?lock=6562046881672518",
        "price": 554.49,
        "stock": 5,
        "category_id": "a21d8a9a-f17c-4521-96aa-a82a5216e6ee",
        "created_at": "2025-08-12T09:59:08Z",
        "updated_at": "2025-08-12T09:59:08Z"
      },
      {
        "pid": "87s9934b-b68d-43c3-a956-02c827b3b2ab",
        "name": "ipad",
        "description": "ipad",
        "detail": "apple ipad",
        "main_image": "https://picsum.photos/seed/DJBziq6OuT/3434/3704",
        "thumbnail": "https://loremflickr.com/3484/1417?lock=6562046881672518",
        "price": 554.49,
        "stock": 5,
        "category_id": "a21d8a9a-f17c-4521-96aa-a82a5216e6ee",
        "created_at": "2025-08-13T06:29:38Z",
        "updated_at": "2025-08-13T06:29:38Z"
      },
      {
        "pid": "b26ca5ca-7727-11f0-943c-0242c0a80002",
        "name": "ddddddd",
        "description": "ipad",
        "detail": "apple ipad",
        "main_image": "https://picsum.photos/seed/b04f4hgjOc/3541/1963",
        "thumbnail": "https://picsum.photos/seed/Zsc50coy/2438/641",
        "price": 800,
        "stock": 68,
        "category_id": "84eb5446-2060-4703-9bec-3b61f3e3e752",
        "created_at": "2025-08-07T09:41:52Z",
        "updated_at": "2025-08-13T11:15:28Z"
      },
      {
        "pid": "be375a12-7727-11f0-943c-0242c0a80002",
        "name": "Laptop",
        "description": "High-performance laptop",
        "detail": "16GB RAM",
        "main_image": "https://example.com/laptop.jpg",
        "thumbnail": "https://example.com/laptop_thumb.jpg",
        "price": 1299.99,
        "stock": 48,
        "category_id": "802fb70e-7728-11f0-943c-0242c0a80002",
        "created_at": "2025-08-07T09:41:52Z",
        "updated_at": "2025-08-13T09:46:30Z"
      }
    ],
    "total": 4
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» products|[object]|true|none||none|
|»»» pid|string|true|none||none|
|»»» name|string|true|none||none|
|»»» description|string|true|none||none|
|»»» detail|string|true|none||none|
|»»» main_image|string|true|none||none|
|»»» thumbnail|string|true|none||none|
|»»» price|number|true|none||none|
|»»» stock|integer|true|none||none|
|»»» category_id|string|true|none||none|
|»»» created_at|string|true|none||none|
|»»» updated_at|string|true|none||none|
|»» total|integer|true|none||none|
|» msg|string|true|none||none|

## GET detail

GET /api/v1/product/detail

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|pid|query|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "product": {
      "pid": "b26ca5ca-7727-11f0-943c-0242c0a80002",
      "name": "ddddddd",
      "description": "ipad",
      "detail": "apple ipad",
      "main_image": "https://picsum.photos/seed/b04f4hgjOc/3541/1963",
      "thumbnail": "https://picsum.photos/seed/Zsc50coy/2438/641",
      "price": 800,
      "stock": 68,
      "category_id": "84eb5446-2060-4703-9bec-3b61f3e3e752",
      "created_at": "2025-08-07T09:41:52Z",
      "updated_at": "2025-08-13T11:15:28Z"
    }
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» product|object|true|none||none|
|»»» pid|string|true|none||none|
|»»» name|string|true|none||none|
|»»» description|string|true|none||none|
|»»» detail|string|true|none||none|
|»»» main_image|string|true|none||none|
|»»» thumbnail|string|true|none||none|
|»»» price|integer|true|none||none|
|»»» stock|integer|true|none||none|
|»»» category_id|string|true|none||none|
|»»» created_at|string|true|none||none|
|»»» updated_at|string|true|none||none|
|» msg|string|true|none||none|

## GET recommended

GET /api/v1/product/recommended

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|limit|query|integer| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "products": [
      {
        "pid": "b26ca5ca-7727-11f0-943c-0242c0a80002",
        "name": "ddddddd",
        "description": "ipad",
        "detail": "apple ipad",
        "main_image": "https://picsum.photos/seed/b04f4hgjOc/3541/1963",
        "thumbnail": "https://picsum.photos/seed/Zsc50coy/2438/641",
        "price": 800,
        "stock": 68,
        "category_id": "84eb5446-2060-4703-9bec-3b61f3e3e752",
        "created_at": "2025-08-07T09:41:52Z",
        "updated_at": "2025-08-13T11:15:28Z"
      },
      {
        "pid": "be375a12-7727-11f0-943c-0242c0a80002",
        "name": "Laptop",
        "description": "High-performance laptop",
        "detail": "16GB RAM",
        "main_image": "https://example.com/laptop.jpg",
        "thumbnail": "https://example.com/laptop_thumb.jpg",
        "price": 1299.99,
        "stock": 48,
        "category_id": "802fb70e-7728-11f0-943c-0242c0a80002",
        "created_at": "2025-08-07T09:41:52Z",
        "updated_at": "2025-08-13T09:46:30Z"
      },
      {
        "pid": "87b9934b-b68d-43c3-a956-02c827b3b2ab",
        "name": "ipad",
        "description": "ipad",
        "detail": "apple ipad",
        "main_image": "https://picsum.photos/seed/DJBziq6OuT/3434/3704",
        "thumbnail": "https://loremflickr.com/3484/1417?lock=6562046881672518",
        "price": 554.49,
        "stock": 5,
        "category_id": "a21d8a9a-f17c-4521-96aa-a82a5216e6ee",
        "created_at": "2025-08-12T09:59:08Z",
        "updated_at": "2025-08-12T09:59:08Z"
      },
      {
        "pid": "87s9934b-b68d-43c3-a956-02c827b3b2ab",
        "name": "ipad",
        "description": "ipad",
        "detail": "apple ipad",
        "main_image": "https://picsum.photos/seed/DJBziq6OuT/3434/3704",
        "thumbnail": "https://loremflickr.com/3484/1417?lock=6562046881672518",
        "price": 554.49,
        "stock": 5,
        "category_id": "a21d8a9a-f17c-4521-96aa-a82a5216e6ee",
        "created_at": "2025-08-13T06:29:38Z",
        "updated_at": "2025-08-13T06:29:38Z"
      }
    ]
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» products|[object]|true|none||none|
|»»» pid|string|true|none||none|
|»»» name|string|true|none||none|
|»»» description|string|true|none||none|
|»»» detail|string|true|none||none|
|»»» main_image|string|true|none||none|
|»»» thumbnail|string|true|none||none|
|»»» price|integer|true|none||none|
|»»» stock|integer|true|none||none|
|»»» category_id|string|true|none||none|
|»»» created_at|string|true|none||none|
|»»» updated_at|string|true|none||none|
|» msg|string|true|none||none|

## POST add

POST /api/v1/product/add

> Body Parameters

```json
{
  "name": "string",
  "description": "string",
  "detail": "string",
  "main_image": "string",
  "thumbnail": "string",
  "price": 0,
  "stock": 0,
  "category_id": "string",
  "is_banner": true,
  "pid": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|
|body|body|object| no |none|
|» name|body|string| yes |none|
|» description|body|string| no |none|
|» detail|body|string| no |none|
|» main_image|body|string| no |none|
|» thumbnail|body|string| no |none|
|» price|body|number| yes |none|
|» stock|body|integer| yes |none|
|» category_id|body|string| yes |none|
|» is_banner|body|boolean| yes |none|
|» pid|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "pid": "64cb39ae-3966-45de-b660-cc8a60e92c34"
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» pid|string|true|none||none|
|» msg|string|true|none||none|

## PUT updata

PUT /api/v1/product/update

> Body Parameters

```json
{
  "name": "string",
  "description": "string",
  "detail": "string",
  "main_image": "string",
  "thumbnail": "string",
  "price": 0,
  "stock": 0,
  "category_id": "string",
  "is_banner": "string",
  "pid": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|
|body|body|object| no |none|
|» name|body|string| yes |none|
|» description|body|string| no |none|
|» detail|body|string| no |none|
|» main_image|body|string| no |none|
|» thumbnail|body|string| no |none|
|» price|body|number| yes |none|
|» stock|body|integer| yes |none|
|» category_id|body|string| yes |none|
|» is_banner|body|string| yes |none|
|» pid|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "success": true
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» success|boolean|true|none||none|
|» msg|string|true|none||none|

# order

## POST create

POST /api/v1/order/create

> Body Parameters

```json
{
  "items": [
    {
      "product_id": "sss",
      "quantity": 1
    }
  ]
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|
|body|body|object| no |none|
|» items|body|[object]| yes |none|
|»» product_id|body|string| no |none|
|»» quantity|body|integer| no |none|
|» use_cart|body|boolean| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "order_id": "acd56c75-1667-4e8e-ad7b-7e8d4ca31402"
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» order_id|string|true|none||none|
|» msg|string|true|none||none|

## GET detail

GET /api/v1/order/detail

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|order_id|query|string| no |none|
|Authorization|header|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "order": {
      "order_id": "acd56c75-1667-4e8e-ad7b-7e8d4ca31402",
      "user_id": "7e1e823a-58ee-4fb6-95e8-e0623383eaf1",
      "items": [
        {
          "pid": "87b9934b-b68d-43c3-a956-02c827b3b2ab",
          "quantity": 1
        }
      ],
      "total_price": 554.49,
      "status": "pending",
      "created_at": "2025-08-18 06:54:32",
      "updated_at": "2025-08-18 06:54:32"
    }
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» order|object|true|none||none|
|»»» order_id|string|true|none||none|
|»»» user_id|string|true|none||none|
|»»» items|[object]|true|none||none|
|»»»» pid|string|false|none||none|
|»»»» quantity|integer|false|none||none|
|»»» total_price|number|true|none||none|
|»»» status|string|true|none||none|
|»»» created_at|string|true|none||none|
|»»» updated_at|string|true|none||none|
|» msg|string|true|none||none|

## GET list

GET /api/v1/order/list

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|integer| no |none|
|page_size|query|integer| no |none|
|Authorization|header|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "orders": [
      {
        "order_id": "acd56c75-1667-4e8e-ad7b-7e8d4ca31402",
        "user_id": "7e1e823a-58ee-4fb6-95e8-e0623383eaf1",
        "items": [
          {
            "pid": "87b9934b-b68d-43c3-a956-02c827b3b2ab",
            "quantity": 1
          }
        ],
        "total_price": 554.49,
        "status": "pending",
        "created_at": "2025-08-18 06:54:32",
        "updated_at": "2025-08-18 06:54:32"
      }
    ],
    "total": 1
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» orders|[object]|true|none||none|
|»»» order_id|string|false|none||none|
|»»» user_id|string|false|none||none|
|»»» items|[object]|false|none||none|
|»»»» pid|string|false|none||none|
|»»»» quantity|integer|false|none||none|
|»»» total_price|number|false|none||none|
|»»» status|string|false|none||none|
|»»» created_at|string|false|none||none|
|»»» updated_at|string|false|none||none|
|»» total|integer|true|none||none|
|» msg|string|true|none||none|

## POST seckill

POST /api/v1/order/seckill

> Body Parameters

```json
{
  "items": [
    {
      "pid": "sss",
      "quantity": 1
    }
  ]
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|
|body|body|object| no |none|
|» pid|body|string| yes |none|
|» quantity|body|integer| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "order_id": "b7119d00-1bd2-4ef5-9bb3-29a0a3e9d46b",
    "success": true
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» order_id|string|true|none||none|
|»» success|boolean|true|none||none|
|» msg|string|true|none||none|

# cart

## GET get list

GET /api/v1/cart

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "cart_id": "",
    "items": [
      {
        "pid": "87b9934b-b68d-43c3-a956-02c827b3b2ab",
        "quantity": 1
      },
      {
        "pid": "b26ca5ca-7727-11f0-943c-0242c0a80002",
        "quantity": 3
      }
    ]
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» cart_id|string|true|none||none|
|»» items|[object]|true|none||none|
|»»» pid|string|true|none||none|
|»»» quantity|integer|true|none||none|
|» msg|string|true|none||none|

## POST add cart

POST /api/v1/cart/add

> Body Parameters

```json
{
  "pid": "string",
  "quantity": 0
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|
|body|body|object| no |none|
|» pid|body|string| yes |none|
|» quantity|body|number| yes |none|

> Response Examples

```json
{
  "code": 0,
  "data": {
    "success": true
  },
  "msg": "success"
}
```

```json
{
  "code": 0,
  "data": {
    "success": true
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» success|boolean|true|none||none|
|» msg|string|true|none||none|

## DELETE clear

DELETE /api/v1/cart/clear

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "success": true
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» success|boolean|true|none||none|
|» msg|string|true|none||none|

## PUT update

PUT /api/v1/cart/update

> Body Parameters

```json
{
  "pid": "string",
  "quantity": 0
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| no |none|
|body|body|object| no |none|
|» pid|body|string| yes |none|
|» quantity|body|number| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "success": true
  },
  "msg": "success"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» data|object|true|none||none|
|»» success|boolean|true|none||none|
|» msg|string|true|none||none|

# Data Schema

