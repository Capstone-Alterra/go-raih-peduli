package helpers

import "strings"

/*
func BuildLikeQuery(column string, keywords []string) string {
	var likeQuery string
	for i, keyword := range keywords {
		if i > 0 {
			likeQuery += " OR "
		}
		keyword = strings.ReplaceAll(keyword, " ", "")
		likeQuery += column + " LIKE '%" + keyword + "%'"
	}
	return "(" + likeQuery + ")"
}

func BuildNotLikeQuery(column string, keywords []string) string {
	var likeQuery string
	for i, keyword := range keywords {
		if i > 0 {
			likeQuery += " AND "
		}
		keyword = strings.ReplaceAll(keyword, " ", "")
		likeQuery += column + " NOT LIKE '%" + keyword + "%'"
	}
	return "(" + likeQuery + ")"
}

*/

func BuildLikeQuery(column string, keywords []string) string {
	var likeQuery string
	for i, keyword := range keywords {
		if i > 0 {
			likeQuery += " OR "
		}
		keyword = strings.ReplaceAll(keyword, " ", "")
		likeQuery += "title " + " LIKE '%" + keyword + "%'"
	}
	for _, keyword := range keywords {
		likeQuery += " OR "
		keyword = strings.ReplaceAll(keyword, " ", "")
		likeQuery += "description" + " LIKE '%" + keyword + "%'"
	}
	return "(" + likeQuery + ")"
}

func BuildNotLikeQuery(column string, keywords []string) string {
	var likeQuery string
	for i, keyword := range keywords {
		if i > 0 {
			likeQuery += " AND "
		}
		keyword = strings.ReplaceAll(keyword, " ", "")
		likeQuery += "title" + " NOT LIKE '%" + keyword + "%'"
	}
	for _, keyword := range keywords {
		likeQuery += " AND "
		keyword = strings.ReplaceAll(keyword, " ", "")
		likeQuery += "description" + " NOT LIKE '%" + keyword + "%'"
	}
	return "(" + likeQuery + ")"
}
