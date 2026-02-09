# NebulaGraph nGQL 构建器使用指南

## 概述

本包提供了一套完整的 nGQL 语句构建器，支持链式调用，方便业务层构建各种图数据库操作语句。

## 文件结构

为了更好的代码组织和维护，构建器按功能拆分为多个文件：

```
├── util.go          # 公共工具函数（formatValue 等）
├── property.go      # 属性类型定义（PropertyType, Property, IndexProperty）
├── pattern.go       # 图模式构建器（VertexPattern, EdgePattern, PathPattern）
├── schema.go        # 图空间管理（CreateSpaceBuilder, DropSpaceBuilder 等）
├── tag.go           # Tag 管理（CreateTagBuilder, AlterTagBuilder 等）
├── edge.go          # Edge 类型管理（CreateEdgeBuilder, AlterEdgeBuilder 等）
├── index.go         # 索引管理（CreateTagIndexBuilder, CreateEdgeIndexBuilder 等）
├── insert.go        # 插入操作（InsertVertexBuilder, InsertEdgeBuilder）
├── query.go         # 查询操作（MatchBuilder, GoBuilder, FetchBuilder, LookupBuilder, GroupByBuilder）
├── update.go        # 更新操作（UpdateVertexBuilder, UpdateEdgeBuilder）
├── upsert.go        # UPSERT 操作（UpsertVertexBuilder, UpsertEdgeBuilder）
├── delete.go        # 删除操作（DeleteVertexBuilder, DeleteEdgeBuilder）
├── pattern_test.go  # 图模式构建器测试
├── pattern_example_test.go  # 图模式构建器示例
├── schema_test.go   # 图空间管理测试
├── tag_test.go      # Tag 管理测试
├── edge_test.go     # Edge 类型管理测试
├── index_test.go    # 索引管理测试
├── insert_test.go   # 插入操作测试
├── query_test.go    # 查询操作测试
├── update_test.go   # 更新操作测试
├── upsert_test.go   # UPSERT 操作测试
├── delete_test.go   # 删除操作测试
└── README.md        # 使用文档
```

## 功能特性

- ✅ **链式调用** - 代码简洁易读
- ✅ **类型安全** - 减少拼写错误
- ✅ **图模式构建器** - 提供 VertexPattern、EdgePattern、PathPattern 用于类型安全的图模式构建
- ✅ **自动转义** - 自动处理值转义和格式化，防止注入
- ✅ **模块化设计** - 每个文件职责单一，易于维护和扩展
- ✅ **完整测试** - 测试覆盖率 78%+

## 设计原则

### 1. 单一职责
每个文件负责一类特定的操作：
- `property.go` - 属性类型定义（PropertyType 枚举、Property 结构体、IndexProperty 结构体）
- `pattern.go` - 图模式构建器（VertexPattern、EdgePattern、PathPattern，用于类型安全的图模式构建）
- `schema.go` - 图空间管理（CREATE/DROP/USE/SHOW/DESCRIBE SPACE）
- `tag.go` - Tag 管理（CREATE/ALTER/DROP/SHOW/DESCRIBE TAG）
- `edge.go` - Edge 类型管理（CREATE/ALTER/DROP/SHOW/DESCRIBE EDGE）
- `index.go` - 索引管理（CREATE/DROP/REBUILD/SHOW/DESCRIBE INDEX）
- `insert.go` - 插入操作（INSERT VERTEX/EDGE）
- `query.go` - 查询操作（MATCH/GO/FETCH/LOOKUP/GROUP BY/FIND PATH/GET SUBGRAPH）
- `update.go` - 更新操作（UPDATE VERTEX/EDGE）
- `upsert.go` - UPSERT 操作（UPSERT VERTEX/EDGE）
- `delete.go` - 删除操作（DELETE VERTEX/EDGE）

### 2. Builder 模式
所有构建器都支持链式调用：
```go
nGQL := ngql.NewInsertVertex("player", "name", "age").
    Values("player100", "Tim Duncan", 42).
    String()
```

### 3. 自动转义
`formatValue()` 函数（位于 util.go）自动处理不同类型的值转换和转义，避免注入风险。

## 图模式构建器（Pattern）

图模式构建器提供了类型安全的方式来构建 nGQL 图模式，特别适用于 MATCH 查询。相比直接使用字符串，Pattern 构建器提供了更好的代码提示和类型检查。

### 顶点模式（VertexPattern）

```go
// 匿名顶点
v1 := ngql.NewVertexPattern("")
// 输出: ()

// 带变量的顶点
v2 := ngql.NewVertexPattern("v")
// 输出: (v)

// 带 Tag 的顶点
v3 := ngql.NewVertexPattern("v").Tag("player")
// 输出: (v:player)

// 带多个 Tag 的顶点
v4 := ngql.NewVertexPattern("v").Tag("player").Tag("person")
// 输出: (v:player:person)

// 带属性过滤的顶点
v5 := ngql.NewVertexPattern("v").
    Tag("player").
    Property("name", "Tim Duncan")
// 输出: (v:player {name: "Tim Duncan"})
```

### 边模式（EdgePattern）

