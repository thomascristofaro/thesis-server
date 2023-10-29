package database

import (
	"context"
	"errors"
	"reflect"
	"strings"
)

type IteratorDB interface {
	Next(dst interface{}) bool
	Close() error
}

type Database interface {
	DBType() DBType
	SetContext(ctx context.Context)
	GetContext() (ctx context.Context)
	Open(connParam ConnectionParameters) error
	Close() error
	Get(value interface{}) (err error)
	Create(value interface{}) (err error)
	Update(value interface{}) (err error)
	Delete(value interface{}) (err error)
	Find(value interface{}, filters []Filter) (it IteratorDB, err error)
	BeginTransaction() (err error)
	CommitTransaction() (err error)
	RollbackTransaction() (err error)
}

type ConnectionParameters interface {
	getURL() (string, error)
	checkParameters() error
}

type DBType int

const (
	SQL DBType = iota
	NOSQL
)

// MODEL

type ModelDB interface {
	DBType() DBType
}

type ModelCtrl struct {
	Name     string
	Model    ModelDB
	ModelArr interface{}
	// xDbModel interface{}
	filters []Filter
	lastErr error
	db      Database
	Conn    ConnectionParameters
	it      IteratorDB
}

func FillNameFields(a interface{}) {
	v := reflect.ValueOf(a).Elem()
	for i := 0; i < v.NumField(); i++ {
		type_field := v.Type().Field(i)
		name := type_field.Name
		if strings.HasPrefix(type_field.Type.String(), "database.Field") {
			f := v.Field(i).FieldByName("Field")
			if f.IsValid() {
				ftoset := reflect.ValueOf(Field{Name: name})
				f.Set(ftoset)
			}
		}
	}
}

func NewModel(model ModelDB) ModelCtrl {
	var db Database
	var dbSQL DatabaseSQL
	var conn ConnectionParameters
	switch model.DBType() {
	case SQL:
		conn = NewConnParamSQLFromEnv()
		dbSQL = DatabaseSQL{}
		db = &dbSQL
	case NOSQL:
		db = NewDatabaseNoSQL()
	}
	name := "TABELLA" // da capire come prenderla dal model
	FillNameFields(model)
	return ModelCtrl{
		Name:  name,
		Model: model,
		Conn:  conn,
		// xDbModel: dbModel,
		db: db,
	}
}

func (m *ModelCtrl) GetLastError() (err error) {
	return m.lastErr
}

func (m *ModelCtrl) Get() bool {
	//TODO check chiave primaria
	m.lastErr = m.db.Get(m.Model)
	return m.lastErr == nil
}

func (m *ModelCtrl) Create() bool {
	m.lastErr = m.db.Create(m.Model)
	return m.lastErr == nil
}

func (m *ModelCtrl) Update() bool {
	//TODO questo però salva tutti i campi,
	// e non so se chackare il fatto che un altro può aver modificato il record
	m.lastErr = m.db.Update(m.Model)
	return m.lastErr == nil
}

func (m *ModelCtrl) Delete() bool {
	//check se la chiave primaria è compilata,
	//sennò cancella TUTTA LA TABELLA
	m.lastErr = m.db.Delete(m.Model)
	return m.lastErr == nil
}

func (m *ModelCtrl) SetFilters(filters map[string][]string) {
	for key, value := range filters {
		m.SetFilter(key, EQUAL, value[0])
	}
}

func (m *ModelCtrl) SetFilterFromStr(filter string) {
	if strings.Contains(filter, ">=") {
		a := strings.Split(filter, ">=")
		m.SetFilter(a[0], ">=", a[1]) // ma non va bene, da capire come fare per il tipo
	} else if strings.Contains(filter, "<=") {
		a := strings.Split(filter, "<=")
		m.SetFilter(a[0], "<=", a[1])
		// } else if strings.Contains(filter, "!=") { // il diverso mi sa che non funziona su docstore
		// 	a := strings.Split(filter, "!=")
		// 	m.SetFilter(a[0], "!=", a[1])
	} else if strings.Contains(filter, "=") {
		a := strings.Split(filter, "=")
		m.SetFilter(a[0], "=", a[1])
	}
}

// sono tutti in AND, capire come gestire anche gli OR più avanti
func (m *ModelCtrl) SetFilter(fieldName string, operator Operator, value interface{}) {
	m.filters = append(m.filters, Filter{
		field: fieldName,
		op:    operator,
		value: value,
	})
}

func (m *ModelCtrl) Find() bool {
	m.it, m.lastErr = m.db.Find(m.Model, m.filters)
	return m.lastErr == nil
}

// manca il close?
func (m *ModelCtrl) Next() bool {
	return m.it.Next(m.Model)
}

func (m *ModelCtrl) Open() bool {
	if m.Conn == nil {
		m.lastErr = errors.New("connection parameters not set")
		return false
	}
	m.lastErr = m.db.Open(m.Conn)
	return m.lastErr == nil
}

func (m *ModelCtrl) Close() bool {
	m.lastErr = m.db.Close()
	return m.lastErr == nil
}

func (m *ModelCtrl) BeginTransaction() bool {
	m.lastErr = m.db.BeginTransaction()
	return m.lastErr == nil
}

func (m *ModelCtrl) CommitTransaction() bool {
	m.lastErr = m.db.CommitTransaction()
	return m.lastErr == nil
}

func (m *ModelCtrl) RollbackTransaction() bool {
	m.lastErr = m.db.RollbackTransaction()
	return m.lastErr == nil
}

func (m *ModelCtrl) SetDB(db Database) {
	m.db = db
}

func (m *ModelCtrl) GetDB() Database {
	return m.db
}

//FIELDS

type Operator string

const (
	EQUAL Operator = "="
	// NOT_EQUAL Operator = "<>"
	GREATER       Operator = ">"
	LESS          Operator = "<"
	GREATER_EQUAL Operator = ">="
	LESS_EQUAL    Operator = "<="
)

type Field struct {
	Name string
}

type FieldString struct {
	Field
	ValueVar string
}

type FieldInt struct {
	Field
	ValueVar int
}

type FieldFloat struct {
	Field
	ValueVar float64
}

// dovranno essere re-implementati con i customFields: capire come modificare il package docstore per implementarli
// TODO APRIRE NUOVO ISSUE SU GO-CLOUD
type Filter struct {
	field string
	op    Operator
	value interface{}
}
