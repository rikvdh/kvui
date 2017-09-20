package types

type KVType int

const (
	KVTypeInvalid KVType = -1
	KVTypeString  KVType = 0
	KVTypeMap     KVType = 1
	KVTypeList    KVType = 2
)

func (k KVType) String() string {
	switch k {
	case KVTypeString:
		return "string"
	case KVTypeMap:
		return "map"
	case KVTypeList:
		return "list"
	}
	return "<invalid>"
}
