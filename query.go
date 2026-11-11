package ngql

import (
	"fmt"
	"strings"
)

// ============ MATCH 查询 ============

// MatchBuilder MATCH 查询构建器
// 子句输出顺序：[OPTIONAL] MATCH pattern [WHERE ...] [WITH ...]... [UNWIND ... AS ...]... [RETURN ...]
type MatchBuilder struct {
	pattern    Pattern
	optional   bool
	whereConds []string
	returns    []string
	distinct   bool
	withs      []string // 预渲染的 WITH 子句，按添加顺序输出
	unwinds    []string // 预渲染的 UNWIND 子句，按添加顺序输出
	orderBy    []string
	limit      int
	skip       int
}

// NewMatch 创建 MATCH 查询构建器
// pattern: 匹配模式，例如 "(v:player{name:\"Tim Duncan\"})"
func NewMatch(pattern Pattern) *MatchBuilder {
	return newMatch(pattern, false)
}

// NewOptionalMatch 创建 OPTIONAL MATCH 查询构建器
// 模式无匹配时不丢弃整行，对应字段以 NULL 填充
func NewOptionalMatch(pattern Pattern) *MatchBuilder {
	return newMatch(pattern, true)
}

func newMatch(pattern Pattern, optional bool) *MatchBuilder {
	return &MatchBuilder{
		pattern:    pattern,
		optional:   optional,
		whereConds: make([]string, 0),
		returns:    make([]string, 0),
		withs:      make([]string, 0),
		unwinds:    make([]string, 0),
		orderBy:    make([]string, 0),
		limit:      -1,
		skip:       -1,
	}
}

// Where 添加 WHERE 条件
func (b *MatchBuilder) Where(condition string) *MatchBuilder {
	b.whereConds = append(b.whereConds, condition)
	return b
}

// With 添加 WITH 子句
// WITH 用于中间投影与范围过滤，可携带自身的 WHERE/ORDER BY/SKIP/LIMIT
func (b *MatchBuilder) With(w *WithBuilder) *MatchBuilder {
	b.withs = append(b.withs, w.String())
	return b
}

// Unwind 添加 UNWIND 子句，将列表表达式展开为多行
// expr: 列表表达式，例如 "names"、"[1, 2, 3]"
// alias: 展开后的变量别名
func (b *MatchBuilder) Unwind(expr, alias string) *MatchBuilder {
	b.unwinds = append(b.unwinds, fmt.Sprintf("UNWIND %s AS %s", expr, alias))
	return b
}

// Return 添加返回字段
func (b *MatchBuilder) Return(fields ...string) *MatchBuilder {
	b.returns = append(b.returns, fields...)
	return b
}

// ReturnDistinct 添加返回字段并去重（RETURN DISTINCT）
func (b *MatchBuilder) ReturnDistinct(fields ...string) *MatchBuilder {
	b.distinct = true
	b.returns = append(b.returns, fields...)
	return b
}

// OrderBy 添加排序
// field: 排序字段
// desc: 是否降序
func (b *MatchBuilder) OrderBy(field string, desc bool) *MatchBuilder {
	order := field
	if desc {
		order += " DESC"
	} else {
		order += " ASC"
	}
	b.orderBy = append(b.orderBy, order)
	return b
}

// Limit 设置返回数量限制
func (b *MatchBuilder) Limit(limit int) *MatchBuilder {
	b.limit = limit
	return b
}

// Skip 设置跳过数量
func (b *MatchBuilder) Skip(skip int) *MatchBuilder {
	b.skip = skip
	return b
}

// String 构建 nGQL 语句
func (b *MatchBuilder) String() string {
	var sb strings.Builder
	if b.optional {
		sb.WriteString("OPTIONAL ")
	}
	sb.WriteString("MATCH ")
	sb.WriteString(b.pattern.String())

	if len(b.whereConds) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(b.whereConds, " AND "))
	}

	for _, with := range b.withs {
		sb.WriteString(" ")
		sb.WriteString(with)
	}

	for _, unwind := range b.unwinds {
		sb.WriteString(" ")
		sb.WriteString(unwind)
	}

	if len(b.returns) > 0 {
		sb.WriteString(" RETURN ")
		if b.distinct {
			sb.WriteString("DISTINCT ")
		}
		sb.WriteString(strings.Join(b.returns, ", "))
	}

	if len(b.orderBy) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.Join(b.orderBy, ", "))
	}

	if b.skip > 0 {
		sb.WriteString(fmt.Sprintf(" SKIP %d", b.skip))
	}

	if b.limit >= 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", b.limit))
	}

	return sb.String()
}

