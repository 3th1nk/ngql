package ngql

import (
	"testing"
)

func TestCreateSpace(t *testing.T) {
	nGQL := NewCreateSpace("my_space").
		IfNotExists().
		PartitionNum(10).
		ReplicaFactor(1).
		VidTypeFixedString(30).
		Comment("测试图空间").
		String()

	expected := `CREATE SPACE IF NOT EXISTS my_space (partition_num = 10, replica_factor = 1, vid_type = FIXED_STRING(30)) COMMENT = "测试图空间"`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestCreateSpaceMinimal(t *testing.T) {
	nGQL := NewCreateSpace("my_space").
		VidTypeInt64().
		String()

	expected := `CREATE SPACE my_space (vid_type = INT64)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestCreateSpaceAs(t *testing.T) {
	nGQL := NewCreateSpaceAs("new_space", "old_space").
		String()

	expected := `CREATE SPACE new_space AS old_space`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestUseSpace(t *testing.T) {
	nGQL := NewUseSpace("my_space").
		String()

	expected := `USE my_space`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestShowSpaces(t *testing.T) {
	nGQL := NewShowSpaces().
		String()

	expected := `SHOW SPACES`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDescribeSpace(t *testing.T) {
	nGQL := NewDescribeSpace("my_space").
		String()

	expected := `DESCRIBE SPACE my_space`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestDropSpace(t *testing.T) {
	nGQL := NewDropSpace("my_space").
		IfExists().
		String()

	expected := `DROP SPACE IF EXISTS my_space`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestClearSpace(t *testing.T) {
	nGQL := NewClearSpace("my_space").
		IfExists().
		String()

	expected := `CLEAR SPACE IF EXISTS my_space`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
