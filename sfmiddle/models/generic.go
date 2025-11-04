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

// Position representa las coordenadas y dimensiones de una firma

// Document representa cada documento en el JSON
type Document struct {
	ActiveDoc       string     `json:"activeDoc"`
	AuthRoleStatus  string     `json:"authRoleStatus"`
	AuthUseStatus   string     `json:"authUseStatus"`
	CreatedAtDoc    string     `json:"createdAtDoc"`
	DocumentExt     string     `json:"documentExt"`
	DocumentHash    string     `json:"documentHash"`
	DocumentName    string     `json:"documentName"`
	DocumentPath    string     `json:"documentPath"`
	DocumentClass   string     `json:"documentClass"`
	Description     string     `json:"description"`
	IdFD            string     `json:"idFD"`
	LastModifiedDoc string     `json:"lastModifiedDoc"`
	Reviewers       []Reviewer `json:"reviewers"`
}
