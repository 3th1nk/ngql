package ngql

import "fmt"

// 示例：基本顶点模式
func ExampleVertexPattern() {
	// 匿名顶点
	v1 := NewVertexPattern("")
	fmt.Println(v1.String())

	// 带变量的顶点
	v2 := NewVertexPattern("v")
	fmt.Println(v2.String())

	// 带 Tag 的顶点
	v3 := NewVertexPattern("v").Tag("player")
	fmt.Println(v3.String())

	// 带多个 Tag 的顶点
	v4 := NewVertexPattern("v").Tag("player").Tag("person")
	fmt.Println(v4.String())

	// 带属性过滤的顶点
	v5 := NewVertexPattern("v").Tag("player").Property("name", "Tim Duncan")
	fmt.Println(v5.String())

	// Output:
	// ()
	// (v)
	// (v:player)
	// (v:player:person)
	// (v:player {name: "Tim Duncan"})
}

// 示例：基本边模式
func ExampleEdgePattern() {
	// 右向边
	e1 := NewEdgePattern(EdgeDirectionRight)
	fmt.Println(e1.String())

	// 左向边
	e2 := NewEdgePattern(EdgeDirectionLeft)
	fmt.Println(e2.String())

	// 双向边
	e3 := NewEdgePattern(EdgeDirectionBoth)
	fmt.Println(e3.String())

	// 带类型的边
	e4 := NewEdgePattern(EdgeDirectionRight).Type("follow")
	fmt.Println(e4.String())

	// 带变量和类型的边
	e5 := NewEdgePattern(EdgeDirectionRight).Variable("e").Type("follow")
	fmt.Println(e5.String())

	// 带多个类型的边
	e6 := NewEdgePattern(EdgeDirectionRight).Type("follow").Type("like")
	fmt.Println(e6.String())

	// Output:
	// ->
	// <-
	// -
	// -[:follow]->
	// -[e:follow]->
	// -[:follow|like]->
}

// 示例：变长边模式
func ExampleEdgePattern_variableLength() {
	// 固定步数
	e1 := NewEdgePattern(EdgeDirectionRight).Type("follow").Hops(2, 2)
	fmt.Println(e1.String())

	// 范围步数
	e2 := NewEdgePattern(EdgeDirectionRight).Type("follow").Hops(1, 3)
	fmt.Println(e2.String())

	// 最小步数
	e3 := NewEdgePattern(EdgeDirectionRight).Type("follow").Hops(1, 0)
	fmt.Println(e3.String())

	// 最大步数
	e4 := NewEdgePattern(EdgeDirectionRight).Type("follow").Hops(0, 3)
	fmt.Println(e4.String())

	// Output:
	// -[:follow*2]->
	// -[:follow*1..3]->
	// -[:follow*1..]->
	// -[:follow*..3]->
}

// 示例：路径模式
func ExamplePathPattern() {
	// 简单路径
	p1 := NewPathPattern().
		Vertex("v1").
		Edge(EdgeDirectionRight).
		Vertex("v2")
	fmt.Println(p1.String())

	// 带变量的路径
	p2 := NewPathPattern().
		Variable("p").
		Vertex("v1").
		Edge(EdgeDirectionRight).
		Vertex("v2")
	fmt.Println(p2.String())

	// 带 Tag 的路径
	p3 := NewPathPattern().
		VertexWithTag("v1", "player").
		EdgeWithType(EdgeDirectionRight, "follow").
		VertexWithTag("v2", "player")
	fmt.Println(p3.String())

	// 复杂路径
	p4 := NewPathPattern().
		Variable("p").
		VertexWithTag("v1", "player").
		EdgeWithType(EdgeDirectionRight, "follow").
		VertexWithTag("v2", "player").
		EdgeWithType(EdgeDirectionRight, "serve").
		VertexWithTag("v3", "team")
	fmt.Println(p4.String())

	// Output:
	// (v1)->(v2)
	// p = (v1)->(v2)
	// (v1:player)-[:follow]->(v2:player)
	// p = (v1:player)-[:follow]->(v2:player)-[:serve]->(v3:team)
}

// 示例：在 MATCH 查询中使用路径模式
func ExamplePathPattern_withMatch() {
	// 创建路径模式
	pattern := NewPathPattern().
		VertexWithTag("v1", "player").
		EdgeWithType(EdgeDirectionRight, "follow").
		VertexWithTag("v2", "player")

	// 在 MATCH 查询中使用
	query := NewMatch(pattern).
		Where("v1.player.age > 30").
		Return("v1.player.name", "v2.player.name").
		String()

	fmt.Println(query)

	// Output:
	// MATCH (v1:player)-[:follow]->(v2:player) WHERE v1.player.age > 30 RETURN v1.player.name, v2.player.name
}

