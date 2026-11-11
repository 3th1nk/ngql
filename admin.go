package ngql

import (
	"fmt"
	"strings"
)

// ============ SUBMIT JOB ============

// SUBMIT JOB 任务类型常量
const (
	JobTypeCompact      = "COMPACT"       // 数据压实
	JobTypeFlush        = "FLUSH"         // 内存刷盘
	JobTypeStats        = "STATS"         // 统计信息收集
	JobTypeRebuildIndex = "REBUILD INDEX" // 重建索引
)

// SubmitJobBuilder SUBMIT JOB 构建器
type SubmitJobBuilder struct {
	jobType    string
	indexNames []string // REBUILD INDEX 的索引名列表
}

// NewSubmitJob 创建 SUBMIT JOB 构建器
// jobType: 任务类型，使用 JobTypeCompact 等常量
func NewSubmitJob(jobType string) *SubmitJobBuilder {
	return &SubmitJobBuilder{
		jobType:    jobType,
		indexNames: make([]string, 0),
	}
}

// Index 指定要重建的索引名列表（仅 REBUILD INDEX 有效，为空则重建所有索引）
func (b *SubmitJobBuilder) Index(indexNames ...string) *SubmitJobBuilder {
	b.indexNames = append(b.indexNames, indexNames...)
	return b
}

// String 构建 nGQL 语句
func (b *SubmitJobBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("SUBMIT JOB ")
	sb.WriteString(b.jobType)

	if len(b.indexNames) > 0 {
		sb.WriteString(":")
		indexStrs := make([]string, 0, len(b.indexNames))
		for _, name := range b.indexNames {
			indexStrs = append(indexStrs, backquoteIfNeed(name))
		}
		sb.WriteString(strings.Join(indexStrs, ", "))
	}

	return sb.String()
}

// ============ JOB 管理 ============

// ShowJobsBuilder SHOW JOBS 构建器
type ShowJobsBuilder struct{}

// NewShowJobs 创建 SHOW JOBS 构建器
func NewShowJobs() *ShowJobsBuilder {
	return &ShowJobsBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowJobsBuilder) String() string {
	return "SHOW JOBS"
}

// ShowJobBuilder SHOW JOB 构建器
type ShowJobBuilder struct {
	jobID int64
}

// NewShowJob 创建 SHOW JOB 构建器
// jobID: 任务 ID
func NewShowJob(jobID int64) *ShowJobBuilder {
	return &ShowJobBuilder{jobID: jobID}
}

// String 构建 nGQL 语句
func (b *ShowJobBuilder) String() string {
	return fmt.Sprintf("SHOW JOB %d", b.jobID)
}

// StopJobBuilder STOP JOB 构建器
type StopJobBuilder struct {
	jobID int64
}

// NewStopJob 创建 STOP JOB 构建器
// jobID: 任务 ID
func NewStopJob(jobID int64) *StopJobBuilder {
	return &StopJobBuilder{jobID: jobID}
}

// String 构建 nGQL 语句
func (b *StopJobBuilder) String() string {
	return fmt.Sprintf("STOP JOB %d", b.jobID)
}

// RecoverJobBuilder RECOVER JOB 构建器
type RecoverJobBuilder struct{}

// NewRecoverJob 创建 RECOVER JOB 构建器
func NewRecoverJob() *RecoverJobBuilder {
	return &RecoverJobBuilder{}
}

// String 构建 nGQL 语句
func (b *RecoverJobBuilder) String() string {
	return "RECOVER JOB"
}

// ============ SHOW 运维信息 ============

// SHOW HOSTS 服务模块常量
const (
	HostModuleGraph   = "GRAPH"   // Graph 服务
	HostModuleMeta    = "META"    // Meta 服务
	HostModuleStorage = "STORAGE" // Storage 服务
)

// ShowHostsBuilder SHOW HOSTS 构建器
type ShowHostsBuilder struct {
	module string // 可选，GRAPH/META/STORAGE，为空显示所有
}

// NewShowHosts 创建 SHOW HOSTS 构建器
// module: 可选服务模块，使用 HostModuleGraph 等常量，不传显示所有服务
func NewShowHosts(module ...string) *ShowHostsBuilder {
	b := &ShowHostsBuilder{}
	if len(module) > 0 {
		b.module = module[0]
	}
	return b
}

// String 构建 nGQL 语句
func (b *ShowHostsBuilder) String() string {
	if b.module == "" {
		return "SHOW HOSTS"
	}
	return "SHOW HOSTS " + b.module
}

// ShowPartsBuilder SHOW PARTS 构建器
type ShowPartsBuilder struct {
	partID int64 // 分片 ID，-1 表示所有
}

// NewShowParts 创建 SHOW PARTS 构建器
// partID: 分片 ID，传 -1 显示所有分片
func NewShowParts(partID int64) *ShowPartsBuilder {
	return &ShowPartsBuilder{partID: partID}
}

// String 构建 nGQL 语句
func (b *ShowPartsBuilder) String() string {
	if b.partID < 0 {
		return "SHOW PARTS"
	}
	return fmt.Sprintf("SHOW PARTS %d", b.partID)
}

// ShowStatsBuilder SHOW STATS 构建器
type ShowStatsBuilder struct{}

// NewShowStats 创建 SHOW STATS 构建器（需先执行 SUBMIT JOB STATS）
func NewShowStats() *ShowStatsBuilder {
	return &ShowStatsBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowStatsBuilder) String() string {
	return "SHOW STATS"
}

