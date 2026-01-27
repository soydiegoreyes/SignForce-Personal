package documentflow

import (
	"fmt"
	"sfmiddle/db"
	"sfmiddle/models"
)

func LoadFolderInfo(idFolder, idUser string, onlyShared, onlyUser bool, page, pageSize int, orderBy, orderDir string) (*models.FolderListResp, error) {

	// Configuración de paginación
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if orderBy == "" {
		orderBy = "creationAt"
	}
	if orderDir != "ASC" && orderDir != "DESC" {
		orderDir = "DESC"
	}

	logic := ""

	// Obtenemos datos del folder
	attrs := []string{"ownerInst_fk", "creatorUser_fk", "creationAt", "lastModified",
		"deletedAt", "deletedReason", "purpose", "description",
		"closedAt", "secuentialSign", "folderBasePath", "expirationDate",
		"numDocs", "numDocsSign", "numSigners", "numReceivers", "completedAt"}

	wheres := map[string][]string{}

	var docWheres map[string][]string // usado para obtener los documentos de idfolderdocuments

	// consulta a invitaciones del usuario para obtener el folder
	if onlyShared {
		// Se obtienen los folders de las invitaciones que ha recibido el usuario ya que se cargan todos los folders en la pantalla
		idsFoldersShared, err := db.DB_con.GenericSelect("invites", "idInvite", []string{"idFolder"}, map[string][]string{"idUserDest_fk": {idUser}})
		if err != nil {
			return nil, fmt.Errorf("error obteniendo datos de invitación: %v", err)
		}
		if len(idsFoldersShared) == 0 {
			return &models.FolderListResp{
				Page:        page,
				PageSize:    pageSize,
				Total:       0,
				FoldersData: make(map[string]interface{}),
			}, nil
		}
		wheresInvDet := map[string][]string{}
		for idInv := range idsFoldersShared {
			wheres["idFolder"] = append(wheres["idFolder"], idsFoldersShared[idInv]["idFolder"])
			wheresInvDet["idInvite"] = append(wheresInvDet["idInvite"], idInv)
		}
		// este mapa tiene la relacion de documentos que se pueden obtener del folder folderdocuments
		idsInvitesdetail, err := db.DB_con.GenericSelect("invitesdetail", "idInviteDetail", []string{"idfolderdocument"}, wheresInvDet)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo datos de invitación: %v", err)
		}

		// filtro para obtener solo el idDocument de cada folder que se obtuvo en el paso previo (WHERE idFolder IN (...) AND idfolderdocument IN (...))
		docWheres = map[string][]string{"idFolder": {}, "idfolderdocument": {}}
		for _, invDetData := range idsInvitesdetail {
			docWheres["idfolderdocument"] = append(docWheres["idfolderdocument"], invDetData["idfolderdocument"])
		}
		docWheres["LOGIC"] = []string{"idFolder AND idfolderdocument"}

		logic = "idFolder"

	} else if onlyUser {
		wheres["creatorUser_fk"] = []string{idUser}
		docWheres = map[string][]string{"idFolder": {}}
		docWheres["LOGIC"] = []string{"idFolder"}
		logic = "creatorUser_fk"

	} else {
		if idFolder != "" {
			wheres["idFolder"] = []string{idFolder}
			logic = "idFolder"

		} else {
			return &models.FolderListResp{}, fmt.Errorf("no hay parámetros de búsqueda")
		}
	}

	logic = fmt.Sprintf("%s ORDER BY %s %s", logic, orderBy, orderDir)
	wheres["LOGIC"] = []string{logic, fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)}

	// Se obtienen los datos de todos los folders
	foldersData, err := db.DB_con.GenericSelect("folders", "idFolder", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos de folders: %v", err)
	}

	if len(foldersData) == 0 {
		return &models.FolderListResp{
			Page:        page,
			PageSize:    pageSize,
			Total:       0,
			FoldersData: make(map[string]interface{}),
		}, nil
	}

	// Obtener documentos de cada folder
	docAttrs := []string{"idDocument", "idFolder"}

	for folderID := range foldersData {
		docWheres["idFolder"] = append(docWheres["idFolder"], folderID)
	}

	folderDocs, err := db.DB_con.GenericSelect("folderdocuments", "idfolderdocument", docAttrs, docWheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo documentos del folder: %v", err)
	}

	// Mapear documentos por folder
	docsByFolder := make(map[string][]models.Document)
	for _, docData := range folderDocs {
		folderID := docData["idFolder"]
		docID := docData["idDocument"]
		docsByFolder[folderID] = append(docsByFolder[folderID], models.Document{IdDocument: docID})
	}

	// Construir el map de folders estructurados
	foldersMap := make(map[string]interface{})
	for folderID, folderData := range foldersData {
		folder := models.Folder{
			IdFolder:       folderID,
			IsSecuential:   folderData["secuentialSign"] == "1",
			ExpirationDate: folderData["expirationDate"],
			ClosedAt:       folderData["closedAt"],
			CompletedAt:    folderData["completedAt"],
			DeletedAt:      folderData["deletedAt"],
			DeletedReason:  folderData["deletedReason"],
			Description:    folderData["description"],
			LastModified:   folderData["lastModified"],
			NnumDocs:       folderData["numDocs"],
			NumDocsSign:    folderData["numDocsSign"],
			NumReceivers:   folderData["numReceivers"],
			NumSigners:     folderData["numSigners"],
			Path:           folderData["folderBasePath"],
			Purpose:        folderData["purpose"],
			Documents:      docsByFolder[folderID],
		}

		// Si no hay documentos para este folder, inicializar como slice vacío
		if folder.Documents == nil {
			folder.Documents = []models.Document{}
		}

		foldersMap[folderID] = folder
	}

	return &models.FolderListResp{
		Page:        page,
		PageSize:    pageSize,
		Total:       len(foldersMap),
		FoldersData: foldersMap,
	}, nil
}
