package ngql

import (
	"fmt"
	"strings"
)

// Pattern 图模式接口
type Pattern interface {
	String() string
}

// ============ 顶点模式 (Vertex Pattern) ============

// VertexPattern 顶点模式
// 语法：(variable[:tag] [{prop: value, ...}])
type VertexPattern struct {
	variable   string                 // 变量名，可选
	tags       []string               // Tag 列表，可选
	properties map[string]interface{} // 属性过滤条件，可选
}

// NewVertexPattern 创建顶点模式
// 示例：
//   - NewVertexPattern("") -> ()
//   - NewVertexPattern("v") -> (v)
//   - NewVertexPattern("v").Tag("player") -> (v:player)
//   - NewVertexPattern("v").Tag("player").Property("name", "Tim") -> (v:player {name: "Tim"})
func NewVertexPattern(variable string) *VertexPattern {
	return &VertexPattern{
		variable:   variable,
		tags:       make([]string, 0),
		properties: make(map[string]interface{}),
	}
}

// Tag 添加 Tag 标签
func (p *VertexPattern) Tag(tag string) *VertexPattern {
	p.tags = append(p.tags, tag)
	return p
}

// Property 添加属性过滤条件
func (p *VertexPattern) Property(key string, value interface{}) *VertexPattern {
	p.properties[key] = value
	return p
}

// String 生成顶点模式字符串
func (p *VertexPattern) String() string {
	var sb strings.Builder
	sb.WriteString("(")

	// 变量名
	if p.variable != "" {
		sb.WriteString(p.variable)
	}

	// Tag 列表
	if len(p.tags) > 0 {
		sb.WriteString(":")
		sb.WriteString(strings.Join(p.tags, ":"))
	}

	// 属性过滤
	if len(p.properties) > 0 {
		sb.WriteString(" {")
		first := true
		for key, value := range p.properties {
			if !first {
				sb.WriteString(", ")
			}
			first = false
			sb.WriteString(fmt.Sprintf("%s: %v", key, formatValue(value)))
		}
		sb.WriteString("}")
	}

	sb.WriteString(")")
	return sb.String()
}

// ============ 边模式 (Edge Pattern) ============

// EdgeDirection 边的方向
type EdgeDirection int

const (
	EdgeDirectionRight EdgeDirection = iota // 右向：->
	EdgeDirectionLeft                       // 左向：<-
	EdgeDirectionBoth                       // 双向：-
)

// EdgePattern 边模式
// 语法：-[variable[:edge_type] [{prop: value, ...}] [*[min..max]]]-{>|<|-}
type EdgePattern struct {
	variable   string                 // 变量名，可选
	edgeTypes  []string               // 边类型列表，可选
	properties map[string]interface{} // 属性过滤条件，可选
	direction  EdgeDirection          // 边的方向
	minHops    int                    // 最小跳数，0 表示不限制
	maxHops    int                    // 最大跳数，0 表示不限制
}

// NewEdgePattern 创建边模式
// 示例：
//   - NewEdgePattern(EdgeDirectionRight) -> ->
//   - NewEdgePattern(EdgeDirectionRight).Type("follow") -> -[:follow]->
//   - NewEdgePattern(EdgeDirectionRight).Variable("e").Type("follow") -> -[e:follow]->
//   - NewEdgePattern(EdgeDirectionRight).Type("follow").Hops(1, 3) -> -[:follow*1..3]->
func NewEdgePattern(direction EdgeDirection) *EdgePattern {
	return &EdgePattern{
		edgeTypes:  make([]string, 0),
		properties: make(map[string]interface{}),
		direction:  direction,
	}
}

// Variable 设置边变量名
func (p *EdgePattern) Variable(variable string) *EdgePattern {
	p.variable = variable
	return p
}

// Type 添加边类型
func (p *EdgePattern) Type(edgeType string) *EdgePattern {
	p.edgeTypes = append(p.edgeTypes, edgeType)
	return p
}

// Property 添加属性过滤条件
func (p *EdgePattern) Property(key string, value interface{}) *EdgePattern {
	p.properties[key] = value
	return p
}

// Hops 设置变长边的跳数范围
// min: 最小跳数（0 表示不限制）
// max: 最大跳数（0 表示不限制）
// 示例：
//   - Hops(2, 2) -> *2
//   - Hops(1, 3) -> *1..3
//   - Hops(0, 0) -> * (任意跳数)
func (p *EdgePattern) Hops(min, max int) *EdgePattern {
	p.minHops = min
	p.maxHops = max
	return p
}

