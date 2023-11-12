package dtos

import "time"

type ResFundraise struct {
	Title 		string `json:"title"`
	Description string `json:"description"`
	Photo	  	string `json:"photo"`
	Target    	int32  `json:"target"`
	StartDate 	time.Time `json:"start_date"`
	EndDate   	time.Time `json:"end_date"`
	Status		string `json:"status"`
	UserID		int	   `json:"user_id"`
}
