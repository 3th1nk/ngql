package ngql

import (
	"fmt"
	"strings"
)

// ============ 更新顶点 ============

// UpdateVertexBuilder 更新顶点构建器
type UpdateVertexBuilder struct {
	tag         string
	vid         string
	sets        []string
	whenCond    string
	yieldFields []string
}

// NewUpdateVertex 创建更新顶点构建器
// tag: 顶点标签
// vid: 顶点ID
func NewUpdateVertex(tag, vid string) *UpdateVertexBuilder {
	return &UpdateVertexBuilder{
		tag:         tag,
		vid:         vid,
		sets:        make([]string, 0),
		yieldFields: make([]string, 0),
	}
}

// Set 设置属性值
// prop: 属性名
// value: 属性值
func (b *UpdateVertexBuilder) Set(prop string, value interface{}) *UpdateVertexBuilder {
	b.sets = append(b.sets, fmt.Sprintf("%s = %s", prop, formatValue(value)))
	return b
}

// SetExpr 设置属性表达式
// expr: 表达式，例如 "age = age + 1"
func (b *UpdateVertexBuilder) SetExpr(expr string) *UpdateVertexBuilder {
	b.sets = append(b.sets, expr)
	return b
}

// When 添加条件
func (b *UpdateVertexBuilder) When(condition string) *UpdateVertexBuilder {
	b.whenCond = condition
	return b
}

// Yield 添加返回字段
func (b *UpdateVertexBuilder) Yield(fields ...string) *UpdateVertexBuilder {
	b.yieldFields = append(b.yieldFields, fields...)
	return b
}

// String 构建 nGQL 语句
func (b *UpdateVertexBuilder) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("UPDATE VERTEX ON %s %s SET ", b.tag, quote(b.vid)))
	sb.WriteString(strings.Join(b.sets, ", "))

	if b.whenCond != "" {
		sb.WriteString(" WHEN ")
		sb.WriteString(b.whenCond)
	}

	if len(b.yieldFields) > 0 {
		sb.WriteString(" YIELD ")
		sb.WriteString(strings.Join(b.yieldFields, ", "))
	}

	return sb.String()
}

// ============ 更新边 ============

// UpdateEdgeBuilder 更新边构建器
type UpdateEdgeBuilder struct {
	edgeType    string
	srcVid      string
	dstVid      string
	rank        *int64
	sets        []string
	whenCond    string
	yieldFields []string
}

// NewUpdateEdge 创建更新边构建器
// edgeType: 边类型
// srcVid: 源顶点ID
// dstVid: 目标顶点ID
func NewUpdateEdge(edgeType, srcVid, dstVid string) *UpdateEdgeBuilder {
	return &UpdateEdgeBuilder{
		edgeType:    edgeType,
		srcVid:      srcVid,
		dstVid:      dstVid,
		sets:        make([]string, 0),
		yieldFields: make([]string, 0),
	}
}

// WithRank 设置边的 rank
func (b *UpdateEdgeBuilder) WithRank(rank int64) *UpdateEdgeBuilder {
	b.rank = &rank
	return b
}

// Set 设置属性值
func (b *UpdateEdgeBuilder) Set(prop string, value interface{}) *UpdateEdgeBuilder {
	b.sets = append(b.sets, fmt.Sprintf("%s = %s", prop, formatValue(value)))
	return b
}

// SetExpr 设置属性表达式
func (b *UpdateEdgeBuilder) SetExpr(expr string) *UpdateEdgeBuilder {
	b.sets = append(b.sets, expr)
	return b
}

// When 添加条件
func (b *UpdateEdgeBuilder) When(condition string) *UpdateEdgeBuilder {
	b.whenCond = condition
	return b
}

// Yield 添加返回字段
func (b *UpdateEdgeBuilder) Yield(fields ...string) *UpdateEdgeBuilder {
	b.yieldFields = append(b.yieldFields, fields...)
	return b
}

// String 构建 nGQL 语句
func (b *UpdateEdgeBuilder) String() string {
	var sb strings.Builder
	rankStr := ""
	if b.rank != nil {
		rankStr = fmt.Sprintf("@%d", *b.rank)
	}
	sb.WriteString(fmt.Sprintf("UPDATE EDGE ON %s %s -> %s%s SET ",
		b.edgeType, quote(b.srcVid), quote(b.dstVid), rankStr))
	sb.WriteString(strings.Join(b.sets, ", "))

	if b.whenCond != "" {
		sb.WriteString(" WHEN ")
		sb.WriteString(b.whenCond)
	}

	if len(b.yieldFields) > 0 {
		sb.WriteString(" YIELD ")
		sb.WriteString(strings.Join(b.yieldFields, ", "))
	}

	return sb.String()
}
