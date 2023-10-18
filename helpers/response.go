package helpers

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