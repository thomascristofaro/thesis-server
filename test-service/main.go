package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"thesis/financial-service/models"
	"thesis/lib/component"
	"thesis/lib/database"
)

func main() {

	ReadConfig()
	os.Setenv("DBName", "financial")

	createTables(database.NewConnParamSQLFromEnv(), []database.ModelDB{models.InvoiceLine{}})

	// page := pages.NewSalesOrderCard()
	// queryParams := map[string][]string{
	// 	"device_id": []string{""},
	// 	"No":        []string{"ORD0001"},
	// }
	// buff, err := pages.PostSalesOrder(page, queryParams)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(string(buff))
	// }

	// utility.SendSQSLogError(context.Background(),
	// 	utility.Message{
	// 		Body: []byte("Test"),
	// 		Metadata: map[string]string{
	// 			"event":    "test",
	// 			"function": "main",
	// 			"attr_1":   "val_1",
	// 		},
	// 	}, errors.New("Errore test"))

	// var filters map[string][]string
	// filters = map[string][]string{
	// 	"SalesOrderNo": []string{"ORD0001"},
	// }
	// TestGetPage(pages.NewSalesOrderLineList(), filters)

}

// config["foo"].(string)
func ReadConfig() map[string]interface{} {
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

	// Set environment variables
	for key, value := range config {
		if err := os.Setenv(key, fmt.Sprintf("%v", value)); err != nil {
			log.Fatalf("Errore impostazione variabile di ambiente " + err.Error())
		}
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

// PAGES Functions
func TestGetSchemaPage(p component.Page) {
	s, err := p.GetSchema()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(string(s))
	}
}

func TestGetPage(p component.Page, filters map[string][]string) {
	s, err := p.Get(filters)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(string(s))
	}
}

func TestPostPage(p component.Page, body string) {
	s, err := p.Post([]byte(body))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(string(s))
	}
}

func TestPatchPage(p component.Page, body string) {
	s, err := p.Patch([]byte(body))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(string(s))
	}
}

func TestDeletePage(p component.Page, filters map[string][]string) {
	s, err := p.Delete(filters)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(string(s))
	}
}

func TestButtonPage(p component.Page, filters map[string][]string) {
	s, err := p.Button(filters)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(string(s))
	}
}
