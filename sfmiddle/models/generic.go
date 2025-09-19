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
