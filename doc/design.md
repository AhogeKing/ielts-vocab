# 单词积累 web 程序

## 系统的核心概念

### 多类型条目

统一使用 **Entry** 表示记录的学习内容。

- WORD 单词
- PHRASE 短语、搭配、习语
- SENTENCE 范文句子、口语表达、写作句型
- ARTICLE 整篇范文、文章

例如：

| 内容                         | 类型       | 词库     |
|----------------------------|----------|--------|
| detrimental                | WORD     | 写作、阅读  |
| pose a threat to           | PHRASE   | 写作     |
| This is largely because... | SENTENCE | 写作     |
| under the hood             | PHASE    | 英文技术文档 |
| I'm not really into        | SENTENCE | 口语     |

**一个条目允许属于多个词库**

例如 *significant* 即可能来自阅读，也可能用于写作。没有必要复制两份记录。

## 第一版 MVP 功能

### 1. 快速记录

新增条目时，最少只要求填写：

- 内容
- 内容类型
- 所属词库 (默认为杂项)

其他字段全部可选：

- 中文释义
- 英文解释
- 词性
- 音标
- 例句
- 笔记
- 来源
- 标签

新增页面提供的两种模式：

- **普通编辑模式**
  
  适合课后认真整理：

  - 完整填写释义
  - 添加例句
  - 添加标签
  - 修改来源
  - 设置掌握状态


- **快速记录模式**

    适合上课和做题时：
    
    ```text
    当前词库：听力
    当前来源：IELTS 课程 2026-07-21
    
    [输入内容]
    [可选中文释义]
    [保存并继续]
    ```
    
    保存后输入框立即清空，并自动沿用上一个词库和来源。

### 2. 条目列表

列表页面需要支持：

- 按六个词库筛选
- 按单词、短语、句子筛选
- 按标签筛选
- 按掌握状态筛选
- 按创建时间排序
- 按更新时间排序
- 收藏筛选
- 分页
- 搜索

采用卡片或紧凑列表形式：

```text
pose a threat to                         [短语]
对……构成威胁

词库：写作、阅读
标签：环境、负面影响、Task 2
来源：雅思写作课 07-21

例句：
Plastic waste poses a serious threat to marine ecosystems.

状态：学习中
```

### 3. 搜索

第一版搜索：

- 英文内容
- 中文释义
- 笔记
- 例句
- 来源
- 标签

初期数据量通常不会特别大，使用 PostgreSQL 的 `ILKE '%keyword%'`。

以后可以增加：

- 拼写错误容忍
- 相似词搜索
- 搜索结果相关度排序

PostgreSQL 的 `pg_trgm` 拓展可以进行基于 trigram 的文本相似度匹配，适合类似 “拼错一点也能搜出来” 的场景；
PostgreSQL 自带全文搜索和 GIN 索引，不过第一版暂时不做这个功能。

### 4. 批量导入和导出

#### 批量导入

允许粘贴：

```text
detrimental
pose a threat to
This is largely because...
be exposed to
```

系统先生成四条草稿，再由用户统一指定。

```text
词库：写作
来源：7 月 21 日雅思课
标签：环境
```

#### CSV 导入

暂定支持：

- CSV：方便用 Excel 查看
- JSON：完整备份和未来迁移
- Markdown：方便生成个人笔记

备份功能对于个人长期项目比复杂的数据统计更重要。

## 数据库设计

整体关系可以表示为：

```text
User
 ├── Collection
 ├── Entry
 │    ├── Example
 │    ├── EntryCollection ── Collection
 │    └── EntryTag ── Tag
 └── Session
```

### 1. users

尽管程序目前只有作者一个用户，还是先做用户表吧。

```text
users
-----
id
username
email
password_hash
created_at
updated_at
```

这样以后可以：

- 在不同设备登陆
- 给其他人使用
- 区分每个用户的词库
- 支持公开分享

### 2. collections

```text
collections
-----------
id
user_id
code
name
description
sort_order
is_system
created_at
updated_at
```

初始化六条数据：

- LISTENING 听力
- SPEAKING  口语
- READING   阅读
- WRITING   写作
- TECH_DOCS 英文技术文档
- MISC      零散内容

`is_system` 用于标记六个默认词库，但用户以后仍然可以新增自定义词库。

### 3. entries

这是最核心的 **条目** 表

```text
entries
-------
id
user_id
content
normalized_content
content_type
translation_zh
definition_en
part_of_speech
phonetic
context_text
notes
source_name
source_detail
learning_status
is_favorite
created_at
updated_at
deleted_at
```

