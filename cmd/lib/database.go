package lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// Variable globale pour la connexion à la base de données
var db *sql.DB

// initialisation de la connexion à la base de données
func InitDB(dataSourceName string) {
	var err error

	// vérifie si base de donnée existe déjà
	_, err = os.Stat(dataSourceName)
	dbExists := !os.IsExist(err)

	// ouvre la connexion à la DB
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	if !dbExists {
		creatTablesFromSQL()
	}
}

// fonction qui lit le fichier SQL et l'exécute
func createTablesFromSQL() {
	// lecture du contenu du fichier SQL
	sqlBytes, err := os.ReadFile("data_base.sql")
	if err != nil {
		fmt.Println("Error reading file", err)
	}

	sqlContent := string(sqlBytes)

	// exécute les commandes sql
	_, err = DB.Exec(sqlContent)
	if err != nil {
		fmt.Println("Error when creating tables", err)
	}

	fmt.Println("Successful creation of file data")
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		fmt.Print(err)
	}
}
