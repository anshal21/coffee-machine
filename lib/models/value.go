package models

// Value is a struct to hold primitive values belonging to one of the DataType
type Value struct {
	Number *float64
	String *string
	Bool   *bool
}

// DataType is a type to represent possible primitive data types in an expression
type DataType int

// Set of different datatypes
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
