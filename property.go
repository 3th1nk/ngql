package ngql

import "fmt"

//go:generate stringer -type=PropertyType -output=property_type_string.go -linecomment
type PropertyType int

const (
	_                         PropertyType = iota
	PropertyTypeInt64                      // int64
	PropertyTypeBoolean                    // bool
	PropertyTypeString                     // string
	PropertyTypeFixedString                // fixed_string
	PropertyTypeDouble                     // double
	PropertyTypeInt32                      // int32
	PropertyTypeInt16                      // int16
	PropertyTypeInt8                       // int8
	PropertyTypeFloat                      // float
	PropertyTypeDate                       // date
	PropertyTypeTime                       // time
	PropertyTypeDateTime                   // datetime
	PropertyTypeTimestamp                  // timestamp
	PropertyTypeGeo                        // geography
	PropertyTypeGeoPoint                   // geography(point)
	PropertyTypeGeoLineString              // geography(linestring)
	PropertyTypeGeoPolygon                 // geography(polygon)
	PropertyTypeDuration                   // duration
)

// Property 属性定义
type Property struct {
	Name     string
	Type     PropertyType // 数据类型
	Nullable bool         // 是否可为空
	Default  interface{}  // 默认值
	Comment  string       // 注释
}

func (p *Property) propDef() string {
	propStr := fmt.Sprintf("%s %s", backquoteIfNeed(p.Name), p.Type)
	if !p.Nullable {
		propStr += " NOT NULL"
	}
	if p.Default != nil {
		propStr += fmt.Sprintf(" DEFAULT %s", formatValue(p.Default))
	}
	if p.Comment != "" {
		propStr += fmt.Sprintf(" COMMENT %s", quote(p.Comment))
	}
	return propStr
}

// IndexProperty 索引属性定义
type IndexProperty struct {
	Name   string
	Length int // 字符串属性的索引长度，0 表示不指定（非字符串类型）
}

func (p *IndexProperty) indexPropDef() string {
	if p.Length > 0 {
		return fmt.Sprintf("%s(%d)", backquoteIfNeed(p.Name), p.Length)
	}
	return fmt.Sprintf("%s", backquoteIfNeed(p.Name))
}
