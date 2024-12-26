package favorite

type OwnerType string

const (
	OwnerTypeUser  OwnerType = "USER"
	OwnerTypeGroup OwnerType = "GROUP"
)

type ObjectType string

const (
	ObjectTypeDocument ObjectType = "DOCUMENT"
	ObjectTypeImage    ObjectType = "IMAGE"
	ObjectTypeVideo    ObjectType = "VIDEO"
)

func IsValidOwnerType(ownerType string) bool {
	switch OwnerType(ownerType) {
	case OwnerTypeUser, OwnerTypeGroup:
		return true
	default:
		return false
	}
}

func IsValidObjectType(objectType string) bool {
	switch ObjectType(objectType) {
	case ObjectTypeDocument, ObjectTypeImage, ObjectTypeVideo:
		return true
	default:
		return false
	}
}
