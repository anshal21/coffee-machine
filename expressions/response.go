package expressions

type Response struct {
	Value Value
	Type  DataType
}

type DataType int

const (
	DataTypeUnknown DataType = iota
	DataTypeBool
	DataTypeNumber
	DataTypeString
)

func (d DataType) String() string {
	switch d {
	case DataTypeBool:
		return "bool"
	case DataTypeNumber:
		return "number"
	case DataTypeString:
		return "string"
	default:
		return "unknown"
	}
}
