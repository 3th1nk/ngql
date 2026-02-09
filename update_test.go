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
