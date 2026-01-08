package models

// Estructura para metadatos de archivo
type FileMetadata struct {
	Id      int    `json:"id"`
	DocType string `json:"documentType"`
	Name    string `json:"name"`
	Ext     string `json:"ext"`
	Size    int64  `json:"size"`
	Hash    string `json:"hash"`
	Path    string `json:"path"`
	Ok      bool   `json:"check"`
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
	IsExternal    bool           `json:"external"`
	User          string         `json:"user"`
	Role          int            `json:"role"`
	AliveProof    bool           `json:"alive_proof"`
	DueDate       string         `json:"due_date"`
	Comment       string         `json:"comment"`
	SignPositions []SignPosition `json:"positions"`
}

// Document representa cada documento en el JSON
type Document struct {
	IdDocument       string `json:"idDocument"`
	ActiveDoc        string `json:"activeDoc"`
	AuthRoleStatus   string `json:"authRoleStatus"`
	AuthUseStatus    string `json:"authUseStatus"`
	CreatedAtDoc     string `json:"createdAtDoc"`
	DocumentExt      string `json:"documentExt"`
	DocumentHash     string `json:"documentHash"`
	DocumentName     string `json:"documentName"`
	DocumentPath     string `json:"documentPath"`
	DocumentClass    string `json:"documentClass"`
	DocumentFullName string `json:"documentFullName"`
	Abstract         string `json:"abstract"`
	LastModifiedDoc  string `json:"lastModifiedDoc"`
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

type FolderInvite struct {
	IdFolder     string   `json:"idFolder"`
	IsSecuential bool     `json:"isSecuential"`
	UserEmisor   UserInfo `json:"userEmisor"`
	Invites      []Invite `json:"invites"`
}

type Folder struct {
	IdFolder       string       `json:"idFolder"`
	IsSecuential   bool         `json:"isSecuential"`
	ExpirationDate string       `json:"expirationDate"`
	UserEmisor     UserInfo     `json:"userEmisor"`
	UserDest       UserDestInfo `json:"userDest"`
	ClosedAt       string       `json:"closedAt"`
	CompletedAt    string       `json:"completedAt"`
	DeletedAt      string       `json:"deletedAt"`
	DeletedReason  string       `json:"deletedReason"`
	Description    string       `json:"description"`
	LastModified   string       `json:"lastModified"`
	NnumDocs       string       `json:"numDocs"`
	NumDocsSign    string       `json:"numDocsSign"`
	NumReceivers   string       `json:"numReceivers"`
	NumSigners     string       `json:"numSigners"`
	Path           string       `json:"path"`
	Purpose        string       `json:"purpose"`
	Documents      []Document   `json:"documents"`
}

type UserInfo struct {
	NameInstEmisor string `json:"nameInstEmisor"`
	NameUserEmisor string `json:"nameUserEmisor"`
	//NameTeamEmisor string `json:"nameTeamEmisor"`
}

type UserDestInfo struct {
	NameInstDest string `json:"nameInstDest"`
	NameUserDest string `json:"nameUserDest"`
	//NameTeamDest string        `json:"nameTeamDest"`
	Keys []*KeysStatus `json:"keys"`
}

type Invite struct {
	IdInvite   string      `json:"idInvite"`
	UserDest   Reviewer    `json:"userDest"`
	SentAt     string      `json:"sentAt"`
	InviteDocs []InviteDoc `json:"invitedocs"`
}

type InviteDoc struct {
	Doc       Document `json:"document"`
	ForSign   bool     `json:"forSign"`
	ExpiresAt string   `json:"expirationDate"`
	Order     int      `json:"order"`
	Comment   string   `json:"comment"`
}

type InviteUser struct {
	IdUserDest string `json:"idInvitado,omitempty"`
	EmailDest  string `json:"emailInvitado,omitempty"`
	RoleApp    string `json:"role,omitempty"`
}