字段说明：

#### content

真正记录的内容：

```text
detrimental
pose a threat to
This is largely because...
```

#### normalized_content

用于搜索和重复检测，例如：

```text
原内容：Pose   a Threat To
标准化：pose a threat to
```

标准化操作可以包括：

- 转小写
- 去除首尾空格
- 连续空格合并
- 同一部份标点

但不修改用户看到的原始 `content`。

#### content_type

使用字符串家数据库约束：

- WORD
- PHRASE
- SENTENCE

相比 PostgreSQL Enum，字符串加 `CHECK` 约束将来增加新类型更方便。

#### content_text

记录第一次看到它的原句：

    The increasing amount of plastic waste poses a threat to marine life.

#### notes

可以支持 Markdown，用来记录：

```markdown
- pose 的第三人称是 poses
- 常用于负面影响
- threat 前面可用 serious / potential / major
```

#### learning_status

第一版暂时只需要：

- NEW
- LEARNING
- MASTERED
- IGNORED

暂时不实现复杂的记忆曲线。

### 4. entry_collections

用于支持一个条目属于多个词库：

```text
entry_collections
-----------------
entry_id
collection_id
created_at
```

联合主键：

```sql
PRIMARY KEY (entry_id, collection_id)
```

### 5. tags and entry_tags

词库负责粗分类，标签负责细分类。

例如：

```text
词库：写作
标签：Task 2、环境、因果、负面影响
```

表结构：

```text
tags
----
id
user_id
name
normalized_name
is_system
created_at
```

```text
entry_tags
----------
entry_id
tag_id
```

标签非常适合区分写作句子的功能：

```text
让步
因果
举例
总结
趋势上升
趋势下降
口语 Part 1
口语 Part 2
计算机网络
Go
PostgreSQL
```

### 6. examples

一条记录可以有多个例句：

```text
examples
--------
id
entry_id
sentence
translation_zh
notes
sort_order
created_at
updated_at
```

例如 `detrimental`：

```text
Excessive screen time can be detrimental to children's health.

The policy may have a detrimental effect on small businesses.
```

对于 `SENTENCE ARTICLE` 类型的条目，其本身就是需要“例句”，因此可以没有额外例句。

### 7. sessions

多设备登陆使用服务端 Session，而不是把长期 JWT 放进 `localStorage`。

```text
sessions
--------
id
user_id
token_hash
expires_at
last_seen_at
created_at
```

浏览器保存一个 `HttpOnly + Secure` Cookie，数据库保存 Token 的哈希。

优点是：

- 可以单独退出某台设备
- 可以查看当前登录设备
- Token 泄漏后能够主动撤销
- 前端不需要操作认证 Token

部署时让前端和 API 使用同一个域名。

## 页面设计

第一版暂时只有少数页面。

### 1. 登录注册页

用户名和密码必填。邮箱选填。

### 2. Dashboard

暂时展示真正有用的信息：

```text
今日新增：12
本周新增：48
待整理：17
学习中：126
已掌握：83

六个词库的条目数量
最近添加的十条内容
快速记录入口
```

### 3. 条目列表页

左侧或顶部显示：

```text
全部
听力
口语
阅读
写作
英文技术文档
零散
```

筛选区域：

```text
搜索框
内容类型
掌握状态
标签
收藏
排序
```

### 4. 条目详情 / 编辑页

分为几个区块：

```text
基本内容
释义与解释
所属词库
标签
语境与例句
来源
学习状态
```

新增和编辑共用一个表单组件。

### 5. 快速记录页

这是移动端也该优化的页面。

```text
当前词库
当前来源
内容输入框
简要释义
保存并继续
```

后期可以增加一些快捷键。

## 5. 后端接口设计

统一前缀：`/api/v1`

### 认证

```http request
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
GET    /api/v1/auth/me
GET    /api/v1/auth/sessions
DELETE /api/v1/auth/sessions/:id
```

### 词库

```http request
GET    /api/v1/collections
POST   /api/v1/collections
PATCH  /api/v1/collections/:id
DELETE /api/v1/collections/:id
```

系统默认词库可以禁止删除，或者删除前要求转移其中的条目。

### 条目

```http request
POST   /api/v1/entries
GET    /api/v1/entries
GET    /api/v1/entries/:id
PATCH  /api/v1/entries/:id
DELETE /api/v1/entries/:id
```