```go
// 右向边
e1 := ngql.NewEdgePattern(ngql.EdgeDirectionRight)
// 输出: ->

// 左向边
e2 := ngql.NewEdgePattern(ngql.EdgeDirectionLeft)
// 输出: <-

// 双向边
e3 := ngql.NewEdgePattern(ngql.EdgeDirectionBoth)
// 输出: -

// 带类型的边
e4 := ngql.NewEdgePattern(ngql.EdgeDirectionRight).Type("follow")
// 输出: -[:follow]->

// 带变量和类型的边
e5 := ngql.NewEdgePattern(ngql.EdgeDirectionRight).
    Variable("e").
    Type("follow")
// 输出: -[e:follow]->

// 带多个类型的边（或关系）
e6 := ngql.NewEdgePattern(ngql.EdgeDirectionRight).
    Type("follow").
    Type("like")
// 输出: -[:follow|like]->

// 变长边 - 固定步数
e7 := ngql.NewEdgePattern(ngql.EdgeDirectionRight).
    Type("follow").
    Hops(2, 2)
// 输出: -[:follow*2]->

// 变长边 - 范围步数
e8 := ngql.NewEdgePattern(ngql.EdgeDirectionRight).
    Type("follow").
    Hops(1, 3)
// 输出: -[:follow*1..3]->
```

### 路径模式（PathPattern）

```go
// 简单路径
p1 := ngql.NewPathPattern().
    Vertex("v1").
    Edge(ngql.EdgeDirectionRight).
    Vertex("v2")
// 输出: (v1)->(v2)

// 带变量的路径
p2 := ngql.NewPathPattern().
    Variable("p").
    Vertex("v1").
    Edge(ngql.EdgeDirectionRight).
    Vertex("v2")
// 输出: p = (v1)->(v2)

// 带 Tag 的路径
p3 := ngql.NewPathPattern().
    VertexWithTag("v1", "player").
    EdgeWithType(ngql.EdgeDirectionRight, "follow").
    VertexWithTag("v2", "player")
// 输出: (v1:player)-[:follow]->(v2:player)

// 复杂路径
p4 := ngql.NewPathPattern().
    Variable("p").
    VertexWithTag("v1", "player").
    EdgeWithType(ngql.EdgeDirectionRight, "follow").
    VertexWithTag("v2", "player").
    EdgeWithType(ngql.EdgeDirectionRight, "serve").
    VertexWithTag("v3", "team")
// 输出: p = (v1:player)-[:follow]->(v2:player)-[:serve]->(v3:team)
```

### 在 MATCH 查询中使用 Pattern

```go
// 使用 Pattern 对象构建 MATCH 查询
pattern := ngql.NewPathPattern().
    VertexWithTag("v1", "player").
    EdgeWithType(ngql.EdgeDirectionRight, "follow").
    VertexWithTag("v2", "player")

query := ngql.NewMatch(pattern).
    Where("v1.player.age > 30").
    Return("v1.player.name", "v2.player.name").
    String()
// 输出: MATCH (v1:player)-[:follow]->(v2:player) WHERE v1.player.age > 30 RETURN v1.player.name, v2.player.name

// 复杂的图模式组合
vertex1 := ngql.NewVertexPattern("v1").
    Tag("player").
    Property("name", "Tim Duncan")

edge := ngql.NewEdgePattern(ngql.EdgeDirectionRight).
    Variable("e").
    Type("follow").
    Hops(1, 3)

path := ngql.NewPathPattern().
    Variable("p").
    VertexPattern(vertex1).
    EdgePattern(edge).
    Vertex("v2")
// 输出: p = (v1:player {name: "Tim Duncan"})-[e:follow*1..3]->(v2)
```

### Pattern 的优势

1. **类型安全**：编译时检查，减少拼写错误
2. **代码提示**：IDE 可以提供完整的方法提示
3. **易于维护**：复杂的图模式更容易理解和修改
4. **可复用**：可以将常用的模式封装成函数
5. **灵活组合**：可以自由组合顶点、边和路径模式

## 快速开始

### 1. 图空间管理

```go
import "github.com/3th1nk/ngql"

// 创建图空间
nGQL := ngql.NewCreateSpace("my_space").
    PartitionNum(10).
    ReplicaFactor(3).
    VidType("FIXED_STRING(32)").
    Comment("我的图空间").
    String()

// 创建图空间（克隆现有空间）
nGQL := ngql.NewCreateSpaceAs("new_space", "old_space").
    String()

// 使用图空间
nGQL := ngql.NewUseSpace("my_space").
    String()

// 显示所有图空间
nGQL := ngql.NewShowSpaces().
    String()

// 查看图空间详情
nGQL := ngql.NewDescribeSpace("my_space").
    String()

// 删除图空间
nGQL := ngql.NewDropSpace("my_space").
    IfExists().
    String()

// 清空图空间
nGQL := ngql.NewClearSpace("my_space").
    String()
```

### 2. Tag 管理

```go
// 创建 Tag
nGQL := ngql.NewCreateTag("player").
    IfNotExists().
    AddProperty(&ngql.Property{
        Name:     "name",
        Type:     "string",
        Nullable: false,
    }).
    AddProperty(&ngql.Property{
        Name:     "age",
        Type:     "int",
        Nullable: true,
        Default:  0,
    }).
    TTL("created_at", 100).
    Comment("球员信息").
    String()

// 简化方式创建 Tag
nGQL := ngql.NewCreateTag("team").
    Property("name", "string").
    Property("city", "string").
    String()

// 修改 Tag
nGQL := ngql.NewAlterTag("player").
    Add(&ngql.Property{Name: "score", Type: "int", Nullable: true}).
    Drop("old_field").
    Change(&ngql.Property{Name: "age", Type: "bigint", Nullable: true}).
    String()

// 删除 Tag
nGQL := ngql.NewDropTag("player").
    IfExists().
    String()

// 显示所有 Tag
nGQL := ngql.NewShowTags().
    String()

// 查看 Tag 详情
nGQL := ngql.NewDescribeTag("player").
    String()

// 显示 Tag 创建语句
nGQL := ngql.NewShowCreateTag("player").
    String()
```

### 3. Edge 类型管理

