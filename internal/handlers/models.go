package handlers

type CheckBody struct {
	Links []string `json:"links"`
}

type GetDataBody struct {
	LinksList []int `json:"links_list"`
}