// ============ WITH 子句 ============

// WithBuilder MATCH 的 WITH 子句构建器
// 输出：WITH field1[, field2...] [WHERE ...] [ORDER BY ...] [SKIP n] [LIMIT n]
type WithBuilder struct {
	fields     []string
	whereConds []string
	orderBy    []string
	limit      int
	skip       int
}

// NewWith 创建 WITH 子句构建器
func NewWith() *WithBuilder {
	return &WithBuilder{
		fields:     make([]string, 0),
		whereConds: make([]string, 0),
		orderBy:    make([]string, 0),
		limit:      -1,
		skip:       -1,
	}
}

// Field 添加字段表达式
func (b *WithBuilder) Field(expr string) *WithBuilder {
	b.fields = append(b.fields, expr)
	return b
}

// FieldAs 添加带别名的字段表达式
func (b *WithBuilder) FieldAs(expr, alias string) *WithBuilder {
	b.fields = append(b.fields, fmt.Sprintf("%s AS %s", expr, alias))
	return b
}

// Where 添加过滤条件，多次调用以 AND 连接
func (b *WithBuilder) Where(condition string) *WithBuilder {
	b.whereConds = append(b.whereConds, condition)
	return b
}

// OrderBy 添加排序
func (b *WithBuilder) OrderBy(field string, desc bool) *WithBuilder {
	order := field
	if desc {
		order += " DESC"
	} else {
		order += " ASC"
	}
	b.orderBy = append(b.orderBy, order)
	return b
}

// Skip 设置跳过数量
func (b *WithBuilder) Skip(skip int) *WithBuilder {
	b.skip = skip
	return b
}

// Limit 设置返回数量限制
func (b *WithBuilder) Limit(limit int) *WithBuilder {
	b.limit = limit
	return b
}

// String 构建 WITH 子句
func (b *WithBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("WITH ")
	sb.WriteString(strings.Join(b.fields, ", "))

	if len(b.whereConds) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(b.whereConds, " AND "))
	}

	if len(b.orderBy) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.Join(b.orderBy, ", "))
	}

	if b.skip > 0 {
		sb.WriteString(fmt.Sprintf(" SKIP %d", b.skip))
	}

	if b.limit >= 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", b.limit))
	}

	return sb.String()
}

// ============ YIELD 语句 ============

// YieldBuilder 独立 YIELD 语句构建器
// 输出：YIELD [DISTINCT] field1[, field2...]
type YieldBuilder struct {
	distinct bool
	fields   []string
}

// NewYield 创建独立 YIELD 语句构建器
func NewYield() *YieldBuilder {
	return &YieldBuilder{
		fields: make([]string, 0),
	}
}

// Distinct 设置 DISTINCT 选项
func (b *YieldBuilder) Distinct() *YieldBuilder {
	b.distinct = true
	return b
}

// Field 添加字段表达式
func (b *YieldBuilder) Field(expr string) *YieldBuilder {
	b.fields = append(b.fields, expr)
	return b
}

// FieldAs 添加带别名的字段表达式
func (b *YieldBuilder) FieldAs(expr, alias string) *YieldBuilder {
	b.fields = append(b.fields, fmt.Sprintf("%s AS %s", expr, alias))
	return b
}

// String 构建 YIELD 语句
func (b *YieldBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("YIELD ")
	if b.distinct {
		sb.WriteString("DISTINCT ")
	}
	sb.WriteString(strings.Join(b.fields, ", "))
	return sb.String()
}

// ============ GO 查询 ============

// GoBuilder GO 查询构建器
type GoBuilder struct {
	fromVids    []string
	overEdges   []string
	steps       int // 步数（StepsRange 设置时为最小步数）
	maxSteps    int // 最大步数，>0 表示变长步数 M TO N STEPS
	reversely   bool
	bidirect    bool
	whereConds  []string
	yieldFields []string
}

