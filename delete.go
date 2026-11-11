package ngql

import (
	"fmt"
	"strings"
)

// ============ 删除 Tag ============

// DeleteTagBuilder DELETE TAG 构建器，删除顶点上指定的 Tag
type DeleteTagBuilder struct {
	tagName string
	vids    []string
}

// NewDeleteTag 创建 DELETE TAG 构建器
// tagName: Tag 名称
// vids: 顶点 ID 列表
func NewDeleteTag(tagName string, vids ...string) *DeleteTagBuilder {
	return &DeleteTagBuilder{
		tagName: tagName,
		vids:    vids,
	}
}

// String 构建 nGQL 语句
func (b *DeleteTagBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("DELETE TAG ")
	sb.WriteString(backquoteIfNeed(b.tagName))
	sb.WriteString(" FROM ")

	vidStrs := make([]string, 0, len(b.vids))
	for _, vid := range b.vids {
		vidStrs = append(vidStrs, quote(vid))
	}
	sb.WriteString(strings.Join(vidStrs, ", "))

	return sb.String()
}

// ============ 删除顶点 ============

// DeleteVertexBuilder 删除顶点构建器
type DeleteVertexBuilder struct {
	vids     []string
	withEdge bool
}

// NewDeleteVertex 创建删除顶点构建器
func NewDeleteVertex(vids ...string) *DeleteVertexBuilder {
	return &DeleteVertexBuilder{
		vids: vids,
	}
}

// WithEdge 同时删除关联的边
func (b *DeleteVertexBuilder) WithEdge() *DeleteVertexBuilder {
	b.withEdge = true
	return b
}

// String 构建 nGQL 语句
func (b *DeleteVertexBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("DELETE VERTEX ")

	vidStrs := make([]string, 0, len(b.vids))
	for _, vid := range b.vids {
		vidStrs = append(vidStrs, quote(vid))
	}
	sb.WriteString(strings.Join(vidStrs, ", "))

	if b.withEdge {
		sb.WriteString(" WITH EDGE")
	}

	return sb.String()
}

// ============ 删除边 ============

// DeleteEdgeBuilder 删除边构建器
type DeleteEdgeBuilder struct {
	edgeType string
	edges    []deleteEdge
}

type deleteEdge struct {
	srcVid string
	dstVid string
	rank   *int64
}

// NewDeleteEdge 创建删除边构建器
func NewDeleteEdge(edgeType string) *DeleteEdgeBuilder {
	return &DeleteEdgeBuilder{
		edgeType: edgeType,
		edges:    make([]deleteEdge, 0),
	}
}

// Edge 添加要删除的边
func (b *DeleteEdgeBuilder) Edge(srcVid, dstVid string) *DeleteEdgeBuilder {
	b.edges = append(b.edges, deleteEdge{
		srcVid: srcVid,
		dstVid: dstVid,
	})
	return b
}

// EdgeWithRank 添加要删除的带 rank 的边
func (b *DeleteEdgeBuilder) EdgeWithRank(srcVid, dstVid string, rank int64) *DeleteEdgeBuilder {
	b.edges = append(b.edges, deleteEdge{
		srcVid: srcVid,
		dstVid: dstVid,
		rank:   &rank,
	})
	return b
}

// String 构建 nGQL 语句
func (b *DeleteEdgeBuilder) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("DELETE EDGE %s ", backquoteIfNeed(b.edgeType)))

	edgeStrs := make([]string, 0, len(b.edges))
	for _, e := range b.edges {
		rankStr := ""
		if e.rank != nil {
			rankStr = fmt.Sprintf("@%d", *e.rank)
		}
		edgeStrs = append(edgeStrs, fmt.Sprintf("%s -> %s%s", quote(e.srcVid), quote(e.dstVid), rankStr))
	}
	sb.WriteString(strings.Join(edgeStrs, ", "))

	return sb.String()
}
