package model

type Dishes struct {
	Id int `json:"id" db:"id"`
	RestaurantId int `json:"restrauntId" db:"restraunt_id"`
	CreatedBy int 	`json:"createdBy" db:"created_by"`
	DishName string	`json:"name" db:"name"`
	Price int `json:"price" db:"price"`
}

type InputDishes struct{	
	RestaurantId int `json:"restrauntId" db:"restraunt_id"`
	CreatedBy int 	`json:"createdBy" db:"created_by"`
	DishName string	`json:"name" db:"name"`
	Price int `json:"price" db:"price"`
}