// NewGo 创建 GO 查询构建器
// fromVids: 起始顶点ID列表
func NewGo(fromVids ...string) *GoBuilder {
	return &GoBuilder{
		fromVids:    fromVids,
		overEdges:   make([]string, 0),
		steps:       1,
		yieldFields: make([]string, 0),
	}
}

// Over 指定要遍历的边类型
func (b *GoBuilder) Over(edges ...string) *GoBuilder {
	b.overEdges = append(b.overEdges, edges...)
	return b
}

// Steps 设置遍历步数
func (b *GoBuilder) Steps(steps int) *GoBuilder {
	b.steps = steps
	return b
}

// StepsRange 设置变长遍历步数范围，输出 M TO N STEPS
func (b *GoBuilder) StepsRange(min, max int) *GoBuilder {
	b.steps = min
	b.maxSteps = max
	return b
}

// Reversely 反向遍历
func (b *GoBuilder) Reversely() *GoBuilder {
	b.reversely = true
	return b
}

// Bidirect 双向遍历
func (b *GoBuilder) Bidirect() *GoBuilder {
	b.bidirect = true
	return b
}

// Where 添加过滤条件，多次调用以 AND 连接
func (b *GoBuilder) Where(condition string) *GoBuilder {
	b.whereConds = append(b.whereConds, condition)
	return b
}

// Yield 添加返回字段
func (b *GoBuilder) Yield(fields ...string) *GoBuilder {
	b.yieldFields = append(b.yieldFields, fields...)
	return b
}

// String 构建 nGQL 语句
func (b *GoBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("GO ")

	if b.maxSteps > 0 {
		sb.WriteString(fmt.Sprintf("%d TO %d STEPS ", b.steps, b.maxSteps))
	} else if b.steps > 1 {
		sb.WriteString(fmt.Sprintf("%d STEPS ", b.steps))
	}

	sb.WriteString("FROM ")
	vidStrs := make([]string, 0, len(b.fromVids))
	for _, vid := range b.fromVids {
		vidStrs = append(vidStrs, quote(vid))
	}
	sb.WriteString(strings.Join(vidStrs, ", "))

	sb.WriteString(" OVER ")
	if len(b.overEdges) == 0 {
		sb.WriteString("*")
	} else {
		edgeStrs := make([]string, 0, len(b.overEdges))
		for _, edge := range b.overEdges {
			edgeStrs = append(edgeStrs, backquoteIfNeed(edge))
		}
		sb.WriteString(strings.Join(edgeStrs, ", "))
	}

	if b.reversely {
		sb.WriteString(" REVERSELY")
	}

	if b.bidirect {
		sb.WriteString(" BIDIRECT")
	}

	if len(b.whereConds) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(b.whereConds, " AND "))
	}

	if len(b.yieldFields) > 0 {
		sb.WriteString(" YIELD ")
		sb.WriteString(strings.Join(b.yieldFields, ", "))
	}

	return sb.String()
}

// ============ FETCH 查询 ============

// FetchBuilder FETCH 查询构建器
type FetchBuilder struct {
	propOn      string // "VERTEX" 或边类型
	tag         string // 顶点标签（仅用于 VERTEX）
	vids        []string
	edges       []fetchEdge
	yieldFields []string
}

type fetchEdge struct {
	srcVid string
	dstVid string
	rank   *int64
}

// NewFetchVertex 创建获取顶点属性的构建器
// tag: 顶点标签
// vids: 顶点ID列表
func NewFetchVertex(tag string, vids ...string) *FetchBuilder {
	return &FetchBuilder{
		propOn:      "VERTEX",
		tag:         tag,
		vids:        vids,
		yieldFields: make([]string, 0),
	}
}

// NewFetchEdge 创建获取边属性的构建器
// edgeType: 边类型
func NewFetchEdge(edgeType string) *FetchBuilder {
	return &FetchBuilder{
		propOn:      edgeType,
		edges:       make([]fetchEdge, 0),
		yieldFields: make([]string, 0),
	}
}