```go
// 创建 Edge 类型
nGQL := ngql.NewCreateEdge("follow").
    IfNotExists().
    AddProperty(&ngql.Property{
        Name:     "degree",
        Type:     "int",
        Nullable: false,
    }).
    AddProperty(&ngql.Property{
        Name:     "created_at",
        Type:     "timestamp",
        Nullable: true,
        Default:  0,
    }).
    TTL("created_at", 100).
    Comment("关注关系").
    String()

// 简化方式创建 Edge 类型
nGQL := ngql.NewCreateEdge("serve").
    Property("start_year", "int").
    Property("end_year", "int").
    String()

// 修改 Edge 类型
nGQL := ngql.NewAlterEdge("follow").
    Add(&ngql.Property{Name: "weight", Type: "double", Nullable: true}).
    Drop("old_field").
    Change(&ngql.Property{Name: "degree", Type: "bigint", Nullable: true}).
    String()

// 删除 Edge 类型
nGQL := ngql.NewDropEdge("follow").
    IfExists().
    String()

// 显示所有 Edge 类型
nGQL := ngql.NewShowEdges().
    String()

// 查看 Edge 类型详情
nGQL := ngql.NewDescribeEdge("follow").
    String()

// 显示 Edge 类型创建语句
nGQL := ngql.NewShowCreateEdge("follow").
    String()
```

### 4. 索引管理

```go
// 创建 Tag 索引
nGQL := ngql.NewCreateTagIndex("idx_player_age", "player").
    IfNotExists().
    Property("age", true).  // true 表示降序
    Comment("年龄索引").
    String()

// 创建复合索引
nGQL := ngql.NewCreateTagIndex("idx_player_name_age", "player").
    Property("name", false).  // false 表示升序
    Property("age", true).
    String()

// 创建 Edge 索引
nGQL := ngql.NewCreateEdgeIndex("idx_follow_degree", "follow").
    IfNotExists().
    Property("degree", false).
    String()

// 删除 Tag 索引
nGQL := ngql.NewDropTagIndex("idx_player_age").
    IfExists().
    String()

// 删除 Edge 索引
nGQL := ngql.NewDropEdgeIndex("idx_follow_degree").
    IfExists().
    String()

// 重建 Tag 索引
nGQL := ngql.NewRebuildTagIndex("idx_player_age", "idx_player_name").
    String()

// 重建所有 Tag 索引
nGQL := ngql.NewRebuildTagIndex().
    String()

// 重建 Edge 索引
nGQL := ngql.NewRebuildEdgeIndex("idx_follow_degree").
    String()

// 显示所有 Tag 索引
nGQL := ngql.NewShowTagIndexes().
    String()

// 显示所有 Edge 索引
nGQL := ngql.NewShowEdgeIndexes().
    String()

// 查看 Tag 索引详情
nGQL := ngql.NewDescribeTagIndex("idx_player_age").
    String()

// 查看 Edge 索引详情
nGQL := ngql.NewDescribeEdgeIndex("idx_follow_degree").
    String()

// 显示 Tag 索引状态
nGQL := ngql.NewShowTagIndexStatus().
    String()

// 显示 Edge 索引状态
nGQL := ngql.NewShowEdgeIndexStatus().
    String()
```

### 5. UPSERT 操作

```go
// UPSERT 顶点（不存在则插入，存在则更新）
nGQL := ngql.NewUpsertVertex("player", "player100").
    Set("age", 30).
    SetExpr("score = score + 10").
    When("name == \"Joe\"").
    Yield("name", "age", "score").
    String()

// 简单 UPSERT 顶点
nGQL := ngql.NewUpsertVertex("player", "player101").
    Set("age", 25).
    String()

// UPSERT 边
nGQL := ngql.NewUpsertEdge("follow", "player100", "player101").
    Set("degree", 90).
    SetExpr("weight = weight + 0.1").
    When("degree < 100").
    Yield("degree", "weight").
    String()

// UPSERT 带 rank 的边
nGQL := ngql.NewUpsertEdge("serve", "player100", "team204").
    WithRank(1).
    Set("end_year", 2020).
    String()
```

### 6. 插入顶点

```go
import "github.com/3th1nk/ngql"

// 插入单个顶点
nGQL := ngql.NewInsertVertex("player", "name", "age").
    Values("player100", "Tim Duncan", 42).
    String()

// 批量插入顶点
nGQL := ngql.NewInsertVertex("player", "name", "age").
    Values("player100", "Tim Duncan", 42).
    Values("player101", "Tony Parker", 36).
    Values("player102", "LaMarcus Aldridge", 33).
    String()

// 使用 IF NOT EXISTS 避免重复插入
nGQL := ngql.NewInsertVertex("player", "name", "age").
    IfNotExists().
    Values("player100", "Tim Duncan", 42).
    String()
```

### 7. 插入边

```go
// 插入单条边
nGQL := ngql.NewInsertEdge("follow", "degree").
    Edge("player100", "player101", 95).
    String()

// 批量插入边
nGQL := ngql.NewInsertEdge("follow", "degree").
    Edge("player100", "player101", 95).
    Edge("player101", "player102", 90).
    Edge("player102", "player100", 75).
    String()

// 插入带 rank 的边（用于同一对顶点间的多条边）
nGQL := ngql.NewInsertEdge("serve", "start_year", "end_year").
    EdgeWithRank("player100", "team204", 1, 1997, 2016).
    String()
```

### 8. MATCH 查询

