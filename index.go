package ngql

import (
	"fmt"
	"strings"
)

// ============ CREATE TAG INDEX ============

// CreateTagIndexBuilder CREATE TAG INDEX 构建器
type CreateTagIndexBuilder struct {
	indexName   string
	tagName     string
	properties  []*IndexProperty
	ifNotExists bool
	comment     string
}

// NewCreateTagIndex 创建 CREATE TAG INDEX 构建器
// indexName: 索引名称
// tagName: Tag 名称
func NewCreateTagIndex(indexName, tagName string) *CreateTagIndexBuilder {
	return &CreateTagIndexBuilder{
		indexName:  indexName,
		tagName:    tagName,
		properties: make([]*IndexProperty, 0),
	}
}

// IfNotExists 设置 IF NOT EXISTS 选项
func (b *CreateTagIndexBuilder) IfNotExists() *CreateTagIndexBuilder {
	b.ifNotExists = true
	return b
}

// Property 添加索引属性
// name: 属性名称
// length: 字符串属性的索引长度（0 表示非字符串类型或使用默认长度）
func (b *CreateTagIndexBuilder) Property(name string, length int) *CreateTagIndexBuilder {
	b.properties = append(b.properties, &IndexProperty{
		Name:   name,
		Length: length,
	})
	return b
}

// Comment 设置注释
func (b *CreateTagIndexBuilder) Comment(comment string) *CreateTagIndexBuilder {
	b.comment = comment
	return b
}

// String 构建 nGQL 语句
func (b *CreateTagIndexBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("CREATE TAG INDEX ")

	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.indexName))
	sb.WriteString(" ON ")
	sb.WriteString(backquoteIfNeed(b.tagName))

	// 属性列表
	if len(b.properties) > 0 {
		sb.WriteString("(")
		propStrs := make([]string, 0, len(b.properties))
		for _, prop := range b.properties {
			propStrs = append(propStrs, prop.indexPropDef())
		}
		sb.WriteString(strings.Join(propStrs, ", "))
		sb.WriteString(")")
	} else {
		sb.WriteString("()")
	}

	// 注释
	if b.comment != "" {
		sb.WriteString(fmt.Sprintf(" COMMENT %s", quote(b.comment)))
	}

	return sb.String()
}

// ============ CREATE EDGE INDEX ============

// CreateEdgeIndexBuilder CREATE EDGE INDEX 构建器
type CreateEdgeIndexBuilder struct {
	indexName   string
	edgeName    string
	properties  []*IndexProperty
	ifNotExists bool
	comment     string
}

// NewCreateEdgeIndex 创建 CREATE EDGE INDEX 构建器
// indexName: 索引名称
// edgeName: Edge 类型名称
func NewCreateEdgeIndex(indexName, edgeName string) *CreateEdgeIndexBuilder {
	return &CreateEdgeIndexBuilder{
		indexName:  indexName,
		edgeName:   edgeName,
		properties: make([]*IndexProperty, 0),
	}
}

// IfNotExists 设置 IF NOT EXISTS 选项
func (b *CreateEdgeIndexBuilder) IfNotExists() *CreateEdgeIndexBuilder {
	b.ifNotExists = true
	return b
}

// Property 添加索引属性
// name: 属性名称
// length: 字符串属性的索引长度（0 表示非字符串类型或使用默认长度）
func (b *CreateEdgeIndexBuilder) Property(name string, length int) *CreateEdgeIndexBuilder {
	b.properties = append(b.properties, &IndexProperty{
		Name:   name,
		Length: length,
	})
	return b
}

// Comment 设置注释
func (b *CreateEdgeIndexBuilder) Comment(comment string) *CreateEdgeIndexBuilder {
	b.comment = comment
	return b
}

