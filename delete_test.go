package ngql

import (
	"testing"
)

func TestDeleteVertex(t *testing.T) {
	nGQL := NewDeleteVertex("player100", "player101").
		String()

	expected := `DELETE VERTEX "player100", "player101"`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDeleteVertexWithEdge(t *testing.T) {
	nGQL := NewDeleteVertex("player100").
		WithEdge().
		String()

	expected := `DELETE VERTEX "player100" WITH EDGE`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDeleteEdge(t *testing.T) {
	nGQL := NewDeleteEdge("follow").
		Edge("player100", "player101").
		Edge("player101", "player102").
		String()

	expected := `DELETE EDGE follow "player100" -> "player101", "player101" -> "player102"`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
