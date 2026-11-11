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

func TestGoWhereAccumulates(t *testing.T) {
	nGQL := NewGo("player100").
		Over("follow").
		Where("follow.degree > 90").
		Where("$$.player.age > 30").
		String()

	expected := `GO FROM "player100" OVER follow WHERE follow.degree > 90 AND $$.player.age > 30`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFindPathWhereAccumulates(t *testing.T) {
	nGQL := NewFindPath("SHORTEST").
		From("player100").
		To("team204").
		Over("follow").
		Where("follow.degree > 90").
		Where("$$.team.name == \"Warriors\"").
		String()

	expected := `FIND SHORTEST PATH FROM "player100" TO "team204" OVER follow WHERE follow.degree > 90 AND $$.team.name == "Warriors"`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGoKeywordOver(t *testing.T) {
	// 关键字 Edge 类型名需自动加反引号
	nGQL := NewGo("player100").
		Over("order").
		String()

	expected := "GO FROM \"player100\" OVER `order`"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFindPathKeywordOver(t *testing.T) {
	nGQL := NewFindPath("SHORTEST").
		From("player100").
		To("team204").
		Over("order").
		String()

	expected := "FIND SHORTEST PATH FROM \"player100\" TO \"team204\" OVER `order`"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGetSubgraphKeywordEdges(t *testing.T) {
	nGQL := NewGetSubgraph(1).
		From("player100").
		In("order").
		Out("case").
		Both("follow").
		String()

	expected := "GET SUBGRAPH FROM \"player100\" IN `order` OUT `case` BOTH follow"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFetchVertexKeywordTag(t *testing.T) {
	nGQL := NewFetchVertex("order", "player100").
		String()

	expected := "FETCH PROP ON `order` \"player100\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFetchEdgeKeywordType(t *testing.T) {
	nGQL := NewFetchEdge("case").
		Edge("player100", "player101").
		String()

	expected := "FETCH PROP ON `case` \"player100\" -> \"player101\""
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestLookupKeywordOn(t *testing.T) {
	nGQL := NewLookup("order").
		String()

	expected := "LOOKUP ON `order`"
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGroupByNoVars(t *testing.T) {
	// 未指定 GROUP BY 变量时不应输出悬空的 GROUP BY 子句
	prevQuery := `GO FROM "player100" OVER follow YIELD $$.player.age AS age`
	nGQL := NewGroupBy(prevQuery).
		Yield("count(*) AS count").
		String()

	expected := prevQuery + ` | YIELD count(*) AS count`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestMatchLimitZero(t *testing.T) {
	// Limit(0) 应生效而非被忽略
	nGQL := NewMatch(NewVertexPattern("v")).Limit(0).String()

	expected := `MATCH (v) LIMIT 0`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestMatchSkip(t *testing.T) {
	pattern := NewVertexPattern("v").Tag("player")
	nGQL := NewMatch(pattern).
		Return("v.player.name").
		Skip(10).
		Limit(5).
		String()

	expected := `MATCH (v:player) RETURN v.player.name SKIP 10 LIMIT 5`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestWithFieldAndSkip(t *testing.T) {
	withClause := NewWith().
		Field("v.player.age").
		Skip(1).
		Limit(2)

	nGQL := NewMatch(NewVertexPattern("v").Tag("player")).
		With(withClause).
		Return("v.player.age").
		String()

	expected := `MATCH (v:player) WITH v.player.age SKIP 1 LIMIT 2 RETURN v.player.age`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestYieldField(t *testing.T) {
	nGQL := NewYield().
		Field("1+1").
		Field("2+2").
		String()

	expected := `YIELD 1+1, 2+2`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGoBidirect(t *testing.T) {
	nGQL := NewGo("player100").
		Over("follow").
		Bidirect().
		Yield("dst(vertex)").
		String()

	expected := `GO FROM "player100" OVER follow BIDIRECT YIELD dst(vertex)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGoOverAllEdges(t *testing.T) {
	// 未指定 Over 时遍历所有边类型（OVER *）
	nGQL := NewGo("player100").
		Yield("dst(vertex)").
		String()

	expected := `GO FROM "player100" OVER * YIELD dst(vertex)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestMatchOrderByAsc(t *testing.T) {
	nGQL := NewMatch(NewVertexPattern("v").Tag("player")).
		Return("v.player.name").
		OrderBy("v.player.name", false).
		String()

	expected := `MATCH (v:player) RETURN v.player.name ORDER BY v.player.name ASC`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestWithOrderByAsc(t *testing.T) {
	withClause := NewWith().
		FieldAs("v.player.name", "name").
		OrderBy("name", false)

	nGQL := NewMatch(NewVertexPattern("v").Tag("player")).
		With(withClause).
		Return("name").
		String()

	expected := `MATCH (v:player) WITH v.player.name AS name ORDER BY name ASC RETURN name`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFindPathOrderByDesc(t *testing.T) {
	nGQL := NewFindPath(PathTypeAll).
		From("player100").
		To("team204").
		Over("follow").
		Yield("path AS p").
		OrderBy("path", true).
		String()

	expected := `FIND ALL PATH FROM "player100" TO "team204" OVER follow YIELD path AS p | ORDER BY path DESC`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFetchEdgeWithRank(t *testing.T) {
	nGQL := NewFetchEdge("serve").
		EdgeWithRank("player100", "team204", 1).
		Yield("properties(edge)").
		String()

	expected := `FETCH PROP ON serve "player100" -> "team204"@1 YIELD properties(edge)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestLookupSkip(t *testing.T) {
	nGQL := NewLookup("player").
		Yield("player.name").
		Skip(5).
		Limit(10).
		String()

	expected := `LOOKUP ON player YIELD player.name | SKIP 5 | LIMIT 10`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGroupByOrderByAscSkip(t *testing.T) {
	prevQuery := `GO FROM "player100" OVER follow YIELD $$.player.age AS age`
	nGQL := NewGroupBy(prevQuery, "$-.age").
		Yield("$-.age", "count(*) AS count").
		OrderBy("count", false).
		Skip(2).
		String()

	expected := prevQuery + ` | GROUP BY $-.age YIELD $-.age, count(*) AS count | ORDER BY count ASC | SKIP 2`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestFindPathReverselyOrderBy(t *testing.T) {
	nGQL := NewFindPath(PathTypeAll).
		From("player100").
		To("team204").
		Over("follow").
		Reversely().
		Yield("path AS p").
		OrderBy("path", false).
		String()

	expected := `FIND ALL PATH FROM "player100" TO "team204" OVER follow REVERSELY YIELD path AS p | ORDER BY path ASC`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestGetSubgraphDefaultSteps(t *testing.T) {
	// steps <= 0 时默认为 1 步，且单步不输出 STEPS
	nGQL := NewGetSubgraph(0).
		From("player100").
		Yield("VERTICES AS nodes").
		String()

	expected := `GET SUBGRAPH FROM "player100" YIELD VERTICES AS nodes`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestPathTypeConstants(t *testing.T) {
	// pathType 常量与字符串字面量等价
	cases := []struct {
		pathType string
		expected string
	}{
		{PathTypeShortest, `FIND SHORTEST PATH FROM "a" TO "b" OVER *`},
		{PathTypeAll, `FIND ALL PATH FROM "a" TO "b" OVER *`},
		{PathTypeNoLoop, `FIND NOLOOP PATH FROM "a" TO "b" OVER *`},
	}
	for _, c := range cases {
		nGQL := NewFindPath(c.pathType).From("a").To("b").String()
		if nGQL != c.expected {
			t.Errorf("Expected: %s\nGot: %s", c.expected, nGQL)
		}
	}
}

func TestGoStepsRange(t *testing.T) {
	// 变长遍历步数 M TO N STEPS
	nGQL := NewGo("player100").
		Over("follow").
		StepsRange(1, 3).
		Yield("$$.player.name").
		String()

	expected := `GO 1 TO 3 STEPS FROM "player100" OVER follow YIELD $$.player.name`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestYield(t *testing.T) {
	// 独立 YIELD 语句
	nGQL := NewYield().
		FieldAs("1+1", "x").
		FieldAs(`"a" == "a"`, "eq").
		Field("count(*)").
		String()

	expected := `YIELD 1+1 AS x, "a" == "a" AS eq, count(*)`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestYieldDistinct(t *testing.T) {
	nGQL := NewYield().
		Distinct().
		FieldAs("2", "x").
		String()

	expected := `YIELD DISTINCT 2 AS x`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestMatchReturnDistinct(t *testing.T) {
	pattern := NewVertexPattern("v").Tag("player")
	nGQL := NewMatch(pattern).
		ReturnDistinct("v.player.name").
		String()

	expected := `MATCH (v:player) RETURN DISTINCT v.player.name`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestMatchOptional(t *testing.T) {
	pattern := NewVertexPattern("v").Tag("player")
	nGQL := NewOptionalMatch(pattern).
		Where("v.player.age > 30").
		Return("v.player.name").
		String()

	expected := `OPTIONAL MATCH (v:player) WHERE v.player.age > 30 RETURN v.player.name`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestMatchWith(t *testing.T) {
	// WITH 子句：字段 + WHERE + ORDER BY + LIMIT
	withClause := NewWith().
		FieldAs("v.player.name", "name").
		Where(`name != "Tim"`).
		OrderBy("name", true).
		Limit(5)

	pattern := NewVertexPattern("v").Tag("player")
	nGQL := NewMatch(pattern).
		Where("v.player.age > 30").
		With(withClause).
		Return("name").
		String()

	expected := `MATCH (v:player) WHERE v.player.age > 30 WITH v.player.name AS name WHERE name != "Tim" ORDER BY name DESC LIMIT 5 RETURN name`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}

func TestMatchUnwind(t *testing.T) {
	// WITH 收集聚合后 UNWIND 展开
	pattern := NewVertexPattern("v").Tag("player")
	nGQL := NewMatch(pattern).
		With(NewWith().FieldAs("collect(v.player.name)", "names")).
		Unwind("names", "name").
		Return("name").
		String()

	expected := `MATCH (v:player) WITH collect(v.player.name) AS names UNWIND names AS name RETURN name`
	if nGQL != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, nGQL)
	}
}
