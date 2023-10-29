package database

import (
	"context"
	"fmt"

	"gocloud.dev/docstore"
	"gocloud.dev/docstore/driver"
)

type ConnParamNoSQL struct {
	Protocol     string
	Table        string
	PartitionKey string
	SortKey      string
}

func NewConnParamNoSQL(protocol string, table string, partitionKey string, sortKey string) ConnectionParameters {
	return &ConnParamNoSQL{
		Protocol:     protocol,
		Table:        table,
		PartitionKey: partitionKey,
		SortKey:      sortKey,
	}
}

func (connParam *ConnParamNoSQL) getURL() (string, error) {
	//controllare i parametri di connessione
	if err := connParam.checkParameters(); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s://%s?partition_key=%s&sort_key=%s",
		connParam.Protocol, connParam.Table, connParam.PartitionKey, connParam.SortKey), nil
}

func (connParam *ConnParamNoSQL) checkParameters() error {
	return nil
}

type DatabaseNoSQL struct {
	Coll *docstore.Collection
	url  string
	ctx  context.Context
}

func NewDatabaseNoSQL() Database {
	return &DatabaseNoSQL{}
}

func NewAndOpenDatabaseNoSQL(protocol string, table string, partitionKey string, sortKey string) (d Database, err error) {
	c := NewConnParamNoSQL(protocol, table, partitionKey, sortKey)
	d = NewDatabaseNoSQL()
	err = d.Open(c)
	return d, err
}

func (d *DatabaseNoSQL) DBType() DBType {
	return NOSQL
}

func (d *DatabaseNoSQL) SetContext(ctx context.Context) {
	d.ctx = ctx
}

func (d *DatabaseNoSQL) GetContext() (ctx context.Context) {
	if d.ctx == nil {
		c := context.TODO()
		d.ctx = c
	}
	return d.ctx
}

func (d *DatabaseNoSQL) Open(connParam ConnectionParameters) (err error) {
	// d.connParam = connParam
	d.url, err = connParam.getURL()
	if err != nil {
		return err
	}
	// si può fare senza url?
	d.Coll, err = docstore.OpenCollection(d.GetContext(), d.url)
	return err
}

func (d *DatabaseNoSQL) Close() error {
	return d.Coll.Close()
}

func (d *DatabaseNoSQL) BeginTransaction() (err error) {
	return nil
}

func (d *DatabaseNoSQL) CommitTransaction() (err error) {
	return nil
}

func (d *DatabaseNoSQL) RollbackTransaction() (err error) {
	return nil
}

func (d *DatabaseNoSQL) Get(value interface{}) (err error) {
	//TODO check chiave primaria
	return d.Coll.Get(d.GetContext(), value)
}

func (d *DatabaseNoSQL) Create(value interface{}) (err error) {
	return d.Coll.Create(d.GetContext(), value)
}

func (d *DatabaseNoSQL) Update(value interface{}) (err error) {
	return d.Coll.Put(d.GetContext(), value)
}

func (d *DatabaseNoSQL) Delete(value interface{}) (err error) {
	return d.Coll.Delete(d.GetContext(), value)
}

func (d *DatabaseNoSQL) Find(value interface{}, filters []Filter) (it IteratorDB, err error) {
	q := d.Coll.Query()
	for _, filter := range filters {
		q.Where(docstore.FieldPath(filter.field), string(filter.op), filter.value)
	}
	rows := q.Get(d.ctx)
	return &IteratorNoSQL{rows: rows, ctx: d.ctx}, nil
}

type IteratorNoSQL struct {
	rows *docstore.DocumentIterator
	ctx  context.Context
}

func (it *IteratorNoSQL) Next(dst interface{}) bool {
	if err := it.rows.Next(it.ctx, dst); err != nil {
		return false
	}
	return true
}

func (it *IteratorNoSQL) Close() error {
	it.rows.Stop()
	return nil
}

// FIELD

func (f FieldString) CustomMarshal(enc driver.Encoder) error {
	enc.EncodeString(f.ValueVar)
	return nil
}

func (f *FieldString) CustomUnmarshal(dec driver.Decoder) bool {
	s, ok := dec.AsString()
	if ok {
		f.ValueVar = s
	}
	return ok
}

func (f FieldInt) CustomMarshal(enc driver.Encoder) error {
	enc.EncodeInt(int64(f.ValueVar))
	return nil
}

func (f *FieldInt) CustomUnmarshal(dec driver.Decoder) bool {
	s, ok := dec.AsInt()
	if ok {
		f.ValueVar = int(s)
	}
	return ok
}

func (f FieldFloat) CustomMarshal(enc driver.Encoder) error {
	enc.EncodeFloat(f.ValueVar)
	return nil
}

func (f *FieldFloat) CustomUnmarshal(dec driver.Decoder) bool {
	s, ok := dec.AsFloat()
	if ok {
		f.ValueVar = s
	}
	return ok
}
