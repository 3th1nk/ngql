package ngql

import (
	"fmt"
	"strings"
)

// ============ UPSERT VERTEX ============

// UpsertVertexBuilder UPSERT VERTEX 构建器
type UpsertVertexBuilder struct {
	tag         string
	vid         string
	setProps    kvPairs // 属性赋值，按调用顺序输出，同名属性后设置覆盖先设置
	setExprs    []string
	whenConds   []string
	yieldFields []string
}

// NewUpsertVertex 创建 UPSERT VERTEX 构建器
// tag: 顶点标签
// vid: 顶点 ID
func NewUpsertVertex(tag, vid string) *UpsertVertexBuilder {
	return &UpsertVertexBuilder{
		tag:         tag,
		vid:         vid,
		setProps:    make(kvPairs, 0),
		setExprs:    make([]string, 0),
		whenConds:   make([]string, 0),
		yieldFields: make([]string, 0),
	}
}

// Set 设置属性值，同名属性后设置覆盖先设置
func (b *UpsertVertexBuilder) Set(prop string, value interface{}) *UpsertVertexBuilder {
	b.setProps.set(prop, value)
	return b
}

// SetExpr 设置属性表达式
func (b *UpsertVertexBuilder) SetExpr(expr string) *UpsertVertexBuilder {
	b.setExprs = append(b.setExprs, expr)
	return b
}

// When 添加条件，多次调用以 AND 连接
func (b *UpsertVertexBuilder) When(condition string) *UpsertVertexBuilder {
	b.whenConds = append(b.whenConds, condition)
	return b
}

// Yield 添加返回字段
func (b *UpsertVertexBuilder) Yield(fields ...string) *UpsertVertexBuilder {
	b.yieldFields = append(b.yieldFields, fields...)
	return b
}

// String 构建 nGQL 语句
func (b *UpsertVertexBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("UPSERT VERTEX ON ")
	sb.WriteString(backquoteIfNeed(b.tag))
	sb.WriteString(" ")
	sb.WriteString(quote(b.vid))

	// SET 子句
	if len(b.setProps) > 0 || len(b.setExprs) > 0 {
		sb.WriteString(" SET ")
		setParts := make([]string, 0, len(b.setProps)+len(b.setExprs))

		// 添加属性赋值
		for _, kv := range b.setProps {
			setParts = append(setParts, fmt.Sprintf("%s = %s", backquoteIfNeed(kv.key), formatValue(kv.value)))
		}

		// 添加表达式
		setParts = append(setParts, b.setExprs...)

		sb.WriteString(strings.Join(setParts, ", "))
	}

	// WHEN 子句
	if len(b.whenConds) > 0 {
		sb.WriteString(" WHEN ")
		sb.WriteString(strings.Join(b.whenConds, " AND "))
	}

	// YIELD 子句
	if len(b.yieldFields) > 0 {
		sb.WriteString(" YIELD ")
		sb.WriteString(strings.Join(b.yieldFields, ", "))
	}

	return sb.String()
}

// ============ UPSERT EDGE ============

// UpsertEdgeBuilder UPSERT EDGE 构建器
type UpsertEdgeBuilder struct {
	edgeType    string
	srcVid      string
	dstVid      string
	rank        *int64
	setProps    kvPairs // 属性赋值，按调用顺序输出，同名属性后设置覆盖先设置
	setExprs    []string
	whenConds   []string
	yieldFields []string
}

// NewUpsertEdge 创建 UPSERT EDGE 构建器
// edgeType: 边类型
// srcVid: 起始顶点 ID
// dstVid: 目标顶点 ID
func NewUpsertEdge(edgeType, srcVid, dstVid string) *UpsertEdgeBuilder {
	return &UpsertEdgeBuilder{
		edgeType:    edgeType,
		srcVid:      srcVid,
		dstVid:      dstVid,
		setProps:    make(kvPairs, 0),
		setExprs:    make([]string, 0),
		whenConds:   make([]string, 0),
		yieldFields: make([]string, 0),
	}
}

// WithRank 设置边的 rank
func (b *UpsertEdgeBuilder) WithRank(rank int64) *UpsertEdgeBuilder {
	b.rank = &rank
	return b
}

// Set 设置属性值，同名属性后设置覆盖先设置
func (b *UpsertEdgeBuilder) Set(prop string, value interface{}) *UpsertEdgeBuilder {
	b.setProps.set(prop, value)
	return b
}

// SetExpr 设置属性表达式
func (b *UpsertEdgeBuilder) SetExpr(expr string) *UpsertEdgeBuilder {
	b.setExprs = append(b.setExprs, expr)
	return b
}

// When 添加条件，多次调用以 AND 连接
func (b *UpsertEdgeBuilder) When(condition string) *UpsertEdgeBuilder {
	b.whenConds = append(b.whenConds, condition)
	return b
}

// Yield 添加返回字段
func (b *UpsertEdgeBuilder) Yield(fields ...string) *UpsertEdgeBuilder {
	b.yieldFields = append(b.yieldFields, fields...)
	return b
}

// String 构建 nGQL 语句
func (b *UpsertEdgeBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("UPSERT EDGE ON ")
	sb.WriteString(backquoteIfNeed(b.edgeType))
	sb.WriteString(fmt.Sprintf(" %s -> %s", quote(b.srcVid), quote(b.dstVid)))

	// Rank
	if b.rank != nil {
		sb.WriteString(fmt.Sprintf("@%d", *b.rank))
	}

	// SET 子句
	if len(b.setProps) > 0 || len(b.setExprs) > 0 {
		sb.WriteString(" SET ")
		setParts := make([]string, 0, len(b.setProps)+len(b.setExprs))

		// 添加属性赋值
		for _, kv := range b.setProps {
			setParts = append(setParts, fmt.Sprintf("%s = %s", backquoteIfNeed(kv.key), formatValue(kv.value)))
		}

		// 添加表达式
		setParts = append(setParts, b.setExprs...)

		sb.WriteString(strings.Join(setParts, ", "))
	}

	// WHEN 子句
	if len(b.whenConds) > 0 {
		sb.WriteString(" WHEN ")
		sb.WriteString(strings.Join(b.whenConds, " AND "))
	}

	// YIELD 子句
	if len(b.yieldFields) > 0 {
		sb.WriteString(" YIELD ")
		sb.WriteString(strings.Join(b.yieldFields, ", "))
	}

	return sb.String()
}
