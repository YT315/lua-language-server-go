package analysis

//唯一类型
var typeNil = &TypeNil{}
var typeTrue = &TypeBool{Value: true}
var typeFalse = &TypeBool{Value: false}
var typeAny = &TypeAny{}

//TypeInfo 类型接口
type TypeInfo interface {
	TypeName() string //类型名称
}

//TypeNil 空类型
type TypeNil struct {
}

//TypeName 类型名称
func (*TypeNil) TypeName() string {
	return "nil"
}

//TypeBool 布尔类型
type TypeBool struct {
	Value bool
}

//TypeName 类型名称
func (*TypeBool) TypeName() string {
	return "bool"
}

//TypeNumber 数字类型
type TypeNumber struct {
	Value float64
}

//TypeName 类型名称
func (*TypeNumber) TypeName() string {
	return "number"
}

//TypeString 字符串类型
type TypeString struct {
	Value string
}

//TypeName 类型名称
func (*TypeString) TypeName() string {
	return "string"
}

//Typelabel 标签类型
type Typelabel struct {
	Value string
}

//TypeName 类型名称
func (*Typelabel) TypeName() string {
	return "label"
}

//TypeAny 任意类型
type TypeAny struct{}

//TypeName 类型名称
func (*TypeAny) TypeName() string {
	return "any"
}

//TypeTable 表类型
type TypeTable struct {
	Fields    map[string]*SymbolInfo  //hash
	Items     map[float64]*SymbolInfo //array
	Metatable *TypeTable              //元表
}

//TypeName 类型名称
func (me *TypeTable) TypeName() string {
	return "table"
}

//TypeFunction 函数类型
type TypeFunction struct {
	Returns [][]TypeInfo //函数可能有多种返回值情况,数字第一索引表示,返回值索引,第二索引,表示此返回值的类型范围
}

//TypeName 类型名称
func (me *TypeFunction) TypeName() string {
	return "function"
}