// 示例：复杂的图模式组合
func ExamplePathPattern_complex() {
	// 创建自定义顶点模式
	vertex1 := NewVertexPattern("v1").
		Tag("player").
		Property("name", "Tim Duncan")

	// 创建自定义边模式（变长边）
	edge := NewEdgePattern(EdgeDirectionRight).
		Variable("e").
		Type("follow").
		Hops(1, 3)

	// 创建路径模式
	path := NewPathPattern().
		Variable("p").
		VertexPattern(vertex1).
		EdgePattern(edge).
		Vertex("v2")

	fmt.Println(path.String())

	// Output:
	// p = (v1:player {name: "Tim Duncan"})-[e:follow*1..3]->(v2)
}

// ============ 多点关联模式示例 ============

// ExamplePathPattern_emptyEdge 演示空边模式
// 空边表示不关心边的类型，匹配任意边
func ExamplePathPattern_emptyEdge() {
	// (a)-[]->(b) - 从 a 到 b 的任意边
	pattern := NewPathPattern().
		Vertex("a").
		Edge(EdgeDirectionRight).
		Vertex("b")

	fmt.Println(pattern.String())
	// Output: (a)->(b)
}

// ExamplePathPattern_threeVertices 演示三点关联模式
// 表示 a 指向 b，c 也指向 b
func ExamplePathPattern_threeVertices() {
	// (a)-[]->(b)<-[]-(c) - a 和 c 都指向 b
	pattern := NewPathPattern().
		Vertex("a").
		Edge(EdgeDirectionRight).
		Vertex("b").
		Edge(EdgeDirectionLeft).
		Vertex("c")

	fmt.Println(pattern.String())
	// Output: (a)->(b)<-(c)
}

// ExamplePathPattern_anonymousMiddle 演示匿名中间顶点
// 用于查找通过某个中间节点连接的两个顶点
func ExamplePathPattern_anonymousMiddle() {
	// (a)-[]->()<-[]-(c) - a 和 c 通过一个匿名顶点连接
	pattern := NewPathPattern().
		Vertex("a").
		Edge(EdgeDirectionRight).
		Vertex(""). // 匿名顶点
		Edge(EdgeDirectionLeft).
		Vertex("c")

	fmt.Println(pattern.String())
	// Output: (a)->()<-(c)
}

// ExamplePathPattern_multipleAnonymous 演示多个匿名顶点
// 用于查找路径，不关心中间节点的标识
func ExamplePathPattern_multipleAnonymous() {
	// ()-[]->()-[]->() - 三个匿名顶点的路径
	pattern := NewPathPattern().
		Vertex("").
		Edge(EdgeDirectionRight).
		Vertex("").
		Edge(EdgeDirectionRight).
		Vertex("")

	fmt.Println(pattern.String())
	// Output: ()->()->()
}

// ============ Tag 模式示例 ============

// ExamplePathPattern_singleTag 演示单个 Tag
func ExamplePathPattern_singleTag() {
	// (a:User)-[]->(b) - a 是 User 类型
	pattern := NewPathPattern().
		VertexWithTag("a", "User").
		Edge(EdgeDirectionRight).
		Vertex("b")

	fmt.Println(pattern.String())
	// Output: (a:User)->(b)
}

// ExamplePathPattern_multipleTags 演示多个 Tag
// 多个 Tag 表示顶点同时具有多个标签
func ExamplePathPattern_multipleTags() {
	// (a:User:Admin)-[]->(b) - a 同时是 User 和 Admin
	pattern := NewPathPattern().
		VertexPattern(
			NewVertexPattern("a").
				Tag("User").
				Tag("Admin"),
		).
		Edge(EdgeDirectionRight).
		Vertex("b")

	fmt.Println(pattern.String())
	// Output: (a:User:Admin)->(b)
}

// ExamplePathPattern_bothWithTags 演示两端都有 Tag
func ExamplePathPattern_bothWithTags() {
	// (a:User)-[]->(b:Post) - 用户指向帖子
	pattern := NewPathPattern().
		VertexWithTag("a", "User").
		Edge(EdgeDirectionRight).
		VertexWithTag("b", "Post")

	fmt.Println(pattern.String())
	// Output: (a:User)->(b:Post)
}

