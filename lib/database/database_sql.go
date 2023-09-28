package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	gosql "gocloud.dev/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConnParamSQL struct {
	Host string
	Port int
	Name string
	User string
	Psw  string
}

func NewConnParamSQL(host string, port int, name string, user string, psw string) ConnectionParameters {
	return &ConnParamSQL{
		Host: host,
		Port: port,
		Name: name,
		User: user,
		Psw:  psw,
	}
}

func (connParam *ConnParamSQL) getURL() (string, error) {
	//controllare i parametri di connessione
	if err := connParam.checkParameters(); err != nil {
		return "", err
	}
	return fmt.Sprintf("mysql://%s:%s@%s:%d/%s",
		connParam.User, connParam.Psw, connParam.Host, connParam.Port, connParam.Name), nil
}

func (connParam *ConnParamSQL) checkParameters() error {
	return nil
}

type DatabaseSQL struct {
	DB     *sql.DB
	GormDB *gorm.DB
	url    string
	ctx    context.Context
}

func NewDatabaseSQL() Database {
	return &DatabaseSQL{}
}

func NewAndOpenDatabaseSQL(host string, port int, name string, user string, psw string) (d Database, err error) {
	c := NewConnParamSQL(host, port, name, user, psw)
	d = NewDatabaseSQL()
	err = d.Open(c)
	return d, err
}

func (d *DatabaseSQL) DBType() DBType {
	return SQL
}

func (d *DatabaseSQL) SetContext(ctx context.Context) {
	d.ctx = ctx
}

func (d *DatabaseSQL) GetContext() (ctx context.Context) {
	if d.ctx == nil {
		c := context.TODO()
		d.ctx = c
	}
	return d.ctx
}

func (d *DatabaseSQL) Open(connParam ConnectionParameters) (err error) {
	// d.connParam = connParam
	d.url, err = connParam.getURL()
	if err != nil {
		return err
	}
	// si può fare senza url?
	d.DB, err = gosql.Open(d.GetContext(), d.url)
	if err != nil {
		return err
	}
	//forse è da rivedere questa istruzione di connessione
	d.GormDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: d.DB,
	}), &gorm.Config{})
	return err
}

func (d *DatabaseSQL) Close() error {
	return d.DB.Close()
}

func (d *DatabaseSQL) Get(value interface{}) (err error) {
	//TODO check chiave primaria
	result := d.GormDB.First(value)
	return result.Error
}

func (d *DatabaseSQL) Create(value interface{}) (err error) {
	result := d.GormDB.Create(value)
	return result.Error
}

func (d *DatabaseSQL) Update(value interface{}) (err error) {
	//TODO questo però salva tutti i campi,
	// e non so se chackare il fatto che un altro può aver modificato il record
	result := d.GormDB.Save(value)
	return result.Error
}

func (d *DatabaseSQL) Delete(value interface{}) (err error) {
	//check se la chiave primaria è compilata,
	//sennò cancella TUTTA LA TABELLA
	result := d.GormDB.Delete(value)
	return result.Error
}

func (d *DatabaseSQL) Find(value interface{}, filters []Filter) (it IteratorDB, err error) {
	tx := d.GormDB.Model(value)
	for _, filter := range filters {
		tx.Where(filter.field+string(filter.op)+"?", filter.value)
	}
	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}
	return &IteratorSQL{rows: rows, gormDB: d.GormDB}, nil
}

type IteratorSQL struct {
	rows   *sql.Rows
	gormDB *gorm.DB
}

func (it *IteratorSQL) Next(dst interface{}) bool {
	if ok := it.rows.Next(); !ok {
		return false
	}
	if err := it.gormDB.ScanRows(it.rows, dst); err != nil {
		return false
	}
	return true
}

func (it *IteratorSQL) Close() error {
	return it.rows.Close()
}

// FIELD

func (f *FieldString) Scan(value interface{}) error {
	if value == nil {
		f.ValueVar = ""
		return nil
	}
	val, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal string value:", value))
	}
	f.ValueVar = val
	return nil
}

func (f FieldString) Value() (driver.Value, error) {
	return f.ValueVar, nil
}

func (f *FieldInt) Scan(value interface{}) error {
	if value == nil {
		f.ValueVar = 0
		return nil
	}
	val, ok := value.(int)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal int value:", value))
	}
	f.ValueVar = val
	return nil
}

func (f FieldInt) Value() (driver.Value, error) {
	return f.ValueVar, nil
}

func (f *FieldFloat) Scan(value interface{}) error {
	if value == nil {
		f.ValueVar = 0
		return nil
	}
	val, ok := value.(float64)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal float value:", value))
	}
	f.ValueVar = val
	return nil
}

func (f FieldFloat) Value() (driver.Value, error) {
	return f.ValueVar, nil
}