// String 构建 nGQL 语句
func (b *CreateEdgeIndexBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("CREATE EDGE INDEX ")

	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.indexName))
	sb.WriteString(" ON ")
	sb.WriteString(backquoteIfNeed(b.edgeName))

	// 属性列表
	if len(b.properties) > 0 {
		sb.WriteString("(")
		propStrs := make([]string, 0, len(b.properties))
		for _, prop := range b.properties {
			propStrs = append(propStrs, prop.indexPropDef())
		}
		sb.WriteString(strings.Join(propStrs, ", "))
		sb.WriteString(")")
	} else {
		sb.WriteString("()")
	}

	// 注释
	if b.comment != "" {
		sb.WriteString(fmt.Sprintf(" COMMENT %s", quote(b.comment)))
	}

	return sb.String()
}

// ============ DROP TAG INDEX ============

// DropTagIndexBuilder DROP TAG INDEX 构建器
type DropTagIndexBuilder struct {
	indexName string
	ifExists  bool
}

// NewDropTagIndex 创建 DROP TAG INDEX 构建器
// indexName: 索引名称
func NewDropTagIndex(indexName string) *DropTagIndexBuilder {
	return &DropTagIndexBuilder{
		indexName: indexName,
	}
}

// IfExists 设置 IF EXISTS 选项
func (b *DropTagIndexBuilder) IfExists() *DropTagIndexBuilder {
	b.ifExists = true
	return b
}

// String 构建 nGQL 语句
func (b *DropTagIndexBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("DROP TAG INDEX ")

	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.indexName))
	return sb.String()
}

// ============ DROP EDGE INDEX ============

// DropEdgeIndexBuilder DROP EDGE INDEX 构建器
type DropEdgeIndexBuilder struct {
	indexName string
	ifExists  bool
}

// NewDropEdgeIndex 创建 DROP EDGE INDEX 构建器
// indexName: 索引名称
func NewDropEdgeIndex(indexName string) *DropEdgeIndexBuilder {
	return &DropEdgeIndexBuilder{
		indexName: indexName,
	}
}

// IfExists 设置 IF EXISTS 选项
func (b *DropEdgeIndexBuilder) IfExists() *DropEdgeIndexBuilder {
	b.ifExists = true
	return b
}

// String 构建 nGQL 语句
func (b *DropEdgeIndexBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("DROP EDGE INDEX ")

	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.indexName))
	return sb.String()
}

// ============ REBUILD TAG INDEX ============

// RebuildTagIndexBuilder REBUILD TAG INDEX 构建器
type RebuildTagIndexBuilder struct {
	indexNames []string
}

// NewRebuildTagIndex 创建 REBUILD TAG INDEX 构建器
// indexNames: 索引名称列表（可选，为空则重建所有）
func NewRebuildTagIndex(indexNames ...string) *RebuildTagIndexBuilder {
	return &RebuildTagIndexBuilder{
		indexNames: indexNames,
	}
}

// String 构建 nGQL 语句
func (b *RebuildTagIndexBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("REBUILD TAG INDEX")

	if len(b.indexNames) > 0 {
		sb.WriteString(" ")
		quotedNames := make([]string, 0, len(b.indexNames))
		for _, name := range b.indexNames {
			quotedNames = append(quotedNames, backquoteIfNeed(name))
		}
		sb.WriteString(strings.Join(quotedNames, ", "))
	}

	return sb.String()
}

// ============ REBUILD EDGE INDEX ============

// RebuildEdgeIndexBuilder REBUILD EDGE INDEX 构建器
type RebuildEdgeIndexBuilder struct {
	indexNames []string
}

// NewRebuildEdgeIndex 创建 REBUILD EDGE INDEX 构建器
// indexNames: 索引名称列表（可选，为空则重建所有）
func NewRebuildEdgeIndex(indexNames ...string) *RebuildEdgeIndexBuilder {
	return &RebuildEdgeIndexBuilder{
		indexNames: indexNames,
	}
}

