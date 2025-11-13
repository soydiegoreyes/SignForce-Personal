package documentflow

import (
	"fmt"
	"sfmiddle/db"
	"sfmiddle/models"
)

func LoadFolderInfos(idFolder, idTeam, idUser string, onlyShared, onlyTeam, onlyUser bool, page, pageSize int, orderBy, orderDir string) (*models.FolderListResp, error) {

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

	// obtenemos datos del folder (compartidos con usuario, solo del equipo, solo creados por usuario)
	attrs := []string{"ownerInst_fk", "ownerInst_fk", "creatorUser_fk", "creationAt", "lastModified", "deletedAt", "deletedReason", "purpose", "description",
		"closedAt", "secuentialSign", "pathSerialized", "expirationDate", "numDocs", "numDocsSign", "numSigners", "numReceivers", "completedAt"}

	wheres := map[string][]string{}

	if onlyShared {
		// Se obtienen los folders de las invitaciones que ha recibido el usuario
		idsFoldersShared, err := db.DB_con.GenericSelect("invites", "idInvite", []string{"idFolder"}, map[string][]string{"idUserDest_fk": {idUser}})
		if err != nil {
			return nil, fmt.Errorf("error obteniendo datos de invitación: %v", err)
		}
		if len(idsFoldersShared) == 0 {
			return nil, fmt.Errorf("invitación no encontrada")
		}

		for idF := range idsFoldersShared {
			wheres["idFolder"] = append(wheres["idFolder"], idF)
		}

		logic = "idFolder"

	} else if onlyTeam {
		wheres["ownerTeam_fk"] = []string{idTeam}
		logic = "ownerTeam_fk"

	} else if onlyUser {
		wheres["creatorUser_fk"] = []string{idUser}
		logic = "creatorUser_fk"

	} else {
		if idFolder != "" {
			wheres["idFolder"] = []string{idFolder}
			logic = "idFolder"
		} else {
			return &models.FolderListResp{}, fmt.Errorf("no hay parametros de busqueda")
		}
	}

	logic = fmt.Sprintf("%s ORDER BY %s", logic, orderBy)
	wheres["LOGIC"] = []string{logic, fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)}
	// se obtienen los datos de todos los folders
	foldersData, err := db.DB_con.GenericSelect("folders", "idFolder", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos de invitación: %v", err)
	}

	// Obtener documentos de cada folder
	attrs = []string{"idDocument", "idFolder"}
	wheres = map[string][]string{
		"idFolder": {},
	}

	foldersWithDocs := make(map[string]interface{})
	for folderID, folderData := range foldersData {
		wheres["idFolder"] = append(wheres["idFolder"], folderID)
		foldersWithDocs[folderID] = folderData
	}

	folderDocs, err := db.DB_con.GenericSelect("folderdocuments", "idfolderdocument", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo documentos del folder: %v", err)
	}
	fmt.Println("folderDocs: ", folderDocs)

	return &models.FolderListResp{
		Page:        page,
		PageSize:    pageSize,
		Total:       len(foldersData),
		FoldersData: foldersWithDocs,
	}, nil
}

func LoadFolderInfo(idFolder, idTeam, idUser string, onlyShared, onlyTeam, onlyUser bool, page, pageSize int, orderBy, orderDir string) (*models.FolderListResp, error) {

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
		"closedAt", "secuentialSign", "pathSerialized", "expirationDate",
		"numDocs", "numDocsSign", "numSigners", "numReceivers", "completedAt"}

	wheres := map[string][]string{}

	if onlyShared {
		// Se obtienen los folders de las invitaciones que ha recibido el usuario
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

		for idF := range idsFoldersShared {
			wheres["idFolder"] = append(wheres["idFolder"], idF)
		}
		logic = "idFolder"

	} else if onlyTeam {
		wheres["ownerTeam_fk"] = []string{idTeam}
		logic = "ownerTeam_fk"

	} else if onlyUser {
		wheres["creatorUser_fk"] = []string{idUser}
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
	docWheres := map[string][]string{
		"idFolder": {},
	}

	for folderID := range foldersData {
		docWheres["idFolder"] = append(docWheres["idFolder"], folderID)
	}

	docWheres["LOGIC"] = []string{"idFolder"}
	folderDocs, err := db.DB_con.GenericSelect("folderdocuments", "idfolderdocument", docAttrs, docWheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo documentos del folder: %v", err)
	}

	// Mapear documentos por folder
	docsByFolder := make(map[string][]string)
	for _, docData := range folderDocs {
		folderID := docData["idFolder"]
		docID := docData["idDocument"]
		docsByFolder[folderID] = append(docsByFolder[folderID], docID)
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
			Path:           folderData["pathSerialized"],
			Purpose:        folderData["purpose"],
			Documents:      docsByFolder[folderID],
		}

		// Si no hay documentos para este folder, inicializar como slice vacío
		if folder.Documents == nil {
			folder.Documents = []string{}
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
