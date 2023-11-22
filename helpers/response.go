package helpers

type ResponseError struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type Pagination struct {
	TotalData    int `json:"total_data"`
	CurrentPage  int `json:"current_page"`
	NextPage     int `json:"next_page"`
	PreviousPage int `json:"previous_page"`
	PageSize     int `json:"page_size"`
	TotalPage    int `json:"total_page"`
}

func Response(message string, datas ...map[string]any) map[string]any {

	var res = map[string]any{
		"message": message,
	}

	if len(datas) > 0 {
		for _, data := range datas {
			for key, value := range data {
				res[key] = value
			}
		}
	}

	return res
}

func BuildErrorResponse(message string) ResponseError {
	res := ResponseError{
		Status:  false,
		Message: message,
	}
	return res
}

func PaginationResponse(page int, pageSize int, totalData int) Pagination {
	var pagination Pagination

	if pageSize >= totalData {
		pagination.PreviousPage = 0
		pagination.NextPage = 0
	} else {
		pagination.PreviousPage = max(page-1, -1)
		pagination.NextPage = min(page+1, (totalData+pageSize-1)/pageSize)
	}

	pagination.TotalData = totalData
	pagination.CurrentPage = page
	pagination.TotalPage = (totalData + pageSize - 1) / pageSize
	pagination.PageSize = pageSize

	return pagination
}