// ExamplePathPattern_complexMultiTag 演示复杂的多 Tag 场景
func ExamplePathPattern_complexMultiTag() {
	// (a:User:Admin:VIP)-[:CREATED]->(b:Post:Article)
	// 具有多个角色的用户创建了一个帖子/文章
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

	fmt.Println(pattern.String())
	// Output: (a:User:Admin:VIP)-[:CREATED]->(b:Post:Article)
}

// ============ 实际应用场景 ============

// ExamplePathPattern_socialNetwork 演示社交网络查询
func ExamplePathPattern_socialNetwork() {
	// 查找用户的朋友的朋友（二度好友）
	pattern := NewPathPattern().
		VertexWithTag("me", "User").
		EdgeWithType(EdgeDirectionRight, "FRIEND").
		VertexWithTag("friend", "User").
		EdgeWithType(EdgeDirectionRight, "FRIEND").
		VertexWithTag("friend_of_friend", "User")

	query := NewMatch(pattern).
		Where("me.User.name == \"Alice\"").
		Return("friend_of_friend.User.name").
		String()

	fmt.Println(query)
	// Output: MATCH (me:User)-[:FRIEND]->(friend:User)-[:FRIEND]->(friend_of_friend:User) WHERE me.User.name == "Alice" RETURN friend_of_friend.User.name
}

// ExamplePathPattern_contentModeration 演示内容审核场景
func ExamplePathPattern_contentModeration() {
	// 查找管理员创建的所有帖子
	pattern := NewPathPattern().
		VertexPattern(
			NewVertexPattern("admin").
				Tag("User").
				Tag("Admin"),
		).
		EdgeWithType(EdgeDirectionRight, "CREATED").
		VertexWithTag("post", "Post")

	query := NewMatch(pattern).
		Return("admin.User.name", "post.Post.title", "post.Post.created_at").
		OrderBy("post.Post.created_at", true).
		String()

	fmt.Println(query)
	// Output: MATCH (admin:User:Admin)-[:CREATED]->(post:Post) RETURN admin.User.name, post.Post.title, post.Post.created_at ORDER BY post.Post.created_at DESC
}

// ExamplePathPattern_recommendation 演示推荐系统场景
func ExamplePathPattern_recommendation() {
	// 查找用户喜欢的帖子的作者创建的其他帖子
	pattern := NewPathPattern().
		VertexWithTag("user", "User").
		EdgeWithType(EdgeDirectionRight, "LIKED").
		VertexWithTag("liked_post", "Post").
		EdgeWithType(EdgeDirectionLeft, "CREATED").
		VertexWithTag("author", "User").
		EdgeWithType(EdgeDirectionRight, "CREATED").
		VertexWithTag("other_post", "Post")

	query := NewMatch(pattern).
		Where("user.User.id == \"user123\" AND other_post.Post.id != liked_post.Post.id").
		Return("other_post.Post.title", "author.User.name").
		Limit(10).
		String()

	fmt.Println(query)
	// Output: MATCH (user:User)-[:LIKED]->(liked_post:Post)<-[:CREATED]-(author:User)-[:CREATED]->(other_post:Post) WHERE user.User.id == "user123" AND other_post.Post.id != liked_post.Post.id RETURN other_post.Post.title, author.User.name LIMIT 10
}

// ExamplePathPattern_anonymousIntermediate 演示匿名中间节点的实际应用
func ExamplePathPattern_anonymousIntermediate() {
	// 查找通过任意中间节点连接的两个用户
	// 例如：用户 A 和用户 B 都关注了同一个话题/群组
	pattern := NewPathPattern().
		VertexWithTag("userA", "User").
		EdgeWithType(EdgeDirectionRight, "FOLLOWS").
		Vertex(""). // 匿名中间节点（可能是话题、群组等）
		EdgeWithType(EdgeDirectionLeft, "FOLLOWS").
		VertexWithTag("userB", "User")

	query := NewMatch(pattern).
		Where("userA.User.id == \"alice\" AND userB.User.id != \"alice\"").
		Return("userB.User.name").
		String()

	fmt.Println(query)
	// Output: MATCH (userA:User)-[:FOLLOWS]->()<-[:FOLLOWS]-(userB:User) WHERE userA.User.id == "alice" AND userB.User.id != "alice" RETURN userB.User.name
}
