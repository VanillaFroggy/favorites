package favorite

type OwnerType string

const (
	OwnerTypeUser  OwnerType = "user"
	OwnerTypeGroup OwnerType = "group"
)

type ObjectType string

const (
	ObjectTypeDocument ObjectType = "document"
	ObjectTypeImage    ObjectType = "image"
	ObjectTypeVideo    ObjectType = "video"
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
