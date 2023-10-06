package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
	"thesis/lib/database"
)

func main() {

	// config := readConfig()

	// dbconn := database.ConnParamSQL{
	// 	Host: config["DBHost"].(string),
	// 	Port: int(config["DBPort"].(float64)),
	// 	Name: "financial",
	// 	User: config["DBUser"].(string),
	// 	Psw:  config["DBPass"].(string),
	// }
	// createTables(&dbconn, []database.ModelDB{models.InvoiceHeader{}, models.InvoiceLine{}})

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
