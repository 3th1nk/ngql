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

func TestDeleteEdgeKeywordType(t *testing.T) {
	// 关键字 Edge 类型名需自动加反引号
	nGQL := NewDeleteEdge("order").
		Edge("player100", "player101").
		String()

	expected := "DELETE EDGE `order` \"player100\" -> \"player101\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDeleteTag(t *testing.T) {
	// 删除顶点上指定的 Tag
	nGQL := NewDeleteTag("player", "player100", "player101").
		String()

	expected := `DELETE TAG player FROM "player100", "player101"`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDeleteTagKeywordName(t *testing.T) {
	nGQL := NewDeleteTag("order", "player100").
		String()

	expected := "DELETE TAG `order` FROM \"player100\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