// Edge 添加要获取的边（仅用于 FetchEdge）
func (b *FetchBuilder) Edge(srcVid, dstVid string) *FetchBuilder {
	b.edges = append(b.edges, fetchEdge{
		srcVid: srcVid,
		dstVid: dstVid,
	})
	return b
}

// EdgeWithRank 添加要获取的带 rank 的边（仅用于 FetchEdge）
func (b *FetchBuilder) EdgeWithRank(srcVid, dstVid string, rank int64) *FetchBuilder {
	b.edges = append(b.edges, fetchEdge{
		srcVid: srcVid,
		dstVid: dstVid,
		rank:   &rank,
	})
	return b
}

// Yield 添加返回字段
func (b *FetchBuilder) Yield(fields ...string) *FetchBuilder {
	b.yieldFields = append(b.yieldFields, fields...)
	return b
}

// String 构建 nGQL 语句
func (b *FetchBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("FETCH PROP ON ")

	if b.propOn == "VERTEX" {
		sb.WriteString(backquoteIfNeed(b.tag))
		sb.WriteString(" ")
		vidStrs := make([]string, 0, len(b.vids))
		for _, vid := range b.vids {
			vidStrs = append(vidStrs, quote(vid))
		}
		sb.WriteString(strings.Join(vidStrs, ", "))
	} else {
		sb.WriteString(backquoteIfNeed(b.propOn))
		sb.WriteString(" ")
		edgeStrs := make([]string, 0, len(b.edges))
		for _, e := range b.edges {
			rankStr := ""
			if e.rank != nil {
				rankStr = fmt.Sprintf("@%d", *e.rank)
			}
			edgeStrs = append(edgeStrs, fmt.Sprintf("%s -> %s%s", quote(e.srcVid), quote(e.dstVid), rankStr))
		}
		sb.WriteString(strings.Join(edgeStrs, ", "))
	}

	if len(b.yieldFields) > 0 {
		sb.WriteString(" YIELD ")
		sb.WriteString(strings.Join(b.yieldFields, ", "))
	}

	return sb.String()
}

// ============ LOOKUP 查询 ============

// LookupBuilder LOOKUP 查询构建器
type LookupBuilder struct {
	onType      string   // 顶点标签或边类型
	whereConds  []string // WHERE 条件列表
	yieldFields []string // YIELD 字段列表
	distinct    bool     // 是否使用 DISTINCT
	limit       int      // LIMIT 数量
	skip        int      // SKIP 数量
}

// NewLookup 创建 LOOKUP 查询构建器
// onType: 顶点标签或边类型，例如 "player" 或 "follow"
func NewLookup(onType string) *LookupBuilder {
	return &LookupBuilder{
		onType:      onType,
		whereConds:  make([]string, 0),
		yieldFields: make([]string, 0),
		limit:       -1,
		skip:        -1,
	}
}

// Where 添加 WHERE 条件
func (b *LookupBuilder) Where(condition string) *LookupBuilder {
	b.whereConds = append(b.whereConds, condition)
	return b
}

// Yield 添加返回字段
func (b *LookupBuilder) Yield(fields ...string) *LookupBuilder {
	b.yieldFields = append(b.yieldFields, fields...)
	return b
}

// Distinct 设置 DISTINCT 选项
func (b *LookupBuilder) Distinct() *LookupBuilder {
	b.distinct = true
	return b
}

// Limit 设置返回数量限制
func (b *LookupBuilder) Limit(limit int) *LookupBuilder {
	b.limit = limit
	return b
}

// Skip 设置跳过数量
func (b *LookupBuilder) Skip(skip int) *LookupBuilder {
	b.skip = skip
	return b
}

// String 构建 nGQL 语句
func (b *LookupBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("LOOKUP ON ")
	sb.WriteString(backquoteIfNeed(b.onType))

	if len(b.whereConds) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(b.whereConds, " AND "))
	}

	if len(b.yieldFields) > 0 {
		sb.WriteString(" YIELD ")
		if b.distinct {
			sb.WriteString("DISTINCT ")
		}
		sb.WriteString(strings.Join(b.yieldFields, ", "))
	}

	if b.skip > 0 {
		sb.WriteString(fmt.Sprintf(" | SKIP %d", b.skip))
	}

	if b.limit >= 0 {
		sb.WriteString(fmt.Sprintf(" | LIMIT %d", b.limit))
	}

	return sb.String()
}

