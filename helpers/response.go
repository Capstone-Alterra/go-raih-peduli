package helpers

type ResponseError struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
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
