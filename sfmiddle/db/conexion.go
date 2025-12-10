package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sfmiddle/configs"
	"sfmiddle/models"
	"sfmiddle/utilities"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// ConexionDB estructura para manejar la conexión
type ConexionDB struct {
	DB *sql.DB
}

// Método para conectar a MySQL
func (cnx *ConexionDB) Conectar() bool {

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbIP := os.Getenv("DB_IP")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPass, dbIP, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al abrir la conexión: %v", err)
		return false
	}

	// Verificar la conexión
	err = db.Ping()
	if err != nil {
		log.Fatalf("No se pudo conectar a MySQL: %v", err)
		return false
	}

	cnx.DB = db
	log.Println("Conexión exitosa a MySQL")
	log.Printf("%s @%s:%s", dbUser, dbIP, dbName)
	return true
}

// Método para desconectar
func (cnx *ConexionDB) Desconectar() {
	if cnx.DB != nil {
		cnx.DB.Close()
		fmt.Println("Conexión a MySQL cerrada correctamente.")
	}
}

// ================================ SELECTS ================================================
// Obtener atributos de una tabla genérica basada en el id principal de la tabla
/*
	EJEMPLO WHERE CON LOGICA:
	wheres := map[string][]string{
		"createdAtDoc":    {"2025-10-21", "2025-10-23"},
		"ownerTeamDoc_fk": {"4"},
		"ownerInstDoc_fk": {"1"},
		"idDocument":      {"ASC"},
		"LOGIC":           {"ownerInstDoc_fk AND NOT ownerTeamDoc_fk AND createdAtDoc BETWEEN ORDER BY idDocument", "LIMIT 5 OFFSET 0"},
	}

	RESULT: SELECT idDocument, documentHash,documentPath FROM documents WHERE ownerInstDoc_fk IN ('1') AND ownerTeamDoc_fk NOT IN ('4') AND createdAtDoc BETWEEN '2025-10-21' AND '2025-10-23' ORDER BY idDocument ASC LIMIT 5 OFFSET 0;
*/
func (cnx *ConexionDB) GenericSelect(tableName string, idColName string, attributes []string, whereMap map[string][]string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	ats := strings.Join(attributes, ",")

	var wheres string
	if logic, exists := whereMap["LOGIC"]; !exists {
		// si no existe lógica significa que no puede haber and, or y not y debe haber solo una columna de atributos
		for k, v := range whereMap {
			wheres += fmt.Sprintf("%s IN ('%s') AND ", k, strings.Join(v, "','"))

		}
		wheres = wheres[:len(wheres)-5]
	} else {
		wheres = logic[0]
		for k, v := range whereMap {
			if strings.Contains(wheres, fmt.Sprintf("NOT %s", k)) {
				wheres = strings.ReplaceAll(wheres, fmt.Sprintf("NOT %s", k), fmt.Sprintf("%s NOT IN ('%s')", k, strings.Join(v, "','")))
			} else if strings.Contains(wheres, fmt.Sprintf("%s BETWEEN", k)) {
				wheres = strings.ReplaceAll(wheres, fmt.Sprintf("%s BETWEEN", k), fmt.Sprintf("%s BETWEEN '%s' AND '%s'", k, v[0], v[1]))
			} else if strings.Contains(wheres, fmt.Sprintf("ORDER BY %s", k)) {
				wheres = strings.ReplaceAll(wheres, fmt.Sprintf("ORDER BY %s", k), fmt.Sprintf("ORDER BY %s %s", k, v[0]))
			} else if strings.Contains(wheres, fmt.Sprintf("REGEXP %s", k)) {
				wheres = strings.ReplaceAll(wheres, fmt.Sprintf("REGEXP %s", k), fmt.Sprintf("%s REGEXP '%s'", k, v[0]))
			} else {
				wheres = strings.ReplaceAll(wheres, k, fmt.Sprintf("%s IN ('%s')", k, strings.Join(v, "','")))
			}
		}
		if len(logic) == 2 {
			wheres = fmt.Sprintf("%s %s", wheres, logic[1])
		}
	}

	query := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s;", idColName, ats, tableName, wheres)
	fmt.Println(query)

	rows, err := cnx.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %s  -> %v", query, err)
	}
	defer rows.Close()

	// preparar slice para scan
	cols := append([]string{idColName}, attributes...)
	values := make([]interface{}, len(cols))
	for i := range values {
		values[i] = new(sql.NullString) // usar NullString
	}

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return nil, fmt.Errorf("error al escanear fila: %v", err)
		}

		colIDns := values[0].(*sql.NullString)
		colID := ""
		if colIDns.Valid {
			colID = colIDns.String
		}

		colData := make(map[string]string)
		for i, attr := range attributes {
			ns := values[i+1].(*sql.NullString)
			if ns.Valid {
				colData[attr] = ns.String
			} else {
				colData[attr] = "" // NULL → string vacío
			}
		}
		result[colID] = colData
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar filas: %v", err)
	}

	return result, nil
}

