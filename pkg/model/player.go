package model

type Player struct {
	ID    string `json:"player"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p Player) Clone() *Player {
	return &Player{
		ID:    p.ID,
		Name:  p.Name,
		Email: p.Email,
	}
}
