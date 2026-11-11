# NebulaGraph 集成测试

本目录包含 NebulaGraph nGQL 构建器的集成测试，用于验证实际数据库操作。

## 测试说明

### 测试内容

#### 1. TestIntegration_FullWorkflow
完整的数据库操作工作流测试，包括：

1. **图空间管理**
   - 创建测试图空间
   - 使用图空间
   - 删除图空间

2. **Schema 管理**
   - 创建 Tag（player, team）
   - 创建 Edge 类型（follow, serve）

3. **数据插入**
   - 批量插入顶点（球员、球队）
   - 批量插入边（关注关系、效力关系）

4. **数据查询**
   - MATCH 查询（条件查询、排序）
   - GO 查询（图遍历）
   - FETCH 查询（获取属性）

5. **数据更新**
   - UPDATE 更新顶点属性
   - 使用表达式更新

6. **UPSERT 操作**
   - 插入或更新顶点

7. **数据删除**
   - 删除边
   - 删除顶点（含关联边）

#### 2. TestIntegration_IndexOperations
索引操作测试，包括：

1. 创建 Tag 索引（单字段、复合索引）
2. 显示所有索引
3. 查看索引详情
4. 重建索引
5. 显示索引状态
6. 删除索引

## 运行测试

### 前提条件

1. **NebulaGraph 服务运行中**（默认连接 127.0.0.1:9669，用户名 root，密码 nebula）

2. **集群初始化**：SessionPool 连接的是 `default_space`，需预先创建：

```ngql
ADD HOSTS "storaged":9779;
CREATE SPACE IF NOT EXISTS default_space(vid_type=FIXED_STRING(32));
```

3. **本地容器快速搭建**（NebulaGraph v3.8.0，单机）：

```bash
# metad/graphd/storaged 三容器，建议 heartbeat_interval_secs=1 加速 schema 传播
docker compose up -d
# 等待 graphd 9669 端口就绪后，通过 console 执行上面的初始化语句
```

### 运行所有集成测试

```bash
go test ./_test -v -count=1
```

> 注意：`_test` 目录以下划线开头，Go 工具在 `./...` 模式下会忽略它，需显式指定路径运行。

### 运行特定测试

```bash
# 运行完整工作流测试
go test ./_test -v -run TestIntegration_FullWorkflow

# 运行索引操作测试
go test ./_test -v -run TestIntegration_IndexOperations
```

### 运行特定测试

```bash
# 运行完整工作流测试
go test -v -run TestIntegration_FullWorkflow

# 运行索引操作测试
go test -v -run TestIntegration_IndexOperations
```

### 测试输出

测试会输出详细的执行日志，包括：
- 执行的 nGQL 语句
- 查询结果行数
- 操作状态

示例输出：
```
=== RUN   TestIntegration_FullWorkflow
    integration_test.go:75: === 清理测试环境 ===
    integration_test.go:82: === 1. 创建测试图空间 ===
    integration_test.go:89: 执行: CREATE SPACE test_space (partition_num = 10, replica_factor = 1, vid_type = FIXED_STRING(32)) COMMENT = "集成测试图空间"
    integration_test.go:99: === 2. 创建 Tag ===
    integration_test.go:107: 执行: CREATE TAG IF NOT EXISTS player(name string, age int32, score int32) COMMENT = "球员信息"
    ...
    integration_test.go:331: === 集成测试完成 ===
--- PASS: TestIntegration_FullWorkflow (26.11s)
PASS
```

## 注意事项

1. **SessionPool 禁止裸 USE 语句**
   - nebula-go 的 SessionPool 在语句执行后会强制切回配置空间
   - 测试统一使用复合执行 `"USE test_space; <stmt>"`，全局语句（CREATE/DROP SPACE）走专用辅助函数

2. **schema 分发存在心跳延迟**
   - CREATE TAG/EDGE 返回成功后，DML 可能短暂报 NotFound（graphd 本地 schema 缓存按心跳周期刷新）
   - 测试内置 `waitForSchema`（轮询 DESCRIBE）与 `rebuildTagIndexWithRetry`（轮询重建）处理
   - 容器环境建议 `heartbeat_interval_secs=1` 加速 schema 传播

3. **测试图空间**
   - 测试使用独立的图空间 `test_space`
   - 测试结束后会自动清理

4. **执行时间**
   - 完整测试约 26 秒、索引测试约 18 秒（含 schema 等待）

5. **并发执行**
   - 不建议并发运行多个集成测试，可能会导致图空间冲突

6. **失败处理**
   - 如果测试失败，检查 NebulaGraph 服务状态与集群初始化是否完成
   - 手动清理测试图空间：`DROP SPACE IF EXISTS test_space`

7. **NebulaGraph 限制**
   - 不支持跨 Tag 或 Edge type 创建复合索引
   - UPDATE/UPSERT 的 YIELD 子句不需要 tag 前缀
   - MATCH 查询的 ORDER BY 应使用别名而非完整路径

## 测试数据

测试使用的示例数据：

**顶点（球员）：**
- player100: Tim Duncan, 42岁, 1000分
- player101: Tony Parker, 36岁, 800分
- player102: LaMarcus Aldridge, 33岁, 700分

**顶点（球队）：**
- team204: Spurs, San Antonio
- team215: Hornets, Charlotte

**边（关注关系）：**
- player100 -> player101 (degree: 95)
- player101 -> player102 (degree: 90)
- player102 -> player100 (degree: 75)

**边（效力关系）：**
- player100 -> team204 (1997-2016)
- player101 -> team204 (2001-2018)

## 扩展测试

如需添加新的测试用例：

1. 在 `integration_test.go` 中添加新的测试函数
2. 函数名以 `TestIntegration_` 开头
3. 使用 `executeNGQL()` 执行语句
4. 使用 `executeNGQLIgnoreError()` 执行清理操作
5. 添加适当的断言验证结果

示例：
```go
func TestIntegration_CustomFeature(t *testing.T) {
    // 使用测试图空间
    useSpaceNGQL := ngql.NewUseSpace(testSpace).String()
    executeNGQL(t, useSpaceNGQL)

    // 执行测试操作
    nGQL := ngql.NewMatch("(v:player)").
        Return("v.player.name").
        String()
    result := executeNGQL(t, nGQL)

    // 验证结果
    if result.GetRowSize() == 0 {
        t.Error("期望有查询结果")
    }
}
```
