package ngql

import (
	"fmt"
	"strings"
)

// 用户角色常量
const (
	RoleGOD   = "GOD"   // 最高权限（仅 God 账号可授予）
	RoleAdmin = "ADMIN" // 图空间管理员
	RoleDBA   = "DBA"   // 数据库管理员
	RoleUser  = "USER"  // 普通用户
	RoleGuest = "GUEST" // 访客
)

// ============ 用户管理 ============

// CreateUserBuilder CREATE USER 构建器
type CreateUserBuilder struct {
	userName    string
	password    string
	ifNotExists bool
}

// NewCreateUser 创建 CREATE USER 构建器
// userName: 用户名
// password: 密码
func NewCreateUser(userName, password string) *CreateUserBuilder {
	return &CreateUserBuilder{
		userName: userName,
		password: password,
	}
}

// IfNotExists 设置 IF NOT EXISTS 选项
func (b *CreateUserBuilder) IfNotExists() *CreateUserBuilder {
	b.ifNotExists = true
	return b
}

// String 构建 nGQL 语句
func (b *CreateUserBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("CREATE USER ")

	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.userName))
	sb.WriteString(" WITH PASSWORD ")
	sb.WriteString(quote(b.password))

	return sb.String()
}

// DropUserBuilder DROP USER 构建器
type DropUserBuilder struct {
	userName string
	ifExists bool
}

// NewDropUser 创建 DROP USER 构建器
// userName: 用户名
func NewDropUser(userName string) *DropUserBuilder {
	return &DropUserBuilder{userName: userName}
}

// IfExists 设置 IF EXISTS 选项
func (b *DropUserBuilder) IfExists() *DropUserBuilder {
	b.ifExists = true
	return b
}

// String 构建 nGQL 语句
func (b *DropUserBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("DROP USER ")

	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}

	sb.WriteString(backquoteIfNeed(b.userName))
	return sb.String()
}

// AlterUserBuilder ALTER USER 构建器（修改密码）
type AlterUserBuilder struct {
	userName string
	password string
}

// NewAlterUser 创建 ALTER USER 构建器
// userName: 用户名
// password: 新密码
func NewAlterUser(userName, password string) *AlterUserBuilder {
	return &AlterUserBuilder{
		userName: userName,
		password: password,
	}
}

// String 构建 nGQL 语句
func (b *AlterUserBuilder) String() string {
	return fmt.Sprintf("ALTER USER %s WITH PASSWORD %s", backquoteIfNeed(b.userName), quote(b.password))
}

// ChangePasswordBuilder CHANGE PASSWORD 构建器
type ChangePasswordBuilder struct {
	userName string
	oldPwd   string
	newPwd   string
}

// NewChangePassword 创建 CHANGE PASSWORD 构建器
// userName: 用户名
// oldPwd: 旧密码
// newPwd: 新密码
func NewChangePassword(userName, oldPwd, newPwd string) *ChangePasswordBuilder {
	return &ChangePasswordBuilder{
		userName: userName,
		oldPwd:   oldPwd,
		newPwd:   newPwd,
	}
}

// String 构建 nGQL 语句
func (b *ChangePasswordBuilder) String() string {
	return fmt.Sprintf("CHANGE PASSWORD %s FROM %s TO %s",
		backquoteIfNeed(b.userName), quote(b.oldPwd), quote(b.newPwd))
}

// ============ 权限管理 ============

// GrantRoleBuilder GRANT ROLE 构建器
type GrantRoleBuilder struct {
	roleType  string
	spaceName string
	userName  string
}

// NewGrantRole 创建 GRANT ROLE 构建器
// roleType: 角色类型，使用 RoleAdmin 等常量
// spaceName: 图空间名称
// userName: 用户名
func NewGrantRole(roleType, spaceName, userName string) *GrantRoleBuilder {
	return &GrantRoleBuilder{
		roleType:  roleType,
		spaceName: spaceName,
		userName:  userName,
	}
}

// String 构建 nGQL 语句
func (b *GrantRoleBuilder) String() string {
	return fmt.Sprintf("GRANT ROLE %s ON %s TO %s",
		b.roleType, backquoteIfNeed(b.spaceName), backquoteIfNeed(b.userName))
}

// RevokeRoleBuilder REVOKE ROLE 构建器
type RevokeRoleBuilder struct {
	roleType  string
	spaceName string
	userName  string
}

// NewRevokeRole 创建 REVOKE ROLE 构建器
// roleType: 角色类型，使用 RoleAdmin 等常量
// spaceName: 图空间名称
// userName: 用户名
func NewRevokeRole(roleType, spaceName, userName string) *RevokeRoleBuilder {
	return &RevokeRoleBuilder{
		roleType:  roleType,
		spaceName: spaceName,
		userName:  userName,
	}
}

// String 构建 nGQL 语句
func (b *RevokeRoleBuilder) String() string {
	return fmt.Sprintf("REVOKE ROLE %s ON %s FROM %s",
		b.roleType, backquoteIfNeed(b.spaceName), backquoteIfNeed(b.userName))
}

// ============ 权限查询 ============

// ShowUsersBuilder SHOW USERS 构建器
type ShowUsersBuilder struct{}

// NewShowUsers 创建 SHOW USERS 构建器
func NewShowUsers() *ShowUsersBuilder {
	return &ShowUsersBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowUsersBuilder) String() string {
	return "SHOW USERS"
}

// ShowRolesBuilder SHOW ROLES 构建器
type ShowRolesBuilder struct {
	spaceName string
}

// NewShowRoles 创建 SHOW ROLES 构建器
// spaceName: 图空间名称
func NewShowRoles(spaceName string) *ShowRolesBuilder {
	return &ShowRolesBuilder{spaceName: spaceName}
}

// String 构建 nGQL 语句
func (b *ShowRolesBuilder) String() string {
	return fmt.Sprintf("SHOW ROLES ON %s", backquoteIfNeed(b.spaceName))
}

// ShowGrantsBuilder SHOW GRANTS 构建器
type ShowGrantsBuilder struct {
	spaceName string // 可选，为空查询当前用户的全部权限
}

// NewShowGrants 创建 SHOW GRANTS 构建器
// spaceName: 可选图空间名称，不传查询当前用户的全部权限
func NewShowGrants(spaceName ...string) *ShowGrantsBuilder {
	b := &ShowGrantsBuilder{}
	if len(spaceName) > 0 {
		b.spaceName = spaceName[0]
	}
	return b
}

// String 构建 nGQL 语句
func (b *ShowGrantsBuilder) String() string {
	if b.spaceName == "" {
		return "SHOW GRANTS"
	}
	return fmt.Sprintf("SHOW GRANTS ON %s", backquoteIfNeed(b.spaceName))
}
