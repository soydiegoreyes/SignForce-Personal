package models

// Estructura para metadatos de archivo
type FileMetadata struct {
	Id      int
	DocType string
	Name    string
	Ext     string
	Size    int64
	Hash    string
	Path    string
	Ok      bool
}

// Estructura para metadatos de llaves enviadas al back
type KeyMetadata struct {
	IdInst  string
	IdUser  string
	PassKey string
	NameKey string
	NameCer string
	HashKey string
	HashCer string
	Path    string
}

// Position representa las coordenadas y dimensiones de una firma
type SignPosition struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Page   int     `json:"page"`
}

// Reviewer representa un firmante o revisor del documento
type Reviewer struct {
	User          string         `json:"user"`
	Role          int            `json:"role"`
	DueDate       string         `json:"due_date"`
	Team          string         `json:"team"`
	Comment       string         `json:"comment"`
	SignPositions []SignPosition `json:"positions"`
}

// Document representa cada documento en el JSON
type Document struct {
	IdDocument       string     `json:"idDocument"`
	ActiveDoc        string     `json:"activeDoc"`
	AuthRoleStatus   string     `json:"authRoleStatus"`
	AuthUseStatus    string     `json:"authUseStatus"`
	CreatedAtDoc     string     `json:"createdAtDoc"`
	DocumentExt      string     `json:"documentExt"`
	DocumentHash     string     `json:"documentHash"`
	DocumentName     string     `json:"documentName"`
	DocumentPath     string     `json:"documentPath"`
	DocumentClass    string     `json:"documentClass"`
	DocumentFullName string     `json:"documentFullName"`
	Abstract         string     `json:"abstract"`
	IdFolder         string     `json:"idFolder"`
	LastModifiedDoc  string     `json:"lastModifiedDoc"`
	Reviewers        []Reviewer `json:"reviewers"`
}

// datos para llenar la plantilla de invitacion por email
type InviteMail struct {
	IdUser        string
	ReviewerName  string
	ReviewerEmail string
	ReviewerInst  string
	SentDate      string
	SenderMessage string
	UrlSignLink   string
}

type Folder struct {
	IdFolder       string       `json:"idFolder"`
	IsSecuential   bool         `json:"isSecuential"`
	ExpirationDate string       `json:"expirationDate"`
	UserEmisor     UserInfo     `json:"userEmisor"`
	UserDest       UserDestInfo `json:"userDest"`
	Invites        []Invite     `json:"invites"`
}

type UserInfo struct {
	NameInstEmisor string `json:"nameInstEmisor"`
	NameUserEmisor string `json:"nameUserEmisor"`
	NameTeamEmisor string `json:"nameTeamEmisor"`
}

type UserDestInfo struct {
	NameInstDest string        `json:"nameInstDest"`
	NameUserDest string        `json:"nameUserDest"`
	NameTeamDest string        `json:"nameTeamDest"`
	Keys         []*KeysStatus `json:"keys"`
}

type Invite struct {
	IdInvite      string   `json:"idInvite"`
	AliveProof    bool     `json:"aliveProof"`
	IsSigner      bool     `json:"isSigner"`
	SentAt        string   `json:"sentAt"`
	MessageEmisor string   `json:"messageEmisor"`
	Document      Document `json:"document"`
}
