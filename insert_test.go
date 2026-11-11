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

func TestInsertVertexKeywordIdentifiers(t *testing.T) {
	// 关键字 Tag 名/属性名需自动加反引号，普通名称不加
	nGQL := NewInsertVertex("order", "from", "status").
		Values("player100", "a", 1).
		String()

	expected := "INSERT VERTEX `order`(`from`, status) VALUES \"player100\": (\"a\", 1)"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestInsertEdgeKeywordIdentifiers(t *testing.T) {
	nGQL := NewInsertEdge("case", "degree").
		Edge("player100", "player101", 95).
		String()

	expected := `INSERT EDGE ` + "`case`" + `(degree) VALUES "player100" -> "player101": (95)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestInsertVertexMultiTag(t *testing.T) {
	// 多 Tag 同 vid：属性值按 Tag 添加顺序拼接
	nGQL := NewInsertVertex("player", "name", "age").
		AddTag("person", "idcard").
		Values("player100", "Tim Duncan", 42, "123456").
		String()

	expected := `INSERT VERTEX player(name, age), person(idcard) VALUES "player100": ("Tim Duncan", 42, "123456")`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestInsertVertexMultiTagBatch(t *testing.T) {
	nGQL := NewInsertVertex("t1", "a").
		AddTag("t2", "b", "c").
		Values("v1", 1, 2, 3).
		Values("v2", 4, 5, 6).
		String()

	expected := `INSERT VERTEX t1(a), t2(b, c) VALUES "v1": (1, 2, 3), "v2": (4, 5, 6)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
