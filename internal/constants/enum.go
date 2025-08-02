package constants

const (
	UserTrung = "Trung"
	UserThang = "Thang"
)

type User int

const (
	UserEnumTrung User = iota
	UserEnumThang
)

func (u User) String() string {
	switch u {
	case UserEnumTrung:
		return UserTrung
	case UserEnumThang:
		return UserThang
	default:
		return "Unknown"
	}
}