// String 生成边模式字符串
func (p *EdgePattern) String() string {
	var sb strings.Builder

	// 边定义
	hasEdgeInfo := p.variable != "" || len(p.edgeTypes) > 0 || len(p.properties) > 0 || p.minHops > 0 || p.maxHops > 0

	if hasEdgeInfo {
		// 左侧连接符
		if p.direction == EdgeDirectionLeft {
			sb.WriteString("<")
		}
		sb.WriteString("-")

		sb.WriteString("[")

		// 变量名
		if p.variable != "" {
			sb.WriteString(p.variable)
		}

		// 边类型
		if len(p.edgeTypes) > 0 {
			sb.WriteString(":")
			sb.WriteString(strings.Join(p.edgeTypes, "|"))
		}

		// 属性过滤
		if len(p.properties) > 0 {
			sb.WriteString(" {")
			first := true
			for key, value := range p.properties {
				if !first {
					sb.WriteString(", ")
				}
				first = false
				sb.WriteString(fmt.Sprintf("%s: %v", key, formatValue(value)))
			}
			sb.WriteString("}")
		}

		// 变长边
		if p.minHops > 0 || p.maxHops > 0 {
			sb.WriteString("*")
			if p.minHops == p.maxHops && p.minHops > 0 {
				// 固定步数：*2
				sb.WriteString(fmt.Sprintf("%d", p.minHops))
			} else if p.minHops > 0 || p.maxHops > 0 {
				// 范围步数：*1..3
				if p.minHops > 0 {
					sb.WriteString(fmt.Sprintf("%d", p.minHops))
				}
				sb.WriteString("..")
				if p.maxHops > 0 {
					sb.WriteString(fmt.Sprintf("%d", p.maxHops))
				}
			}
			// 如果 minHops 和 maxHops 都是 0，则为 * (任意跳数)
		}

		sb.WriteString("]")

		// 右侧连接符
		sb.WriteString("-")
		if p.direction == EdgeDirectionRight {
			sb.WriteString(">")
		}
	} else {
		// 简化形式：没有边信息时
		if p.direction == EdgeDirectionLeft {
			sb.WriteString("<-")
		} else if p.direction == EdgeDirectionRight {
			sb.WriteString("->")
		} else {
			sb.WriteString("-")
		}
	}

	return sb.String()
}

// ============ 路径模式 (Path Pattern) ============

// PathPattern 路径模式
// 由顶点和边组成的完整路径
type PathPattern struct {
	variable string    // 路径变量名，可选
	elements []Pattern // 路径元素（顶点和边的交替序列）
}

// NewPathPattern 创建路径模式
// 示例：
//   - NewPathPattern().Vertex("v1").Edge(EdgeDirectionRight).Vertex("v2")
//     -> (v1)->(v2)
//   - NewPathPattern().Variable("p").Vertex("v1").Edge(EdgeDirectionRight).Vertex("v2")
//     -> p = (v1)->(v2)
func NewPathPattern() *PathPattern {
	return &PathPattern{
		elements: make([]Pattern, 0),
	}
}

// Variable 设置路径变量名
func (p *PathPattern) Variable(variable string) *PathPattern {
	p.variable = variable
	return p
}

// Vertex 添加顶点模式
func (p *PathPattern) Vertex(variable string) *PathPattern {
	p.elements = append(p.elements, NewVertexPattern(variable))
	return p
}

// VertexWithTag 添加带 Tag 的顶点模式
func (p *PathPattern) VertexWithTag(variable string, tag string) *PathPattern {
	p.elements = append(p.elements, NewVertexPattern(variable).Tag(tag))
	return p
}

// VertexPattern 添加自定义顶点模式
func (p *PathPattern) VertexPattern(vertex *VertexPattern) *PathPattern {
	p.elements = append(p.elements, vertex)
	return p
}

// Edge 添加边模式
func (p *PathPattern) Edge(direction EdgeDirection) *PathPattern {
	p.elements = append(p.elements, NewEdgePattern(direction))
	return p
}

// EdgeWithType 添加带类型的边模式
func (p *PathPattern) EdgeWithType(direction EdgeDirection, edgeType string) *PathPattern {
	p.elements = append(p.elements, NewEdgePattern(direction).Type(edgeType))
	return p
}

// EdgePattern 添加自定义边模式
func (p *PathPattern) EdgePattern(edge *EdgePattern) *PathPattern {
	p.elements = append(p.elements, edge)
	return p
}

// String 生成路径模式字符串
func (p *PathPattern) String() string {
	var sb strings.Builder

	// 路径变量
	if p.variable != "" {
		sb.WriteString(p.variable)
		sb.WriteString(" = ")
	}

	// 路径元素
	for _, element := range p.elements {
		sb.WriteString(element.String())
	}

	return sb.String()
}
