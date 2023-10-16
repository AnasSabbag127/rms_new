package model

type Restaurant struct{
	Id int		`json:"id" db:"id"`
	Name string	 `json:"name" db:"name"`
	Address string `json:"address" db:"address"`
	CreatedBy int	`json:"createdBy" db:"created_by"`
	Stars int `json:"stars" db:"stars"`
}

type InputRestraunt struct{
	Name string	 `json:"name" db:"name"`
	Address string `json:"address" db:"address"`
	CreatedBy int	`json:"createdBy" db:"created_by"`
	Stars int `json:"stars" db:"stars"`
}
