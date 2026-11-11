//go:build integration

// 集成测试：需连接真实 NebulaGraph 实例运行（见 README.md），
// 通过 build tag 与单元测试隔离，保证 go mod tidy / go test ./... 不会误触
package integration_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/3th1nk/ngql"
	nebula "github.com/vesoft-inc/nebula-go/v3"
)

// 连接参数支持环境变量覆盖（CI 中 graphd 运行在 dind 容器内，通过服务名访问）
var (
	host         = envOr("NEBULA_HOST", "127.0.0.1")
	port         = envIntOr("NEBULA_PORT", 9669)
	storageHost  = envOr("NEBULA_STORAGE_HOST", "storaged")
	storagePort  = envIntOr("NEBULA_STORAGE_PORT", 9779)
	defaultSpace = envOr("NEBULA_DEFAULT_SPACE", "default_space")
	username     = envOr("NEBULA_USERNAME", "root")
	password     = envOr("NEBULA_PASSWORD", "nebula")
	testSpace    = "test_space"
)

var (
	sessionPool *nebula.SessionPool
)

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envIntOr(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// TestMain 测试主函数：自动初始化集群（幂等）后创建 SessionPool
func TestMain(m *testing.M) {
	if err := initCluster(); err != nil {
		fmt.Println("集群初始化失败: ", err)
		os.Exit(1)
	}

	poolConf, err := nebula.NewSessionPoolConf(
		username, password,
		[]nebula.HostAddress{
			{Host: host, Port: port},
		},
		defaultSpace,
	)
	if err != nil {
		fmt.Println("创建 SessionPoolConf 失败: ", err)
		os.Exit(1)
	}
	// default_space 的分区 leader 就绪前 USE 可能失败，重试等待
	for i := 0; i < 15; i++ {
		sessionPool, err = nebula.NewSessionPool(*poolConf, nebula.DefaultLogger{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		fmt.Println("创建 SessionPool 失败: ", err)
		os.Exit(1)
	}
	code := m.Run()
	sessionPool.Close()
	os.Exit(code)
}

// initCluster 初始化集群（幂等）：注册存储节点、确保 default_space 存在
func initCluster() error {
	// 等待 graphd 就绪后建立连接
	var pool *nebula.ConnectionPool
	var session *nebula.Session
	var err error
	for i := 0; i < 30; i++ {
		pool, err = nebula.NewConnectionPool(
			[]nebula.HostAddress{{Host: host, Port: port}},
			nebula.GetDefaultConf(),
			nebula.DefaultLogger{},
		)
		if err == nil {
			session, err = pool.GetSession(username, password)
			if err == nil {
				break
			}
			pool.Close()
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("连接 graphd %s:%d 失败: %w", host, port, err)
	}
	defer pool.Close()
	defer session.Release()

	exec := func(stmt string) error {
		rs, err := session.Execute(stmt)
		if err != nil {
			return err
		}
		if !rs.IsSucceed() {
			return fmt.Errorf("%s 执行失败: %s", stmt, rs.GetErrorMsg())
		}
		return nil
	}

	// 注册存储节点（重复注册会报错，忽略即可）
	if err := exec(fmt.Sprintf(`ADD HOSTS "%s":%d`, storageHost, storagePort)); err != nil {
		fmt.Printf("ADD HOSTS（已注册则忽略）: %v\n", err)
	}

	// 等待存储节点被 meta 心跳感知（ONLINE）：未感知时创建图空间会报 Host not enough
	registered := false
	for i := 0; i < 30; i++ {
		rs, err := session.Execute("SHOW HOSTS")
		if err == nil && rs.IsSucceed() {
			for _, row := range rs.AsStringTable() {
				joined := strings.Join(row, " ")
				if strings.Contains(joined, storageHost) && strings.Contains(joined, "ONLINE") {
					registered = true
					break
				}
			}
		}
		if registered {
			break
		}
		time.Sleep(time.Second)
	}
	if !registered {
		return fmt.Errorf("等待存储节点 %s:%d 注册超时", storageHost, storagePort)
	}

	// 创建 default_space（冷启动时存储节点心跳未活跃会暂时报 Host not enough，重试等待）
	var createErr error
	for i := 0; i < 30; i++ {
		if err := exec(fmt.Sprintf("CREATE SPACE IF NOT EXISTS %s(vid_type=FIXED_STRING(32))", defaultSpace)); err == nil {
			createErr = nil
			break
		} else {
			createErr = err
		}
		time.Sleep(time.Second)
	}
	if createErr != nil {
		return createErr
	}

	// 等待 default_space 的分区分配完成（分区分配依赖空间创建，需在 CREATE 之后等待）
	spaceReady := false
	for i := 0; i < 30; i++ {
		rs, err := session.Execute("SHOW HOSTS")
		if err == nil && rs.IsSucceed() {
			for _, row := range rs.AsStringTable() {
				if strings.Contains(strings.Join(row, " "), defaultSpace) {
					spaceReady = true
					break
				}
			}
		}
		if spaceReady {
			break
		}
		time.Sleep(time.Second)
	}
	if !spaceReady {
		return fmt.Errorf("等待 %s 分区分配超时", defaultSpace)
	}
	return probeStorage(session)
}

// probeStorage 写入探针：真实验证 graphd schema 缓存与 storage RPC 写入链路就绪
// 冷启动后首次 DML 可能报 RPC timeout，就绪后其余语句不再受影响；探针数据自动清理
func probeStorage(session *nebula.Session) error {
	execOK := func(stmt string) bool {
		rs, err := session.Execute(stmt)
		return err == nil && rs.IsSucceed()
	}
	// 创建探针 Tag
	for i := 0; i < 15; i++ {
		if execOK(fmt.Sprintf("USE %s; CREATE TAG IF NOT EXISTS probe_tag(c int)", defaultSpace)) {
			break
		}
		time.Sleep(time.Second)
	}
	// 轮询写入直到 storage RPC 就绪
	for i := 0; i < 30; i++ {
		if execOK(fmt.Sprintf(`USE %s; INSERT VERTEX probe_tag(c) VALUES "probe_tag":(1)`, defaultSpace)) {
			execOK(fmt.Sprintf(`USE %s; DELETE VERTEX "probe_tag"`, defaultSpace))
			execOK(fmt.Sprintf("USE %s; DROP TAG probe_tag", defaultSpace))
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("存储写入链路探测超时")
}

// executeNGQL 执行 nGQL 语句并返回结果
// 注意：SessionPool 禁止裸 USE 语句（执行后会强制切回 pool 配置的空间），
// 空间上下文语句统一使用 "USE test_space; <stmt>" 复合执行，确保落点正确
func executeNGQL(t *testing.T, nGQL string) *nebula.ResultSet {
	result, err := sessionPool.Execute("USE " + testSpace + "; " + nGQL)
	if err != nil {
		t.Fatalf("执行 nGQL 失败: %s\n错误: %v", nGQL, err)
	}
	if !result.IsSucceed() {
		t.Fatalf("nGQL 执行失败: %s\n错误码: %d, 错误信息: %s",
			nGQL, result.GetErrorCode(), result.GetErrorMsg())
	}
	return result
}

// executeNGQLIgnoreError 执行 nGQL 语句，忽略错误（用于清理操作）
func executeNGQLIgnoreError(t *testing.T, nGQL string) {
	result, err := sessionPool.Execute("USE " + testSpace + "; " + nGQL)
	if err != nil {
		t.Logf("执行 nGQL 出错（已忽略）: %s\n错误: %v", nGQL, err)
		return
	}
	if !result.IsSucceed() {
		t.Logf("nGQL 执行失败（已忽略）: %s\n错误码: %d, 错误信息: %s",
			nGQL, result.GetErrorCode(), result.GetErrorMsg())
	}
}

// executeGlobalNGQL 执行全局 nGQL 语句（不带空间上下文，如 CREATE/DROP SPACE）
func executeGlobalNGQL(t *testing.T, nGQL string) *nebula.ResultSet {
	result, err := sessionPool.Execute(nGQL)
	if err != nil {
		t.Fatalf("执行 nGQL 失败: %s\n错误: %v", nGQL, err)
	}
	if !result.IsSucceed() {
		t.Fatalf("nGQL 执行失败: %s\n错误码: %d, 错误信息: %s",
			nGQL, result.GetErrorCode(), result.GetErrorMsg())
	}
	return result
}

// executeGlobalNGQLIgnoreError 执行全局 nGQL 语句并忽略错误（用于清理操作）
func executeGlobalNGQLIgnoreError(t *testing.T, nGQL string) {
	result, err := sessionPool.Execute(nGQL)
	if err != nil {
		t.Logf("执行 nGQL 出错（已忽略）: %s\n错误: %v", nGQL, err)
		return
	}
	if !result.IsSucceed() {
		t.Logf("nGQL 执行失败（已忽略）: %s\n错误码: %d, 错误信息: %s",
			nGQL, result.GetErrorCode(), result.GetErrorMsg())
	}
}

// waitForSchema 等待 Tag/Edge schema 对 DML 可见
// DESCRIBE 直查 meta 立即可见，但 INSERT 等 DML 走 graphd 本地 schema 缓存（按心跳周期刷新），
// 因此 DESCRIBE 通过后需再等待一个缓存刷新周期
func waitForSchema(t *testing.T, stmt string) {
	for i := 0; i < 30; i++ {
		result, err := sessionPool.Execute("USE " + testSpace + "; " + stmt)
		if err == nil && result.IsSucceed() {
			time.Sleep(2 * time.Second) // 等待 graphd schema 缓存刷新
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("等待 schema 可见超时: %s", stmt)
}

// rebuildTagIndexWithRetry 重建 Tag 索引并等待其就绪
// REBUILD 对索引可见性的要求高于 DESCRIBE（需 schema 分发至 storaged），就绪前会报 SemanticError
func rebuildTagIndexWithRetry(t *testing.T, indexName string) {
	stmt := "REBUILD TAG INDEX " + indexName
	for i := 0; i < 30; i++ {
		result, err := sessionPool.Execute("USE " + testSpace + "; " + stmt)
		if err == nil && result.IsSucceed() {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("重建索引 %s 超时", indexName)
}

// createTestSpaceWithRetry 创建测试图空间，冷启动时 meta 心跳未活跃会暂时报 Host not enough
func createTestSpaceWithRetry(t *testing.T, stmt string) {
	for i := 0; i < 30; i++ {
		result, err := sessionPool.Execute(stmt)
		if err == nil && result.IsSucceed() {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("创建测试图空间超时: %s", stmt)
}

// TestIntegration_FullWorkflow 完整的集成测试工作流
func TestIntegration_FullWorkflow(t *testing.T) {
	// 清理可能存在的测试图空间
	t.Log("=== 清理测试环境 ===")
	cleanupNGQL := ngql.NewDropSpace(testSpace).
		IfExists().
		String()
	executeGlobalNGQLIgnoreError(t, cleanupNGQL)
	time.Sleep(2 * time.Second) // 等待删除完成

	// 1. 创建测试图空间
	t.Log("=== 1. 创建测试图空间 ===")
	createSpaceNGQL := ngql.NewCreateSpace(testSpace).
		PartitionNum(10).
		ReplicaFactor(1).
		VidTypeFixedString(32).
		Comment("集成测试图空间").
		String()
	t.Logf("执行: %s", createSpaceNGQL)
	createTestSpaceWithRetry(t, createSpaceNGQL)
	time.Sleep(10 * time.Second) // 等待图空间创建完成（NebulaGraph 异步创建）

	// 2. 创建 Tag
	t.Log("=== 2. 创建 Tag ===")
	createPlayerTagNGQL := ngql.NewCreateTag("player").
		IfNotExists().
		Property("name", ngql.PropertyTypeString).
		Property("age", ngql.PropertyTypeInt32).
		Property("score", ngql.PropertyTypeInt32).
		Comment("球员信息").
		String()
	t.Logf("执行: %s", createPlayerTagNGQL)
	executeNGQL(t, createPlayerTagNGQL)

	createTeamTagNGQL := ngql.NewCreateTag("team").
		IfNotExists().
		Property("name", ngql.PropertyTypeString).
		Property("city", ngql.PropertyTypeString).
		String()
	t.Logf("执行: %s", createTeamTagNGQL)
	executeNGQL(t, createTeamTagNGQL)
	waitForSchema(t, "DESCRIBE TAG player")
	waitForSchema(t, "DESCRIBE TAG team")

	// 3. 创建 Edge 类型
	t.Log("=== 3. 创建 Edge 类型 ===")
	createFollowEdgeNGQL := ngql.NewCreateEdge("follow").
		IfNotExists().
		Property("degree", ngql.PropertyTypeInt32).
		Comment("关注关系").
		String()
	t.Logf("执行: %s", createFollowEdgeNGQL)
	executeNGQL(t, createFollowEdgeNGQL)

	createServeEdgeNGQL := ngql.NewCreateEdge("serve").
		IfNotExists().
		Property("start_year", ngql.PropertyTypeInt32).
		Property("end_year", ngql.PropertyTypeInt32).
		String()
	t.Logf("执行: %s", createServeEdgeNGQL)
	executeNGQL(t, createServeEdgeNGQL)
	waitForSchema(t, "DESCRIBE EDGE follow")
	waitForSchema(t, "DESCRIBE EDGE serve")

	// 4. 插入顶点数据
	t.Log("=== 4. 插入顶点数据 ===")
	insertPlayersNGQL := ngql.NewInsertVertex("player", "name", "age", "score").
		Values("player100", "Tim Duncan", 42, 1000).
		Values("player101", "Tony Parker", 36, 800).
		Values("player102", "LaMarcus Aldridge", 33, 700).
		String()
	t.Logf("执行: %s", insertPlayersNGQL)
	executeNGQL(t, insertPlayersNGQL)

	insertTeamsNGQL := ngql.NewInsertVertex("team", "name", "city").
		Values("team204", "Spurs", "San Antonio").
		Values("team215", "Hornets", "Charlotte").
		String()
	t.Logf("执行: %s", insertTeamsNGQL)
	executeNGQL(t, insertTeamsNGQL)

	// 5. 插入边数据
	t.Log("=== 5. 插入边数据 ===")
	insertFollowEdgesNGQL := ngql.NewInsertEdge("follow", "degree").
		Edge("player100", "player101", 95).
		Edge("player101", "player102", 90).
		Edge("player102", "player100", 75).
		String()
	t.Logf("执行: %s", insertFollowEdgesNGQL)
	executeNGQL(t, insertFollowEdgesNGQL)

	insertServeEdgesNGQL := ngql.NewInsertEdge("serve", "start_year", "end_year").
		Edge("player100", "team204", 1997, 2016).
		Edge("player101", "team204", 2001, 2018).
		String()
	t.Logf("执行: %s", insertServeEdgesNGQL)
	executeNGQL(t, insertServeEdgesNGQL)

	// 创建索引（MATCH 查询需要索引）
	t.Log("=== 创建索引 ===")
	createIndexNGQL := ngql.NewCreateTagIndex("idx_player_age", "player").
		IfNotExists().
		Property("age", 0).
		String()
	t.Logf("执行: %s", createIndexNGQL)
	executeNGQL(t, createIndexNGQL)

	// 重建索引
	rebuildTagIndexWithRetry(t, "idx_player_age")
	time.Sleep(5 * time.Second) // 等待索引构建完成

	// 6. 查询数据 - MATCH
	t.Log("=== 6. 查询数据 - MATCH ===")
	p := ngql.NewVertexPattern("v").Tag("player")
	matchNGQL := ngql.NewMatch(p).
		Where("v.player.age > 35").
		Return("v.player.name AS name", "v.player.age AS age").
		OrderBy("age", true). // 使用别名而不是完整路径
		String()
	t.Logf("执行: %s", matchNGQL)
	result := executeNGQL(t, matchNGQL)
	t.Logf("查询结果行数: %d", result.GetRowSize())
	if result.GetRowSize() < 2 {
		t.Errorf("期望至少 2 行结果，实际: %d", result.GetRowSize())
	}

	// 7. 查询数据 - GO
	t.Log("=== 7. 查询数据 - GO ===")
	goNGQL := ngql.NewGo("player100").
		Over("follow").
		Yield("$$.player.name AS name", "$$.player.age AS age").
		String()
	t.Logf("执行: %s", goNGQL)
	result = executeNGQL(t, goNGQL)
	t.Logf("查询结果行数: %d", result.GetRowSize())
	if result.GetRowSize() < 1 {
		t.Errorf("期望至少 1 行结果，实际: %d", result.GetRowSize())
	}

	// 8. 查询数据 - FETCH
	t.Log("=== 8. 查询数据 - FETCH ===")
	fetchNGQL := ngql.NewFetchVertex("player", "player100").
		Yield("player.name", "player.age", "player.score").
		String()
	t.Logf("执行: %s", fetchNGQL)
	result = executeNGQL(t, fetchNGQL)
	t.Logf("查询结果行数: %d", result.GetRowSize())

	// 9. 更新顶点数据
	t.Log("=== 9. 更新顶点数据 ===")
	updateNGQL := ngql.NewUpdateVertex("player", "player100").
		Set("age", 43).
		SetExpr("score = score + 100").
		Yield("name", "age", "score"). // UPDATE 语句的 YIELD 不需要 tag 前缀
		String()
	t.Logf("执行: %s", updateNGQL)
	result = executeNGQL(t, updateNGQL)
	t.Logf("更新结果行数: %d", result.GetRowSize())

	// 验证更新结果
	fetchAfterUpdateNGQL := ngql.NewFetchVertex("player", "player100").
		Yield("player.age", "player.score").
		String()
	t.Logf("执行: %s", fetchAfterUpdateNGQL)
	result = executeNGQL(t, fetchAfterUpdateNGQL)
	t.Logf("验证更新后的数据行数: %d", result.GetRowSize())

	// 10. UPSERT 操作
	t.Log("=== 10. UPSERT 操作 ===")
	upsertNGQL := ngql.NewUpsertVertex("player", "player103").
		Set("name", "Manu Ginobili").
		Set("age", 41).
		Set("score", 900).
		Yield("name", "age"). // UPSERT 语句的 YIELD 不需要 tag 前缀
		String()
	t.Logf("执行: %s", upsertNGQL)
	result = executeNGQL(t, upsertNGQL)
	t.Logf("UPSERT 结果行数: %d", result.GetRowSize())

	// 11. 删除边
	t.Log("=== 11. 删除边 ===")
	deleteEdgeNGQL := ngql.NewDeleteEdge("follow").
		Edge("player102", "player100").
		String()
	t.Logf("执行: %s", deleteEdgeNGQL)
	executeNGQL(t, deleteEdgeNGQL)

	// 12. 删除顶点
	t.Log("=== 12. 删除顶点 ===")
	deleteVertexNGQL := ngql.NewDeleteVertex("player103").
		WithEdge().
		String()
	t.Logf("执行: %s", deleteVertexNGQL)
	executeNGQL(t, deleteVertexNGQL)

	// 13. 清理测试图空间
	t.Log("=== 13. 清理测试图空间 ===")
	dropSpaceNGQL := ngql.NewDropSpace(testSpace).
		IfExists().
		String()
	t.Logf("执行: %s", dropSpaceNGQL)
	executeGlobalNGQL(t, dropSpaceNGQL)

	t.Log("=== 集成测试完成 ===")
}

// TestIntegration_IndexOperations 索引操作集成测试
func TestIntegration_IndexOperations(t *testing.T) {
	// 清理可能存在的测试图空间
	t.Log("=== 清理测试环境 ===")
	cleanupNGQL := ngql.NewDropSpace(testSpace).
		IfExists().
		String()
	executeGlobalNGQLIgnoreError(t, cleanupNGQL)
	time.Sleep(2 * time.Second)

	// 创建测试图空间
	t.Log("=== 创建测试图空间 ===")
	createSpaceNGQL := ngql.NewCreateSpace(testSpace).
		PartitionNum(10).
		ReplicaFactor(1).
		VidTypeFixedString(32).
		String()
	createTestSpaceWithRetry(t, createSpaceNGQL)
	time.Sleep(10 * time.Second)

	// 创建 player Tag
	t.Log("=== 创建 player Tag ===")
	createPlayerTagNGQL := ngql.NewCreateTag("player").
		IfNotExists().
		Property("name", ngql.PropertyTypeString).
		Property("age", ngql.PropertyTypeInt32).
		Property("score", ngql.PropertyTypeInt32).
		String()
	executeNGQL(t, createPlayerTagNGQL)
	time.Sleep(5 * time.Second)

	// 1. 创建 Tag 索引
	t.Log("=== 1. 创建 Tag 索引 ===")
	createTagIndexNGQL := ngql.NewCreateTagIndex("idx_player_age", "player").
		IfNotExists().
		Property("age", 0).
		String()
	t.Logf("执行: %s", createTagIndexNGQL)
	executeNGQL(t, createTagIndexNGQL)

	// 2. 创建另一个 Tag 索引
	t.Log("=== 2. 创建另一个 Tag 索引 ===")
	createScoreIndexNGQL := ngql.NewCreateTagIndex("idx_player_score", "player").
		IfNotExists().
		Property("score", 0).
		String()
	t.Logf("执行: %s", createScoreIndexNGQL)
	executeNGQL(t, createScoreIndexNGQL)
	t.Log("=== 3. 显示所有 Tag 索引 ===")
	showIndexesNGQL := ngql.NewShowTagIndexes().String()
	t.Logf("执行: %s", showIndexesNGQL)
	result := executeNGQL(t, showIndexesNGQL)
	t.Logf("索引数量: %d", result.GetRowSize())

	// 4. 查看索引详情
	t.Log("=== 4. 查看索引详情 ===")
	describeIndexNGQL := ngql.NewDescribeTagIndex("idx_player_age").String()
	t.Logf("执行: %s", describeIndexNGQL)
	result = executeNGQL(t, describeIndexNGQL)
	t.Logf("索引详情行数: %d", result.GetRowSize())

	// 5. 重建索引
	t.Log("=== 5. 重建索引 ===")
	rebuildTagIndexWithRetry(t, "idx_player_age")

	// 6. 显示索引状态
	t.Log("=== 6. 显示索引状态 ===")
	showIndexStatusNGQL := ngql.NewShowTagIndexStatus().String()
	t.Logf("执行: %s", showIndexStatusNGQL)
	result = executeNGQL(t, showIndexStatusNGQL)
	t.Logf("索引状态行数: %d", result.GetRowSize())

	// 7. 删除索引
	t.Log("=== 7. 删除索引 ===")
	dropIndexNGQL := ngql.NewDropTagIndex("idx_player_score").
		IfExists().
		String()
	t.Logf("执行: %s", dropIndexNGQL)
	executeNGQL(t, dropIndexNGQL)

	// 清理测试图空间
	t.Log("=== 清理测试图空间 ===")
	dropSpaceNGQL := ngql.NewDropSpace(testSpace).
		IfExists().
		String()
	t.Logf("执行: %s", dropSpaceNGQL)
	executeGlobalNGQL(t, dropSpaceNGQL)

	t.Log("=== 索引操作测试完成 ===")
}
