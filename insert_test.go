package ngql

import (
	"testing"
)

func TestInsertVertex(t *testing.T) {
	nGQL := NewInsertVertex("player", "name", "age").
		Values("player100", "Tim Duncan", 42).
		Values("player101", "Tony Parker", 36).
		String()

	expected := `INSERT VERTEX player(name, age) VALUES "player100": ("Tim Duncan", 42), "player101": ("Tony Parker", 36)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestInsertVertexIfNotExists(t *testing.T) {
	nGQL := NewInsertVertex("player", "name", "age").
		IfNotExists().
		Values("player100", "Tim Duncan", 42).
		String()

	expected := `INSERT VERTEX IF NOT EXISTS player(name, age) VALUES "player100": ("Tim Duncan", 42)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestInsertEdge(t *testing.T) {
	nGQL := NewInsertEdge("follow", "degree").
		Edge("player100", "player101", 95).
		Edge("player101", "player102", 90).
		String()

	expected := `INSERT EDGE follow(degree) VALUES "player100" -> "player101": (95), "player101" -> "player102": (90)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestInsertEdgeWithRank(t *testing.T) {
	nGQL := NewInsertEdge("follow", "degree").
		EdgeWithRank("player100", "player101", 1, 95).
		String()

	expected := `INSERT EDGE follow(degree) VALUES "player100" -> "player101"@1: (95)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
