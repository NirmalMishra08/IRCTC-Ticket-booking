// types.go
package train

type CoachType string
type BerthType string

const (
	CoachType3A CoachType = "3A"
	CoachType2A CoachType = "2A"
	CoachType1A CoachType = "1A"
	CoachTypeSL CoachType = "SL"
	CoachTypeGN CoachType = "GN"
)

const (
	BerthUpper  BerthType = "UP"
	BerthLower  BerthType = "DOWN"
	BerthMiddle BerthType = "MID"
)

// SeatConfiguration defines how seats should be arranged in a coach
type SeatConfiguration struct {
	CoachType       string `json:"coach_type" validate:"required,oneof=3A 2A 1A SL GN"`
	NumberOfCoaches int    `json:"number_of_coaches" validate:"required,min=1"`
	SeatsPerCoach   int    `json:"seats_per_coach" validate:"required,min=1"`
	// Berth configuration (only applicable for sleeper coaches)
	LowerBerthRatio     float64 `json:"lower_berth_ratio" validate:"min=0,max=1"`
	MiddleBerthRatio    float64 `json:"middle_berth_ratio" validate:"min=0,max=1"`
	UpperBerthRatio     float64 `json:"upper_berth_ratio" validate:"min=0,max=1"`
	SideUpperBerthRatio float64 `json:"side_upper_berth_ratio" validate:"min=0,max=1"`
	SideLowerBerthRatio float64 `json:"side_lower_berth_ratio" validate:"min=0,max=1"`
}

type CreateCoachSeatRequest struct {
	TrainID        int                 `json:"train_id" validate:"required"`
	Configurations []SeatConfiguration `json:"configurations" validate:"required,min=1,dive"`
}

type CoachSeatResponse struct {
	TrainID        int                `json:"train_id"`
	TrainNumber    int                `json:"train_number"`
	TrainName      string             `json:"train_name"`
	TotalCoaches   int                `json:"total_coaches"`
	TotalSeats     int                `json:"total_seats"`
	CoachBreakdown []CoachTypeSummary `json:"coach_breakdown"`
	Message        string             `json:"message"`
}

type CoachTypeSummary struct {
	CoachType       string `json:"coach_type"`
	NumberOfCoaches int    `json:"number_of_coaches"`
	SeatsPerCoach   int    `json:"seats_per_coach"`
	TotalSeats      int    `json:"total_seats"`
}

type BerthAllocation struct {
	BerthType BerthType
	Count     int
}