```go
// 基本查询 - 使用 VertexPattern
pattern := ngql.NewVertexPattern("v").Tag("player")
nGQL := ngql.NewMatch(pattern).
    Return("v.player.name", "v.player.age").
    Limit(10).
    String()

// 或者使用 PathPattern（功能相同）
pattern := ngql.NewPathPattern().
    VertexWithTag("v", "player")
nGQL := ngql.NewMatch(pattern).
    Return("v.player.name", "v.player.age").
    Limit(10).
    String()

// 带条件查询
pattern := ngql.NewVertexPattern("v").Tag("player")
nGQL := ngql.NewMatch(pattern).
    Where("v.player.age > 30").
    Where("v.player.name != \"\"").
    Return("v.player.name AS name", "v.player.age AS age").
    OrderBy("v.player.age", true).  // true 表示降序
    Limit(20).
    String()

// 带分页查询
pattern := ngql.NewVertexPattern("v").Tag("player")
nGQL := ngql.NewMatch(pattern).
    Return("v.player.name", "v.player.age").
    OrderBy("v.player.age", false).  // false 表示升序
    Skip(10).
    Limit(10).
    String()

// 复杂图模式查询
pattern := ngql.NewPathPattern().
    VertexWithTag("v", "player").
    EdgeWithType(ngql.EdgeDirectionRight, "follow").
    VertexWithTag("v2", "player")

nGQL := ngql.NewMatch(pattern).
    Where("v.player.age > 30").
    Return("v.player.name AS follower", "v2.player.name AS followee").
    String()

// 变长路径查询
pattern := ngql.NewPathPattern().
    VertexWithTag("v1", "player").
    EdgePattern(
        ngql.NewEdgePattern(ngql.EdgeDirectionRight).
            Type("follow").
            Hops(1, 3),
    ).
    VertexWithTag("v2", "player")

nGQL := ngql.NewMatch(pattern).
    Return("v1.player.name", "v2.player.name").
    String()
```

### 9. GO 查询（图遍历）

```go
// 基本遍历
nGQL := ngql.NewGo("player100").
    Over("follow").
    Yield("$$.player.name", "$$.player.age").
    String()

// 多步遍历
nGQL := ngql.NewGo("player100").
    Over("follow").
    Steps(2).
    Yield("$$.player.name").
    String()

// 反向遍历
nGQL := ngql.NewGo("player100").
    Over("follow").
    Reversely().
    Yield("$$.player.name AS follower").
    String()

// 双向遍历
nGQL := ngql.NewGo("player100").
    Over("follow").
    Bidirect().
    Yield("$$.player.name").
    String()

// 遍历所有边类型
nGQL := ngql.NewGo("player100").
    Over("*").  // 或者不调用 Over()
    Yield("$$.player.name", "$$").
    String()

// 带条件的遍历
nGQL := ngql.NewGo("player100").
    Over("follow").
    Where("$$.player.age > 30").
    Yield("$$.player.name", "$$.player.age").
    String()

// 从多个起点遍历
nGQL := ngql.NewGo("player100", "player101", "player102").
    Over("follow").
    Yield("$$.player.name").
    String()
```

### 10. 更新顶点

```go
// 基本更新
nGQL := ngql.NewUpdateVertex("player", "player100").
    Set("age", 43).
    String()

// 更新多个属性
nGQL := ngql.NewUpdateVertex("player", "player100").
    Set("age", 43).
    Set("name", "Tim Duncan Jr.").
    String()

// 使用表达式更新
nGQL := ngql.NewUpdateVertex("player", "player100").
    SetExpr("age = age + 1").
    SetExpr("score = score + 10").
    String()

// 条件更新
nGQL := ngql.NewUpdateVertex("player", "player100").
    Set("age", 43).
    When("age > 40").
    String()

// 更新并返回结果
nGQL := ngql.NewUpdateVertex("player", "player100").
    Set("age", 43).
    Yield("player.name", "player.age").
    String()
```

### 11. 更新边

```go
// 基本更新
nGQL := ngql.NewUpdateEdge("follow", "player100", "player101").
    Set("degree", 96).
    String()

// 更新带 rank 的边
nGQL := ngql.NewUpdateEdge("serve", "player100", "team204").
    WithRank(1).
    Set("end_year", 2017).
    String()

// 条件更新
nGQL := ngql.NewUpdateEdge("follow", "player100", "player101").
    Set("degree", 96).
    When("degree < 100").
    Yield("properties(edge)").
    String()
```

### 12. 删除顶点

```go
// 删除单个顶点（不删除关联的边）
nGQL := ngql.NewDeleteVertex("player100").
    String()

// 删除多个顶点
nGQL := ngql.NewDeleteVertex("player100", "player101", "player102").
    String()

// 删除顶点及其关联的边
nGQL := ngql.NewDeleteVertex("player100").
    WithEdge().
    String()
```

### 13. 删除边

```go
// 删除单条边
nGQL := ngql.NewDeleteEdge("follow").
    Edge("player100", "player101").
    String()

// 删除多条边
nGQL := ngql.NewDeleteEdge("follow").
    Edge("player100", "player101").
    Edge("player101", "player102").
    String()

// 删除带 rank 的边
nGQL := ngql.NewDeleteEdge("serve").
    EdgeWithRank("player100", "team204", 1).
    String()
```

### 14. FETCH 查询（获取属性）

```go
// 获取顶点属性
nGQL := ngql.NewFetchVertex("player", "player100", "player101").
    Yield("properties(vertex)").
    String()

// 获取特定属性
nGQL := ngql.NewFetchVertex("player", "player100").
    Yield("player.name", "player.age").
    String()

// 获取边属性
nGQL := ngql.NewFetchEdge("follow").
    Edge("player100", "player101").
    Yield("properties(edge)").
    String()

// 获取带 rank 的边属性
nGQL := ngql.NewFetchEdge("serve").
    EdgeWithRank("player100", "team204", 1).
    Yield("serve.start_year", "serve.end_year").
    String()
```

### 15. LOOKUP 查询（索引查询）

