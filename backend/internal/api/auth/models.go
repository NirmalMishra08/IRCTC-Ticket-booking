package auth

type authRequest struct {
	Provider string `json:"provider"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

type authResponse struct {
	UserID   string `json:"user_id"`
	Provider string `json:"provider"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	// Phone            string `json:"phone"`
	// ActiveMembership bool `json:"active_membership"`
}
