package lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

// Variable globale pour la connexion à la base de données
var DB *sql.DB

// initialisation de la connexion à la base de données
func InitDB(dataSourceName string) {
	var err error

	// vérifie si base de donnée existe déjà
	_, err = os.Stat(dataSourceName)
	dbExists := !os.IsNotExist(err)

	// ouvre la connexion à la DB
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	if !dbExists {
		createTablesFromSQL()
	}
}

// fonction qui lit le fichier SQL et l'exécute
func createTablesFromSQL() {
	// lecture du contenu du fichier SQL
	sqlBytes, err := os.ReadFile("data_base.sql")
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}

	sqlContent := string(sqlBytes)

	// exécute les commandes sql
	_, err = DB.Exec(sqlContent)
	if err != nil {
		fmt.Println("Error when creating tables", err)
		return
	}

	fmt.Println("Successful creation of file data")
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		fmt.Print(err)
	}
}

func InsertUser(nickname, age int, gender string, firstName, lastName, email, password string) error {

}

// fonction qui vérifie si la colonne Status existe et l'ajoute si nécessaire
func EnsureStatusColumn() error {
	// Vérifie si la colonne Status existe dans la table USERS
	rows, err := DB.Query("PRAGMA table_info(USERS)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasStatusColumn := false
	for rows.Next() {
		var cid, notnull, pk int
		var name, type_name string
		var dflt_value interface{}

		err := rows.Scan(&cid, &name, &type_name, &notnull, &dflt_value, &pk)
		if err != nil {
			return err
		}

		if name == "Status" {
			hasStatusColumn = true
			break
		}
	}

	// Si la colonne Status n'existe pas, l'ajouter
	if !hasStatusColumn {
		_, err := DB.Exec(`ALTER TABLE USERS ADD COLUMN Status TEXT DEFAULT 'offline'`)
		if err != nil {
			return err
		}
		fmt.Println("Status column added to USERS table")
	}

	return nil
}

// fonction qui met à jour si le user est en ligne
func UpdateUserStatus(userId int, status string) error {
	err := EnsureStatusColumn()

	if err != nil {
		return err
	}

	// MAJ du status en ligne ou non
	_, err = DB.Exec("UPDATE USERS SET Status = ? WHERE UserId = ?", status, userId)
	return err

}

// fonction qui vérifie si une chaîne contient une sous-chaîne
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
