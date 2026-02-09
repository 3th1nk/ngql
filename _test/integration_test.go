package _test

import (
	"fmt"
	"testing"
	"time"

	"github.com/3th1nk/ngql"
	nebula "github.com/vesoft-inc/nebula-go/v3"
)

const (
	host      = "127.0.0.1"
	port      = 9669
	username  = "root"
	password  = "nebula"
	testSpace = "test_space"
)

var (
	sessionPool *nebula.SessionPool
)

// TestMain 测试主函数，初始化 SessionPool
func TestMain(m *testing.M) {
	poolConf, err := nebula.NewSessionPoolConf(
		username, password,
		[]nebula.HostAddress{
			{Host: host, Port: port},
		},
		"default_space", // 需要先手动创建图空间
	)
	if err != nil {
		fmt.Println("创建 SessionPoolConf 失败: ", err)
		return
	}
	sessionPool, err = nebula.NewSessionPool(*poolConf, nebula.DefaultLogger{})
	if err != nil {
		fmt.Println("创建 SessionPool 失败: ", err)
		return
	}
	m.Run()
}

// executeNGQL 执行 nGQL 语句并返回结果
func executeNGQL(t *testing.T, nGQL string) *nebula.ResultSet {
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

// executeNGQLIgnoreError 执行 nGQL 语句，忽略错误（用于清理操作）
func executeNGQLIgnoreError(t *testing.T, nGQL string) {
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

// TestIntegration_FullWorkflow 完整的集成测试工作流
func TestIntegration_FullWorkflow(t *testing.T) {
	// 清理可能存在的测试图空间
	t.Log("=== 清理测试环境 ===")
	cleanupNGQL := ngql.NewDropSpace(testSpace).
		IfExists().
		String()
	executeNGQLIgnoreError(t, cleanupNGQL)
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
	executeNGQL(t, createSpaceNGQL)
	time.Sleep(10 * time.Second) // 等待图空间创建完成（NebulaGraph 异步创建）

	// 使用测试图空间
	useSpaceNGQL := ngql.NewUseSpace(testSpace).String()
	t.Logf("执行: %s", useSpaceNGQL)
	executeNGQL(t, useSpaceNGQL)

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
	time.Sleep(5 * time.Second) // 等待 Tag 创建完成（NebulaGraph 异步创建）

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
	time.Sleep(5 * time.Second) // 等待 Edge 类型创建完成（NebulaGraph 异步创建）

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
	time.Sleep(3 * time.Second) // 等待索引创建完成

	// 重建索引
	rebuildIndexNGQL := ngql.NewRebuildTagIndex("idx_player_age").String()
	t.Logf("执行: %s", rebuildIndexNGQL)
	executeNGQL(t, rebuildIndexNGQL)
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
	executeNGQL(t, dropSpaceNGQL)

	t.Log("=== 集成测试完成 ===")
}

// TestIntegration_IndexOperations 索引操作集成测试
func TestIntegration_IndexOperations(t *testing.T) {
	// 清理可能存在的测试图空间
	t.Log("=== 清理测试环境 ===")
	cleanupNGQL := ngql.NewDropSpace(testSpace).
		IfExists().
		String()
	executeNGQLIgnoreError(t, cleanupNGQL)
	time.Sleep(2 * time.Second)

	// 创建测试图空间
	t.Log("=== 创建测试图空间 ===")
	createSpaceNGQL := ngql.NewCreateSpace(testSpace).
		PartitionNum(10).
		ReplicaFactor(1).
		VidTypeFixedString(32).
		String()
	executeNGQL(t, createSpaceNGQL)
	time.Sleep(10 * time.Second)

	// 使用测试图空间
	useSpaceNGQL := ngql.NewUseSpace(testSpace).String()
	executeNGQL(t, useSpaceNGQL)

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

	// 3. 显示所有 Tag 索引
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
	rebuildIndexNGQL := ngql.NewRebuildTagIndex("idx_player_age").String()
	t.Logf("执行: %s", rebuildIndexNGQL)
	executeNGQL(t, rebuildIndexNGQL)

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
	executeNGQL(t, dropSpaceNGQL)

	t.Log("=== 索引操作测试完成 ===")
}