```go
// 基本 LOOKUP 查询
nGQL := ngql.NewLookup("player").
    Where("player.age > 30").
    Yield("id(vertex)", "player.name", "player.age").
    String()

// 使用 DISTINCT 去重
nGQL := ngql.NewLookup("player").
    Where("player.name STARTS WITH \"T\"").
    Distinct().
    Yield("player.age").
    String()

// 多条件查询
nGQL := ngql.NewLookup("player").
    Where("player.age > 30").
    Where("player.name != \"\"").
    Yield("id(vertex)", "player.name").
    Limit(10).
    String()

// 查询边
nGQL := ngql.NewLookup("follow").
    Where("follow.degree > 90").
    Yield("properties(edge)").
    String()

// 带分页的查询
nGQL := ngql.NewLookup("player").
    Where("player.age BETWEEN 20 AND 40").
    Yield("id(vertex)", "player.name", "player.age").
    Skip(10).
    Limit(20).
    String()
```

### 16. GROUP BY 分组聚合

```go
// 基本分组统计
prevQuery := ngql.NewGo("player100").
    Over("follow").
    Yield("$$.player.age AS age").
    String()

nGQL := ngql.NewGroupBy(prevQuery, "$-.age").
    Yield("$-.age", "count(*) AS count").
    String()

// 多字段分组
prevQuery := ngql.NewGo("player100").
    Over("follow").
    Yield("$$.player.age AS age", "$$.player.name AS name").
    String()

nGQL := ngql.NewGroupBy(prevQuery, "$-.age").
    Yield("$-.age", "count(*) AS count", "collect($-.name) AS names").
    String()

// 带排序和限制的分组
prevQuery := ngql.NewGo("player100").
    Over("follow").
    Yield("$$.player.age AS age").
    String()

nGQL := ngql.NewGroupBy(prevQuery, "$-.age").
    Yield("$-.age", "count(*) AS count").
    OrderBy("count", true).  // 按数量降序
    Limit(5).
    String()

// 使用聚合函数
prevQuery := ngql.NewMatch("(v:player)-[:follow]->(v2:player)").
    Return("v2.player.age AS age", "v2.player.name AS name").
    String()

nGQL := ngql.NewGroupBy(prevQuery, "$-.age").
    Yield("$-.age",
          "count(*) AS count",
          "avg($-.age) AS avg_age",
          "min($-.age) AS min_age",
          "max($-.age) AS max_age").
    String()
```

### 17. FIND PATH 路径查询

```go
// 查找最短路径
nGQL := ngql.NewFindPath("SHORTEST").
    From("player100").
    To("team204").
    Over("*").
    Yield("path AS p").
    String()

// 查找所有路径
nGQL := ngql.NewFindPath("ALL").
    From("player100").
    To("player101", "player102").
    Over("follow", "serve").
    UpTo(5).  // 最多5步
    Yield("path AS p").
    Limit(10).
    String()

// 查找无环路径
nGQL := ngql.NewFindPath("NOLOOP").
    From("player100").
    To("team204").
    Over("*").
    Yield("path AS p").
    String()

// 带属性的路径查询
nGQL := ngql.NewFindPath("SHORTEST").
    WithProp().
    From("player100").
    To("team204").
    Over("follow", "serve").
    Bidirect().  // 双向查找
    Yield("path AS p").
    String()

// 带条件的路径查询
nGQL := ngql.NewFindPath("SHORTEST").
    From("player100").
    To("team204").
    Over("*").
    Where("follow.degree > 80").
    UpTo(10).
    Yield("path AS p").
    String()
```

### 18. GET SUBGRAPH 子图查询

```go
// 基本子图查询
nGQL := ngql.NewGetSubgraph(2).  // 2步子图
    From("player100").
    Yield("VERTICES AS nodes", "EDGES AS relationships").
    String()

// 带属性的子图查询
nGQL := ngql.NewGetSubgraph(1).
    WithProp().
    From("player100", "player101").
    Yield("VERTICES AS nodes", "EDGES AS edges").
    String()

// 指定边方向的子图查询
nGQL := ngql.NewGetSubgraph(2).
    From("player100").
    In("follow").      // 入边
    Out("serve").      // 出边
    Yield("VERTICES AS nodes").
    String()

// 双向边的子图查询
nGQL := ngql.NewGetSubgraph(3).
    From("player100").
    Both("follow", "serve").  // 双向边
    Yield("VERTICES AS nodes", "EDGES AS edges").
    String()

// 带条件的子图查询
nGQL := ngql.NewGetSubgraph(2).
    WithProp().
    From("player100").
    Both("follow").
    Where("$$.player.age > 30").
    Where("$$.player.name != \"\"").
    Yield("VERTICES AS nodes", "EDGES AS edges").
    String()

// 多起点子图查询
nGQL := ngql.NewGetSubgraph(1).
    From("player100", "player101", "player102").
    Both("*").
    Yield("VERTICES AS nodes", "EDGES AS edges").
    String()
```

## 最佳实践

### 1. 使用链式调用

```go
// ✅ 推荐
nGQL := ngql.NewInsertVertex("player", "name", "age").
    Values("player100", "Tim Duncan", 42).
    String()

// ❌ 不推荐
builder := ngql.NewInsertVertex("player", "name", "age")
builder.Values("player100", "Tim Duncan", 42)
nGQL := builder.String()
```

### 2. 批量操作

```go
// ✅ 推荐：一次插入多条
nGQL := ngql.NewInsertVertex("player", "name", "age").
    Values("player100", "Tim Duncan", 42).
    Values("player101", "Tony Parker", 36).
    Values("player102", "LaMarcus Aldridge", 33).
    String()

// ❌ 不推荐：多次插入
for _, player := range players {
    nGQL := ngql.NewInsertVertex("player", "name", "age").
        Values(player.ID, player.Name, player.Age).
        String()
    session.Execute(nGQL)
}
```

