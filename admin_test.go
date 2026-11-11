package ngql

import (
	"testing"
)

// ============ SUBMIT JOB 测试 ============

func TestSubmitJob(t *testing.T) {
	cases := []struct {
		nGQL     string
		expected string
	}{
		{NewSubmitJob(JobTypeCompact).String(), "SUBMIT JOB COMPACT"},
		{NewSubmitJob(JobTypeFlush).String(), "SUBMIT JOB FLUSH"},
		{NewSubmitJob(JobTypeStats).String(), "SUBMIT JOB STATS"},
	}
	for _, c := range cases {
		if c.nGQL != c.expected {
			t.Errorf("Expected: %s, Got: %s", c.expected, c.nGQL)
		}
	}
}

func TestSubmitJobRebuildIndex(t *testing.T) {
	cases := []struct {
		nGQL     string
		expected string
	}{
		{NewSubmitJob(JobTypeRebuildIndex).String(), "SUBMIT JOB REBUILD INDEX"},
		{NewSubmitJob(JobTypeRebuildIndex).Index("idx1").String(), "SUBMIT JOB REBUILD INDEX:idx1"},
		{NewSubmitJob(JobTypeRebuildIndex).Index("idx1", "idx2").String(), "SUBMIT JOB REBUILD INDEX:idx1, idx2"},
	}
	for _, c := range cases {
		if c.nGQL != c.expected {
			t.Errorf("Expected: %s, Got: %s", c.expected, c.nGQL)
		}
	}
}

func TestJobManagement(t *testing.T) {
	cases := []struct {
		nGQL     string
		expected string
	}{
		{NewShowJobs().String(), "SHOW JOBS"},
		{NewShowJob(42).String(), "SHOW JOB 42"},
		{NewStopJob(42).String(), "STOP JOB 42"},
		{NewRecoverJob().String(), "RECOVER JOB"},
	}
	for _, c := range cases {
		if c.nGQL != c.expected {
			t.Errorf("Expected: %s, Got: %s", c.expected, c.nGQL)
		}
	}
}

// ============ SHOW 家族测试 ============

func TestShowAdmin(t *testing.T) {
	cases := []struct {
		nGQL     string
		expected string
	}{
		{NewShowHosts().String(), "SHOW HOSTS"},
		{NewShowHosts(HostModuleGraph).String(), "SHOW HOSTS GRAPH"},
		{NewShowHosts(HostModuleStorage).String(), "SHOW HOSTS STORAGE"},
		{NewShowParts(-1).String(), "SHOW PARTS"},
		{NewShowParts(5).String(), "SHOW PARTS 5"},
		{NewShowStats().String(), "SHOW STATS"},
		{NewShowSnapshots().String(), "SHOW SNAPSHOTS"},
		{NewShowSessions().String(), "SHOW SESSIONS"},
		{NewShowCharset().String(), "SHOW CHARSET"},
		{NewShowCollation().String(), "SHOW COLLATION"},
	}
	for _, c := range cases {
		if c.nGQL != c.expected {
			t.Errorf("Expected: %s, Got: %s", c.expected, c.nGQL)
		}
	}
}

func TestShowQueries(t *testing.T) {
	cases := []struct {
		nGQL     string
		expected string
	}{
		{NewShowQueries("").String(), "SHOW QUERIES"},
		{NewShowQueries(QueryScopeCurrent).String(), "SHOW QUERIES CURRENT"},
		{NewShowQueries(QueryScopeAll).String(), "SHOW QUERIES ALL"},
	}
	for _, c := range cases {
		if c.nGQL != c.expected {
			t.Errorf("Expected: %s, Got: %s", c.expected, c.nGQL)
		}
	}
}

func TestShowCreateSpace(t *testing.T) {
	nGQL := NewShowCreateSpace("basketball").String()
	expected := "SHOW CREATE SPACE basketball"
	if nGQL != expected {
		t.Errorf("Expected: %s, Got: %s", expected, nGQL)
	}
}

func TestShowCreateSpaceKeywordName(t *testing.T) {
	nGQL := NewShowCreateSpace("order").String()
	expected := "SHOW CREATE SPACE `order`"
	if nGQL != expected {
		t.Errorf("Expected: %s, Got: %s", expected, nGQL)
	}
}

// ============ CONFIGS 测试 ============

func TestShowConfigs(t *testing.T) {
	cases := []struct {
		nGQL     string
		expected string
	}{
		{NewShowConfigs("").String(), "SHOW CONFIGS"},
		{NewShowConfigs(ConfigModuleGraph).String(), "SHOW CONFIGS graph"},
		{NewShowConfigs(ConfigModuleStorage).String(), "SHOW CONFIGS storage"},
	}
	for _, c := range cases {
		if c.nGQL != c.expected {
			t.Errorf("Expected: %s, Got: %s", c.expected, c.nGQL)
		}
	}
}

func TestUpdateConfig(t *testing.T) {
	cases := []struct {
		nGQL     string
		expected string
	}{
		{NewUpdateConfig("", "heartbeat_interval_secs", 10).String(), "UPDATE CONFIGS heartbeat_interval_secs=10"},
		{NewUpdateConfig("graph", "heartbeat_interval_secs", 10).String(), "UPDATE CONFIGS graph:heartbeat_interval_secs=10"},
		{NewUpdateConfig("storage", "rocksdb_db_path", "/data/storage").String(), `UPDATE CONFIGS storage:rocksdb_db_path="/data/storage"`},
		{NewUpdateConfig("", "enable_optimizer", true).String(), "UPDATE CONFIGS enable_optimizer=true"},
	}
	for _, c := range cases {
		if c.nGQL != c.expected {
			t.Errorf("Expected: %s, Got: %s", c.expected, c.nGQL)
		}
	}
}

func TestReloadConfigs(t *testing.T) {
	if got := NewReloadConfigs().String(); got != "RELOAD CONFIGS" {
		t.Errorf("Expected: RELOAD CONFIGS, Got: %s", got)
	}
}