// ============ GROUP BY 操作 ============

// GroupByBuilder GROUP BY 构建器（管道操作）
type GroupByBuilder struct {
	prevQuery   string   // 前置查询语句
	groupByVars []string // GROUP BY 变量列表
	yieldFields []string // YIELD 字段列表（包含聚合函数）
	orderBy     []string // ORDER BY 子句
	limit       int      // LIMIT 数量
	skip        int      // SKIP 数量
}

// NewGroupBy 创建 GROUP BY 构建器
// prevQuery: 前置查询语句
// groupByVars: GROUP BY 的变量，例如 "$-.age"
func NewGroupBy(prevQuery string, groupByVars ...string) *GroupByBuilder {
	return &GroupByBuilder{
		prevQuery:   prevQuery,
		groupByVars: groupByVars,
		yieldFields: make([]string, 0),
		orderBy:     make([]string, 0),
		limit:       -1,
		skip:        -1,
	}
}

// Yield 添加返回字段（通常包含聚合函数）
// 例如: "$-.age", "count(*) AS count", "avg($-.score) AS avg_score"
func (b *GroupByBuilder) Yield(fields ...string) *GroupByBuilder {
	b.yieldFields = append(b.yieldFields, fields...)
	return b
}

// OrderBy 添加排序
func (b *GroupByBuilder) OrderBy(field string, desc bool) *GroupByBuilder {
	order := field
	if desc {
		order += " DESC"
	} else {
		order += " ASC"
	}
	b.orderBy = append(b.orderBy, order)
	return b
}

// Limit 设置返回数量限制
func (b *GroupByBuilder) Limit(limit int) *GroupByBuilder {
	b.limit = limit
	return b
}

// Skip 设置跳过数量
func (b *GroupByBuilder) Skip(skip int) *GroupByBuilder {
	b.skip = skip
	return b
}

// String 构建 nGQL 语句
func (b *GroupByBuilder) String() string {
	var sb strings.Builder
	sb.WriteString(b.prevQuery)

	if len(b.groupByVars) > 0 {
		sb.WriteString(" | GROUP BY ")
		sb.WriteString(strings.Join(b.groupByVars, ", "))
		if len(b.yieldFields) > 0 {
			sb.WriteString(" YIELD ")
			sb.WriteString(strings.Join(b.yieldFields, ", "))
		}
	} else if len(b.yieldFields) > 0 {
		// 未指定分组变量时不输出 GROUP BY，避免生成非法语句
		sb.WriteString(" | YIELD ")
		sb.WriteString(strings.Join(b.yieldFields, ", "))
	}

	if len(b.orderBy) > 0 {
		sb.WriteString(" | ORDER BY ")
		sb.WriteString(strings.Join(b.orderBy, ", "))
	}

	if b.skip > 0 {
		sb.WriteString(fmt.Sprintf(" | SKIP %d", b.skip))
	}

	if b.limit >= 0 {
		sb.WriteString(fmt.Sprintf(" | LIMIT %d", b.limit))
	}

	return sb.String()
}

// ============ FIND PATH 路径查询 ============

// FIND PATH 路径类型常量
const (
	PathTypeShortest = "SHORTEST" // 最短路径
	PathTypeAll      = "ALL"      // 所有路径
	PathTypeNoLoop   = "NOLOOP"   // 不含环的路径
)

// FindPathBuilder FIND PATH 查询构建器
type FindPathBuilder struct {
	pathType    string   // SHORTEST, ALL, NOLOOP（见 PathTypeShortest 等常量）
	withProp    bool     // 是否包含属性
	fromVids    []string // 起始顶点列表
	toVids      []string // 目标顶点列表
	overEdges   []string // 边类型列表
	reversely   bool     // 反向遍历
	bidirect    bool     // 双向遍历
	whereConds  []string // WHERE 条件列表
	upToSteps   int      // 最大步数
	yieldFields []string // YIELD 字段
	orderBy     []string // ORDER BY 子句
	limit       int      // LIMIT 数量
}

