# NebulaGraph 变量引用指南

NebulaGraph nGQL 提供了多种变量引用机制，用于在查询中引用顶点、边、属性和中间结果。

## 目录

- [变量类型概览](#变量类型概览)
- [管道引用变量 `$-`](#管道引用变量--)
- [用户自定义变量 `$var`](#用户自定义变量-var)
- [顶点属性引用 `$^` 和 `$$`](#顶点属性引用--和-)
- [边属性引用](#边属性引用)
- [MATCH 模式变量](#match-模式变量)
- [变量对比](#变量对比)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

## 变量类型概览

| 变量类型 | 语法 | 适用场景 | 示例 |
|---------|------|---------|------|
| 管道引用变量 | `$-` | 管道操作符中引用前一个查询结果 | `$-.id` |
| 用户自定义变量 | `$var` | 存储查询结果供后续使用 | `$friends.id` |
| 源顶点属性 | `$^` | GO 语句中引用起始顶点属性 | `$^.player.name` |
| 目标顶点属性 | `$$` | GO 语句中引用目标顶点属性 | `$$.player.age` |
| 边属性 | `<edge>.<prop>` | 引用边的属性 | `follow.degree` |
| MATCH 模式变量 | `(v)`, `[e]`, `p` | MATCH 语句中的顶点、边、路径变量 | `(v:player)` |

## 管道引用变量 `$-`

### 概述

`$-` 是管道操作符的隐式引用变量，用于引用**前一个查询的输出结果**。

### 语法

```ngql
查询1 | 查询2 使用 $-.字段名
```

### 使用场景

#### 1. 简单链式查询

```ngql
-- 查找 player100 关注的人，然后查找这些人关注的人（二度关系）
GO FROM "player100" OVER follow
YIELD dst(edge) AS dstid |
GO FROM $-.dstid OVER follow
YIELD dst(edge);
```

#### 2. 多字段引用

```ngql
-- 引用前一个查询的多个字段
GO FROM "player100" OVER follow
YIELD dst(edge) AS dstid, properties($$).name AS name |
GO FROM $-.dstid OVER serve
WHERE $-.name != "Tony Parker"
YIELD properties($$).name AS team;
```

#### 3. 属性过滤

```ngql
-- 使用前一个查询的结果进行过滤
GO FROM "player100" OVER follow
YIELD dst(edge) AS id, properties($$).age AS age |
GO FROM $-.id OVER serve
WHERE $-.age > 30
YIELD properties($$).name AS team;
```

### 注意事项

- `$-` 只能在管道符 `|` 后的**下一个查询**中使用
- 不能跨越多个管道操作符引用
- 字段名必须在前一个查询的 YIELD 子句中定义

## 用户自定义变量 `$var`

### 概述

用户自定义变量用于**显式存储查询结果**，可以在后续多个查询中重复使用。

### 语法

```ngql
$变量名 = 查询语句;
后续查询使用 $变量名.字段名;
```

### 变量命名规则

- 必须以 `$` 开头
- 由字母、数字、下划线组成
- 区分大小写
- 示例：`$result`, `$player_list`, `$temp123`

### 使用场景

#### 1. 存储中间结果

```ngql
-- 存储关注列表
$friends = GO FROM "player100" OVER follow
           YIELD dst(edge) AS id;

-- 查询朋友的球队
GO FROM $friends.id OVER serve
YIELD properties($$).name AS team;
```

#### 2. 多次使用同一结果

```ngql
-- 存储查询结果
$players = GO FROM "team200" OVER serve
           YIELD src(edge) AS id;

-- 查询1：获取球员年龄
GO FROM $players.id
YIELD properties($$).age AS age;

-- 查询2：获取球员关注关系
GO FROM $players.id OVER follow
YIELD dst(edge) AS friend;
```

#### 3. 复杂多步查询

```ngql
-- 第一步：找到所有关注者
$followers = GO FROM "player100" OVER follow
             YIELD dst(edge) AS id;

-- 第二步：找到关注者的球队
$teams = GO FROM $followers.id OVER serve
         YIELD DISTINCT properties($$).name AS team_name;

-- 第三步：使用球队信息进行进一步查询
GO FROM $teams.team_name OVER serve
YIELD src(edge) AS player_id;
```

### 注意事项

- 变量在**当前会话**中有效
- 复合查询中必须使用分号 `;` 终止
- 变量会占用内存，大结果集需注意性能

## 顶点属性引用 `$^` 和 `$$`

### 概述

在 GO 语句中，`$^` 和 `$$` 用于引用遍历过程中的顶点属性。

| 符号 | 含义 | 说明 |
|------|------|------|
| `$^` | 源顶点（Source Vertex） | GO 语句的起始顶点 |
| `$$` | 目标顶点（Destination Vertex） | 通过边到达的目标顶点 |

### 语法

```ngql
$^.<tag_name>.<prop_name>  -- 源顶点属性
$$.<tag_name>.<prop_name>  -- 目标顶点属性
```

### 使用场景

#### 1. 基本属性引用

```ngql
-- 查询关注关系，返回源顶点和目标顶点的名字
GO FROM "player100" OVER follow
YIELD $^.player.name AS from_name,
      $$.player.name AS to_name;
```

#### 2. 属性过滤

```ngql
-- 查找年龄大于 30 的目标顶点
GO FROM "player100" OVER follow
WHERE $$.player.age > 30
YIELD $$.player.name AS name,
      $$.player.age AS age;
```

#### 3. 多 Tag 属性引用

```ngql
-- 当顶点有多个 Tag 时，指定 Tag 名称
GO FROM "player100" OVER follow
YIELD $^.player.name AS player_name,
      $^.person.birthday AS birthday,
      $$.player.age AS friend_age;
```

#### 4. 与函数结合使用

```ngql
-- 使用 properties() 函数获取所有属性
GO FROM "player100" OVER follow
YIELD properties($^) AS source_props,
      properties($$) AS dest_props;
```

### 新旧语法对比

```ngql
-- 旧语法（仍然支持）
GO FROM "player100" OVER follow
YIELD $^.player.name, $$.player.age;

-- 新语法（推荐）
GO FROM "player100" OVER follow
YIELD properties($^).name, properties($$).age;
```

### 注意事项

- 只能在 **GO 语句**中使用
- 必须指定 Tag 名称（除非使用 properties() 函数）
- MATCH 语句中使用模式变量，不使用 `$^` 和 `$$`

## 边属性引用

### 概述

边属性引用用于访问边（关系）的属性，包括用户自定义属性和内置属性。

### 用户自定义属性

#### 语法

```ngql
<edge_type>.<prop_name>  -- 旧语法
properties(edge).<prop_name>  -- 新语法（推荐）
```

#### 示例

```ngql
-- 查询关注关系的亲密度
GO FROM "player100" OVER follow
YIELD follow.degree;

-- 使用新语法
GO FROM "player100" OVER follow
YIELD properties(edge).degree;
```

### 内置属性

边有 4 个内置属性：

| 属性 | 说明 | 旧语法 | 新语法 |
|------|------|--------|--------|
| `_src` | 起始顶点 ID | `follow._src` | `src(edge)` |
| `_dst` | 目标顶点 ID | `follow._dst` | `dst(edge)` |
| `_type` | 边类型 | `follow._type` | `type(edge)` |
| `_rank` | 边的 rank 值 | `follow._rank` | `rank(edge)` |

#### 示例

```ngql
-- 旧语法
GO FROM "player100" OVER follow
YIELD follow._src, follow._dst, follow._type, follow._rank;

-- 新语法（推荐）
GO FROM "player100" OVER follow
YIELD src(edge), dst(edge), type(edge), rank(edge);
```

### 使用场景

#### 1. 边属性过滤

```ngql
-- 查找亲密度大于 90 的关注关系
GO FROM "player100" OVER follow
WHERE properties(edge).degree > 90
YIELD dst(edge), properties(edge).degree;
```

#### 2. 获取所有边属性

```ngql
-- 使用 properties() 获取边的所有属性
GO FROM "player100" OVER follow
YIELD properties(edge);
```

#### 3. 边属性与顶点属性结合

```ngql
-- 同时返回顶点和边的属性
GO FROM "player100" OVER follow
YIELD $^.player.name AS from_name,
      $$.player.name AS to_name,
      properties(edge).degree AS intimacy;
```

## MATCH 模式变量

### 概述

MATCH 语句使用模式变量来表示图中的顶点、边和路径。

### 变量类型

#### 1. 顶点变量

```ngql
-- 基本顶点变量
MATCH (v) RETURN v;

-- 带 Tag 的顶点变量
MATCH (v:player) RETURN v;

-- 带属性过滤的顶点变量
MATCH (v:player {name: "Tim Duncan"}) RETURN v;
```

#### 2. 边变量

```ngql
-- 基本边变量
MATCH (v1)-[e]->(v2) RETURN e;

-- 带类型的边变量
MATCH (v1)-[e:follow]->(v2) RETURN e;

-- 带属性过滤的边变量
MATCH (v1)-[e:follow {degree: 95}]->(v2) RETURN e;
```

#### 3. 路径变量

```ngql
-- 路径变量
MATCH p = (v1:player)-[:follow]->(v2:player)
RETURN p;

-- 变长路径
MATCH p = (v1:player)-[:follow*1..3]->(v2:player)
RETURN p;
```

### 使用场景

#### 1. 基本查询

```ngql
-- 查找球员及其关注关系
MATCH (v1:player)-[e:follow]->(v2:player)
WHERE v1.player.name == "Tim Duncan"
RETURN v1.player.name AS from,
       v2.player.name AS to,
       e.degree AS degree;
```

#### 2. 多跳查询

```ngql
-- 查找二度关系
MATCH (v1:player)-[:follow]->(v2:player)-[:follow]->(v3:player)
WHERE v1.player.name == "Tim Duncan"
RETURN v3.player.name AS friend_of_friend;
```

#### 3. 变长路径查询

```ngql
-- 查找 1-3 跳内的所有关注关系
MATCH p = (v1:player {name: "Tim Duncan"})-[:follow*1..3]->(v2:player)
RETURN v2.player.name AS reachable_players,
       length(p) AS hops;
```

#### 4. 路径函数

```ngql
-- 使用路径变量和路径函数
MATCH p = (v1:player)-[:follow*]->(v2:player)
WHERE v1.player.name == "Tim Duncan"
  AND v2.player.name == "Tony Parker"
RETURN p, length(p) AS path_length
ORDER BY length(p) ASC
LIMIT 1;
```

### 变量作用域

```ngql
-- 变量在整个 MATCH 语句中可用
MATCH (v1:player)-[e:follow]->(v2:player)
WHERE v1.player.age > 30
  AND e.degree > 90
  AND v2.player.age < 35
RETURN v1, e, v2;
```

## 变量对比

### GO 语句 vs MATCH 语句

| 特性 | GO 语句 | MATCH 语句 |
|------|---------|-----------|
| 顶点引用 | `$^`, `$$` | 模式变量 `(v)` |
| 边引用 | `edge`, `<edge_type>.<prop>` | 模式变量 `[e]` |
| 路径引用 | 不支持 | 路径变量 `p` |
| 管道操作 | 支持 `$-` | 不支持 |
| 用户变量 | 支持 `$var` | 支持 `$var` |

### 示例对比

#### GO 语句风格

```ngql
-- 使用 GO 语句查询
GO FROM "player100" OVER follow
YIELD $^.player.name AS from_name,
      $$.player.name AS to_name,
      properties(edge).degree AS degree;
```

#### MATCH 语句风格

```ngql
-- 使用 MATCH 语句查询（等价）
MATCH (v1:player)-[e:follow]->(v2:player)
WHERE id(v1) == "player100"
RETURN v1.player.name AS from_name,
       v2.player.name AS to_name,
       e.degree AS degree;
```

## 最佳实践

### 1. 选择合适的变量类型

```ngql
-- ✅ 简单链式查询使用管道
GO FROM "player100" OVER follow YIELD dst(edge) AS id |
GO FROM $-.id OVER serve YIELD properties($$).name;

-- ✅ 复杂多步查询使用自定义变量
$friends = GO FROM "player100" OVER follow YIELD dst(edge) AS id;
GO FROM $friends.id OVER serve YIELD properties($$).name AS team;
GO FROM $friends.id YIELD properties($$).age AS age;
```

### 2. 使用新语法

```ngql
-- ❌ 旧语法
GO FROM "player100" OVER follow
YIELD $^.player.name, $$.player.age, follow.degree;

-- ✅ 新语法（推荐）
GO FROM "player100" OVER follow
YIELD properties($^).name, properties($$).age, properties(edge).degree;
```

### 3. 明确字段别名

```ngql
-- ✅ 使用清晰的别名
GO FROM "player100" OVER follow
YIELD dst(edge) AS friend_id, properties($$).name AS friend_name |
GO FROM $-.friend_id OVER serve
YIELD properties($$).name AS team_name;
```

### 4. MATCH vs GO 选择

```ngql
-- ✅ 简单遍历使用 GO（性能更好）
GO FROM "player100" OVER follow YIELD dst(edge);

-- ✅ 复杂模式匹配使用 MATCH（更直观）
MATCH (v1:player)-[:follow]->(v2:player)-[:serve]->(t:team)
WHERE v1.player.name == "Tim Duncan"
RETURN v2.player.name, t.team.name;
```

### 5. 避免变量名冲突

```ngql
-- ❌ 变量名可能冲突
$id = GO FROM "player100" OVER follow YIELD dst(edge) AS id;
GO FROM $id.id OVER serve YIELD properties($$).name;

-- ✅ 使用描述性变量名
$friends = GO FROM "player100" OVER follow YIELD dst(edge) AS friend_id;
GO FROM $friends.friend_id OVER serve YIELD properties($$).name;
```

## 常见问题

### Q1: `$-` 和 `$var` 有什么区别？

**A:**
- `$-` 是隐式的，只能在管道符后的下一个查询中使用
- `$var` 是显式的，可以在会话中多次使用
- `$-` 适合简单链式查询，`$var` 适合复杂多步查询

### Q2: 为什么 MATCH 语句中不能使用 `$^` 和 `$$`？

**A:** `$^` 和 `$$` 是 GO 语句特有的语法。MATCH 语句使用模式变量（如 `(v1)`, `(v2)`）来表示顶点，这种方式更符合图模式匹配的语义。

### Q3: 如何在 MATCH 语句中引用边属性？

**A:** 使用边变量：

```ngql
MATCH (v1)-[e:follow]->(v2)
WHERE e.degree > 90
RETURN v1, v2, e.degree;
```

### Q4: 变量的作用域是什么？

**A:**
- `$-`: 仅在管道符后的下一个查询中
- `$var`: 当前会话中
- MATCH 模式变量: 当前 MATCH 语句中

### Q5: 如何获取顶点或边的所有属性？

**A:** 使用 `properties()` 函数：

```ngql
-- GO 语句
GO FROM "player100" OVER follow
YIELD properties($^), properties($$), properties(edge);

-- MATCH 语句
MATCH (v1)-[e]->(v2)
RETURN properties(v1), properties(e), properties(v2);
```

### Q6: 管道操作可以嵌套吗？

**A:** 可以，但 `$-` 只能引用**直接前一个**查询的结果：

```ngql
-- ✅ 正确：每个 $- 引用前一个查询
GO FROM "player100" OVER follow YIELD dst(edge) AS id1 |
GO FROM $-.id1 OVER follow YIELD dst(edge) AS id2 |
GO FROM $-.id2 OVER serve YIELD properties($$).name;

-- ❌ 错误：不能跨越管道引用
GO FROM "player100" OVER follow YIELD dst(edge) AS id1 |
GO FROM $-.id1 OVER follow YIELD dst(edge) AS id2 |
GO FROM $-.id1 OVER serve YIELD properties($$).name;  -- 错误！id1 不在前一个查询中
```

### Q7: 如何在 WHERE 子句中使用变量？

**A:**

```ngql
-- GO 语句中使用 $^, $$
GO FROM "player100" OVER follow
WHERE $$.player.age > 30 AND properties(edge).degree > 90
YIELD $$.player.name;

-- MATCH 语句中使用模式变量
MATCH (v1:player)-[e:follow]->(v2:player)
WHERE v1.player.name == "Tim Duncan"
  AND v2.player.age > 30
  AND e.degree > 90
RETURN v2.player.name;

-- 管道中使用 $-
GO FROM "player100" OVER follow YIELD dst(edge) AS id, properties($$).age AS age |
GO FROM $-.id OVER serve WHERE $-.age > 30 YIELD properties($$).name;
```

## 参考文档

- [NebulaGraph 官方文档 - 变量和复合查询](https://docs.nebula-graph.com.cn/3.8.0/3.ngql-guide/4.variable-and-composite-queries/)
- [NebulaGraph 官方文档 - 属性引用](https://docs.nebula-graph.com.cn/3.8.0/3.ngql-guide/4.variable-and-composite-queries/3.property-reference/)
- [NebulaGraph 官方文档 - MATCH 语句](https://docs.nebula-graph.com.cn/3.8.0/3.ngql-guide/7.general-query-statements/2.match/)
- [NebulaGraph 官方文档 - GO 语句](https://docs.nebula-graph.com.cn/3.8.0/3.ngql-guide/7.general-query-statements/3.go/)
