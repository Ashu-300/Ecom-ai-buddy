package dto

type Address struct {
    Street     string `json:"street"`
    City       string `json:"city"`
    State      string `json:"state"`
    PostalCode string `json:"postal_code"`
    Country    string `json:"country"`
}

type User struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Role      string    `json:"role"`
    Addresses []Address `json:"Addresses"`
}

// Response wrapper for the API
type AuthResponse struct {
    Message  string `json:"message"`
    UserInfo User   `json:"userInfo"`
}
