package ngql

import (
	"fmt"
	"strings"
)

// ============ CREATE SPACE ============

// CreateSpaceBuilder CREATE SPACE 构建器
type CreateSpaceBuilder struct {
	spaceName     string
	ifNotExists   bool
	partitionNum  int
	replicaFactor int
	vidType       string // "FIXED_STRING(N)" or "INT64"
	comment       string
}

// NewCreateSpace 创建 CREATE SPACE 构建器
// spaceName: 图空间名称
func NewCreateSpace(spaceName string) *CreateSpaceBuilder {
	return &CreateSpaceBuilder{
		spaceName:     spaceName,
		partitionNum:  -1,
		replicaFactor: -1,
	}
}

// IfNotExists 设置 IF NOT EXISTS 选项
func (b *CreateSpaceBuilder) IfNotExists() *CreateSpaceBuilder {
	b.ifNotExists = true
	return b
}

// PartitionNum 设置分区数量
func (b *CreateSpaceBuilder) PartitionNum(num int) *CreateSpaceBuilder {
	b.partitionNum = num
	return b
}

// ReplicaFactor 设置副本因子
func (b *CreateSpaceBuilder) ReplicaFactor(factor int) *CreateSpaceBuilder {
	b.replicaFactor = factor
	return b
}

func (b *CreateSpaceBuilder) VidTypeInt64() *CreateSpaceBuilder {
	b.vidType = "INT64"
	return b
}

func (b *CreateSpaceBuilder) VidTypeFixedString(length int) *CreateSpaceBuilder {
	b.vidType = fmt.Sprintf("FIXED_STRING(%d)", length)
	return b
}

// Comment 设置注释
func (b *CreateSpaceBuilder) Comment(comment string) *CreateSpaceBuilder {
	b.comment = comment
	return b
}

// String 构建 nGQL 语句
func (b *CreateSpaceBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("CREATE SPACE ")

	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}

	sb.WriteString(b.spaceName)

	// 构建选项列表
	options := make([]string, 0)
	if b.partitionNum > 0 {
		options = append(options, fmt.Sprintf("partition_num = %d", b.partitionNum))
	}
	if b.replicaFactor > 0 {
		options = append(options, fmt.Sprintf("replica_factor = %d", b.replicaFactor))
	}
	if b.vidType != "" {
		options = append(options, fmt.Sprintf("vid_type = %s", b.vidType))
	}

	if len(options) > 0 {
		sb.WriteString(" (")
		sb.WriteString(strings.Join(options, ", "))
		sb.WriteString(")")
	}

	if b.comment != "" {
		sb.WriteString(fmt.Sprintf(" COMMENT = %s", quote(b.comment)))
	}

	return sb.String()
}

// ============ CREATE SPACE AS (Clone) ============

// CreateSpaceAsBuilder CREATE SPACE AS 构建器（克隆图空间）
type CreateSpaceAsBuilder struct {
	newSpaceName string
	oldSpaceName string
}

// NewCreateSpaceAs 创建 CREATE SPACE AS 构建器
// newSpaceName: 新图空间名称
// oldSpaceName: 要克隆的图空间名称
func NewCreateSpaceAs(newSpaceName, oldSpaceName string) *CreateSpaceAsBuilder {
	return &CreateSpaceAsBuilder{
		newSpaceName: newSpaceName,
		oldSpaceName: oldSpaceName,
	}
}

// String 构建 nGQL 语句
func (b *CreateSpaceAsBuilder) String() string {
	return fmt.Sprintf("CREATE SPACE %s AS %s", backquoteIfNeed(b.newSpaceName), backquoteIfNeed(b.oldSpaceName))
}

// ============ USE SPACE ============

// UseSpaceBuilder USE SPACE 构建器
type UseSpaceBuilder struct {
	spaceName string
}

// NewUseSpace 创建 USE SPACE 构建器
// spaceName: 图空间名称
func NewUseSpace(spaceName string) *UseSpaceBuilder {
	return &UseSpaceBuilder{
		spaceName: spaceName,
	}
}

// String 构建 nGQL 语句
func (b *UseSpaceBuilder) String() string {
	return fmt.Sprintf("USE %s", backquoteIfNeed(b.spaceName))
}

// ============ SHOW SPACES ============

// ShowSpacesBuilder SHOW SPACES 构建器
type ShowSpacesBuilder struct{}

// NewShowSpaces 创建 SHOW SPACES 构建器
func NewShowSpaces() *ShowSpacesBuilder {
	return &ShowSpacesBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowSpacesBuilder) String() string {
	return "SHOW SPACES"
}

// ============ DESCRIBE SPACE ============

// DescribeSpaceBuilder DESCRIBE SPACE 构建器
type DescribeSpaceBuilder struct {
	spaceName string
}

// NewDescribeSpace 创建 DESCRIBE SPACE 构建器
// spaceName: 图空间名称
func NewDescribeSpace(spaceName string) *DescribeSpaceBuilder {
	return &DescribeSpaceBuilder{
		spaceName: spaceName,
	}
}

// String 构建 nGQL 语句
func (b *DescribeSpaceBuilder) String() string {
	return fmt.Sprintf("DESCRIBE SPACE %s", backquoteIfNeed(b.spaceName))
}

// ============ DROP SPACE ============

// DropSpaceBuilder DROP SPACE 构建器
type DropSpaceBuilder struct {
	spaceName string
	ifExists  bool
}

// NewDropSpace 创建 DROP SPACE 构建器
// spaceName: 图空间名称
func NewDropSpace(spaceName string) *DropSpaceBuilder {
	return &DropSpaceBuilder{
		spaceName: spaceName,
	}
}

// IfExists 设置 IF EXISTS 选项
func (b *DropSpaceBuilder) IfExists() *DropSpaceBuilder {
	b.ifExists = true
	return b
}

// String 构建 nGQL 语句
func (b *DropSpaceBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("DROP SPACE ")

	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.spaceName))
	return sb.String()
}

// ============ CLEAR SPACE ============

// ClearSpaceBuilder CLEAR SPACE 构建器
type ClearSpaceBuilder struct {
	spaceName string
	ifExists  bool
}

// NewClearSpace 创建 CLEAR SPACE 构建器
// spaceName: 图空间名称
func NewClearSpace(spaceName string) *ClearSpaceBuilder {
	return &ClearSpaceBuilder{
		spaceName: spaceName,
	}
}

// IfExists 设置 IF EXISTS 选项
func (b *ClearSpaceBuilder) IfExists() *ClearSpaceBuilder {
	b.ifExists = true
	return b
}

// String 构建 nGQL 语句
func (b *ClearSpaceBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("CLEAR SPACE ")

	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.spaceName))
	return sb.String()
}