// ShowSnapshotsBuilder SHOW SNAPSHOTS 构建器
type ShowSnapshotsBuilder struct{}

// NewShowSnapshots 创建 SHOW SNAPSHOTS 构建器
func NewShowSnapshots() *ShowSnapshotsBuilder {
	return &ShowSnapshotsBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowSnapshotsBuilder) String() string {
	return "SHOW SNAPSHOTS"
}

// ShowSessionsBuilder SHOW SESSIONS 构建器
type ShowSessionsBuilder struct{}

// NewShowSessions 创建 SHOW SESSIONS 构建器
func NewShowSessions() *ShowSessionsBuilder {
	return &ShowSessionsBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowSessionsBuilder) String() string {
	return "SHOW SESSIONS"
}

// SHOW QUERIES 范围常量
const (
	QueryScopeCurrent = "CURRENT" // 当前 session 的查询
	QueryScopeAll     = "ALL"     // 所有 session 的查询
)

// ShowQueriesBuilder SHOW QUERIES 构建器
type ShowQueriesBuilder struct {
	scope string // 可选，CURRENT/ALL，为空显示当前 session
}

// NewShowQueries 创建 SHOW QUERIES 构建器
// scope: 可选范围，使用 QueryScopeCurrent/QueryScopeAll 常量，不传显示当前 session 查询
func NewShowQueries(scope ...string) *ShowQueriesBuilder {
	b := &ShowQueriesBuilder{}
	if len(scope) > 0 {
		b.scope = scope[0]
	}
	return b
}

// String 构建 nGQL 语句
func (b *ShowQueriesBuilder) String() string {
	if b.scope == "" {
		return "SHOW QUERIES"
	}
	return "SHOW QUERIES " + b.scope
}

// ShowCharsetBuilder SHOW CHARSET 构建器
type ShowCharsetBuilder struct{}

// NewShowCharset 创建 SHOW CHARSET 构建器
func NewShowCharset() *ShowCharsetBuilder {
	return &ShowCharsetBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowCharsetBuilder) String() string {
	return "SHOW CHARSET"
}

// ShowCollationBuilder SHOW COLLATION 构建器
type ShowCollationBuilder struct{}

// NewShowCollation 创建 SHOW COLLATION 构建器
func NewShowCollation() *ShowCollationBuilder {
	return &ShowCollationBuilder{}
}

// String 构建 nGQL 语句
func (b *ShowCollationBuilder) String() string {
	return "SHOW COLLATION"
}

// ShowCreateSpaceBuilder SHOW CREATE SPACE 构建器
type ShowCreateSpaceBuilder struct {
	spaceName string
}

// NewShowCreateSpace 创建 SHOW CREATE SPACE 构建器
// spaceName: 图空间名称
func NewShowCreateSpace(spaceName string) *ShowCreateSpaceBuilder {
	return &ShowCreateSpaceBuilder{spaceName: spaceName}
}

// String 构建 nGQL 语句
func (b *ShowCreateSpaceBuilder) String() string {
	return fmt.Sprintf("SHOW CREATE SPACE %s", backquoteIfNeed(b.spaceName))
}

// ============ CONFIGS ============

// SHOW/UPDATE CONFIGS 服务模块常量
const (
	ConfigModuleGraph   = "graph"   // Graph 服务配置
	ConfigModuleMeta    = "meta"    // Meta 服务配置
	ConfigModuleStorage = "storage" // Storage 服务配置
)

// ShowConfigsBuilder SHOW CONFIGS 构建器
type ShowConfigsBuilder struct {
	module string // 可选模块，为空显示所有
}

// NewShowConfigs 创建 SHOW CONFIGS 构建器
// module: 可选模块，使用 ConfigModuleGraph 等常量，不传显示所有配置
func NewShowConfigs(module ...string) *ShowConfigsBuilder {
	b := &ShowConfigsBuilder{}
	if len(module) > 0 {
		b.module = module[0]
	}
	return b
}

// String 构建 nGQL 语句
func (b *ShowConfigsBuilder) String() string {
	if b.module == "" {
		return "SHOW CONFIGS"
	}
	return "SHOW CONFIGS " + b.module
}

// UpdateConfigBuilder UPDATE CONFIGS 构建器
type UpdateConfigBuilder struct {
	module string // 可选模块，为空时更新所有模块中的同名配置
	key    string
	value  interface{}
}

// NewUpdateConfig 创建 UPDATE CONFIGS 构建器
// module: 可选模块（ConfigModuleGraph 等），传空时不限定模块
// key: 配置项名称
// value: 配置项的值
func NewUpdateConfig(module, key string, value interface{}) *UpdateConfigBuilder {
	return &UpdateConfigBuilder{
		module: module,
		key:    key,
		value:  value,
	}
}

// String 构建 nGQL 语句
func (b *UpdateConfigBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("UPDATE CONFIGS ")
	if b.module != "" {
		sb.WriteString(b.module)
		sb.WriteString(":")
	}
	sb.WriteString(b.key)
	sb.WriteString("=")
	sb.WriteString(formatValue(b.value))
	return sb.String()
}

// ReloadConfigsBuilder RELOAD CONFIGS 构建器
type ReloadConfigsBuilder struct{}

// NewReloadConfigs 创建 RELOAD CONFIGS 构建器
func NewReloadConfigs() *ReloadConfigsBuilder {
	return &ReloadConfigsBuilder{}
}

// String 构建 nGQL 语句
func (b *ReloadConfigsBuilder) String() string {
	return "RELOAD CONFIGS"
}