### 3. 错误处理

```go
session := config.Default().MustGetNebulaSessionPool()
result, err := session.Execute(nGQL)
if err != nil {
    return gError.DbError(err, "操作失败")
}
```

## 注意事项

1. **值转义**：构建器会自动处理字符串值的转义，无需手动处理
2. **VID 格式**：顶点 ID 会自动添加引号，传入时使用原始字符串即可
3. **属性名称**：属性名称需要包含标签前缀，如 `player.name`
4. **表达式**：使用 `SetExpr` 时，表达式不会被转义，可以使用 nGQL 函数和运算符
5. **性能考虑**：批量操作时，尽量使用单个语句插入多条数据，而不是多次调用

## API 参考

### 图模式构建器 (pattern.go)

| 构建器 | 方法 | 说明 |
|--------|------|------|
| VertexPattern | NewVertexPattern(variable) | 创建顶点模式构建器 |
| | Tag(tag) | 添加 Tag 标签 |
| | Property(key, value) | 添加属性过滤条件 |
| | String() | 生成顶点模式字符串 |
| EdgePattern | NewEdgePattern(direction) | 创建边模式构建器 |
| | Variable(variable) | 设置边变量名 |
| | Type(edgeType) | 添加边类型 |
| | Property(key, value) | 添加属性过滤条件 |
| | Hops(min, max) | 设置变长边的跳数范围 |
| | String() | 生成边模式字符串 |
| PathPattern | NewPathPattern() | 创建路径模式构建器 |
| | Variable(variable) | 设置路径变量名 |
| | Vertex(variable) | 添加顶点模式 |
| | VertexWithTag(variable, tag) | 添加带 Tag 的顶点模式 |
| | VertexPattern(vertex) | 添加自定义顶点模式 |
| | Edge(direction) | 添加边模式 |
| | EdgeWithType(direction, edgeType) | 添加带类型的边模式 |
| | EdgePattern(edge) | 添加自定义边模式 |
| | String() | 生成路径模式字符串 |

**EdgeDirection 常量**：
- `EdgeDirectionRight` - 右向边 (->)
- `EdgeDirectionLeft` - 左向边 (<-)
- `EdgeDirectionBoth` - 双向边 (-)

### 图空间管理 (schema.go)

| 构建器 | 方法 | 说明 |
|--------|------|------|
| CreateSpaceBuilder | NewCreateSpace(spaceName) | 创建图空间构建器 |
| | IfNotExists() | 设置 IF NOT EXISTS 选项 |
| | PartitionNum(num) | 设置分区数 |
| | ReplicaFactor(num) | 设置副本数 |
| | VidType(vidType) | 设置 VID 类型 |
| | Comment(comment) | 设置注释 |
| | String() | 生成 nGQL 语句 |
| CreateSpaceAsBuilder | NewCreateSpaceAs(newSpace, oldSpace) | 创建克隆图空间构建器 |
| | String() | 生成 nGQL 语句 |
| UseSpaceBuilder | NewUseSpace(spaceName) | 创建使用图空间构建器 |
| | String() | 生成 nGQL 语句 |
| ShowSpacesBuilder | NewShowSpaces() | 创建显示所有图空间构建器 |
| | String() | 生成 nGQL 语句 |
| DescribeSpaceBuilder | NewDescribeSpace(spaceName) | 创建查看图空间详情构建器 |
| | String() | 生成 nGQL 语句 |
| DropSpaceBuilder | NewDropSpace(spaceName) | 创建删除图空间构建器 |
| | IfExists() | 设置 IF EXISTS 选项 |
| | String() | 生成 nGQL 语句 |
| ClearSpaceBuilder | NewClearSpace(spaceName) | 创建清空图空间构建器 |
| | String() | 生成 nGQL 语句 |

### Tag 管理 (tag.go)

| 构建器 | 方法                          | 说明 |
|--------|-----------------------------|------|
| CreateTagBuilder | NewCreateTag(tagName)       | 创建 Tag 构建器 |
| | IfNotExists()               | 设置 IF NOT EXISTS 选项 |
| | AddProperty(prop *Property) | 添加属性（完整定义） |
| | Property(name, type)        | 添加属性（简化方式） |
| | TTL(col, duration)                  | 设置 TTL |
| | Comment(comment)            | 设置注释 |
| | String()                    | 生成 nGQL 语句 |
| AlterTagBuilder | NewAlterTag(tagName)        | 创建修改 Tag 构建器 |
| | Add(prop *Property)         | 添加属性 |
| | Drop(propName)              | 删除属性 |
| | Change(prop *Property)      | 修改属性 |
| | String()                    | 生成 nGQL 语句 |
| DropTagBuilder | NewDropTag(tagName)         | 创建删除 Tag 构建器 |
| | IfExists()                  | 设置 IF EXISTS 选项 |
| | String()                    | 生成 nGQL 语句 |
| ShowTagsBuilder | NewShowTags()               | 创建显示所有 Tag 构建器 |
| | String()                    | 生成 nGQL 语句 |
| DescribeTagBuilder | NewDescribeTag(tagName)     | 创建查看 Tag 详情构建器 |
| | String()                    | 生成 nGQL 语句 |
| ShowCreateTagBuilder | NewShowCreateTag(tagName)   | 创建显示 Tag 创建语句构建器 |
| | String()                    | 生成 nGQL 语句 |

### Edge 类型管理 (edge.go)

