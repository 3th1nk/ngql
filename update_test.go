package ngql

import (
	"testing"
)

func TestUpdateVertex(t *testing.T) {
	nGQL := NewUpdateVertex("player", "player100").
		Set("age", 43).
		SetExpr("score = score + 10").
		String()

	expected := `UPDATE VERTEX ON player "player100" SET age = 43, score = score + 10`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpdateVertexWithWhen(t *testing.T) {
	nGQL := NewUpdateVertex("player", "player100").
		Set("age", 43).
		When("age > 40").
		Yield("player.name", "player.age").
		String()

	expected := `UPDATE VERTEX ON player "player100" SET age = 43 WHEN age > 40 YIELD player.name, player.age`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpdateEdge(t *testing.T) {
	nGQL := NewUpdateEdge("follow", "player100", "player101").
		Set("degree", 96).
		String()

	expected := `UPDATE EDGE ON follow "player100" -> "player101" SET degree = 96`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpdateEdgeFull(t *testing.T) {
	nGQL := NewUpdateEdge("follow", "player100", "player101").
		WithRank(1).
		Set("degree", 96).
		SetExpr("weight = weight + 0.1").
		When("degree > 90").
		Yield("follow.degree").
		String()

	expected := `UPDATE EDGE ON follow "player100" -> "player101"@1 SET degree = 96, weight = weight + 0.1 WHEN degree > 90 YIELD follow.degree`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpdateVertexKeywordIdentifiers(t *testing.T) {
	// 关键字 Tag 名/属性名需自动加反引号
	nGQL := NewUpdateVertex("order", "player100").
		Set("from", "a").
		String()

	expected := "UPDATE VERTEX ON `order` \"player100\" SET `from` = \"a\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpdateEdgeKeywordIdentifiers(t *testing.T) {
	nGQL := NewUpdateEdge("case", "player100", "player101").
		Set("from", "a").
		String()

	expected := "UPDATE EDGE ON `case` \"player100\" -> \"player101\" SET `from` = \"a\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUpdateVertexWhenAccumulates(t *testing.T) {
	nGQL := NewUpdateVertex("player", "player100").
		Set("age", 43).
		When("age > 40").
		When("age < 50").
		String()

	expected := `UPDATE VERTEX ON player "player100" SET age = 43 WHEN age > 40 AND age < 50`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
