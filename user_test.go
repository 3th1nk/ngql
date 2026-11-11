package ngql

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	cases := []struct {
		nGQL     string
		expected string
	}{
		{NewCreateUser("user1", "pwd1").String(), `CREATE USER user1 WITH PASSWORD "pwd1"`},
		{NewCreateUser("user1", "pwd1").IfNotExists().String(), `CREATE USER IF NOT EXISTS user1 WITH PASSWORD "pwd1"`},
	}
	for _, c := range cases {
		if c.nGQL != c.expected {
			t.Errorf("Expected: %s, Got: %s", c.expected, c.nGQL)
		}
	}
}

func TestDropUser(t *testing.T) {
	cases := []struct {
		nGQL     string
		expected string
	}{
		{NewDropUser("user1").String(), "DROP USER user1"},
		{NewDropUser("user1").IfExists().String(), "DROP USER IF EXISTS user1"},
	}
	for _, c := range cases {
		if c.nGQL != c.expected {
			t.Errorf("Expected: %s, Got: %s", c.expected, c.nGQL)
		}
	}
}

func TestAlterUser(t *testing.T) {
	nGQL := NewAlterUser("user1", "newpwd").String()
	expected := `ALTER USER user1 WITH PASSWORD "newpwd"`
	if nGQL != expected {
		t.Errorf("Expected: %s, Got: %s", expected, nGQL)
	}
}

func TestChangePassword(t *testing.T) {
	nGQL := NewChangePassword("user1", "oldpwd", "newpwd").String()
	expected := `CHANGE PASSWORD user1 FROM "oldpwd" TO "newpwd"`
	if nGQL != expected {
		t.Errorf("Expected: %s, Got: %s", expected, nGQL)
	}
}

func TestGrantRole(t *testing.T) {
	nGQL := NewGrantRole(RoleAdmin, "basketball", "user1").String()
	expected := "GRANT ROLE ADMIN ON basketball TO user1"
	if nGQL != expected {
		t.Errorf("Expected: %s, Got: %s", expected, nGQL)
	}
}

func TestRevokeRole(t *testing.T) {
	nGQL := NewRevokeRole(RoleUser, "basketball", "user1").String()
	expected := "REVOKE ROLE USER ON basketball FROM user1"
	if nGQL != expected {
		t.Errorf("Expected: %s, Got: %s", expected, nGQL)
	}
}

func TestShowUserAccess(t *testing.T) {
	cases := []struct {
		nGQL     string
		expected string
	}{
		{NewShowUsers().String(), "SHOW USERS"},
		{NewShowRoles("basketball").String(), "SHOW ROLES ON basketball"},
		{NewShowGrants("").String(), "SHOW GRANTS"},
		{NewShowGrants("basketball").String(), "SHOW GRANTS ON basketball"},
	}
	for _, c := range cases {
		if c.nGQL != c.expected {
			t.Errorf("Expected: %s, Got: %s", c.expected, c.nGQL)
		}
	}
}
