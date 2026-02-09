package ngql

import (
	"fmt"
	"strings"
)

// ============ 插入顶点 ============

// InsertVertexBuilder 插入顶点构建器
type InsertVertexBuilder struct {
	tag        string
	props      []string
	values     []VertexValue
	ifNotExist bool
}

// VertexValue 顶点值
type VertexValue struct {
	Vid    string
	Values []interface{}
}

// NewInsertVertex 创建插入顶点构建器
// tag: 顶点标签名称
// props: 属性名称列表
func NewInsertVertex(tag string, props ...string) *InsertVertexBuilder {
	return &InsertVertexBuilder{
		tag:    tag,
		props:  props,
		values: make([]VertexValue, 0),
	}
}

// IfNotExists 设置 IF NOT EXISTS 选项
func (b *InsertVertexBuilder) IfNotExists() *InsertVertexBuilder {
	b.ifNotExist = true
	return b
}

// Values 添加顶点数据
// vid: 顶点ID
// values: 属性值列表，顺序需与 props 一致
func (b *InsertVertexBuilder) Values(vid string, values ...interface{}) *InsertVertexBuilder {
	b.values = append(b.values, VertexValue{
		Vid:    vid,
		Values: values,
	})
	return b
}

// String 构建 nGQL 语句
func (b *InsertVertexBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("INSERT VERTEX ")

	if b.ifNotExist {
		sb.WriteString("IF NOT EXISTS ")
	}

	sb.WriteString(b.tag)
	sb.WriteString("(")
	sb.WriteString(strings.Join(b.props, ", "))
	sb.WriteString(") VALUES ")

	valueStrs := make([]string, 0, len(b.values))
	for _, v := range b.values {
		vals := make([]string, 0, len(v.Values))
		for _, val := range v.Values {
			vals = append(vals, formatValue(val))
		}
		valueStrs = append(valueStrs, fmt.Sprintf("%s: (%s)", quote(v.Vid), strings.Join(vals, ", ")))
	}
	sb.WriteString(strings.Join(valueStrs, ", "))

	return sb.String()
}

// ============ 插入边 ============

// InsertEdgeBuilder 插入边构建器
type InsertEdgeBuilder struct {
	edgeType   string
	props      []string
	edges      []EdgeValue
	ifNotExist bool
}

// EdgeValue 边值
type EdgeValue struct {
	SrcVid string
	DstVid string
	Rank   *int64
	Values []interface{}
}

// NewInsertEdge 创建插入边构建器
// edgeType: 边类型名称
// props: 属性名称列表
func NewInsertEdge(edgeType string, props ...string) *InsertEdgeBuilder {
	return &InsertEdgeBuilder{
		edgeType: edgeType,
		props:    props,
		edges:    make([]EdgeValue, 0),
	}
}

// IfNotExists 设置 IF NOT EXISTS 选项
func (b *InsertEdgeBuilder) IfNotExists() *InsertEdgeBuilder {
	b.ifNotExist = true
	return b
}

// Edge 添加边数据
// srcVid: 源顶点ID
// dstVid: 目标顶点ID
// values: 属性值列表，顺序需与 props 一致
func (b *InsertEdgeBuilder) Edge(srcVid, dstVid string, values ...interface{}) *InsertEdgeBuilder {
	b.edges = append(b.edges, EdgeValue{
		SrcVid: srcVid,
		DstVid: dstVid,
		Values: values,
	})
	return b
}

// EdgeWithRank 添加带 rank 的边数据
// srcVid: 源顶点ID
// dstVid: 目标顶点ID
// rank: 边的排序值
// values: 属性值列表，顺序需与 props 一致
func (b *InsertEdgeBuilder) EdgeWithRank(srcVid, dstVid string, rank int64, values ...interface{}) *InsertEdgeBuilder {
	b.edges = append(b.edges, EdgeValue{
		SrcVid: srcVid,
		DstVid: dstVid,
		Rank:   &rank,
		Values: values,
	})
	return b
}

// String 构建 nGQL 语句
func (b *InsertEdgeBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("INSERT EDGE ")

	if b.ifNotExist {
		sb.WriteString("IF NOT EXISTS ")
	}

	sb.WriteString(b.edgeType)
	sb.WriteString("(")
	sb.WriteString(strings.Join(b.props, ", "))
	sb.WriteString(") VALUES ")

	edgeStrs := make([]string, 0, len(b.edges))
	for _, e := range b.edges {
		vals := make([]string, 0, len(e.Values))
		for _, val := range e.Values {
			vals = append(vals, formatValue(val))
		}

		rankStr := ""
		if e.Rank != nil {
			rankStr = fmt.Sprintf("@%d", *e.Rank)
		}

		edgeStrs = append(edgeStrs, fmt.Sprintf("%s -> %s%s: (%s)",
			quote(e.SrcVid), quote(e.DstVid), rankStr, strings.Join(vals, ", ")))
	}
	sb.WriteString(strings.Join(edgeStrs, ", "))

	return sb.String()
}
