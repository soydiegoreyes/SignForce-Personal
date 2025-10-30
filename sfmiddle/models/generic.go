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
	IdFD            string     `json:"idFD"`
	LastModifiedDoc string     `json:"lastModifiedDoc"`
	Reviewers       []Reviewer `json:"reviewers"`
}

type InviteRequest map[string]Document

/*

{
  "21": {
    "activeDoc": "1",
    "authRoleStatus": "1",
    "authUseStatus": "1",
    "createdAtDoc": "2025-10-23T18:03:56Z",
    "documentExt": "pdf",
    "documentHash": "xGuJQvdKRhjRKoy5imPHL3I8UY2Vli7+ESLGjWcHu7A=",
    "documentName": "Acta Matrimonio",
    "documentPath": "uploaded/documents/2/4/2/",
    "idFD": "f41ab4a4-2d62-45aa-8ad7-c472f18260a0",
    "lastModifiedDoc": "2025-10-23T18:03:56Z",
    "reviewers": [
      {
        "user": "2",
        "role": "Signer",
        "due_date": "2025-10-30",
        "team": "1",
		"inst": "4",
        "comment": "Firme por favor",
		"positions": [
			{
				"page": 1,
				"x": 76.60537190082644,
				"y": 107.87959866220736,
				"width": 168.59504132231405,
				"height": 50.5685618729097
			}
		]
      },
      {
        "user": "3",
        "role": "Signer",
        "due_date": "2025-10-30",
        "team": "2",
		"inst": "4",
        "comment": "Favor de firmar",
		"positions": [
			{
				"page": 1,
				"x": 74.9194214876033,
				"y": 178.67558528428094,
				"width": 168.59504132231405,
				"height": 50.5685618729097

			}
		]
      }
    ],
  },
  "22": {
    "activeDoc": "1",
    "authRoleStatus": "1",
    "authUseStatus": "1",
    "createdAtDoc": "2025-10-23T18:34:21Z",
    "documentExt": "pdf",
    "documentHash": "AS3AmpsLCWqCAp0LzI5rBq2YCYLUDYVHfG1xQz1U16k=",
    "documentName": "ComprobanteDomicilio",
    "documentPath": "uploaded/documents/2/4/2/",
    "idFD": "02f5f106-058d-4a55-a92a-de37dae795b5",
    "lastModifiedDoc": "2025-10-23T18:34:21Z",
    "reviewers": [
      {
        "user": "5",
        "role": "Signer",
        "due_date": "2025-10-30",
        "team": "3",
		"inst": "6",
        "comment": "Firme este documento"
      },
      {
        "user": "7",
        "role": "Signer",
        "due_date": "2025-10-30",
        "team": "4",
		"inst": "6",
        "comment": "Firme"
      }
    ],
    "signPositions": [
      {
        "user": "5",
        "page": 1,
        "position": {
          "x": 130.94792967044214,
          "y": 325.7126361204013,
          "width": 133.18502183406113,
          "height": 39.88317993311037
        }
      },
      {
        "user": "7",
        "page": 1,
        "position": {
          "x": 289.4381056529749,
          "y": 327.042075451505,
          "width": 133.18502183406113,
          "height": 39.88317993311037
        }
      }
    ]
  },
  "23": {
    "activeDoc": "1",
    "authRoleStatus": "1",
    "authUseStatus": "1",
    "createdAtDoc": "2025-10-24T15:02:35Z",
    "documentExt": "pdf",
    "documentHash": "3F/kMZlWS15PwMqeDNeJB54gR32c9d/9U3ANpILbIC4=",
    "documentName": "AcuseCita",
    "documentPath": "uploaded/documents/2/4/2/",
    "idFD": "4f3c5e5f-718d-4151-8aea-0da426065593",
    "lastModifiedDoc": "2025-10-24T15:02:35Z",
    "reviewers": [
      {
        "user": "11",
        "role": "Signer",
        "due_date": "2025-10-30",
        "team": "1",
		"inst": "7",
        "comment": "Tiene que firmar ahora"
      },
      {
        "user": "21",
        "role": "Viewer",
        "due_date": "2025-10-30",
        "team": "3",
		"inst": "7",
        "comment": "Favor de firmar"
      }
    ],
    "signPositions": [
      {
        "user": "11",
        "page": 1,
        "position": {
          "x": 63.5119724025974,
          "y": 662.2073578595317,
          "width": 132.46753246753246,
          "height": 39.7324414715719
        }
      },
      {
        "user": "21",
        "page": 1,
        "position": {
          "x": 431.7717126623377,
          "y": 670.1538461538462,
          "width": 132.46753246753246,
          "height": 39.7324414715719
        }
      }
    ]
  }
}



amigo tengo este json que es una peticion para una api de firma de documentos que se traduce en un folder que contiene documentos los cuales tienen usuarios que fungen como firmantes o revisores (en todo caso ambos aprueban o desaprueban)

donde el id principal es el id del documento a firmar y idFD es el id de la tabla pivote entre idFolder y idDocumento, solo se hace una firma por usuario pero esa firma puede quedar estampada en una o muchas paginas del mismo documento y eso queda en una tabla donde se ponen sus coordenadas y dimensiones y la pagina que representa. me gustaría que generaras una estructura para golang para volcar el json en ella y que me permita tener acceso para hacer inserts en estas 3 tablas con sus respectivos campos:



invites:





"idFolder_fk" (de tabla pivote usando el idFD)



"idDocument_fk" (cada indice del json)



"idUserIssuer_fk" (se obtiene de token)



"idInstitution_fk" (se obtiene de token)



"idUser_fk" (de reviewers user)



"idTeam_fk" (de reviewers team)



"expirationDate" (de reviewers due_date)



"signerViewer" (de reviewers role)



"descriptionText" (de reviewers  comment)



signatures:





idInvite (una vez creado el invite se obtendrá este dato)



idFolder_fk (de tabla pivote usando el idFD)



digestValueSign (documentHash)



idInstitution_fk (se consultará en db usando el id del reviewer y viendo su institucion)



idUser_fk (de reviewers user)



idTeam_fk (de reviewers team)



signstamps:





idSignature (una vez creado el registro en signature)



xSign (de signPositions)



ySign (de signPositions)



wSign (de signPositions)



hSign (de signPositions)



pageSign (de signPositions)
*/
