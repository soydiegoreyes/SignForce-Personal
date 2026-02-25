package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"sfmiddle/configs"
	"sfmiddle/models"
	"sfmiddle/utilities"

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

// Función auxiliar para generar los placeholders (?, ?, ?) de un IN
func placeholders(n int) string {
	ps := make([]string, n)
	for i := range ps {
		ps[i] = "?"
	}
	return strings.Join(ps, ",")
}

func (cnx *ConexionDB) GenericSelect(tableName string, idColName string, attributes []string, whereMap map[string][]string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)
	var args []interface{}
	ats := strings.Join(attributes, ",")

	var wheres string
	logic, hasLogic := whereMap["LOGIC"]

	// 1. Extraer y filtrar las llaves (quitando "LOGIC")
	var keys []string
	for k := range whereMap {
		if k != "LOGIC" {
			keys = append(keys, k) // <--- Aquí faltaba la 'k'
		}
	}

	if !hasLogic {
		// CASO SIN LOGICA: Ordenamos alfabéticamente para que siempre sea igual
		sort.Strings(keys)
		var parts []string
		for _, k := range keys {
			v := whereMap[k]
			parts = append(parts, fmt.Sprintf("%s IN (%s)", k, placeholders(len(v))))
			for _, val := range v {
				args = append(args, val)
			}
		}
		wheres = strings.Join(parts, " AND ")
	} else {
		// CASO CON LOGICA: Ordenamos las llaves según su aparición en el string logic[0]
		wheres = logic[0]

		// Ordenamos las llaves basándonos en su posición en el string de lógica
		sort.Slice(keys, func(i, j int) bool {
			return strings.Index(wheres, keys[i]) < strings.Index(wheres, keys[j])
		})

		// Ahora que están ordenadas según aparecen en el SQL, procesamos
		for _, k := range keys {
			v := whereMap[k]

			// Prioridad de reemplazo (de más complejo a más simple)
			if strings.Contains(wheres, fmt.Sprintf("NOT %s", k)) {
				wheres = strings.ReplaceAll(wheres, fmt.Sprintf("NOT %s", k), fmt.Sprintf("%s NOT IN (%s)", k, placeholders(len(v))))
				for _, val := range v {
					args = append(args, val)
				}

			} else if strings.Contains(wheres, fmt.Sprintf("%s BETWEEN", k)) {
				wheres = strings.ReplaceAll(wheres, fmt.Sprintf("%s BETWEEN", k), fmt.Sprintf("%s BETWEEN ? AND ?", k))
				args = append(args, v[0], v[1])

			} else if strings.Contains(wheres, fmt.Sprintf("REGEXP %s", k)) {
				wheres = strings.ReplaceAll(wheres, fmt.Sprintf("REGEXP %s", k), fmt.Sprintf("%s REGEXP ?", k))
				args = append(args, v[0])

			} else if strings.Contains(wheres, fmt.Sprintf("LIKE %s", k)) {
				wheres = strings.ReplaceAll(wheres, fmt.Sprintf("LIKE %s", k), fmt.Sprintf("%s LIKE ?", k))
				args = append(args, "%"+v[0]+"%")

			} else if strings.Contains(wheres, k) {
				// Reemplazo para el nombre de la columna simple (usando IN por defecto)
				// Usamos un reemplazo cuidadoso para no romper nombres de columnas similares
				wheres = strings.ReplaceAll(wheres, k, fmt.Sprintf("%s IN (%s)", k, placeholders(len(v))))
				for _, val := range v {
					args = append(args, val)
				}
			}
		}

		if len(logic) == 2 {
			wheres = fmt.Sprintf("%s %s", wheres, logic[1])
		}
	}

	query := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s;", idColName, ats, tableName, wheres)
	fmt.Println(query, args)
	// IMPORTANTE: Ahora pasamos los 'args' a la query
	rows, err := cnx.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar: %v", err)
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

/*
func (cnx *ConexionDB) GenericJoinSelect(

	mainTable string,
	joinTable string,
	joinCondition string,
	idColName string,
	wVals []string, // Estos son los valores peligrosos
	mainAttributes []string,
	joinAttributes []string,

	) (map[string]map[string]string, error) {
		result := make(map[string]map[string]string)

		// 1. Preparar los argumentos para la ejecución
		var args []interface{}
		for _, v := range wVals {
			args = append(args, v)
		}

		// 2. Construir la lista de atributos (Estructura - Segura si es interna)
		mainAttrs := ""
		if len(mainAttributes) > 0 {
			mainAttrs = mainTable + "." + strings.Join(mainAttributes, ", "+mainTable+".")
		}

		joinAttrs := ""
		if len(joinAttributes) > 0 {
			joinAttrs = joinTable + "." + strings.Join(joinAttributes, ", "+joinTable+".")
		}

		allAttrs := []string{}
		if mainAttrs != "" {
			allAttrs = append(allAttrs, mainAttrs)
		}
		if joinAttrs != "" {
			allAttrs = append(allAttrs, joinAttrs)
		}
		selectClause := strings.Join(allAttrs, ", ")

		// 3. Crear los placeholders (?, ?, ?) basado en la cantidad de wVals
		pHolders := make([]string, len(wVals))
		for i := range pHolders {
			pHolders[i] = "?"
		}
		whereList := strings.Join(pHolders, ", ")

		// 4. Construir la consulta SQL usando los placeholders
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
		fmt.Println(query, args)
		// 5. Ejecutar pasando los argumentos por separado
		rows, err := cnx.DB.Query(query, args...)
		if err != nil {
			return nil, fmt.Errorf("error al ejecutar la consulta: %v", err)
		}
		defer rows.Close()

		// --- El resto del procesamiento de filas se mantiene igual ---
		columns, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("error al obtener columnas: %v", err)
		}

		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(sql.NullString)
		}

		for rows.Next() {
			err := rows.Scan(values...)
			if err != nil {
				return nil, fmt.Errorf("error al escanear fila: %v", err)
			}

			// El ID es el último valor según tu SELECT
			idNS := values[len(values)-1].(*sql.NullString)
			id := ""
			if idNS.Valid {
				id = idNS.String
			}

			rowData := make(map[string]string)
			for i, colName := range columns {
				// Usamos el nombre de la columna para el mapa, pero saltamos el ID final
				// para no duplicarlo si ya está en los atributos
				if i == len(columns)-1 {
					continue
				}
				ns := values[i].(*sql.NullString)
				if ns.Valid {
					rowData[colName] = ns.String
				} else {
					rowData[colName] = ""
				}
			}
			result[id] = rowData
		}

		return result, nil
	}
*/
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
			if x, ok := val.(bool); ok {
				values = append(values, utilities.Bool2Int(x))
			} else {
				values = append(values, val)
			}
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
		hash, err := utilities.GetHash(v, configs.HashConf, true)
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
		if x, ok := v.(bool); ok {
			vals += "'" + fmt.Sprintf("%v", utilities.Bool2Int(x)) + "'"
		} else {
			vals += "'" + fmt.Sprintf("%v", v) + "'"
		}
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

func Update_algos() {
	var q = "SELECT * FROM signforceapi.algos;"
	algosData, err := DB_con.ExecuteSelect(q)
	if err != nil {
		fmt.Println("Error al actualizar algoritmos de base de datos")
		return
	}
	utilities.AlgosMap = algosData
}

func NewConn() ConexionDB {
	DB_cnx := ConexionDB{}
	if DB_cnx.Conectar() {
		fmt.Println("Conexion exitosa")
	}

	return DB_cnx
}

var DB_con ConexionDB
