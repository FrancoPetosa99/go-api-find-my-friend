package pagination

type FilterPet struct {
	UserID        *int    `json:"user_id"`
	Type          *string `json:"type"`
	Breed         *string `json:"breed"`
	LastSeenPlace *string `json:"last_seen_place"`
}
