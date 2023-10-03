package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
	"thesis/lib/database"
	"thesis/utility-service/models"
)

func main() {

	config := readConfig()

	dbconn := database.ConnParamSQL{
		Host: config["RDSHost"].(string),
		Port: int(config["RDSPort"].(float64)),
		Name: "utility",
		User: config["RDSUser"].(string),
		Psw:  config["RDSPass"].(string),
	}
	createTables(&dbconn, []database.ModelDB{models.Navigation{}, models.Log{}})

}

// config["foo"].(string)
func readConfig() map[string]interface{} {
	// Read the JSON file
	data, err := os.ReadFile("../config.json")
	if err != nil {
		log.Fatalf("Errore lettura file di configurazione " + err.Error())
	}

	// Unmarshal the JSON data into a map[string]interface{}
	var config map[string]interface{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Errore lettura file di configurazione " + err.Error())
	}
	return config
}

func createTables(dbconn database.ConnectionParameters, models []database.ModelDB) {
	var db database.DatabaseSQL
	err := db.Open(dbconn)
	if err != nil {
		log.Fatalf("Errore connessione")
		return
	}
	log.Println("Connesso")

	for _, model := range models {
		err = db.GormDB.Migrator().CreateTable(model)
		if err != nil {
			log.Printf("Errore creazione tabella %s", err.Error())
		} else {
			log.Printf("Creata tabella %s", strings.ToLower(reflect.TypeOf(model).Name()))
		}
	}
	err = db.Close()
	if err != nil {
		log.Fatalf("Errore chiusura")
	}
	log.Println("Disconnesso")
}
