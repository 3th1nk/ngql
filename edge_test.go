package ngql

import (
	"testing"
)

func TestCreateEdge(t *testing.T) {
	nGQL := NewCreateEdge("follow").
		IfNotExists().
		AddProperty(&Property{Name: "degree", Type: PropertyTypeInt32, Nullable: false}).
		AddProperty(&Property{Name: "created_at", Type: PropertyTypeTimestamp, Nullable: true, Default: 0}).
		TTL("created_at", 100).
		Comment("关注关系").
		String()

	expected := "CREATE EDGE IF NOT EXISTS follow(degree int32 NOT NULL, created_at timestamp DEFAULT 0) TTL_DURATION = 100, TTL_COL = \"created_at\" COMMENT = \"关注关系\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestCreateEdgeSimple(t *testing.T) {
	nGQL := NewCreateEdge("serve").
		Property("start_year", PropertyTypeInt32).
		Property("end_year", PropertyTypeInt32).
		String()

	expected := "CREATE EDGE serve(start_year int32, end_year int32)"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestAlterEdge(t *testing.T) {
	nGQL := NewAlterEdge("follow").
		Add(&Property{Name: "weight", Type: PropertyTypeDouble, Nullable: true}).
		Drop("old_field").
		Change(&Property{Name: "degree", Type: PropertyTypeInt64, Nullable: true}).
		String()

	expected := "ALTER EDGE follow ADD (weight double), DROP (old_field), CHANGE (degree int64)"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDropEdge(t *testing.T) {
	nGQL := NewDropEdge("follow").
		IfExists().
		String()

	expected := "DROP EDGE IF EXISTS follow"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestShowEdges(t *testing.T) {
	nGQL := NewShowEdges().
		String()

	expected := `SHOW EDGES`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDescribeEdge(t *testing.T) {
	nGQL := NewDescribeEdge("follow").
		String()

	expected := "DESCRIBE EDGE follow"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestShowCreateEdge(t *testing.T) {
	nGQL := NewShowCreateEdge("follow").
		String()

	expected := "SHOW CREATE EDGE follow"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
