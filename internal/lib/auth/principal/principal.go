package principal

// Principal is current entity, parsed from jwt token
type Principal struct {
	UserID string
	Role   string
}
