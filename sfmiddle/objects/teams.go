package objects

import (
	"fmt"
	"sfmiddle/db"
)

//type Team struct{}

func CreateTeam(idInst, idUser, nameTeam, description string) string {
	columns := []string{
		"idInstitution_fk", "creatorUser_fk", "nameTeam", "description",
	}
	values := []interface{}{
		idInst, idUser, nameTeam, description,
	}
	idTeam, err := db.DB_con.GenericInsert("teams", columns, values)
	if err != nil {
		fmt.Println("Error: al insertar datos", err)
		return ""
	}
	return idTeam
}
