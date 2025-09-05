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