// String 构建 nGQL 语句
func (b *RebuildEdgeIndexBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("REBUILD EDGE INDEX")

	if len(b.indexNames) > 0 {
		sb.WriteString(" ")
		quotedNames := make([]string, 0, len(b.indexNames))
		for _, name := range b.indexNames {
			quotedNames = append(quotedNames, backquoteIfNeed(name))
		}
		sb.WriteString(strings.Join(quotedNames, ", "))
	}

	return sb.String()
}

// ============ SHOW TAG INDEXES ============

// ShowTagIndexesBuilder SHOW TAG INDEXES 构建器
type ShowTagIndexesBuilder struct{}

// NewShowTagIndexes 创建 SHOW TAG INDEXES 构建器
func NewShowTagIndexes() *ShowTagIndexesBuilder {
	return &ShowTagIndexesBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowTagIndexesBuilder) String() string {
	return "SHOW TAG INDEXES"
}

// ============ SHOW EDGE INDEXES ============

// ShowEdgeIndexesBuilder SHOW EDGE INDEXES 构建器
type ShowEdgeIndexesBuilder struct{}

// NewShowEdgeIndexes 创建 SHOW EDGE INDEXES 构建器
func NewShowEdgeIndexes() *ShowEdgeIndexesBuilder {
	return &ShowEdgeIndexesBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowEdgeIndexesBuilder) String() string {
	return "SHOW EDGE INDEXES"
}

// ============ DESCRIBE TAG INDEX ============

// DescribeTagIndexBuilder DESCRIBE TAG INDEX 构建器
type DescribeTagIndexBuilder struct {
	indexName string
}

// NewDescribeTagIndex 创建 DESCRIBE TAG INDEX 构建器
// indexName: 索引名称
func NewDescribeTagIndex(indexName string) *DescribeTagIndexBuilder {
	return &DescribeTagIndexBuilder{
		indexName: indexName,
	}
}

// String 构建 nGQL 语句
func (b *DescribeTagIndexBuilder) String() string {
	return fmt.Sprintf("DESCRIBE TAG INDEX %s", backquoteIfNeed(b.indexName))
}

// ============ DESCRIBE EDGE INDEX ============

// DescribeEdgeIndexBuilder DESCRIBE EDGE INDEX 构建器
type DescribeEdgeIndexBuilder struct {
	indexName string
}

// NewDescribeEdgeIndex 创建 DESCRIBE EDGE INDEX 构建器
// indexName: 索引名称
func NewDescribeEdgeIndex(indexName string) *DescribeEdgeIndexBuilder {
	return &DescribeEdgeIndexBuilder{
		indexName: indexName,
	}
}

// String 构建 nGQL 语句
func (b *DescribeEdgeIndexBuilder) String() string {
	return fmt.Sprintf("DESCRIBE EDGE INDEX %s", backquoteIfNeed(b.indexName))
}

// ============ SHOW TAG INDEX STATUS ============

// ShowTagIndexStatusBuilder SHOW TAG INDEX STATUS 构建器
type ShowTagIndexStatusBuilder struct{}

// NewShowTagIndexStatus 创建 SHOW TAG INDEX STATUS 构建器
func NewShowTagIndexStatus() *ShowTagIndexStatusBuilder {
	return &ShowTagIndexStatusBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowTagIndexStatusBuilder) String() string {
	return "SHOW TAG INDEX STATUS"
}

// ============ SHOW EDGE INDEX STATUS ============

// ShowEdgeIndexStatusBuilder SHOW EDGE INDEX STATUS 构建器
type ShowEdgeIndexStatusBuilder struct{}

// NewShowEdgeIndexStatus 创建 SHOW EDGE INDEX STATUS 构建器
func NewShowEdgeIndexStatus() *ShowEdgeIndexStatusBuilder {
	return &ShowEdgeIndexStatusBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowEdgeIndexStatusBuilder) String() string {
	return "SHOW EDGE INDEX STATUS"
}