| 构建器 | 方法                          | 说明 |
|--------|-----------------------------|------|
| CreateEdgeBuilder | NewCreateEdge(edgeName)     | 创建 Edge 类型构建器 |
| | IfNotExists()               | 设置 IF NOT EXISTS 选项 |
| | AddProperty(prop *Property) | 添加属性（完整定义） |
| | Property(name, type)        | 添加属性（简化方式） |
| | TTL(col,duration)                   | 设置 TTL |
| | Comment(comment)            | 设置注释 |
| | String()                    | 生成 nGQL 语句 |
| AlterEdgeBuilder | NewAlterEdge(edgeName)      | 创建修改 Edge 类型构建器 |
| | Add(prop *Property)         | 添加属性 |
| | Drop(propName)              | 删除属性 |
| | Change(prop *Property)      | 修改属性 |
| | String()                    | 生成 nGQL 语句 |
| DropEdgeBuilder | NewDropEdge(edgeName)       | 创建删除 Edge 类型构建器 |
| | IfExists()                  | 设置 IF EXISTS 选项 |
| | String()                    | 生成 nGQL 语句 |
| ShowEdgesBuilder | NewShowEdges()              | 创建显示所有 Edge 类型构建器 |
| | String()                    | 生成 nGQL 语句 |
| DescribeEdgeBuilder | NewDescribeEdge(edgeName)   | 创建查看 Edge 类型详情构建器 |
| | String()                    | 生成 nGQL 语句 |
| ShowCreateEdgeBuilder | NewShowCreateEdge(edgeName) | 创建显示 Edge 类型创建语句构建器 |
| | String()                    | 生成 nGQL 语句 |

### 索引管理 (index.go)

| 构建器 | 方法 | 说明 |
|--------|------|------|
| CreateTagIndexBuilder | NewCreateTagIndex(indexName, tagName) | 创建 Tag 索引构建器 |
| | IfNotExists() | 设置 IF NOT EXISTS 选项 |
| | Property(name, desc) | 添加索引属性（desc: 是否降序） |
| | Comment(comment) | 设置注释 |
| | String() | 生成 nGQL 语句 |
| CreateEdgeIndexBuilder | NewCreateEdgeIndex(indexName, edgeName) | 创建 Edge 索引构建器 |
| | IfNotExists() | 设置 IF NOT EXISTS 选项 |
| | Property(name, desc) | 添加索引属性（desc: 是否降序） |
| | Comment(comment) | 设置注释 |
| | String() | 生成 nGQL 语句 |
| DropTagIndexBuilder | NewDropTagIndex(indexName) | 创建删除 Tag 索引构建器 |
| | IfExists() | 设置 IF EXISTS 选项 |
| | String() | 生成 nGQL 语句 |
| DropEdgeIndexBuilder | NewDropEdgeIndex(indexName) | 创建删除 Edge 索引构建器 |
| | IfExists() | 设置 IF EXISTS 选项 |
| | String() | 生成 nGQL 语句 |
| RebuildTagIndexBuilder | NewRebuildTagIndex(indexNames...) | 创建重建 Tag 索引构建器 |
| | String() | 生成 nGQL 语句 |
| RebuildEdgeIndexBuilder | NewRebuildEdgeIndex(indexNames...) | 创建重建 Edge 索引构建器 |
| | String() | 生成 nGQL 语句 |
| ShowTagIndexesBuilder | NewShowTagIndexes() | 创建显示所有 Tag 索引构建器 |
| | String() | 生成 nGQL 语句 |
| ShowEdgeIndexesBuilder | NewShowEdgeIndexes() | 创建显示所有 Edge 索引构建器 |
| | String() | 生成 nGQL 语句 |
| DescribeTagIndexBuilder | NewDescribeTagIndex(indexName) | 创建查看 Tag 索引详情构建器 |
| | String() | 生成 nGQL 语句 |
| DescribeEdgeIndexBuilder | NewDescribeEdgeIndex(indexName) | 创建查看 Edge 索引详情构建器 |
| | String() | 生成 nGQL 语句 |
| ShowTagIndexStatusBuilder | NewShowTagIndexStatus() | 创建显示 Tag 索引状态构建器 |
| | String() | 生成 nGQL 语句 |
| ShowEdgeIndexStatusBuilder | NewShowEdgeIndexStatus() | 创建显示 Edge 索引状态构建器 |
| | String() | 生成 nGQL 语句 |

### UPSERT 操作 (upsert.go)

| 构建器 | 方法 | 说明 |
|--------|------|------|
| UpsertVertexBuilder | NewUpsertVertex(tag, vid) | 创建 UPSERT 顶点构建器 |
| | Set(prop, value) | 设置属性值 |
| | SetExpr(expr) | 设置属性表达式 |
| | When(condition) | 添加条件 |
| | Yield(fields...) | 添加返回字段 |
| | String() | 生成 nGQL 语句 |
| UpsertEdgeBuilder | NewUpsertEdge(edgeType, srcVid, dstVid) | 创建 UPSERT 边构建器 |
| | WithRank(rank) | 设置边的 rank |
| | Set(prop, value) | 设置属性值 |
| | SetExpr(expr) | 设置属性表达式 |
| | When(condition) | 添加条件 |
| | Yield(fields...) | 添加返回字段 |
| | String() | 生成 nGQL 语句 |

### 插入操作 (insert.go)

| 构建器 | 方法 | 说明 |
|--------|------|------|
| InsertVertexBuilder | NewInsertVertex(tag, props...) | 创建插入顶点构建器 |
| | IfNotExists() | 设置 IF NOT EXISTS 选项 |
| | Values(vid, values...) | 添加顶点数据 |
| | String() | 生成 nGQL 语句 |
| InsertEdgeBuilder | NewInsertEdge(edgeType, props...) | 创建插入边构建器 |
| | IfNotExists() | 设置 IF NOT EXISTS 选项 |
| | Edge(srcVid, dstVid, values...) | 添加边数据 |
| | EdgeWithRank(srcVid, dstVid, rank, values...) | 添加带 rank 的边数据 |
| | String() | 生成 nGQL 语句 |

