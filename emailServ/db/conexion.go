package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// ConexionDB estructura para manejar la conexión
type ConexionDB struct {
	DB *sql.DB
}

// Cargar las variables de entorno desde .env
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}
}

// Método para conectar a MySQL
func (cnx *ConexionDB) Conectar() bool {
	loadEnv()

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
	fmt.Println("Conexión exitosa a MySQL")
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
func (cnx *ConexionDB) GenericSelect(tableName string, idColName string, attributes []string, whereMap map[string][]string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	ats := strings.Join(attributes, ",")

	var wheres string
	if logic, exists := whereMap["LOGIC"]; !exists {
		// si no existe lógica significa que no puede haber and, or y not y debe haber solo una columna de atributos, se t
		for k, v := range whereMap {
			wheres += fmt.Sprintf("%s IN ('%s') AND ", k, strings.Join(v, "','"))

		}
		wheres = wheres[:len(wheres)-5]
	} else {
		wheres = logic[0]
		for k, v := range whereMap {
			if strings.Contains(wheres, fmt.Sprintf("NOT %s", k)) {
				wheres = strings.ReplaceAll(wheres, fmt.Sprintf("NOT %s", k), fmt.Sprintf("%s NOT IN ('%s')", k, strings.Join(v, "','")))
			} else {
				wheres = strings.ReplaceAll(wheres, k, fmt.Sprintf("%s IN ('%s')", k, strings.Join(v, "','")))
			}
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

// =========================================================================================

// Obtener atributos de usuario
func (cnx *ConexionDB) GetEmailAttributes(whereAttr string, whEqAttr string, attributes []string) (map[string]string, error) {
	data := make(map[string]string)
	ats := ""
	for _, attr := range attributes[:len(attributes)-1] {
		ats += attr
		ats += ","
	}
	ats += attributes[len(attributes)-1]

	query := fmt.Sprintf("SELECT %s FROM emails WHERE %s = %s", ats, whEqAttr, whereAttr)

	//fmt.Println(query)

	row := cnx.DB.QueryRow(query)
	values := make([]interface{}, len(attributes))
	for i := range values {
		values[i] = new(string)
	}

	err := row.Scan(values...)
	if err != nil {
		return nil, err
	}

	for i, attr := range attributes {
		data[attr] = *(values[i].(*string))
	}

	return data, nil
}

// Actualizar atributos de usuario
func (cnx *ConexionDB) UpdateEmailAttributes(idEmail string, attributes []string, newValues []string) bool {
	if len(attributes) != len(newValues) {
		log.Println("Error: La cantidad de atributos y valores no coinciden")
		return false
	}

	query := "UPDATE emails SET "
	args := []interface{}{}

	for i, attr := range attributes {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("`%s` = ?", attr)
		args = append(args, newValues[i])
	}

	query += " WHERE idEmail = ?"
	args = append(args, idEmail)

	result, err := cnx.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Error al actualizar usuario %s: %v", idEmail, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0
}

// InsertEmailDynamic inserta un nuevo email con atributos dinámicos
func (cnx *ConexionDB) InsertEmailAttributes(attributes []string, values []interface{}) (int64, error) {
	// Iniciar transacción
	tx, err := cnx.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("error al iniciar transacción: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback() // Si hay error, hacemos rollback
		}
	}()

	// Construir y ejecutar query (similar a tu función original)
	query := fmt.Sprintf("INSERT INTO emails (%s) VALUES (%s)",
		strings.Join(attributes, ", "),
		strings.Repeat("?, ", len(values)-1)+"?")

	result, err := tx.Exec(query, values...)
	if err != nil {
		return 0, fmt.Errorf("error en ejecución: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener ID: %v", err)
	}

	// Commit explícito para persistir cambios
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("error en commit: %v", err)
	}

	return id, nil
}

// Obtener atributos de usuario
func (cnx *ConexionDB) GetUserAttributes(idUsers []string, attributes []string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	// Construir la lista de atributos para la consulta SQL
	ats := strings.Join(attributes, ",")

	// Construir la lista de usuarios para la cláusula IN
	usrs := ""
	for i, idUser := range idUsers {
		usrs += "'" + idUser + "'" // Agregar comillas para valores string en SQL
		if i != len(idUsers)-1 {
			usrs += ","
		}
	}

	query := fmt.Sprintf("SELECT idUser, %s FROM users WHERE idUser IN (%s);", ats, usrs)

	rows, err := cnx.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %v", err)
	}
	defer rows.Close()

	// Preparar slice para almacenar los valores escaneados
	cols := append([]string{"idUser"}, attributes...)
	values := make([]interface{}, len(cols))
	for i := range values {
		values[i] = new(string)
	}

	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("error al escanear fila: %v", err)
		}

		userID := *(values[0].(*string))
		userData := make(map[string]string)

		for i, attr := range attributes {
			userData[attr] = *(values[i+1].(*string)) // +1 porque el primer valor es usuClave
		}

		result[userID] = userData
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar filas: %v", err)
	}

	return result, nil
}

func NewConn() ConexionDB {
	DB_cnx := ConexionDB{}
	if DB_cnx.Conectar() {
		fmt.Println("Conexion exitosa")
	}

	return DB_cnx
}

var DB_con ConexionDB = NewConn()
