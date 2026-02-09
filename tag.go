package ngql

import (
	"fmt"
	"strings"
)

// ============ CREATE TAG ============

// CreateTagBuilder CREATE TAG 构建器
type CreateTagBuilder struct {
	tagName     string
	ifNotExists bool
	properties  []*Property
	ttlDuration int
	ttlCol      string
	comment     string
}

// NewCreateTag 创建 CREATE TAG 构建器
// tagName: Tag 名称
func NewCreateTag(tagName string) *CreateTagBuilder {
	return &CreateTagBuilder{
		tagName:     tagName,
		properties:  make([]*Property, 0),
		ttlDuration: -1,
	}
}

// IfNotExists 设置 IF NOT EXISTS 选项
func (b *CreateTagBuilder) IfNotExists() *CreateTagBuilder {
	b.ifNotExists = true
	return b
}

// AddProperty 添加属性
func (b *CreateTagBuilder) AddProperty(prop *Property) *CreateTagBuilder {
	b.properties = append(b.properties, prop)
	return b
}

// Property 添加属性（简化方法）
// name: 属性名
// typ: 数据类型
func (b *CreateTagBuilder) Property(name string, typ PropertyType) *CreateTagBuilder {
	b.properties = append(b.properties, &Property{
		Name:     name,
		Type:     typ,
		Nullable: true,
	})
	return b
}

// TTL 设置 TTL 配置
// col: TTL 列名，数据类型必须是int或者timestamp
// duration: TTL 持续时间（秒）
func (b *CreateTagBuilder) TTL(col string, duration int) *CreateTagBuilder {
	b.ttlDuration = duration
	b.ttlCol = col
	return b
}

// Comment 设置注释
func (b *CreateTagBuilder) Comment(comment string) *CreateTagBuilder {
	b.comment = comment
	return b
}

// String 构建 nGQL 语句
func (b *CreateTagBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("CREATE TAG ")

	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.tagName))

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

	// TTL 配置
	if b.ttlDuration > 0 && b.ttlCol != "" {
		sb.WriteString(fmt.Sprintf(" TTL_DURATION = %d, TTL_COL = %s", b.ttlDuration, quote(b.ttlCol)))
	}

	// 注释
	if b.comment != "" {
		sb.WriteString(fmt.Sprintf(" COMMENT = %s", quote(b.comment)))
	}

	return sb.String()
}

// ============ ALTER TAG ============

// AlterOperation ALTER 操作类型
type AlterOperation struct {
	Op   string    // "ADD", "DROP", "CHANGE"
	Prop *Property // 属性定义（DROP 时只需要 Name）
}

// AlterTagBuilder ALTER TAG 构建器
type AlterTagBuilder struct {
	tagName     string
	operations  []*AlterOperation
	ttlDuration int
	ttlCol      string
	comment     string
}

// NewAlterTag 创建 ALTER TAG 构建器
// tagName: Tag 名称
func NewAlterTag(tagName string) *AlterTagBuilder {
	return &AlterTagBuilder{
		tagName:     tagName,
		operations:  make([]*AlterOperation, 0),
		ttlDuration: -1,
	}
}

// Add 添加属性
func (b *AlterTagBuilder) Add(prop *Property) *AlterTagBuilder {
	b.operations = append(b.operations, &AlterOperation{
		Op:   "ADD",
		Prop: prop,
	})
	return b
}

// Drop 删除属性
func (b *AlterTagBuilder) Drop(propName string) *AlterTagBuilder {
	b.operations = append(b.operations, &AlterOperation{
		Op:   "DROP",
		Prop: &Property{Name: propName},
	})
	return b
}

// Change 修改属性
func (b *AlterTagBuilder) Change(prop *Property) *AlterTagBuilder {
	b.operations = append(b.operations, &AlterOperation{
		Op:   "CHANGE",
		Prop: prop,
	})
	return b
}

// TTL 设置 TTL 配置
func (b *AlterTagBuilder) TTL(duration int, col string) *AlterTagBuilder {
	b.ttlDuration = duration
	b.ttlCol = col
	return b
}

// Comment 设置注释
func (b *AlterTagBuilder) Comment(comment string) *AlterTagBuilder {
	b.comment = comment
	return b
}

// String 构建 nGQL 语句
func (b *AlterTagBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("ALTER TAG ")
	sb.WriteString(backquoteIfNeed(b.tagName))

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

	// TTL 配置
	if b.ttlDuration > 0 && b.ttlCol != "" {
		sb.WriteString(fmt.Sprintf(" TTL_DURATION = %d, TTL_COL = %s", b.ttlDuration, quote(b.ttlCol)))
	}

	// 注释
	if b.comment != "" {
		sb.WriteString(fmt.Sprintf(" COMMENT = %s", quote(b.comment)))
	}

	return sb.String()
}

// ============ DROP TAG ============

// DropTagBuilder DROP TAG 构建器
type DropTagBuilder struct {
	tagName  string
	ifExists bool
}

// NewDropTag 创建 DROP TAG 构建器
// tagName: Tag 名称
func NewDropTag(tagName string) *DropTagBuilder {
	return &DropTagBuilder{
		tagName: tagName,
	}
}

// IfExists 设置 IF EXISTS 选项
func (b *DropTagBuilder) IfExists() *DropTagBuilder {
	b.ifExists = true
	return b
}

// String 构建 nGQL 语句
func (b *DropTagBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("DROP TAG ")

	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.tagName))
	return sb.String()
}

// ============ SHOW TAGS ============

// ShowTagsBuilder SHOW TAGS 构建器
type ShowTagsBuilder struct{}

// NewShowTags 创建 SHOW TAGS 构建器
func NewShowTags() *ShowTagsBuilder {
	return &ShowTagsBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowTagsBuilder) String() string {
	return "SHOW TAGS"
}

// ============ DESCRIBE TAG ============

// DescribeTagBuilder DESCRIBE TAG 构建器
type DescribeTagBuilder struct {
	tagName string
}

// NewDescribeTag 创建 DESCRIBE TAG 构建器
// tagName: Tag 名称
func NewDescribeTag(tagName string) *DescribeTagBuilder {
	return &DescribeTagBuilder{
		tagName: tagName,
	}
}

// String 构建 nGQL 语句
func (b *DescribeTagBuilder) String() string {
	return fmt.Sprintf("DESCRIBE TAG %s", backquoteIfNeed(b.tagName))
}

// ============ SHOW CREATE TAG ============

// ShowCreateTagBuilder SHOW CREATE TAG 构建器
type ShowCreateTagBuilder struct {
	tagName string
}

// NewShowCreateTag 创建 SHOW CREATE TAG 构建器
// tagName: Tag 名称
func NewShowCreateTag(tagName string) *ShowCreateTagBuilder {
	return &ShowCreateTagBuilder{
		tagName: tagName,
	}
}

// String 构建 nGQL 语句
func (b *ShowCreateTagBuilder) String() string {
	return fmt.Sprintf("SHOW CREATE TAG %s", backquoteIfNeed(b.tagName))
}