// NewFindPath 创建 FIND PATH 查询构建器
// pathType: 路径类型 "SHORTEST", "ALL", "NOLOOP"
func NewFindPath(pathType string) *FindPathBuilder {
	return &FindPathBuilder{
		pathType:    pathType,
		fromVids:    make([]string, 0),
		toVids:      make([]string, 0),
		overEdges:   make([]string, 0),
		yieldFields: make([]string, 0),
		orderBy:     make([]string, 0),
		upToSteps:   -1,
		limit:       -1,
	}
}

// WithProp 设置包含属性
func (b *FindPathBuilder) WithProp() *FindPathBuilder {
	b.withProp = true
	return b
}

// From 设置起始顶点
func (b *FindPathBuilder) From(vids ...string) *FindPathBuilder {
	b.fromVids = append(b.fromVids, vids...)
	return b
}

// To 设置目标顶点
func (b *FindPathBuilder) To(vids ...string) *FindPathBuilder {
	b.toVids = append(b.toVids, vids...)
	return b
}

// Over 指定要遍历的边类型
func (b *FindPathBuilder) Over(edges ...string) *FindPathBuilder {
	b.overEdges = append(b.overEdges, edges...)
	return b
}

// Reversely 反向遍历
func (b *FindPathBuilder) Reversely() *FindPathBuilder {
	b.reversely = true
	return b
}

// Bidirect 双向遍历
func (b *FindPathBuilder) Bidirect() *FindPathBuilder {
	b.bidirect = true
	return b
}

// Where 添加过滤条件，多次调用以 AND 连接
func (b *FindPathBuilder) Where(condition string) *FindPathBuilder {
	b.whereConds = append(b.whereConds, condition)
	return b
}

// UpTo 设置最大步数
func (b *FindPathBuilder) UpTo(steps int) *FindPathBuilder {
	b.upToSteps = steps
	return b
}

// Yield 添加返回字段
func (b *FindPathBuilder) Yield(fields ...string) *FindPathBuilder {
	b.yieldFields = append(b.yieldFields, fields...)
	return b
}

// OrderBy 添加排序
func (b *FindPathBuilder) OrderBy(field string, desc bool) *FindPathBuilder {
	order := field
	if desc {
		order += " DESC"
	} else {
		order += " ASC"
	}
	b.orderBy = append(b.orderBy, order)
	return b
}

// Limit 设置返回数量限制
func (b *FindPathBuilder) Limit(limit int) *FindPathBuilder {
	b.limit = limit
	return b
}

// String 构建 nGQL 语句
func (b *FindPathBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("FIND ")
	sb.WriteString(b.pathType)
	sb.WriteString(" PATH")

	if b.withProp {
		sb.WriteString(" WITH PROP")
	}

	sb.WriteString(" FROM ")
	fromStrs := make([]string, 0, len(b.fromVids))
	for _, vid := range b.fromVids {
		fromStrs = append(fromStrs, quote(vid))
	}
	sb.WriteString(strings.Join(fromStrs, ", "))

	sb.WriteString(" TO ")
	toStrs := make([]string, 0, len(b.toVids))
	for _, vid := range b.toVids {
		toStrs = append(toStrs, quote(vid))
	}
	sb.WriteString(strings.Join(toStrs, ", "))

	sb.WriteString(" OVER ")
	if len(b.overEdges) == 0 {
		sb.WriteString("*")
	} else {
		edgeStrs := make([]string, 0, len(b.overEdges))
		for _, edge := range b.overEdges {
			edgeStrs = append(edgeStrs, backquoteIfNeed(edge))
		}
		sb.WriteString(strings.Join(edgeStrs, ", "))
	}

	if b.reversely {
		sb.WriteString(" REVERSELY")
	}

	if b.bidirect {
		sb.WriteString(" BIDIRECT")
	}

	if len(b.whereConds) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(b.whereConds, " AND "))
	}

	if b.upToSteps > 0 {
		sb.WriteString(fmt.Sprintf(" UPTO %d STEPS", b.upToSteps))
	}

	if len(b.yieldFields) > 0 {
		sb.WriteString(" YIELD ")
		sb.WriteString(strings.Join(b.yieldFields, ", "))
	}

	if len(b.orderBy) > 0 {
		sb.WriteString(" | ORDER BY ")
		sb.WriteString(strings.Join(b.orderBy, ", "))
	}

	if b.limit >= 0 {
		sb.WriteString(fmt.Sprintf(" | LIMIT %d", b.limit))
	}

	return sb.String()
}

