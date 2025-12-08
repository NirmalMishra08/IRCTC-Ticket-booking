package util

type Role int

const (
	EM_USER Role = iota
	NON_EM_USER
	AMIGO_USER
	SCOOTER_USER
	END_OF_LINE_USER
)

func ConvertIntToRole(role Role) string {
	switch role {
	case 0:
		return "EM_USER"
	case 1:
		return "NON_EM_USER"
	case 2:
		return "AMIGO_USER"
	case 3:
		return "SCOOTER_USER"
	case 4:
		return "END_OF_LINE_USER"
	default:
		return "UNKNOWN"
	}
}

func ConvertRoleToInt(role string) int32 {
	switch role {
	case "EM_USER":
		return 0
	case "NON_EM_USER":
		return 1
	case "AMIGO_USER":
		return 2
	case "SCOOTER_USER":
		return 3
	case "END_OF_LINE_USER":
		return 4
	default:
		return -1
	}
}
