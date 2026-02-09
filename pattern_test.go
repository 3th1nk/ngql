package ngql

import (
	"testing"
)

// ============ 顶点模式测试 ============

func TestVertexPattern_Anonymous(t *testing.T) {
	pattern := NewVertexPattern("")
	expected := "()"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestVertexPattern_WithVariable(t *testing.T) {
	pattern := NewVertexPattern("v")
	expected := "(v)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestVertexPattern_WithTag(t *testing.T) {
	pattern := NewVertexPattern("v").Tag("player")
	expected := "(v:player)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestVertexPattern_WithMultipleTags(t *testing.T) {
	pattern := NewVertexPattern("v").Tag("player").Tag("person")
	expected := "(v:player:person)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestVertexPattern_WithProperty(t *testing.T) {
	pattern := NewVertexPattern("v").Tag("player").Property("name", "Tim Duncan")
	result := pattern.String()
	// 属性顺序可能不固定，所以只检查关键部分
	if result != `(v:player {name: "Tim Duncan"})` {
		t.Errorf("Expected: (v:player {name: \"Tim Duncan\"}), Got: %s", result)
	}
}

func TestVertexPattern_WithMultipleProperties(t *testing.T) {
	pattern := NewVertexPattern("v").
		Tag("player").
		Property("name", "Tim Duncan").
		Property("age", 42)
	result := pattern.String()
	// 检查是否包含所有必要的部分
	if !(result == `(v:player {name: "Tim Duncan", age: 42})` ||
		result == `(v:player {age: 42, name: "Tim Duncan"})`) {
		t.Errorf("Unexpected result: %s", result)
	}
}

// ============ 边模式测试 ============

func TestEdgePattern_RightDirection(t *testing.T) {
	pattern := NewEdgePattern(EdgeDirectionRight)
	expected := "->"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestEdgePattern_LeftDirection(t *testing.T) {
	pattern := NewEdgePattern(EdgeDirectionLeft)
	expected := "<-"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestEdgePattern_BothDirection(t *testing.T) {
	pattern := NewEdgePattern(EdgeDirectionBoth)
	expected := "-"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestEdgePattern_WithType(t *testing.T) {
	pattern := NewEdgePattern(EdgeDirectionRight).Type("follow")
	expected := "-[:follow]->"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestEdgePattern_WithVariable(t *testing.T) {
	pattern := NewEdgePattern(EdgeDirectionRight).Variable("e").Type("follow")
	expected := "-[e:follow]->"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestEdgePattern_WithMultipleTypes(t *testing.T) {
	pattern := NewEdgePattern(EdgeDirectionRight).Type("follow").Type("like")
	expected := "-[:follow|like]->"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestEdgePattern_WithProperty(t *testing.T) {
	pattern := NewEdgePattern(EdgeDirectionRight).
		Type("follow").
		Property("degree", 90)
	expected := "-[:follow {degree: 90}]->"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestEdgePattern_VariableLength_Fixed(t *testing.T) {
	pattern := NewEdgePattern(EdgeDirectionRight).Type("follow").Hops(2, 2)
	expected := "-[:follow*2]->"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestEdgePattern_VariableLength_Range(t *testing.T) {
	pattern := NewEdgePattern(EdgeDirectionRight).Type("follow").Hops(1, 3)
	expected := "-[:follow*1..3]->"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestEdgePattern_VariableLength_MinOnly(t *testing.T) {
	pattern := NewEdgePattern(EdgeDirectionRight).Type("follow").Hops(1, 0)
	expected := "-[:follow*1..]->"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestEdgePattern_VariableLength_MaxOnly(t *testing.T) {
	pattern := NewEdgePattern(EdgeDirectionRight).Type("follow").Hops(0, 3)
	expected := "-[:follow*..3]->"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestEdgePattern_VariableLength_Any(t *testing.T) {
	// 任意跳数应该使用 Hops(1, 0) 表示 *1.. 或者不设置 Hops
	pattern := NewEdgePattern(EdgeDirectionRight).Type("follow").Hops(1, 0)
	expected := "-[:follow*1..]->"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

// ============ 路径模式测试 ============

func TestPathPattern_Simple(t *testing.T) {
	pattern := NewPathPattern().
		Vertex("v1").
		Edge(EdgeDirectionRight).
		Vertex("v2")
	expected := "(v1)->(v2)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestPathPattern_WithVariable(t *testing.T) {
	pattern := NewPathPattern().
		Variable("p").
		Vertex("v1").
		Edge(EdgeDirectionRight).
		Vertex("v2")
	expected := "p = (v1)->(v2)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestPathPattern_WithTags(t *testing.T) {
	pattern := NewPathPattern().
		VertexWithTag("v1", "player").
		EdgeWithType(EdgeDirectionRight, "follow").
		VertexWithTag("v2", "player")
	expected := "(v1:player)-[:follow]->(v2:player)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestPathPattern_Complex(t *testing.T) {
	pattern := NewPathPattern().
		Variable("p").
		VertexWithTag("v1", "player").
		EdgeWithType(EdgeDirectionRight, "follow").
		VertexWithTag("v2", "player").
		EdgeWithType(EdgeDirectionRight, "serve").
		VertexWithTag("v3", "team")
	expected := "p = (v1:player)-[:follow]->(v2:player)-[:serve]->(v3:team)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestPathPattern_WithCustomPatterns(t *testing.T) {
	vertex := NewVertexPattern("v1").Tag("player").Property("name", "Tim Duncan")
	edge := NewEdgePattern(EdgeDirectionRight).Variable("e").Type("follow").Hops(1, 3)

	pattern := NewPathPattern().
		VertexPattern(vertex).
		EdgePattern(edge).
		Vertex("v2")

	result := pattern.String()
	expected := `(v1:player {name: "Tim Duncan"})-[e:follow*1..3]->(v2)`
	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}
}

func TestPathPattern_BidirectionalEdge(t *testing.T) {
	pattern := NewPathPattern().
		Vertex("v1").
		Edge(EdgeDirectionBoth).
		Vertex("v2")
	expected := "(v1)-(v2)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestPathPattern_LeftDirectionEdge(t *testing.T) {
	pattern := NewPathPattern().
		Vertex("v1").
		EdgeWithType(EdgeDirectionLeft, "follow").
		Vertex("v2")
	expected := "(v1)<-[:follow]-(v2)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

// ============ 多点关联模式测试 ============

func TestPathPattern_EmptyEdge(t *testing.T) {
	// 测试空边：(a)-[]->(b)
	pattern := NewPathPattern().
		Vertex("a").
		Edge(EdgeDirectionRight).
		Vertex("b")
	expected := "(a)->(b)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestPathPattern_ThreeVertices_TwoEdges(t *testing.T) {
	// 测试三点两边：(a)-[]->(b)<-[]-(c)
	// 注意：这需要分两个路径来表示，或者使用更复杂的结构
	// 这里先测试线性路径：(a)-[]->(b)-[]->(c)
	pattern := NewPathPattern().
		Vertex("a").
		Edge(EdgeDirectionRight).
		Vertex("b").
		Edge(EdgeDirectionLeft).
		Vertex("c")
	expected := "(a)->(b)<-(c)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestPathPattern_AnonymousMiddleVertex(t *testing.T) {
	// 测试匿名中间顶点：(a)-[]->()<-[]-(c)
	pattern := NewPathPattern().
		Vertex("a").
		Edge(EdgeDirectionRight).
		Vertex("").
		Edge(EdgeDirectionLeft).
		Vertex("c")
	expected := "(a)->()<-(c)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestPathPattern_MultipleAnonymousVertices(t *testing.T) {
	// 测试多个匿名顶点：()-[]->()-[]-()
	pattern := NewPathPattern().
		Vertex("").
		Edge(EdgeDirectionRight).
		Vertex("").
		Edge(EdgeDirectionRight).
		Vertex("")
	expected := "()->()->()"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

// ============ Tag 模式测试 ============

func TestPathPattern_SingleTag(t *testing.T) {
	// 测试单个 Tag：(a:User)-[]->(b)
	pattern := NewPathPattern().
		VertexWithTag("a", "User").
		Edge(EdgeDirectionRight).
		Vertex("b")
	expected := "(a:User)->(b)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestPathPattern_MultipleTags(t *testing.T) {
	// 测试多个 Tag：(a:User:Admin)-[]->(b)
	pattern := NewPathPattern().
		VertexPattern(
			NewVertexPattern("a").
				Tag("User").
				Tag("Admin"),
		).
		Edge(EdgeDirectionRight).
		Vertex("b")
	expected := "(a:User:Admin)->(b)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestPathPattern_BothVerticesWithTags(t *testing.T) {
	// 测试两个顶点都有 Tag：(a:User)-[]->(b:Post)
	pattern := NewPathPattern().
		VertexWithTag("a", "User").
		Edge(EdgeDirectionRight).
		VertexWithTag("b", "Post")
	expected := "(a:User)->(b:Post)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}

func TestPathPattern_ComplexMultiTag(t *testing.T) {
	// 测试复杂的多 Tag 场景：(a:User:Admin:VIP)-[:CREATED]->(b:Post:Article)
	pattern := NewPathPattern().
		VertexPattern(
			NewVertexPattern("a").
				Tag("User").
				Tag("Admin").
				Tag("VIP"),
		).
		EdgeWithType(EdgeDirectionRight, "CREATED").
		VertexPattern(
			NewVertexPattern("b").
				Tag("Post").
				Tag("Article"),
		)
	expected := "(a:User:Admin:VIP)-[:CREATED]->(b:Post:Article)"
	if pattern.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, pattern.String())
	}
}