// ============ GET SUBGRAPH 子图查询 ============

// GetSubgraphBuilder GET SUBGRAPH 查询构建器
type GetSubgraphBuilder struct {
	withProp    bool     // 是否包含属性
	steps       int      // 遍历步数
	fromVids    []string // 起始顶点列表
	inEdges     []string // IN 方向的边类型
	outEdges    []string // OUT 方向的边类型
	bothEdges   []string // BOTH 方向的边类型
	whereConds  []string // WHERE 条件列表
	yieldFields []string // YIELD 字段列表
}

// NewGetSubgraph 创建 GET SUBGRAPH 查询构建器
// steps: 遍历步数，默认为1
func NewGetSubgraph(steps int) *GetSubgraphBuilder {
	if steps <= 0 {
		steps = 1
	}
	return &GetSubgraphBuilder{
		steps:       steps,
		fromVids:    make([]string, 0),
		inEdges:     make([]string, 0),
		outEdges:    make([]string, 0),
		bothEdges:   make([]string, 0),
		whereConds:  make([]string, 0),
		yieldFields: make([]string, 0),
	}
}

// WithProp 设置包含属性
func (b *GetSubgraphBuilder) WithProp() *GetSubgraphBuilder {
	b.withProp = true
	return b
}

// From 设置起始顶点
func (b *GetSubgraphBuilder) From(vids ...string) *GetSubgraphBuilder {
	b.fromVids = append(b.fromVids, vids...)
	return b
}

// In 指定 IN 方向的边类型
func (b *GetSubgraphBuilder) In(edges ...string) *GetSubgraphBuilder {
	b.inEdges = append(b.inEdges, edges...)
	return b
}

// Out 指定 OUT 方向的边类型
func (b *GetSubgraphBuilder) Out(edges ...string) *GetSubgraphBuilder {
	b.outEdges = append(b.outEdges, edges...)
	return b
}

// Both 指定 BOTH 方向的边类型
func (b *GetSubgraphBuilder) Both(edges ...string) *GetSubgraphBuilder {
	b.bothEdges = append(b.bothEdges, edges...)
	return b
}

// Where 添加 WHERE 条件
func (b *GetSubgraphBuilder) Where(condition string) *GetSubgraphBuilder {
	b.whereConds = append(b.whereConds, condition)
	return b
}

// Yield 添加返回字段
func (b *GetSubgraphBuilder) Yield(fields ...string) *GetSubgraphBuilder {
	b.yieldFields = append(b.yieldFields, fields...)
	return b
}

// String 构建 nGQL 语句
func (b *GetSubgraphBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("GET SUBGRAPH")

	if b.withProp {
		sb.WriteString(" WITH PROP")
	}

	if b.steps > 1 {
		sb.WriteString(fmt.Sprintf(" %d STEPS", b.steps))
	}

	sb.WriteString(" FROM ")
	vidStrs := make([]string, 0, len(b.fromVids))
	for _, vid := range b.fromVids {
		vidStrs = append(vidStrs, quote(vid))
	}
	sb.WriteString(strings.Join(vidStrs, ", "))

	// 处理边类型和方向
	if len(b.inEdges) > 0 {
		sb.WriteString(" IN ")
		sb.WriteString(joinBackquoted(b.inEdges))
	}

	if len(b.outEdges) > 0 {
		sb.WriteString(" OUT ")
		sb.WriteString(joinBackquoted(b.outEdges))
	}

	if len(b.bothEdges) > 0 {
		sb.WriteString(" BOTH ")
		sb.WriteString(joinBackquoted(b.bothEdges))
	}

	if len(b.whereConds) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(b.whereConds, " AND "))
	}

	if len(b.yieldFields) > 0 {
		sb.WriteString(" YIELD ")
		sb.WriteString(strings.Join(b.yieldFields, ", "))
	}

	return sb.String()
}