列表查询示例：

`GET /api/v1/entries?collection_id=4&type=PHRASE&status=LEARNING&q=threat&page=1&page_size=20`

### 批量操作

```http request
POST   /api/v1/entries/batch
PATCH  /api/v1/entries/batch/status
DELETE /api/v1/entries/batch
```

### 标签

```http request
GET    /api/v1/tags?q=env
POST   /api/v1/tags
DELETE /api/v1/tags/:id
```

### 导入导出

```http request
POST /api/v1/imports/csv
GET  /api/v1/exports/csv
GET  /api/v1/exports/loginJson
GET  /api/v1/exports/markdown
```

## 新增条目请求示例

```JSON
{
  "content": "pose a threat to",
  "contentType": "PHRASE",
  "translationZh": "对……构成威胁",
  "definitionEn": "to create a risk or danger to something",
  "partOfSpeech": "verb phrase",
  "contextText": "Plastic waste poses a threat to marine ecosystems.",
  "notes": "常用于雅思写作中的负面影响论证。",
  "source": {
    "name": "IELTS writing class",
    "detail": "2026-07-21 environmental topics"
  },
  "collectionIds": [3, 4],
  "tags": [
    "environment",
    "negative impact",
    "Task 2"
  ],
  "examples": [
    {
      "sentence": "Climate change poses a serious threat to coastal communities.",
      "translationZh": "气候变化对沿海社区构成严重威胁。"
    }
  ]
}
```

响应使用统一结构：

```JSON
{
  "success": true,
  "data": {
    "id": 1024
  },
  "error": null
}
```

## 技术栈

### 后端

- Gin
- PostgreSQL
- GORM
- 数据库迁移工具

### 前端

- Vue 3
- TypeScript
- Vite
- Vue Router
- Pinia
- Axios

## Go 项目目录结构

```text
ielts-vault/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   ├── entry/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── dto.go
│   │   └── model.go
│   ├── collection/
│   ├── tag/
│   ├── importexport/
│   ├── middleware/
│   ├── database/
│   │   ├── db.go
│   │   └── generated/
│   ├── config/
│   └── server/
│       └── router.go
├── db/
│   ├── migrations/
│   └── queries/
├── web/
│   └── 前端项目
├── sqlc.yaml
├── compose.yaml
├── go.mod
└── README.md
```

### 请求流程：

```text
Handler
   ↓
Service
   ↓
Repository / sqlc
   ↓
PostgreSQL
```

### 各层职责

#### Handler

- 从 Gin 获取参数
- Binding 和 Validation
- 调用 Service
- 转换 HTTP 状态码
- 返回 JSON

#### Service

- 权限和业务规则
- 重复内容检测
- 创建标签
- 处理多词库关系
- 开启事务
- 组织返回结果

#### Repository

- 执行 SQL
- 不包含 HTTP 逻辑
- 尽量不包含复杂业务

## 开发顺序

### 最小版本

先只完成：

1. 数据库连接
2. entries 表
3. collections 表
4. 新增条目接口
5. 条目列表接口
6. 一个最简单的前端新增页面
7. 一个最简单的前端列表页面

目标：

```text
浏览器输入
→ Gin 接收
→ Service 处理
→ PostgreSQL 保存
→ 列表查询
→ 浏览器展示
```

### 核心整理功能

加入：

```text
编辑和删除
多个词库
标签
例句
来源
分页
筛选
搜索
重复内容提醒
```

重复内容不直接禁止，例如同一个单词可能有不同语境。创建时可以返回：

```text
发现 2 条相似内容：
1. threat
2. pose a threat to

仍然创建 / 查看已有记录
```

### 真正投入使用

加入：

```text
用户登录
移动端适配
快速记录模式
批量粘贴
CSV 导入
JSON 备份
部署和 HTTPS
```

### 复习功能

等积累一段时间的数据后再实现：

```text
今日复习
随机抽查
隐藏中文看英文
隐藏英文看中文
熟悉度评分
错词本
复习历史
简单间隔复习
```

此时再增加：

```text
review_logs
entry_learning_states
review_sessions
```

### AI 辅助

AI 应当辅助整理内容：

```text
自动识别单词、短语或句子
生成中文释义草稿
补充英文解释
生成 IELTS 风格例句
识别词性
推荐标签
判断适合听说读写中的哪些词库
发现重复或相似表达
将范文拆分为可积累表达
```

每次由 AI 产生的内容都先显示为草稿，由你确认后保存。
