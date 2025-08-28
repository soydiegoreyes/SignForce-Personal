package models

// ==================== REQUESTS ======================//
type RegisterRequest struct {
	RegisterSignupName     string `json:"firstName"`
	RegisterSignupLastname string `json:"lastName"`
	LegalSignupName        string `json:"legalName"`
	LegalSignupLastname    string `json:"legalLastName"`
	NameInst               string `json:"businessName"`
	AliasNameInst          string `json:"businessAlias"`
	TaxNumInst             string `json:"businessTaxNum"`
	ContactEmailInst       string `json:"email"`
	ContactPhoneInst       string `json:"phone"`
	City                   string `json:"city"`
	Country                string `json:"country"`
}

type GetUserRequest struct {
	IdUser string   `json:"iduser"`
	Fields []string `json:"fields"`
}

// Decodificar el cuerpo (ej: {"id": "123", "password": "secret"})
type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type EmailRequest struct {
	IdUser   string   `json:"iduser"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
	Dest     []string `json:"dest"`
	MimeType string   `json:"mimetype"`
}

// ==================== RESPONSES ====================//
// RegisterResponse estructura para respuesta a register
type RegisterResponse struct {
	Check  bool   `json:"check"`
	InstId string `json:"instId"`
	Error  string `json:"error"`
}

// LoginResponse estructura para la respuesta del login
type LoginResponse struct {
	//Token      string `json:"token"`
	Error      string `json:"error"`
	RedirectTo string `json:"redirectTo"`
}

// UserResponse devuelve datos del usuario (sin sensibles)
type UserResponse struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

type ValidationResponse struct {
	LegalName           string `json:"legalName"`
	AliasName           string `json:"aliasName"`
	TaxNum              string `json:"taxNum"`
	LegalSignupName     string `json:"legalSignupName"`
	LegalSignupLastname string `json:"legalSignupLastname"`
	StreetAddress       string `json:"streetAddr"`
	AddressLine         string `json:"lineAddr"`
	PostalCode          string `json:"postalCode"`
	Neighborhood        string `json:"neighborhood"`
	Locality            string `json:"locality"`
	ActaConstitutiva    string `json:"docActa"`
	PoderRepresentante  string `json:"docPoder"`
	IdentidadOficial    string `json:"docIdentidad"`
	PruebaResidencia    string `json:"docResidencia"`
}
