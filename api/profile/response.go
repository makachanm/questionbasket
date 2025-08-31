package profile

type ProfileResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BasedOn     int    `json:"based_on"`
}