### 查询操作 (query.go)

| 构建器 | 方法 | 说明 |
|--------|------|------|
| MatchBuilder | NewMatch(pattern) | 创建 MATCH 查询构建器（pattern 必须是 Pattern 对象，如 VertexPattern 或 PathPattern） |
| | Where(condition) | 添加 WHERE 条件 |
| | Return(fields...) | 添加返回字段 |
| | OrderBy(field, desc) | 添加排序 |
| | Limit(limit) | 设置返回数量限制 |
| | Skip(skip) | 设置跳过数量 |
| | String() | 生成 nGQL 语句 |
| GoBuilder | NewGo(fromVids...) | 创建 GO 查询构建器 |
| | Over(edges...) | 指定要遍历的边类型 |
| | Steps(steps) | 设置遍历步数 |
| | Reversely() | 反向遍历 |
| | Bidirect() | 双向遍历 |
| | Where(condition) | 添加过滤条件 |
| | Yield(fields...) | 添加返回字段 |
| | String() | 生成 nGQL 语句 |
| FetchBuilder | NewFetchVertex(tag, vids...) | 创建获取顶点属性的构建器 |
| | NewFetchEdge(edgeType) | 创建获取边属性的构建器 |
| | Edge(srcVid, dstVid) | 添加要获取的边 |
| | EdgeWithRank(srcVid, dstVid, rank) | 添加要获取的带 rank 的边 |
| | Yield(fields...) | 添加返回字段 |
| | String() | 生成 nGQL 语句 |
| LookupBuilder | NewLookup(onType) | 创建 LOOKUP 查询构建器 |
| | Where(condition) | 添加 WHERE 条件 |
| | Yield(fields...) | 添加返回字段 |
| | Distinct() | 设置 DISTINCT 选项 |
| | Limit(limit) | 设置返回数量限制 |
| | Skip(skip) | 设置跳过数量 |
| | String() | 生成 nGQL 语句 |
| GroupByBuilder | NewGroupBy(prevQuery, groupByVars...) | 创建 GROUP BY 构建器 |
| | Yield(fields...) | 添加返回字段（含聚合函数） |
| | OrderBy(field, desc) | 添加排序 |
| | Limit(limit) | 设置返回数量限制 |
| | Skip(skip) | 设置跳过数量 |
| | String() | 生成 nGQL 语句 |
| FindPathBuilder | NewFindPath(pathType) | 创建 FIND PATH 查询构建器 |
| | WithProp() | 设置包含属性 |
| | From(vids...) | 设置起始顶点 |
| | To(vids...) | 设置目标顶点 |
| | Over(edges...) | 指定要遍历的边类型 |
| | Reversely() | 反向遍历 |
| | Bidirect() | 双向遍历 |
| | Where(condition) | 添加过滤条件 |
| | UpTo(steps) | 设置最大步数 |
| | Yield(fields...) | 添加返回字段 |
| | OrderBy(field, desc) | 添加排序 |
| | Limit(limit) | 设置返回数量限制 |
| | String() | 生成 nGQL 语句 |
| GetSubgraphBuilder | NewGetSubgraph(steps) | 创建 GET SUBGRAPH 查询构建器 |
| | WithProp() | 设置包含属性 |
| | From(vids...) | 设置起始顶点 |
| | In(edges...) | 指定 IN 方向的边类型 |
| | Out(edges...) | 指定 OUT 方向的边类型 |
| | Both(edges...) | 指定 BOTH 方向的边类型 |
| | Where(condition) | 添加 WHERE 条件 |
| | Yield(fields...) | 添加返回字段 |
| | String() | 生成 nGQL 语句 |

### 更新操作 (update.go)

| 构建器 | 方法 | 说明 |
|--------|------|------|
| UpdateVertexBuilder | NewUpdateVertex(tag, vid) | 创建更新顶点构建器 |
| | Set(prop, value) | 设置属性值 |
| | SetExpr(expr) | 设置属性表达式 |
| | When(condition) | 添加条件 |
| | Yield(fields...) | 添加返回字段 |
| | String() | 生成 nGQL 语句 |
| UpdateEdgeBuilder | NewUpdateEdge(edgeType, srcVid, dstVid) | 创建更新边构建器 |
| | WithRank(rank) | 设置边的 rank |
| | Set(prop, value) | 设置属性值 |
| | SetExpr(expr) | 设置属性表达式 |
| | When(condition) | 添加条件 |
| | Yield(fields...) | 添加返回字段 |
| | String() | 生成 nGQL 语句 |

### 删除操作 (delete.go)

| 构建器 | 方法 | 说明 |
|--------|------|------|
| DeleteVertexBuilder | NewDeleteVertex(vids...) | 创建删除顶点构建器 |
| | WithEdge() | 同时删除关联的边 |
| | String() | 生成 nGQL 语句 |
| DeleteEdgeBuilder | NewDeleteEdge(edgeType) | 创建删除边构建器 |
| | Edge(srcVid, dstVid) | 添加要删除的边 |
| | EdgeWithRank(srcVid, dstVid, rank) | 添加要删除的带 rank 的边 |
| | String() | 生成 nGQL 语句 |

## 参考文档

- [NebulaGraph 官方文档](https://docs.nebula-graph.com.cn/)
- [NebulaGraph nGQL 指南](https://docs.nebula-graph.com.cn/3.8.0/3.ngql-guide/1.nGQL-overview/ngql-style-guide)
- [NebulaGraph Go 客户端](https://github.com/vesoft-inc/nebula-go)
