package util

type Level int

const (
	SMS Level = iota
	EMOTORAD_USER
	VIEWER
	ANALYST
	MODERATOR_1
	MODERATOR_2
	MODERATOR_3
	MODERATOR_4
	ADMIN
	OWNER
	EMOTORAD_ADMIN
	BACKEND_ADMIN
)

func ConvertIntToLevel(level Level) string {
	switch level {
	case 0:
		return "SMS"
	case 1:
		return "EMOTORAD_USER"
	case 2:
		return "VIEWER"
	case 3:
		return "ANALYST"
	case 4:
		return "MODERATOR_1"
	case 5:
		return "MODERATOR_2"
	case 6:
		return "MODERATOR_3"
	case 7:
		return "MODERATOR_4"
	case 8:
		return "ADMIN"
	case 9:
		return "OWNER"
	case 10:
		return "EMOTORAD_ADMIN"
	case 11:
		return "BACKEND_ADMIN"
	default:
		return "UNKNOWN"
	}
}

func ConvertLevelToInt(level string) int32 {
	switch level {
	case "SMS":
		return 0
	case "EMOTORAD_USER":
		return 1
	case "VIEWER":
		return 2
	case "ANALYST":
		return 3
	case "MODERATOR_1":
		return 4
	case "MODERATOR_2":
		return 5
	case "MODERATOR_3":
		return 6
	case "MODERATOR_4":
		return 7
	case "ADMIN":
		return 8
	case "OWNER":
		return 9
	case "EMOTORAD_ADMIN":
		return 10
	case "BACKEND_ADMIN":
		return 11
	default:
		return -1
	}
}
