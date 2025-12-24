package models

import "mime/multipart"

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

type UserDataReq struct {
	IdInvite string `json:"idInvite"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Alias    string `json:"alias"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	TaxNum   string `json:"rfc"`
	PobUid   string `json:"curp"`
	Role     string `json:"role"`
}

type GetUserRequest struct {
	IdUser string   `json:"iduser"`
	Fields []string `json:"fields"`
}

type GetTeamUsersRequest struct {
	IdTeam string   `json:"idteam"`
	Fields []string `json:"fields,omitempty"`
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

type ValidationRequest struct {
	StreetAddress      string `json:"streetAddress,omitempty"`
	PostalCode         string `json:"postalCode,omitempty"`
	Neighborhood       string `json:"neighborhood,omitempty"`
	Locality           string `json:"locality,omitempty"`
	ActaConstitutiva   string `json:"docActa,omitempty"`
	PoderRepresentante string `json:"docPoder,omitempty"`
	IdentidadOficial   string `json:"docIdentidad,omitempty"`
	PruebaResidencia   string `json:"docResidencia,omitempty"`
	// Campos para manejar archivos (no se serializan a JSON)
	ActaFile       multipart.File `json:"-"`
	PoderFile      multipart.File `json:"-"`
	IdentidadFile  multipart.File `json:"-"`
	ResidenciaFile multipart.File `json:"-"`
}

type PaymentReq struct {
	Plan       string `json:"plan"`
	CardNum    string `json:"cardnumber"`
	Expiration string `json:"exp"`
	CVV        string `json:"cvv"`
	NameOwner  string `json:"nameowner"`
}

type DocDataRequest struct {
	IdDocs   []string `json:"idDocs,omitempty"`
	Type     string   `json:"type,omitempty"`
	PathDocs []string `json:"paths,omitempty"`
	DateFrom string   `json:"date_from,omitempty"`
	DateTo   string   `json:"date_to,omitempty"`
	Page     int      `json:"page,omitempty"`
	PageSize int      `json:"page_size,omitempty"`
	OrderBy  string   `json:"order_by,omitempty"`  // opcional, default: createdAtDoc
	OrderDir string   `json:"order_dir,omitempty"` // ASC o DESC
}

type DocRev struct {
	IdFolder  string     `json:"idfolder"`
	Document  Document   `json:"document"`
	Reviewers []Reviewer `json:"reviewers"`
}

type InviteRequest map[string]DocRev

type GetByInviteRequest struct {
	IdInvite string `json:"idInvite"`
}

type FolderRequest struct {
	IdFolder   string `json:"idFolder,omitempty"`
	OnlyShared bool   `json:"onlyShared,omitempty"`
	//OnlyTeam   bool   `json:"onlyTeam,omitempty"`
	OnlyUser bool   `json:"onlyUser,omitempty"`
	DateFrom string `json:"dateFrom,omitempty"`
	DateTo   string `json:"dateTo,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"pageSize,omitempty"`
	OrderBy  string `json:"orderBy,omitempty"`
	OrderDir string `json:"orderDir,omitempty"`
}

type LLMrequest struct {
	Path   string `json:"path"`
	Action int    `json:"action,omitempty"`
	Prompt string `json:"prompt,omitempty"`
}

type SignDoc struct {
	Aut       string     `json:"aut"`
	IdInvite  string     `json:"inviteId"`
	IdKey     string     `json:"keyId"`
	IdFolder  string     `json:"folderId"`
	Documents []Document `json:"signDocuments"`
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
	Error      string `json:"error"`
	RedirectTo string `json:"redirectTo"`
}

// UserResponse devuelve datos del usuario (sin sensibles)
type UserResponse struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Active bool   `json:"active"`
}

type InstResponse struct {
	LegalName           string `json:"legalName"`
	AliasName           string `json:"aliasName"`
	TaxNum              string `json:"taxNum"`
	LegalSignupName     string `json:"legalSignupName"`
	LegalSignupLastname string `json:"legalSignupLastname"`
	StreetAddress       string `json:"streetAddress"`
	PostalCode          string `json:"postalCode"`
	Neighborhood        string `json:"neighborhood"`
	Locality            string `json:"locality"`
	State               string `json:"state"`
	Country             string `json:"country"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	Logo                string `json:"logoUrl"`
	IsActive            bool   `json:"isActive"`
}

type ValidationResponse struct {
	LegalName           string `json:"legalName"`
	AliasName           string `json:"aliasName"`
	TaxNum              string `json:"taxNum"`
	LegalSignupName     string `json:"legalSignupName"`
	LegalSignupLastname string `json:"legalSignupLastname"`
	StreetAddress       string `json:"streetAddress"`
	//AddressLine         string `json:"lineAddr"`
	PostalCode         string `json:"postalCode"`
	Neighborhood       string `json:"neighborhood"`
	Locality           string `json:"locality"`
	ActaConstitutiva   string `json:"docActa"`
	PoderRepresentante string `json:"docPoder"`
	IdentidadOficial   string `json:"docIdentidad"`
	PruebaResidencia   string `json:"docResidencia"`
}

// Estructura para la respuesta -> MODELS
type UploadResponse struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	DocumentIDs []string `json:"document_ids,omitempty"`
	Errors      []string `json:"errors,omitempty"`
}

type PaymentResp struct {
	CurrentSatat string `json:"current"`
	Success      bool   `json:"success"`
	Status       string `json:"status"`
	Plan         string `json:"plan"`
	Expiration   string `json:"expiration"`
}

type UserStatus struct {
	CurrentSatat string `json:"current"`
	Success      bool   `json:"success"`
	Status       string `json:"status"`
	Active       string `json:"plan"`
	Expiration   string `json:"expiration"`
}

type KeysStatus struct {
	IdKey           string `json:"idKey"`
	NameKey         string `json:"nameKey"`
	NameCer         string `json:"nameCer"`
	Expiration      string `json:"expiration"`
	Owner           string `json:"owner"`
	SubjectUniqueId string `json:"subjectUniqueId"`
	IssuerRFC4514   string `json:"issuerData"`
	UploadedAt      string `json:"uploadedAt"`
	Selected        bool   `json:"selected"`
}

type UserDataResp struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	LastName     string `json:"lastname"`
	Alias        string `json:"alias"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Active       string `json:"active"`
	Role         string `json:"role"`
	SignedDocs   int    `json:"signed,omitempty"`
	TotalSigns   int    `json:"totalsigns,omitempty"`
	Kyc          string `json:"kyc,omitempty"`
	IsAlive      bool   `json:"isAlive"`
	CreatedAt    string `json:"createdat"`
	DeletedAt    string `json:"deletedat"`
	LastModified string `json:"lastmodified"`
}

type TeamDataResp struct {
	Id           string `json:"id"`
	CreatorUser  string `json:"creator"`
	Name         string `json:"name"`
	LimitSigners string `json:"limitsigners"`
	LimitUsers   string `json:"limitusers"`
	DeletedAt    string `json:"deletedAt"`
	Description  string `json:"description"`
}

type InviteInfoResp struct {
	Folder FolderInvite `json:"folder"`
}

type SignResponse struct {
	Message string                       `json:"message"`
	Signed  map[string]map[string]string `json:"signed"`
}
type FolderListResp struct {
	Page        int                    `json:"page"`
	PageSize    int                    `json:"page_size"`
	Total       int                    `json:"total"`
	FoldersData map[string]interface{} `json:"folders"`
}

type LLMresp struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Date    string `json:"date"`
}