// Realiza un JOIN entre dos tablas y devuelve los resultados
func (cnx *ConexionDB) GenericJoinSelect(
	mainTable string, // Tabla principal (ej: "users")
	joinTable string, // Tabla a unir (ej: "profiles")
	joinCondition string, // Condición de JOIN (ej: "users.id = profiles.user_id")
	idColName string, // Columna de ID para filtrar (ej: "users.id")
	wVals []string, // IDs a buscar (ej: ["1", "2"])
	mainAttributes []string, // Atributos de la tabla principal (ej: ["name", "email"])
	joinAttributes []string, // Atributos de la tabla secundaria (ej: ["bio", "avatar"])
) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	// Construir la lista de atributos para el SELECT
	mainAttrs := ""
	if len(mainAttributes) > 0 {
		mainAttrs = mainTable + "." + strings.Join(mainAttributes, ", "+mainTable+".")
	}

	joinAttrs := ""
	if len(joinAttributes) > 0 {
		joinAttrs = joinTable + "." + strings.Join(joinAttributes, ", "+joinTable+".")
	}

	// Combinar atributos de ambas tablas
	allAttrs := []string{}
	if mainAttrs != "" {
		allAttrs = append(allAttrs, mainAttrs)
	}
	if joinAttrs != "" {
		allAttrs = append(allAttrs, joinAttrs)
	}
	selectClause := strings.Join(allAttrs, ", ")

	// Construir la lista de IDs para el IN
	wherePlaceholders := make([]string, len(wVals))
	for i, whereV := range wVals {
		wherePlaceholders[i] = "'" + whereV + "'"
	}
	whereList := strings.Join(wherePlaceholders, ", ")

	// Construir la consulta SQL
	query := fmt.Sprintf(
		"SELECT %s, %s.%s FROM %s JOIN %s ON %s WHERE %s.%s IN (%s);",
		selectClause,
		mainTable,
		idColName,
		mainTable,
		joinTable,
		joinCondition,
		mainTable,
		idColName,
		whereList,
	)
	fmt.Println(query)
	rows, err := cnx.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %s -> %v", query, err)
	}
	defer rows.Close()

	// Obtener nombres de columnas
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error al obtener columnas: %v", err)
	}

	// Preparar valores para Scan
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(sql.NullString) // usar NullString
	}

	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("error al escanear fila: %v", err)
		}

		// Obtener el ID principal (asumimos que es la última columna)
		idNS := values[len(values)-1].(*sql.NullString)
		id := ""
		if idNS.Valid {
			id = idNS.String
		}

		rowData := make(map[string]string)

		// Mapear cada columna a su valor
		for i, colName := range columns {
			if colName == idColName {
				continue // Saltar la columna de ID (ya lo tenemos)
			}
			ns := values[i].(*sql.NullString)
			if ns.Valid {
				rowData[colName] = ns.String
			} else {
				rowData[colName] = "" // representar NULL como string vacío
			}
		}

		result[id] = rowData
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar filas: %v", err)
	}

	return result, nil
}

func (cnx *ConexionDB) ExecuteSelect(query string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)
	fmt.Println(query)
	// Ejecutar el SELECT
	rows, err := cnx.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %s -> %v", query, err)
	}
	defer rows.Close()

	// Obtener nombres de columnas dinámicamente
	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error al obtener columnas: %v", err)
	}
	if len(cols) == 0 {
		return nil, fmt.Errorf("la consulta no devolvió columnas")
	}

	// Preparar valores y punteros
	values := make([]interface{}, len(cols))
	for i := range values {
		values[i] = new(sql.NullString)
	}

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return nil, fmt.Errorf("error al escanear fila: %v", err)
		}

		// La primera columna será usada como clave (igual que en tu función original)
		idVal := ""
		if ns, ok := values[0].(*sql.NullString); ok && ns.Valid {
			idVal = ns.String
		}

		rowData := make(map[string]string)
		for i, colName := range cols {
			ns := values[i].(*sql.NullString)
			if ns.Valid {
				rowData[colName] = ns.String
			} else {
				rowData[colName] = ""
			}
		}

		// Si no hay valor en la primera columna, genera una clave numérica
		if idVal == "" {
			idVal = fmt.Sprintf("row_%d", len(result)+1)
		}

		result[idVal] = rowData
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar filas: %v", err)
	}

	return result, nil
}

// =========================================================================================
// ======================================== UPDATES ========================================
// GenericUpdate actualiza registros en una tabla basada en el ID principal
/*
UPDATE Users
SET userPass = 'MiPassword', userActive = '1'
WHERE idUser = 3;
*/

