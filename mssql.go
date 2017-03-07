package pop

import (
	"fmt"
	"os/exec"

	_ "github.com/go-sql-driver/mssql"
	"github.com/markbates/going/clam"
	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
	"github.com/pkg/errors"
)

type mssql struct {
	ConnectionDetails *ConnectionDetails
}

func (m *mssql) Details() *ConnectionDetails {
	return m.ConnectionDetails
}

func (m *mssql) URL() string {
	c := m.ConnectionDetails
	if c.URL != "" {
		return c.URL
	}

	s := "%s:%s@(%s:%s)/%s?parseTime=true&multiStatements=true&readTimeout=1s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database)
}

func (m *mssql) MigrationURL() string {
	return m.URL()
}

func (m *mssql) Create(s store, model *Model, cols Columns) error {
	return errors.Wrap(genericCreate(s, model, cols), "mssql create")
}

func (m *mssql) Update(s store, model *Model, cols Columns) error {
	return errors.Wrap(genericUpdate(s, model, cols), "mssql update")
}

func (m *mssql) Destroy(s store, model *Model) error {
	return errors.Wrap(genericDestroy(s, model), "mssql destroy")
}

func (m *mssql) SelectOne(s store, model *Model, query Query) error {
	return errors.Wrap(genericSelectOne(s, model, query), "mssql select one")
}

func (m *mssql) SelectMany(s store, models *Model, query Query) error {
	return errors.Wrap(genericSelectMany(s, models, query), "mssql select many")
}

func (m *mssql) CreateDB() error {
	c := m.ConnectionDetails
	cmd := exec.Command("mssql", "-u", c.User, "-p"+c.Password, "-h", c.Host, "-P", c.Port, "-e", fmt.Sprintf("create database %s", c.Database))
	err := clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
	return errors.Wrapf(err, "error creating mssql database %s", c.Database)
}

func (m *mssql) DropDB() error {
	c := m.ConnectionDetails
	cmd := exec.Command("mssql", "-u", c.User, "-p"+c.Password, "-h", c.Host, "-P", c.Port, "-e", fmt.Sprintf("drop database %s", c.Database))
	err := clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
	return errors.Wrapf(err, "error dropping mssql database %s", c.Database)
}

func (m *mssql) TranslateSQL(sql string) string {
	return sql
}

func (m *mssql) FizzTranslator() fizz.Translator {
	t := translators.Newmssql(m.URL(), m.Details().Database)
	return t
}

func (m *mssql) Lock(fn func() error) error {
	return fn()
}

func newMsSQL(deets *ConnectionDetails) dialect {
	cd := &mssql{
		ConnectionDetails: deets,
	}

	return cd
}
