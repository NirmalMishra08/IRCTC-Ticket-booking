package tatkal

type TatkalJob struct {
	UserID string         `json:"user_id"`
	Data   TatkalRequest `json:"data"`
}
