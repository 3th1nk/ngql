package ngql

import (
	"testing"
)

func TestCreateTagIndex(t *testing.T) {
	nGQL := NewCreateTagIndex("idx_player_age", "player").
		IfNotExists().
		Property("age", 0).
		Comment("年龄索引").
		String()

	expected := "CREATE TAG INDEX IF NOT EXISTS idx_player_age ON player(age) COMMENT \"年龄索引\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestCreateTagIndexComposite(t *testing.T) {
	nGQL := NewCreateTagIndex("idx_player_name_age", "player").
		Property("name", 20).
		Property("age", 0).
		String()

	expected := "CREATE TAG INDEX idx_player_name_age ON player(name(20), age)"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestCreateTagIndexEmpty(t *testing.T) {
	nGQL := NewCreateTagIndex("idx_player", "player").
		String()

	expected := "CREATE TAG INDEX idx_player ON player()"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestCreateEdgeIndex(t *testing.T) {
	nGQL := NewCreateEdgeIndex("idx_follow_degree", "follow").
		IfNotExists().
		Property("degree", 0).
		String()

	expected := "CREATE EDGE INDEX IF NOT EXISTS idx_follow_degree ON follow(degree)"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDropTagIndex(t *testing.T) {
	nGQL := NewDropTagIndex("idx_player_age").
		IfExists().
		String()

	expected := "DROP TAG INDEX IF EXISTS idx_player_age"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDropEdgeIndex(t *testing.T) {
	nGQL := NewDropEdgeIndex("idx_follow_degree").
		IfExists().
		String()

	expected := "DROP EDGE INDEX IF EXISTS idx_follow_degree"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestRebuildTagIndex(t *testing.T) {
	nGQL := NewRebuildTagIndex("idx_player_age", "idx_player_name").
		String()

	expected := "REBUILD TAG INDEX idx_player_age, idx_player_name"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestRebuildTagIndexAll(t *testing.T) {
	nGQL := NewRebuildTagIndex().
		String()

	expected := `REBUILD TAG INDEX`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestRebuildEdgeIndex(t *testing.T) {
	nGQL := NewRebuildEdgeIndex("idx_follow_degree").
		String()

	expected := "REBUILD EDGE INDEX idx_follow_degree"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestShowTagIndexes(t *testing.T) {
	nGQL := NewShowTagIndexes().
		String()

	expected := `SHOW TAG INDEXES`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestShowEdgeIndexes(t *testing.T) {
	nGQL := NewShowEdgeIndexes().
		String()

	expected := `SHOW EDGE INDEXES`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDescribeTagIndex(t *testing.T) {
	nGQL := NewDescribeTagIndex("idx_player_age").
		String()

	expected := "DESCRIBE TAG INDEX idx_player_age"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDescribeEdgeIndex(t *testing.T) {
	nGQL := NewDescribeEdgeIndex("idx_follow_degree").
		String()

	expected := "DESCRIBE EDGE INDEX idx_follow_degree"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestShowTagIndexStatus(t *testing.T) {
	nGQL := NewShowTagIndexStatus().
		String()

	expected := `SHOW TAG INDEX STATUS`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestShowEdgeIndexStatus(t *testing.T) {
	nGQL := NewShowEdgeIndexStatus().
		String()

	expected := `SHOW EDGE INDEX STATUS`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestCreateEdgeIndexComment(t *testing.T) {
	nGQL := NewCreateEdgeIndex("idx_follow_degree", "follow").
		Property("degree", 0).
		Comment("关注度索引").
		String()

	expected := "CREATE EDGE INDEX idx_follow_degree ON follow(degree) COMMENT \"关注度索引\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestCreateEdgeIndexEmpty(t *testing.T) {
	nGQL := NewCreateEdgeIndex("idx_follow", "follow").
		String()

	expected := "CREATE EDGE INDEX idx_follow ON follow()"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
