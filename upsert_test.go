package ngql

import (
	"testing"
)

func TestUpsertVertex(t *testing.T) {
	nGQL := NewUpsertVertex("player", "player100").
		Set("age", 30).
		SetExpr("score = score + 10").
		When("name == \"Joe\"").
		Yield("name", "age", "score").
		String()

	expected := `UPSERT VERTEX ON player "player100" SET age = 30, score = score + 10 WHEN name == "Joe" YIELD name, age, score`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpsertVertexSimple(t *testing.T) {
	nGQL := NewUpsertVertex("player", "player101").
		Set("age", 25).
		String()

	expected := `UPSERT VERTEX ON player "player101" SET age = 25`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpsertEdge(t *testing.T) {
	nGQL := NewUpsertEdge("follow", "player100", "player101").
		Set("degree", 90).
		SetExpr("weight = weight + 0.1").
		When("degree < 100").
		Yield("degree", "weight").
		String()

	expected := `UPSERT EDGE ON follow "player100" -> "player101" SET degree = 90, weight = weight + 0.1 WHEN degree < 100 YIELD degree, weight`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpsertEdgeWithRank(t *testing.T) {
	nGQL := NewUpsertEdge("serve", "player100", "team204").
		WithRank(1).
		Set("end_year", 2020).
		String()

	expected := `UPSERT EDGE ON serve "player100" -> "team204"@1 SET end_year = 2020`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpsertVertexSetOrderPreserved(t *testing.T) {
	// 多次 Set 按调用顺序输出
	nGQL := NewUpsertVertex("player", "player100").
		Set("b", 2).
		Set("a", 1).
		String()

	expected := `UPSERT VERTEX ON player "player100" SET b = 2, a = 1`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpsertVertexSetOverwrite(t *testing.T) {
	// 同名属性后设置覆盖先设置
	nGQL := NewUpsertVertex("player", "player100").
		Set("age", 30).
		Set("age", 42).
		String()

	expected := `UPSERT VERTEX ON player "player100" SET age = 42`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpsertVertexWhenAccumulates(t *testing.T) {
	nGQL := NewUpsertVertex("player", "player100").
		Set("age", 42).
		When("age > 30").
		When("age < 50").
		String()

	expected := `UPSERT VERTEX ON player "player100" SET age = 42 WHEN age > 30 AND age < 50`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpsertEdgeKeywordIdentifiers(t *testing.T) {
	// 关键字标识符需自动加反引号
	nGQL := NewUpsertEdge("order", "player100", "player101").
		Set("from", "a").
		String()

	expected := "UPSERT EDGE ON `order` \"player100\" -> \"player101\" SET `from` = \"a\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
