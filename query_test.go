package ngql

import (
	"testing"
)

func TestMatch(t *testing.T) {
	pattern := NewVertexPattern("v").Tag("player")
	nGQL := NewMatch(pattern).
		Where("v.player.age > 30").
		Return("v.player.name", "v.player.age").
		OrderBy("v.player.age", true).
		Limit(10).
		String()

	expected := `MATCH (v:player) WHERE v.player.age > 30 RETURN v.player.name, v.player.age ORDER BY v.player.age DESC LIMIT 10`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGo(t *testing.T) {
	nGQL := NewGo("player100").
		Over("follow").
		Steps(2).
		Yield("$$.player.name", "$$.player.age").
		String()

	expected := `GO 2 STEPS FROM "player100" OVER follow YIELD $$.player.name, $$.player.age`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGoReversely(t *testing.T) {
	nGQL := NewGo("player100").
		Over("follow").
		Reversely().
		Yield("$$.player.name").
		String()

	expected := `GO FROM "player100" OVER follow REVERSELY YIELD $$.player.name`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFetchVertex(t *testing.T) {
	nGQL := NewFetchVertex("player", "player100", "player101").
		Yield("properties(vertex)").
		String()

	expected := `FETCH PROP ON player "player100", "player101" YIELD properties(vertex)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFetchEdge(t *testing.T) {
	nGQL := NewFetchEdge("follow").
		Edge("player100", "player101").
		Yield("properties(edge)").
		String()

	expected := `FETCH PROP ON follow "player100" -> "player101" YIELD properties(edge)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestLookup(t *testing.T) {
	nGQL := NewLookup("player").
		Where("player.age > 30").
		Yield("id(vertex)", "player.name", "player.age").
		String()

	expected := `LOOKUP ON player WHERE player.age > 30 YIELD id(vertex), player.name, player.age`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestLookupWithDistinct(t *testing.T) {
	nGQL := NewLookup("player").
		Where("player.name STARTS WITH \"T\"").
		Distinct().
		Yield("player.age").
		Limit(10).
		String()

	expected := `LOOKUP ON player WHERE player.name STARTS WITH "T" YIELD DISTINCT player.age | LIMIT 10`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestLookupEdge(t *testing.T) {
	nGQL := NewLookup("follow").
		Where("follow.degree > 90").
		Yield("properties(edge)").
		String()

	expected := `LOOKUP ON follow WHERE follow.degree > 90 YIELD properties(edge)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGroupBy(t *testing.T) {
	prevQuery := `GO FROM "player100" OVER follow YIELD $$.player.age AS age`
	nGQL := NewGroupBy(prevQuery, "$-.age").
		Yield("$-.age", "count(*) AS count").
		String()

	expected := `GO FROM "player100" OVER follow YIELD $$.player.age AS age | GROUP BY $-.age YIELD $-.age, count(*) AS count`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGroupByWithOrderAndLimit(t *testing.T) {
	prevQuery := `GO FROM "player100" OVER follow YIELD $$.player.age AS age, $$.player.name AS name`
	nGQL := NewGroupBy(prevQuery, "$-.age").
		Yield("$-.age", "count(*) AS count", "collect($-.name) AS names").
		OrderBy("count", true).
		Limit(5).
		String()

	expected := `GO FROM "player100" OVER follow YIELD $$.player.age AS age, $$.player.name AS name | GROUP BY $-.age YIELD $-.age, count(*) AS count, collect($-.name) AS names | ORDER BY count DESC | LIMIT 5`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFindShortestPath(t *testing.T) {
	nGQL := NewFindPath("SHORTEST").
		From("player100").
		To("team204").
		Over("*").
		Yield("path AS p").
		String()

	expected := `FIND SHORTEST PATH FROM "player100" TO "team204" OVER * YIELD path AS p`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFindAllPath(t *testing.T) {
	nGQL := NewFindPath("ALL").
		From("player100").
		To("player101", "player102").
		Over("follow", "serve").
		UpTo(5).
		Yield("path AS p").
		Limit(10).
		String()

	expected := `FIND ALL PATH FROM "player100" TO "player101", "player102" OVER follow, serve UPTO 5 STEPS YIELD path AS p | LIMIT 10`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFindPathWithProp(t *testing.T) {
	nGQL := NewFindPath("SHORTEST").
		WithProp().
		From("player100").
		To("team204").
		Over("*").
		Bidirect().
		Yield("path AS p").
		String()

	expected := `FIND SHORTEST PATH WITH PROP FROM "player100" TO "team204" OVER * BIDIRECT YIELD path AS p`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGetSubgraph(t *testing.T) {
	nGQL := NewGetSubgraph(2).
		From("player100").
		Yield("VERTICES AS nodes", "EDGES AS relationships").
		String()

	expected := `GET SUBGRAPH 2 STEPS FROM "player100" YIELD VERTICES AS nodes, EDGES AS relationships`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGetSubgraphWithProp(t *testing.T) {
	nGQL := NewGetSubgraph(1).
		WithProp().
		From("player100", "player101").
		In("follow").
		Out("serve").
		Yield("VERTICES AS nodes").
		String()

	expected := `GET SUBGRAPH WITH PROP FROM "player100", "player101" IN follow OUT serve YIELD VERTICES AS nodes`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGetSubgraphBoth(t *testing.T) {
	nGQL := NewGetSubgraph(3).
		From("player100").
		Both("follow", "serve").
		Where("$$.player.age > 30").
		Yield("VERTICES AS nodes", "EDGES AS edges").
		String()

	expected := `GET SUBGRAPH 3 STEPS FROM "player100" BOTH follow, serve WHERE $$.player.age > 30 YIELD VERTICES AS nodes, EDGES AS edges`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
