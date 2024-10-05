---
title: 个人项目
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

# 个人项目

Base URLs:

# Authentication

# Default

## POST 注册

POST /api/user/register

输入username ,password,并在sql数据库中查找是否有该账号

> Body 请求参数

```yaml
username: test
password: test

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» username|body|string| 否 |none|
|» password|body|string| 否 |none|

> 返回示例

```json
{
  "error": "用户名已存在"
}
```

```json
{
  "error": "注册失败"
}
```

```json
{
  "message": "注册成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 登出

POST /api/user/logout

设置会话过期，并保存

> 返回示例

```json
{
  "message": "登出成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 提问

POST /api/question

需要登录，否则不能使用该api，需要传入question，本函数会插入问题，并记录提问者，同时会默认由ai 先产生一个回答，作为默认回答

> Body 请求参数

```yaml
question: say"1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» question|body|string| 否 |none|

> 返回示例

```json
{
  "message": "提问成功"
}
```

```json
{
  "error": "无法生成答案"
}
```

```json
{
  "error": "未登录"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 回答

POST /api/answer

需要登录，否则不能使用该api,传入想回答的问题id，并插入答案，同时记录回答者

> Body 请求参数

```yaml
id: "1"
answer: HDUHELP is a computer learning community

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» id|body|string| 否 |none|
|» answer|body|string| 否 |none|

> 返回示例

```json
{
  "message": "回答成功"
}
```

```json
{
  "error": "回答失败"
}
```

```json
{
  "error": "未登录"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 点赞

POST /api/user/like

需要登录，否则不能使用该api，传入要点点赞的答案id，完成点赞功能

> Body 请求参数

```yaml
id: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» id|body|string| 否 |none|

> 返回示例

```json
"{\n    \"message\": \"点赞成功\"\n}\n{\"likes_count\": count}"
```

```json
{
  "error": "err"
}
```

```json
{
  "error": "未登录"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 查看问题和回答

GET /api/user/show

需要登录，通过审核的问题和回答，会加入该页面，以问题和答案的对应格式输出（可以一条问题对应多个回答）

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

```json
[
  {
    "id": 1,
    "question_text": "what is hduhelp",
    "total_likes": 1,
    "answers": [
      {
        "answer_id": 1,
        "answer_text": "HDUhelp could refer to different things depending on the context.\n\nIn the context of Hangzhou Dianzi University (HDU), it might be some kind of assistance or support related to the university. For example, it could be an online platform, a student - run help service, or an internal system within the university designed to provide help with various aspects such as academic support (course - related questions, study resources), campus life guidance (dormitory issues, campus facilities), or administrative assistance.\n\nIt could also potentially be the name of a self - developed software or service by HDU students or faculty for the benefit of the HDU community.",
        "respondent": "ai",
        "likes_count": 1
      }
    ],
    "questioner": "111"
  }
]
```

```json
[]
```

```json
{
  "error": "未登录"
}
```

```json
{
  "error": "err.Error()"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 管理员查看所有未审核的问题和回答

GET /api/admin

只有登录并且用户名为admin 的用户才能查看，会显示所有问题和答案，不管审核还是未审核

> 返回示例

```json
[
  {
    "id": 1,
    "question_text": "what is hduhelp",
    "total_likes": 1,
    "answers": [
      {
        "answer_id": 1,
        "answer_text": "HDUhelp could refer to different things depending on the context.\n\nIn the context of Hangzhou Dianzi University (HDU), it might be some kind of assistance or support related to the university. For example, it could be an online platform, a student - run help service, or an internal system within the university designed to provide help with various aspects such as academic support (course - related questions, study resources), campus life guidance (dormitory issues, campus facilities), or administrative assistance.\n\nIt could also potentially be the name of a self - developed software or service by HDU students or faculty for the benefit of the HDU community.",
        "respondent": "ai",
        "likes_count": 1
      }
    ],
    "questioner": "111"
  }
]
```

```json
{
  "error": "未登录"
}
```

```json
{
  "error": "管理员才可以访问！"
}
```

```json
{
  "error": "err.Error()"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 审核问题

POST /api/admin/checkQuestion

只有登录并且用户名为admin 的用户才能查看，输入问题id，和审核参数（同意为1）

> Body 请求参数

```yaml
id: "1"
check: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» id|body|string| 否 |ID 编号|
|» check|body|string| 否 |none|

> 返回示例

```json
{
  "message": "审核状态更新成功"
}
```

```json
{
  "error": "无效的审核状态"
}
```

```json
{
  "error": "管理员才可以访问！"
}
```

```json
{
  "error": "未登录"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 审核回答

POST /api/admin/checkAnswer

只有登录并且用户名为admin 的用户才能查看，输入答案id，和审核参数（同意为1）

> Body 请求参数

```yaml
id: "1"
check: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» id|body|string| 否 |ID 编号|
|» check|body|string| 否 |none|

> 返回示例

```json
{
  "message": "审核状态更新成功"
}
```

```json
{
  "error": "无效的审核状态"
}
```

```json
{
  "error": "err.Error()"
}
```

```json
{
  "error": "管理员才可以访问！"
}
```

```json
{
  "error": "未登录"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 按点赞数量排序

GET /api/user/like_sort

需要登录，会根据题目和答案的数量由高到低来排序，同时要求审核通过

> 返回示例

```json
[
  {
    "id": 1,
    "question_text": "what is hduhelp",
    "total_likes": 1,
    "answers": [
      {
        "answer_id": 1,
        "answer_text": "HDUhelp could refer to different things depending on the context.\n\nIn the context of Hangzhou Dianzi University (HDU), it might be some kind of assistance or support related to the university. For example, it could be an online platform, a student - run help service, or an internal system within the university designed to provide help with various aspects such as academic support (course - related questions, study resources), campus life guidance (dormitory issues, campus facilities), or administrative assistance.\n\nIt could also potentially be the name of a self - developed software or service by HDU students or faculty for the benefit of the HDU community.",
        "respondent": "ai",
        "likes_count": 1
      },
      {
        "answer_id": 2,
        "answer_text": "HDUHELP is a computer learning community",
        "respondent": "lili",
        "likes_count": 0
      }
    ],
    "questioner": "111"
  }
]
```

```json
{
  "error": "未登录"
}
```

```json
{
  "error": "err.Error()"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 修改问题

POST /api/question/changeQuestion

需要登录，只有本作者才能修改，传入问题的id，和更新后的内容

> Body 请求参数

```yaml
id: "1"
question: what is vidar_team

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» id|body|string| 否 |ID 编号|
|» question|body|string| 否 |none|

> 返回示例

```json
{
  "message": "修改成功"
}
```

```json
{
  "error": "err.Error()"
}
```

```json
{
  "error": "只有本作者才可以修改！"
}
```

```json
{
  "error": "未登录"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 修改答案

POST /api/answer/changeAnswer

需要登录，只有本作者才能修改，传入答案的id，和更新后的内容

> Body 请求参数

```yaml
id: "1"
answer: I do not know

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» id|body|string| 否 |ID 编号|
|» answer|body|string| 否 |none|

> 返回示例

```json
{
  "message": "修改成功"
}
```

```json
{
  "error": "err"
}
```

```json
{
  "error": "未登录"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 删除问题

POST /api/question/deleteQuestion

需要登录，只有本作者才能删除，传入问题的id

> Body 请求参数

```yaml
id: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» id|body|string| 否 |none|

> 返回示例

```json
{
  "message": "删除成功"
}
```

```json
{
  "error": "只有本作者才可以删除！"
}
```

```json
{
  "error": "err.Error()"
}
```

```json
{
  "error": "未登录"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 删除答案

POST /api/answer/deleteAnswer

需要登录，只有本作者才能删除，传入问题的id

> Body 请求参数

```yaml
id: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» id|body|string| 否 |none|

> 返回示例

```json
{
  "message": "删除成功"
}
```

```json
{
  "error": "err.Error()"
}
```

```json
{
  "error": "只有本作者才可以删除！"
}
```

```json
{
  "error": "未登录"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 主界面

## GET 主界面

GET /127.0.0.1:8000

接入了index.html，可以选择登录或者注册

> 返回示例

```json
"index.html"
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 登录

## POST 登录

POST /api/user/login

输入账号和密码，在sql里面插入对应的数据

> Body 请求参数

```yaml
username: test
password: testtest

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» username|body|string| 否 |none|
|» password|body|string| 否 |none|

> 返回示例

```json
{
  "message": "登录成功"
}
```

```json
{
  "error": "用户名或密码错误"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 数据模型

