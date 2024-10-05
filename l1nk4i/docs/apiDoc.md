---
title: hduhelp_backend_task
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
generator: "@tarslib/widdershins v4.0.23"

---

# hduhelp_backend_task

Base URLs:

# Authentication

# users

## POST 用户注册

POST /api/users/register

注册

> Body 请求参数

```json
{
  "username": "l1nkQAQ",
  "password": "vidarteam"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» username|body|string| 是 |none|
|» password|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|none|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» msg|string|true|none||none|

状态码 **400**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» error|string|true|none||none|

状态码 **500**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» error|string|true|none||none|

## GET 获取当前用户信息

GET /api/users/userinfo

查询当前用户信息

> 返回示例

> 200 Response

```json
{
  "role": "string",
  "user_id": "string",
  "username": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» role|string|true|none||none|
|» user_id|string|true|none||none|
|» username|string|true|none||none|

状态码 **400**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» error|string|true|none||none|

## POST 用户登录

POST /api/users/login

登录

> Body 请求参数

```json
{
  "username": "l1nkQAQ",
  "password": "vidarteam"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» username|body|string| 是 |none|
|» password|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» msg|string|true|none||none|

## POST 用户登出

POST /api/users/logout

登出

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» msg|string|true|none||none|

# questions

## POST 创建问题

POST /api/questions

创建问题

> Body 请求参数

```json
{
  "title": "man",
  "content": "what can i say"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» title|body|string| 是 |none|
|» content|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» msg|string|true|none||none|

## GET 列出用户发表的所有问题

GET /api/questions

列出用户发表的所有问题

> 返回示例

> 200 Response

```json
{
  "question_id": [
    "string"
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» question_id|[string]|true|none||none|

## POST 创建指定问题的回答

POST /api/questions/{question-id}/answers

发送回答

> Body 请求参数

```json
{
  "content": "manman"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|question-id|path|string| 是 |none|
|body|body|object| 否 |none|
|» question_id|body|string| 是 |none|
|» content|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "answer_id": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» answer_id|string|true|none||none|

## GET 获取指定问题的所有回答

GET /api/questions/{question-id}/answers

获得问题对应的所有回答

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|question-id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "data": [
    {
      "ID": 0,
      "CreatedAt": "string",
      "UpdatedAt": "string",
      "DeletedAt": null,
      "AnswerID": "string",
      "UserID": "string",
      "QuestionID": "string",
      "Content": "string"
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» data|[object]|true|none||none|
|»» ID|integer|true|none||none|
|»» CreatedAt|string|true|none||none|
|»» UpdatedAt|string|true|none||none|
|»» DeletedAt|null|true|none||none|
|»» AnswerID|string|true|none||none|
|»» UserID|string|true|none||none|
|»» QuestionID|string|true|none||none|
|»» Content|string|true|none||none|

## DELETE 删除指定问题

DELETE /api/questions/{question-id}

删除问题

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|question-id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» msg|string|true|none||none|

## PUT 更新指定问题

PUT /api/questions/{question-id}

更改问题的标题或者内容

> Body 请求参数

```json
{
  "title": "manba",
  "content": "out"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|question-id|path|string| 是 |none|
|body|body|object| 否 |none|
|» question_id|body|string| 是 |none|
|» title|body|string| 是 |none|
|» content|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» msg|string|true|none||none|

## GET 获取指定问题

GET /api/questions/{question-id}

通过question_id查找对应问题对象

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|question-id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "data": {
    "ID": 0,
    "CreatedAt": "string",
    "UpdatedAt": "string",
    "DeletedAt": null,
    "QuestionID": "string",
    "UserID": "string",
    "BestAnswerID": "string",
    "Title": "string",
    "Content": "string"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» data|object|true|none||none|
|»» ID|integer|true|none||none|
|»» CreatedAt|string|true|none||none|
|»» UpdatedAt|string|true|none||none|
|»» DeletedAt|null|true|none||none|
|»» QuestionID|string|true|none||none|
|»» UserID|string|true|none||none|
|»» BestAnswerID|string|true|none||none|
|»» Title|string|true|none||none|
|»» Content|string|true|none||none|

## GET 搜索问题

GET /api/questions/search

通过输入匹配问题标题和内容，搜索问题（无鉴权）

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|content|query|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "data": [
    {
      "ID": 0,
      "CreatedAt": "string",
      "UpdatedAt": "string",
      "DeletedAt": null,
      "QuestionID": "string",
      "UserID": "string",
      "BestAnswerID": "string",
      "Title": "string",
      "Content": "string"
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» data|[object]|true|none||none|
|»» ID|integer|true|none||none|
|»» CreatedAt|string|true|none||none|
|»» UpdatedAt|string|true|none||none|
|»» DeletedAt|null|true|none||none|
|»» QuestionID|string|true|none||none|
|»» UserID|string|true|none||none|
|»» BestAnswerID|string|true|none||none|
|»» Title|string|true|none||none|
|»» Content|string|true|none||none|

## PUT 更新最佳答案

PUT /api/questions/{question-id}/best/{answer-id}

更新最佳答案

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|question-id|path|string| 是 |none|
|answer-id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» msg|string|true|none||none|

# answers

## DELETE 删除指定回答

DELETE /api/answers/{answer-id}

删除回答

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|answer-id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» msg|string|true|none||none|

## PUT 更新指定回答

PUT /api/answers/{answer-id}

更改回答内容

> Body 请求参数

```json
{
  "content": "manba out"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|answer-id|path|string| 是 |none|
|body|body|object| 否 |none|
|» answer_id|body|string| 是 |none|
|» content|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» msg|string|true|none||none|

# admin

## DELETE 删除问题

DELETE /api/admin/questions/{question-id}

删除问题

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|question-id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» msg|string|true|none||none|

## DELETE 删除答案

DELETE /api/admin/answers/{answer-id}

删除答案

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|answer-id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» msg|string|true|none||none|

# 数据模型

<h2 id="tocS_User">User</h2>

<a id="schemauser"></a>
<a id="schema_User"></a>
<a id="tocSuser"></a>
<a id="tocsuser"></a>

```json
{
  "username": "string",
  "password": "string",
  "role": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "deleted_at": "2019-08-24T14:15:22Z",
  "id": 1,
  "updated_at": "2019-08-24T14:15:22Z",
  "user_id": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|username|string|true|none||none|
|password|string|true|none||none|
|role|string|true|none||none|
|created_at|string(date-time)|false|none||none|
|deleted_at|string(date-time)|false|none||none|
|id|integer|true|none||none|
|updated_at|string(date-time)|false|none||none|
|user_id|string|true|none||none|

<h2 id="tocS_Question">Question</h2>

<a id="schemaquestion"></a>
<a id="schema_Question"></a>
<a id="tocSquestion"></a>
<a id="tocsquestion"></a>

```json
{
  "title": "string",
  "content": "string",
  "best_answer_id": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "deleted_at": "2019-08-24T14:15:22Z",
  "id": 1,
  "question_id": "string",
  "updated_at": "2019-08-24T14:15:22Z",
  "user_id": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|title|string|true|none||none|
|content|string|true|none||none|
|best_answer_id|string|false|none||none|
|created_at|string(date-time)|false|none||none|
|deleted_at|string(date-time)|false|none||none|
|id|integer|true|none||none|
|question_id|string|true|none||none|
|updated_at|string(date-time)|false|none||none|
|user_id|string|true|none||none|

<h2 id="tocS_Answer">Answer</h2>

<a id="schemaanswer"></a>
<a id="schema_Answer"></a>
<a id="tocSanswer"></a>
<a id="tocsanswer"></a>

```json
{
  "content": "string",
  "answer_id": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "deleted_at": "2019-08-24T14:15:22Z",
  "id": 1,
  "question_id": "string",
  "updated_at": "2019-08-24T14:15:22Z",
  "user_id": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|content|string|true|none||none|
|answer_id|string|true|none||none|
|created_at|string(date-time)|false|none||none|
|deleted_at|string(date-time)|false|none||none|
|id|integer|true|none||none|
|question_id|string|true|none||none|
|updated_at|string(date-time)|false|none||none|
|user_id|string|true|none||none|