func (cnx *ConexionDB) GenericBatchUpdate(tableName string, whereColumn string, updates map[string]map[string]interface{}) error {
	if tableName == "" || whereColumn == "" || len(updates) == 0 {
		return fmt.Errorf("parámetros inválidos")
	}

	for whereValue, updateData := range updates {
		if len(updateData) == 0 {
			continue
		}

		setClauses := make([]string, 0, len(updateData))
		values := make([]interface{}, 0, len(updateData)+1)

		for field, val := range updateData {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
			values = append(values, val)
		}

		// Agregar el valor del WHERE al final
		values = append(values, whereValue)

		query := fmt.Sprintf(
			"UPDATE %s SET %s WHERE %s = ?",
			tableName,
			strings.Join(setClauses, ", "),
			whereColumn,
		)
		fmt.Println(query, values)
		// Ejecutar la consulta
		if _, err := cnx.DB.Exec(query, values...); err != nil {
			return fmt.Errorf("error actualizando %s=%v: %w", whereColumn, whereValue, err)
		}

		fmt.Println("Valores actualizados.")
	}

	return nil
}

func (cnx *ConexionDB) UpdateValData(idInst string, idUser string, valReq *models.ValidationRequest) error {

	instFields := map[string]string{
		"streetAddress": valReq.StreetAddress,
		"postalCode":    valReq.PostalCode,
		"neighborhood":  valReq.Neighborhood,
		"locality":      valReq.Locality,
		//"lineAddr":valReq.AddressLine,
	}
	kycFields := map[string]string{
		"ActaConstitutiva":   valReq.ActaConstitutiva,
		"PoderRepresentante": valReq.PoderRepresentante,
		"IdentidadOficial":   valReq.IdentidadOficial,
		"PruebaResidencia":   valReq.PruebaResidencia,
	}

	// se llenan los campos de institutions
	var query, updt string
	for k, v := range instFields {
		if v == "" {
			continue
		}
		updt += k + " = '" + v + "',"
	}
	if len(updt) > 0 {
		updt = updt[:len(updt)-1]
		query = fmt.Sprintf("UPDATE institutions SET %s WHERE idInstitution = '%s';", updt, idInst)

		// Ejecutar la consulta
		if _, err := cnx.DB.Exec(query); err != nil {
			return fmt.Errorf("error actualizando %s: %w", query, err)
		}
		fmt.Println("Valores actualizados.")
	}

	// se llenan los campos de kyc idInsttitution_fk, idUser_fk, documentHash, documentExt, documentName, documentClass, documentPath, expirationDate
	for k, v := range kycFields {
		if v == "" {
			continue
		}

		// se sacan los datos para cada columna
		ls := strings.LastIndex(v, "/")
		ld := strings.LastIndex(v, ".")

		path := v[:ls+1]
		name := v[ls+1 : ld]

		ext := v[ld+1:]
		hash, err := utilities.GetHash(v, configs.HashConf)
		if err != nil {
			return fmt.Errorf("error al obtener hash del documento: %w", err)
		}

		cols := []string{"idInsttitution_fk", "idUser_fk", "documentHash", "documentExt", "documentName", "documentClass", "documentPath"}
		idx, err := cnx.GenericInsert("kyc", cols, []interface{}{idInst, idUser, hash, ext, name, k, path})
		if err != nil {
			return fmt.Errorf("error al insertar valores de documento: %w", err)
		}
		fmt.Println("Valores actualizados en indice: ", idx)

	}

	return nil
}

// =========================================================================================
// ======================================== INSERTS ========================================
// GenericInsert añade registros en una tabla
/*
INSERT INTO table (COLUMNS) VALUES (VALUES);
*/
func (cnx *ConexionDB) GenericInsert(tableName string, columns []string, values []interface{}) (string, error) {
	if tableName == "" || len(columns) == 0 || len(values) == 0 {
		return "", fmt.Errorf("parámetros inválidos")
	}

	var keys, vals string

	// constuimos las columnas que son keys
	keys = strings.Join(columns, ",")

	// se construyen los valores que son los atributos
	for i, v := range values {
		vals += "'" + fmt.Sprintf("%v", v) + "'"
		if i < len(values)-1 {
			vals += ", "
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, keys, vals)
	fmt.Println(query)
	// Ejecutar la consulta
	if result, err := cnx.DB.Exec(query); err != nil {
		return "", fmt.Errorf("error insertando datos: %s\nERROR: %w", query, err)
	} else {
		// Obtener información útil del resultado
		if lastID, err := result.LastInsertId(); err == nil {
			fmt.Printf("Registro insertado con ID: %d\n", lastID)
			return fmt.Sprintf("%v", lastID), nil
		} else {
			return "", nil
		}
	}
}

func NewConn() ConexionDB {
	DB_cnx := ConexionDB{}
	if DB_cnx.Conectar() {
		fmt.Println("Conexion exitosa")
	}

	return DB_cnx
}

var DB_con ConexionDB
