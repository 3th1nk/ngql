package ngql

import (
	"fmt"
	"strings"
)

// ============ CREATE EDGE ============

// CreateEdgeBuilder CREATE EDGE 构建器
type CreateEdgeBuilder struct {
	edgeName    string
	ifNotExists bool
	properties  []*Property
	ttlDuration int
	ttlCol      string
	comment     string
}

// NewCreateEdge 创建 CREATE EDGE 构建器
// edgeName: Edge 类型名称
func NewCreateEdge(edgeName string) *CreateEdgeBuilder {
	return &CreateEdgeBuilder{
		edgeName:    edgeName,
		properties:  make([]*Property, 0),
		ttlDuration: -1,
	}
}

// IfNotExists 设置 IF NOT EXISTS 选项
func (b *CreateEdgeBuilder) IfNotExists() *CreateEdgeBuilder {
	b.ifNotExists = true
	return b
}

// AddProperty 添加属性
func (b *CreateEdgeBuilder) AddProperty(prop *Property) *CreateEdgeBuilder {
	b.properties = append(b.properties, prop)
	return b
}

// Property 添加属性（简化方法）
// name: 属性名
// typ: 数据类型
func (b *CreateEdgeBuilder) Property(name string, typ PropertyType) *CreateEdgeBuilder {
	b.properties = append(b.properties, &Property{
		Name:     name,
		Type:     typ,
		Nullable: true,
	})
	return b
}

// TTL 设置 TTL 配置
// col: TTL 属性，数据类型必须是int或者timestamp
// duration: TTL 持续时间（秒）
func (b *CreateEdgeBuilder) TTL(col string, duration int) *CreateEdgeBuilder {
	b.ttlDuration = duration
	b.ttlCol = col
	return b
}

// Comment 设置注释
func (b *CreateEdgeBuilder) Comment(comment string) *CreateEdgeBuilder {
	b.comment = comment
	return b
}

// String 构建 nGQL 语句
func (b *CreateEdgeBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("CREATE EDGE ")

	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.edgeName))

	// 构建属性列表
	if len(b.properties) > 0 {
		sb.WriteString("(")
		propStrs := make([]string, 0, len(b.properties))
		for _, prop := range b.properties {
			propStrs = append(propStrs, prop.propDef())
		}
		sb.WriteString(strings.Join(propStrs, ", "))
		sb.WriteString(")")
	}

	appendTTLAndComment(&sb, b.ttlDuration, b.ttlCol, b.comment)

	return sb.String()
}

// ============ ALTER EDGE ============

// AlterEdgeBuilder ALTER EDGE 构建器
type AlterEdgeBuilder struct {
	edgeName    string
	operations  []*AlterOperation
	ttlDuration int
	ttlCol      string
	comment     string
}

// NewAlterEdge 创建 ALTER EDGE 构建器
// edgeName: Edge 类型名称
func NewAlterEdge(edgeName string) *AlterEdgeBuilder {
	return &AlterEdgeBuilder{
		edgeName:    edgeName,
		operations:  make([]*AlterOperation, 0),
		ttlDuration: -1,
	}
}

// Add 添加属性
func (b *AlterEdgeBuilder) Add(prop *Property) *AlterEdgeBuilder {
	b.operations = append(b.operations, &AlterOperation{
		Op:   "ADD",
		Prop: prop,
	})
	return b
}

// Drop 删除属性
func (b *AlterEdgeBuilder) Drop(propName string) *AlterEdgeBuilder {
	b.operations = append(b.operations, &AlterOperation{
		Op:   "DROP",
		Prop: &Property{Name: propName},
	})
	return b
}

// Change 修改属性
func (b *AlterEdgeBuilder) Change(prop *Property) *AlterEdgeBuilder {
	b.operations = append(b.operations, &AlterOperation{
		Op:   "CHANGE",
		Prop: prop,
	})
	return b
}

// TTL 设置 TTL 配置
func (b *AlterEdgeBuilder) TTL(duration int, col string) *AlterEdgeBuilder {
	b.ttlDuration = duration
	b.ttlCol = col
	return b
}

// Comment 设置注释
func (b *AlterEdgeBuilder) Comment(comment string) *AlterEdgeBuilder {
	b.comment = comment
	return b
}

// String 构建 nGQL 语句
func (b *AlterEdgeBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("ALTER EDGE ")
	sb.WriteString(backquoteIfNeed(b.edgeName))

	// 构建操作列表
	if len(b.operations) > 0 {
		opStrs := make([]string, 0, len(b.operations))
		for _, op := range b.operations {
			var opStr string
			switch op.Op {
			case "ADD":
				opStr = fmt.Sprintf("ADD (%s)", op.Prop.propDef())
			case "DROP":
				opStr = fmt.Sprintf("DROP (%s)", backquoteIfNeed(op.Prop.Name))
			case "CHANGE":
				opStr = fmt.Sprintf("CHANGE (%s)", op.Prop.propDef())
			}
			opStrs = append(opStrs, opStr)
		}
		sb.WriteString(" ")
		sb.WriteString(strings.Join(opStrs, ", "))
	}

	appendTTLAndComment(&sb, b.ttlDuration, b.ttlCol, b.comment)

	return sb.String()
}

// ============ DROP EDGE ============

// DropEdgeBuilder DROP EDGE 构建器
type DropEdgeBuilder struct {
	edgeName string
	ifExists bool
}

// NewDropEdge 创建 DROP EDGE 构建器
// edgeName: Edge 类型名称
func NewDropEdge(edgeName string) *DropEdgeBuilder {
	return &DropEdgeBuilder{
		edgeName: edgeName,
	}
}

// IfExists 设置 IF EXISTS 选项
func (b *DropEdgeBuilder) IfExists() *DropEdgeBuilder {
	b.ifExists = true
	return b
}

// String 构建 nGQL 语句
func (b *DropEdgeBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("DROP EDGE ")

	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.edgeName))
	return sb.String()
}

// ============ SHOW EDGES ============

// ShowEdgesBuilder SHOW EDGES 构建器
type ShowEdgesBuilder struct{}

// NewShowEdges 创建 SHOW EDGES 构建器
func NewShowEdges() *ShowEdgesBuilder {
	return &ShowEdgesBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowEdgesBuilder) String() string {
	return "SHOW EDGES"
}

// ============ DESCRIBE EDGE ============

// DescribeEdgeBuilder DESCRIBE EDGE 构建器
type DescribeEdgeBuilder struct {
	edgeName string
}

// NewDescribeEdge 创建 DESCRIBE EDGE 构建器
// edgeName: Edge 类型名称
func NewDescribeEdge(edgeName string) *DescribeEdgeBuilder {
	return &DescribeEdgeBuilder{
		edgeName: edgeName,
	}
}

// String 构建 nGQL 语句
func (b *DescribeEdgeBuilder) String() string {
	return fmt.Sprintf("DESCRIBE EDGE %s", backquoteIfNeed(b.edgeName))
}

// ============ SHOW CREATE EDGE ============

// ShowCreateEdgeBuilder SHOW CREATE EDGE 构建器
type ShowCreateEdgeBuilder struct {
	edgeName string
}

// NewShowCreateEdge 创建 SHOW CREATE EDGE 构建器
// edgeName: Edge 类型名称
func NewShowCreateEdge(edgeName string) *ShowCreateEdgeBuilder {
	return &ShowCreateEdgeBuilder{
		edgeName: edgeName,
	}
}

// String 构建 nGQL 语句
func (b *ShowCreateEdgeBuilder) String() string {
	return fmt.Sprintf("SHOW CREATE EDGE %s", backquoteIfNeed(b.edgeName))
}
