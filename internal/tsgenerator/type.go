package tsgenerator

// TODO Add support for deprecated models properties and relationships

type TSModel struct {
	ClassName     string
	Properties    []TSProperty
	Relationships []TSRelationship
	ApiFilePath   string
}

type TSInterface struct {
	InterfaceName string
	Properties    []TSProperty
	ApiFilePath   string
}

type TSEnum struct {
	EnumName    string
	Values      []string
	ApiFilePath string
}

type TSProperty struct {
	Name string
	Type TSTypeToken
}

type TSRelationship struct {
	Name string
	Type APIModelToken
}

type TSTypeToken string

const (
	STRING_TOKEN        TSTypeToken = "string"
	NUMBER_TOKEN        TSTypeToken = "number"
	BOOLEAN_TOKEN       TSTypeToken = "boolean"
	BIGINT_TOKEN        TSTypeToken = "bigint"
	DATE_TOKEN          TSTypeToken = "Date"
	OBJECT_TOKEN        TSTypeToken = "Object"
	STRING_ARRAY_TOKEN  TSTypeToken = "string[]"
	NUMBER_ARRAY_TOKEN  TSTypeToken = "number[]"
	BOOLEAN_ARRAY_TOKEN TSTypeToken = "boolean[]"
	DATE_ARRAY_TOKEN    TSTypeToken = "Date[]"
	OBJECT_ARRAY_TOKEN  TSTypeToken = "Object[]"
)

type TSStructToken string

const (
	ENUM_TOKEN      TSTypeToken = "enum"
	INTERFACE_TOKEN TSTypeToken = "interface"
	TYPE_TOKEN      TSTypeToken = "type"
	CLASS_TOKEN     TSTypeToken = "class"
)

type APIModelToken string

const (
	COLUMN_TOKEN APIModelToken = "@Column"
)

type APIModelRelationshipToken string

const (
	BELONGS_TO_TOKEN APIModelRelationshipToken = "@BelongsTo"
	HAS_MANY_TOKEN   APIModelRelationshipToken = "@HasMany"
)

type APIModelConstraintToken string

const (
	AllowNullToken APIModelConstraintToken = "@AllowNull"
)
