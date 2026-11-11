package ngql

import (
	"testing"
)

func TestCreateTag(t *testing.T) {
	nGQL := NewCreateTag("player").
		IfNotExists().
		AddProperty(&Property{Name: "name", Type: PropertyTypeString, Nullable: false}).
		AddProperty(&Property{Name: "age", Type: PropertyTypeInt32, Nullable: true, Default: 0}).
		TTL("age", 100).
		Comment("玩家标签").
		String()

	expected := "CREATE TAG IF NOT EXISTS player(name string NOT NULL, age int32 DEFAULT 0) TTL_DURATION = 100, TTL_COL = \"age\", COMMENT = \"玩家标签\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestCreateTagSimple(t *testing.T) {
	nGQL := NewCreateTag("team").
		Property("name", PropertyTypeString).
		Property("founded", PropertyTypeInt8).
		String()

	expected := "CREATE TAG team(name string, founded int8)"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestAlterTag(t *testing.T) {
	nGQL := NewAlterTag("player").
		Add(&Property{Name: "score", Type: PropertyTypeInt32, Nullable: true}).
		Drop("old_field").
		Change(&Property{Name: "age", Type: PropertyTypeInt64, Nullable: true}).
		String()

	expected := "ALTER TAG player ADD (score int32), DROP (old_field), CHANGE (age int64)"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestAlterTagTTLComment(t *testing.T) {
	nGQL := NewAlterTag("player").
		TTL(100, "age").
		Comment("球员标签").
		String()

	// TTL 子句与 COMMENT 之间以逗号分隔
	expected := "ALTER TAG player TTL_DURATION = 100, TTL_COL = \"age\", COMMENT = \"球员标签\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDropTag(t *testing.T) {
	nGQL := NewDropTag("player").
		IfExists().
		String()

	expected := "DROP TAG IF EXISTS player"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestShowTags(t *testing.T) {
	nGQL := NewShowTags().
		String()

	expected := `SHOW TAGS`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDescribeTag(t *testing.T) {
	nGQL := NewDescribeTag("player").
		String()

	expected := "DESCRIBE TAG player"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestShowCreateTag(t *testing.T) {
	nGQL := NewShowCreateTag("player").
		String()

	expected := "SHOW CREATE TAG player"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
