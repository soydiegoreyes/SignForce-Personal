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

// Obtener atributos de una tabla genérica basada en el id
func (cnx *ConexionDB) GenericSelect(tableName string, idColName string, ids []string, attributes []string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	// construir la lista de los atributos para la consulta SQL
	ats := strings.Join(attributes, ",")

	// construir la lista de ids para la cláusula IN
	idstring := ""
	for i, id := range ids {
		idstring += "'" + id + "'"
		if i != len(ids)-1 {
			idstring += ","
		}
	}

	query := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s IN (%s);", idColName, ats, tableName, idColName, idstring)
	log.Println(query)
	rows, err := cnx.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %s  -> %v", query, err)
	}
	defer rows.Close()

	// Preparar slice para almacenar los valores escaneados
	cols := append([]string{idColName}, attributes...)
	values := make([]interface{}, len(cols))
	for i := range values {
		values[i] = new(string)
	}
	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("error al escanear fila: %v", err)
		}

		colID := *(values[0].(*string))
		colData := make(map[string]string)

		for i, attr := range attributes {
			colData[attr] = *(values[i+1].(*string)) // +1 porque el primer valor es usuClave
		}

		result[colID] = colData
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

var DB_con ConexionDB
