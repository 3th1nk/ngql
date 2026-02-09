# NebulaGraph 图模式 (Graph Patterns)

图模式是 NebulaGraph nGQL 中用于描述图结构的语法，主要用于 MATCH 查询语句中。

## 目录

- [顶点模式 (Vertex Pattern)](#顶点模式-vertex-pattern)
- [边模式 (Edge Pattern)](#边模式-edge-pattern)
- [路径模式 (Path Pattern)](#路径模式-path-pattern)
- [多点关联模式](#多点关联模式)
- [Tag 模式](#tag-模式)
- [在查询中使用](#在查询中使用)

## 顶点模式 (Vertex Pattern)

顶点模式用于匹配图中的顶点。

### 基本语法

```
(variable[:tag1[:tag2...]] [{prop: value, ...}])
```

### 使用示例

```go
// 匿名顶点（匹配任意顶点）
v1 := NewVertexPattern("")
// 输出: ()

// 带变量的顶点
v2 := NewVertexPattern("v")
// 输出: (v)

// 带 Tag 的顶点
v3 := NewVertexPattern("v").Tag("player")
// 输出: (v:player)

// 带多个 Tag 的顶点
v4 := NewVertexPattern("v").Tag("player").Tag("person")
// 输出: (v:player:person)

// 带属性过滤的顶点
v5 := NewVertexPattern("v").
    Tag("player").
    Property("name", "Tim Duncan").
    Property("age", 42)
// 输出: (v:player {name: "Tim Duncan", age: 42})
```

## 边模式 (Edge Pattern)

边模式用于匹配图中的边（关系）。

### 基本语法

```
-[variable[:edge_type1[|edge_type2...]] [{prop: value, ...}] [*[min..max]]]-{>|<|-}
```

### 边的方向

- `->` : 右向边（从左到右）
- `<-` : 左向边（从右到左）
- `-` : 双向边（不限方向）

### 使用示例

#### 基本边模式

```go
// 右向边
e1 := NewEdgePattern(EdgeDirectionRight)
// 输出: ->

// 左向边
e2 := NewEdgePattern(EdgeDirectionLeft)
// 输出: <-

// 双向边
e3 := NewEdgePattern(EdgeDirectionBoth)
// 输出: -

// 带类型的边
e4 := NewEdgePattern(EdgeDirectionRight).Type("follow")
// 输出: -[:follow]->

// 带变量和类型的边
e5 := NewEdgePattern(EdgeDirectionRight).
    Variable("e").
    Type("follow")
// 输出: -[e:follow]->

// 带多个类型的边（匹配任意一种类型）
e6 := NewEdgePattern(EdgeDirectionRight).
    Type("follow").
    Type("like")
// 输出: -[:follow|like]->

// 带属性过滤的边
e7 := NewEdgePattern(EdgeDirectionRight).
    Type("follow").
    Property("degree", 90)
// 输出: -[:follow {degree: 90}]->
```

#### 变长边模式

变长边用于匹配多跳关系。

```go
// 固定步数（恰好 2 跳）
e1 := NewEdgePattern(EdgeDirectionRight).
    Type("follow").
    Hops(2, 2)
// 输出: -[:follow*2]->

// 范围步数（1 到 3 跳）
e2 := NewEdgePattern(EdgeDirectionRight).
    Type("follow").
    Hops(1, 3)
// 输出: -[:follow*1..3]->

// 最小步数（至少 1 跳）
e3 := NewEdgePattern(EdgeDirectionRight).
    Type("follow").
    Hops(1, 0)
// 输出: -[:follow*1..]->

// 最大步数（最多 3 跳）
e4 := NewEdgePattern(EdgeDirectionRight).
    Type("follow").
    Hops(0, 3)
// 输出: -[:follow*..3]->
```

## 路径模式 (Path Pattern)

路径模式由顶点和边组成，描述完整的图遍历路径。

### 基本语法

```
[variable =] (vertex)-[edge]-(vertex)[-[edge]-(vertex)...]
```

### 使用示例

#### 简单路径

```go
// 两个顶点之间的简单路径
p1 := NewPathPattern().
    Vertex("v1").
    Edge(EdgeDirectionRight).
    Vertex("v2")
// 输出: (v1)->(v2)

// 带路径变量
p2 := NewPathPattern().
    Variable("p").
    Vertex("v1").
    Edge(EdgeDirectionRight).
    Vertex("v2")
// 输出: p = (v1)->(v2)
```

#### 带 Tag 的路径

```go
// 使用便捷方法
p1 := NewPathPattern().
    VertexWithTag("v1", "player").
    EdgeWithType(EdgeDirectionRight, "follow").
    VertexWithTag("v2", "player")
// 输出: (v1:player)-[:follow]->(v2:player)
```

#### 复杂路径

```go
// 多跳路径
p1 := NewPathPattern().
    Variable("p").
    VertexWithTag("v1", "player").
    EdgeWithType(EdgeDirectionRight, "follow").
    VertexWithTag("v2", "player").
    EdgeWithType(EdgeDirectionRight, "serve").
    VertexWithTag("v3", "team")
// 输出: p = (v1:player)-[:follow]->(v2:player)-[:serve]->(v3:team)

// 使用自定义模式
vertex := NewVertexPattern("v1").
    Tag("player").
    Property("name", "Tim Duncan")

edge := NewEdgePattern(EdgeDirectionRight).
    Variable("e").
    Type("follow").
    Hops(1, 3)

path := NewPathPattern().
    Variable("p").
    VertexPattern(vertex).
    EdgePattern(edge).
    Vertex("v2")
// 输出: p = (v1:player {name: "Tim Duncan"})-[e:follow*1..3]->(v2)
```

## 多点关联模式

多点关联模式用于描述多个顶点之间的复杂关系，可以包含多条边和多个顶点。

### 空边模式

当不需要指定边的类型、属性或变量时，可以使用空边模式。

```go
// 简单的空边连接
p1 := NewPathPattern().
    Vertex("a").
    Edge(EdgeDirectionRight).
    Vertex("b")
// 输出: (a)->(b)

// 带 Tag 的空边连接
p2 := NewPathPattern().
    VertexWithTag("a", "User").
    Edge(EdgeDirectionRight).
    VertexWithTag("b", "Post")
// 输出: (a:User)->(b:Post)
```

### 三点两边模式

描述三个顶点通过两条边连接的关系。

```go
// 三个顶点的线性连接
p1 := NewPathPattern().
    Vertex("a").
    Edge(EdgeDirectionRight).
    Vertex("b").
    Edge(EdgeDirectionLeft).
    Vertex("c")
// 输出: (a)->(b)<-(c)

// 带边类型的三点连接
p2 := NewPathPattern().
    VertexWithTag("user", "User").
    EdgeWithType(EdgeDirectionRight, "CREATED").
    VertexWithTag("post", "Post").
    EdgeWithType(EdgeDirectionLeft, "LIKED").
    VertexWithTag("other", "User")
// 输出: (user:User)-[:CREATED]->(post:Post)<-[:LIKED]-(other:User)
```

### 匿名顶点模式

在路径中使用匿名顶点（不指定变量名）来表示中间节点。

```go
// 匿名中间顶点
p1 := NewPathPattern().
    Vertex("a").
    Edge(EdgeDirectionRight).
    Vertex("").  // 匿名顶点
    Edge(EdgeDirectionLeft).
    Vertex("c")
// 输出: (a)->()<-(c)

// 多个匿名顶点
p2 := NewPathPattern().
    Vertex("").
    Edge(EdgeDirectionRight).
    Vertex("").
    Edge(EdgeDirectionRight).
    Vertex("")
// 输出: ()->()->()

// 带类型的匿名顶点路径
p3 := NewPathPattern().
    VertexWithTag("start", "User").
    EdgeWithType(EdgeDirectionRight, "FOLLOWS").
    Vertex("").  // 中间用户（不关心具体是谁）
    EdgeWithType(EdgeDirectionRight, "FOLLOWS").
    VertexWithTag("end", "User")
// 输出: (start:User)-[:FOLLOWS]->()-[:FOLLOWS]->(end:User)
```

### 实际应用场景

```go
// 场景1: 查找共同好友
// 用户 A 和用户 C 都关注了某个中间用户
pattern := NewPathPattern().
    VertexWithTag("userA", "User").
    EdgeWithType(EdgeDirectionRight, "FOLLOWS").
    Vertex("").  // 共同关注的用户
    EdgeWithType(EdgeDirectionLeft, "FOLLOWS").
    VertexWithTag("userC", "User")

// 场景2: 内容审核链路
// 用户创建内容，内容被审核员审核
pattern := NewPathPattern().
    VertexWithTag("creator", "User").
    EdgeWithType(EdgeDirectionRight, "CREATED").
    VertexWithTag("content", "Post").
    EdgeWithType(EdgeDirectionLeft, "REVIEWED").
    VertexWithTag("moderator", "Moderator")
```

## Tag 模式

Tag 模式用于指定顶点的类型标签，一个顶点可以有多个 Tag。

### 单 Tag 模式

为顶点指定单个 Tag 标签。

```go
// 基本单 Tag
p1 := NewPathPattern().
    VertexWithTag("a", "User").
    Edge(EdgeDirectionRight).
    Vertex("b")
// 输出: (a:User)->(b)

// 带边类型的单 Tag
p2 := NewPathPattern().
    VertexWithTag("user", "User").
    EdgeWithType(EdgeDirectionRight, "POSTED").
    VertexWithTag("content", "Post")
// 输出: (user:User)-[:POSTED]->(content:Post)
```

### 多 Tag 模式

一个顶点可以同时拥有多个 Tag，用于表示多重身份或类型。

```go
// 双 Tag 顶点
v1 := NewVertexPattern("a").
    Tag("User").
    Tag("Admin")
p1 := NewPathPattern().
    VertexPattern(v1).
    Edge(EdgeDirectionRight).
    Vertex("b")
// 输出: (a:User:Admin)->(b)

// 三 Tag 顶点
v2 := NewVertexPattern("a").
    Tag("User").
    Tag("Admin").
    Tag("VIP")
p2 := NewPathPattern().
    VertexPattern(v2).
    EdgeWithType(EdgeDirectionRight, "CREATED").
    VertexWithTag("b", "Post")
// 输出: (a:User:Admin:VIP)-[:CREATED]->(b:Post)
```

### 复杂 Tag 组合

在实际应用中，可以组合使用多个 Tag 来精确匹配特定类型的顶点。

```go
// 多 Tag 顶点与多 Tag 目标
admin := NewVertexPattern("admin").
    Tag("User").
    Tag("Admin").
    Tag("VIP")

article := NewVertexPattern("content").
    Tag("Post").
    Tag("Article")

edge := NewEdgePattern(EdgeDirectionRight).
    Type("CREATED")

pattern := NewPathPattern().
    VertexPattern(admin).
    EdgePattern(edge).
    VertexPattern(article)
// 输出: (admin:User:Admin:VIP)-[:CREATED]->(content:Post:Article)
```

### 实际应用场景

```go
// 场景1: 权限管理
// 查找具有管理员和审核员双重身份的用户创建的内容
pattern := NewPathPattern().
    VertexPattern(
        NewVertexPattern("user").
            Tag("User").
            Tag("Admin").
            Tag("Moderator"),
    ).
    EdgeWithType(EdgeDirectionRight, "CREATED").
    VertexWithTag("content", "Post")

// 场景2: 内容分类
// 查找技术博客文章（同时具有 Post 和 TechArticle 标签）
pattern := NewPathPattern().
    VertexWithTag("author", "User").
    EdgeWithType(EdgeDirectionRight, "PUBLISHED").
    VertexPattern(
        NewVertexPattern("article").
            Tag("Post").
            Tag("TechArticle").
            Tag("Featured"),
    )

// 场景3: 用户推荐
// 查找 VIP 用户关注的其他 VIP 用户
pattern := NewPathPattern().
    VertexPattern(
        NewVertexPattern("vip1").
            Tag("User").
            Tag("VIP"),
    ).
    EdgeWithType(EdgeDirectionRight, "FOLLOWS").
    VertexPattern(
        NewVertexPattern("vip2").
            Tag("User").
            Tag("VIP"),
    )
```

### Tag 模式注意事项

1. **Tag 顺序**：多个 Tag 的顺序会影响匹配结果，建议按照重要性排序
2. **Tag 数量**：虽然支持多个 Tag，但过多的 Tag 可能影响查询性能
3. **Tag 命名**：Tag 名称应遵循 NebulaGraph 的命名规范，建议使用大驼峰命名
4. **Tag 组合**：合理使用 Tag 组合可以实现更精确的查询，避免在 WHERE 子句中过滤

## 在查询中使用

### MATCH 查询

```go
// 创建路径模式
pattern := NewPathPattern().
    VertexWithTag("v1", "player").
    EdgeWithType(EdgeDirectionRight, "follow").
    VertexWithTag("v2", "player")

// 在 MATCH 查询中使用
query := NewMatch(pattern).
    Where("v1.player.age > 30").
    Return("v1.player.name", "v2.player.name").
    OrderBy("v1.player.age", true).
    Limit(10).
    String()

// 输出:
// MATCH (v1:player)-[:follow]->(v2:player)
// WHERE v1.player.age > 30
// RETURN v1.player.name, v2.player.name
// ORDER BY v1.player.age DESC
// LIMIT 10
```

### 复杂查询示例

```go
// 查找朋友的朋友
pattern := NewPathPattern().
    VertexWithTag("me", "player").
    EdgePattern(
        NewEdgePattern(EdgeDirectionRight).
            Type("follow").
            Hops(2, 2),
    ).
    VertexWithTag("friend_of_friend", "player")

query := NewMatch(pattern).
    Where("me.player.name == \"Tim Duncan\"").
    Return("friend_of_friend.player.name").
    String()

// 查找最短路径
pattern := NewPathPattern().
    Variable("p").
    VertexWithTag("start", "player").
    EdgePattern(
        NewEdgePattern(EdgeDirectionRight).
            Type("follow").
            Hops(1, 5),
    ).
    VertexWithTag("end", "player")

query := NewMatch(pattern).
    Where("start.player.name == \"Tim Duncan\" AND end.player.name == \"Tony Parker\"").
    Return("p").
    OrderBy("length(p)", false).
    Limit(1).
    String()
```

## API 参考

### VertexPattern

| 方法 | 说明 | 示例 |
|------|------|------|
| `NewVertexPattern(variable string)` | 创建顶点模式 | `NewVertexPattern("v")` |
| `Tag(tag string)` | 添加 Tag 标签 | `.Tag("player")` |
| `Property(key string, value interface{})` | 添加属性过滤 | `.Property("age", 30)` |
| `String()` | 生成模式字符串 | `.String()` |

### EdgePattern

| 方法 | 说明 | 示例 |
|------|------|------|
| `NewEdgePattern(direction EdgeDirection)` | 创建边模式 | `NewEdgePattern(EdgeDirectionRight)` |
| `Variable(variable string)` | 设置边变量名 | `.Variable("e")` |
| `Type(edgeType string)` | 添加边类型 | `.Type("follow")` |
| `Property(key string, value interface{})` | 添加属性过滤 | `.Property("degree", 90)` |
| `Hops(min, max int)` | 设置变长边跳数 | `.Hops(1, 3)` |
| `String()` | 生成模式字符串 | `.String()` |

### PathPattern

| 方法 | 说明 | 示例 |
|------|------|------|
| `NewPathPattern()` | 创建路径模式 | `NewPathPattern()` |
| `Variable(variable string)` | 设置路径变量名 | `.Variable("p")` |
| `Vertex(variable string)` | 添加顶点 | `.Vertex("v1")` |
| `VertexWithTag(variable, tag string)` | 添加带 Tag 的顶点 | `.VertexWithTag("v1", "player")` |
| `VertexPattern(vertex *VertexPattern)` | 添加自定义顶点模式 | `.VertexPattern(vertex)` |
| `Edge(direction EdgeDirection)` | 添加边 | `.Edge(EdgeDirectionRight)` |
| `EdgeWithType(direction EdgeDirection, edgeType string)` | 添加带类型的边 | `.EdgeWithType(EdgeDirectionRight, "follow")` |
| `EdgePattern(edge *EdgePattern)` | 添加自定义边模式 | `.EdgePattern(edge)` |
| `String()` | 生成模式字符串 | `.String()` |

## 注意事项

1. **属性值格式化**：属性值会自动根据类型进行格式化（字符串加引号，数字直接输出等）
2. **变量命名**：变量名应遵循 NebulaGraph 的命名规范
3. **Tag 和边类型**：可以添加多个 Tag 或边类型，用于匹配多种可能
4. **变长边性能**：使用变长边时注意性能影响，建议设置合理的跳数范围
5. **路径变量**：路径变量可用于后续的 RETURN 或 WHERE 子句中

## 参考文档

- [NebulaGraph 官方文档 - 图模式](https://docs.nebula-graph.com.cn/3.8.0/3.ngql-guide/1.nGQL-overview/3.graph-patterns/)
- [NebulaGraph 官方文档 - MATCH 语句](https://docs.nebula-graph.com.cn/3.8.0/3.ngql-guide/7.general-query-statements/2.match/)
