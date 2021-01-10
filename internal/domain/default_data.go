package domain

const (
	SystemClassUser = "_User"
)

var SystemClasses = map[string]string{SystemClassUser: SystemClassUser}

var DefaultField = Fields{
	FieldObjectID:        Field{Type: FieldTypeUUID},
	FieldReadPermission:  Field{Type: FieldTypeArray},
	FieldWritePermission: Field{Type: FieldTypeArray},
	FieldCreatedAt:       Field{Type: FieldTypeDate},
	FieldUpdatedAt:       Field{Type: FieldTypeDate},
}

var DefaultFieldUser = Fields{
	"username":      Field{Type: FieldTypeString, Required: true},
	"password":      Field{Type: FieldTypeString, Required: true},
	"email":         Field{Type: FieldTypeString, Required: true},
	"emailVerified": Field{Type: FieldTypeBoolean, Required: true},
	"authData":      Field{Type: FieldTypeObject, Required: true},
}

var DefaultFields = map[string]Fields{
	"_DEFAULT":      DefaultField,
	SystemClassUser: DefaultFieldUser,
}
