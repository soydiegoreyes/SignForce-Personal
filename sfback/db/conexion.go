package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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
		values[i] = new(string)
	}

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return nil, fmt.Errorf("error al escanear fila: %v", err)
		}

		colID := *(values[0].(*string))
		colData := make(map[string]string)
		for i, attr := range attributes {
			colData[attr] = *(values[i+1].(*string))
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
		values[i] = new(string)
	}

	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("error al escanear fila: %v", err)
		}

		// Obtener el ID principal (asumimos que es el último valor)
		id := *(values[len(values)-1].(*string))
		rowData := make(map[string]string)

		// Mapear cada columna a su valor
		for i, colName := range columns {
			if colName == idColName {
				continue // Saltar la columna de ID (ya la tenemos)
			}
			rowData[colName] = *(values[i].(*string))
		}

		result[id] = rowData
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

		// Ejecutar la consulta
		if _, err := cnx.DB.Exec(query, values...); err != nil {
			return fmt.Errorf("error actualizando %s=%v: %w", whereColumn, whereValue, err)
		}

		fmt.Println("Valores actualizados:")
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
	// Ejecutar la consulta
	if result, err := cnx.DB.Exec(query); err != nil {
		return "", fmt.Errorf("error insertando datos: %s\nERROR: %w", query, err)
	} else {
		// Obtener información útil del resultado
		if lastID, err := result.LastInsertId(); err == nil {
			log.Printf("Registro insertado con ID: %d", lastID)
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
