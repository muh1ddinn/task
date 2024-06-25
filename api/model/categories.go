package model

type Categories struct {
	Name string `json:"name"`
}

type UpdateCategories struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetAllcategoriesResponse struct {
	Categories []Getcategoriest `json:"categories"`
	Count      int64            `json:"count"`
}

type GetAllCategoriestRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type Getcategoriest struct {
	Id         string `json:"id"`
	Name       string `json:"Name"`
	Created_at string `json:"created_at"`
}

type CategoriesLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Patchcategories struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
