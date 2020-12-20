package schema

const (
	SystemClassUser = "_User"
)

var SystemClasses = map[string]string{SystemClassUser: SystemClassUser}

var DefaultColumn = Fields{
	FieldObjectID:        Field{Type: FieldTypeUUID},
	FieldReadPermission:  Field{Type: FieldTypeArray},
	FieldWritePermission: Field{Type: FieldTypeArray},
	FieldCreatedAt:       Field{Type: FieldTypeDate},
	FieldUpdatedAt:       Field{Type: FieldTypeDate},
}

var DefaultColumnUser = Fields{
	"username":      Field{Type: FieldTypeString, Required: true},
	"password":      Field{Type: FieldTypeString, Required: true},
	"email":         Field{Type: FieldTypeString, Required: true},
	"emailVerified": Field{Type: FieldTypeBoolean, Required: true},
	"authData":      Field{Type: FieldTypeObject, Required: true},
}

var DefaultColumns = map[string]Fields{
	"_DEFAULT":      DefaultColumn,
	SystemClassUser: DefaultColumnUser,
}
