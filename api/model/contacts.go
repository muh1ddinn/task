package model

type Contact struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Category string `json:"category"`
}

type GetAllContactResponse struct {
	Contact []GetAllContact `json:"contact"`
	Count   int64           `json:"count"`
}

type GetAllContactRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllContact struct {
	Id         string `json:"id"`
	Name       string `json:"Name"`
	Email      string `json:"Email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	Category   string `json:"category"`
	Created_at string `json:"createdAt"`
}

type ContactLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ContactLoginRequest struct {
	Email    string `json:"mail"`
	Password string `json:"password"`
}

type PatchContact struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Check struct {
	Email string `json:"email"`
}

type ContactHistory struct {
	ID         string `json:"id"`
	ContactID  string `json:"contact_id"`
	Phone      string `json:"phone"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	Category   string `json:"category"`
	Changed_at string `json:"changed_at"`
	ChangeType string `json:"change_type"`
}

type GetAllhistoryResponse struct {
	ContactHistory []ContactHistory `json:"contactHistory"`
	Count          int64            `json:"count"`
}
