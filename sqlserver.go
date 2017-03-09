package pop

import (
	"fmt"
	"os/exec"
	_ "github.com/minus5/gofreetds"
	"github.com/markbates/going/clam"
	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
	"github.com/pkg/errors"
)

type sqlserver struct {
	ConnectionDetails *ConnectionDetails
}

func (m *sqlserver) Details() *ConnectionDetails {
	return m.ConnectionDetails
}

func (m *sqlserver) URL() string {
	c := m.ConnectionDetails
	if c.URL != "" {
		return c.URL
	}
	// TODO JW -> insert url code
	s := "%s:%s@(%s:%s)/%s?parseTime=true&multiStatements=true&readTimeout=1s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database)
}

func (m *sqlserver) MigrationURL() string {
	return m.URL()
}

func (m *sqlserver) Create(s store, model *Model, cols Columns) error {
	return errors.Wrap(genericCreate(s, model, cols), "sqlserver create")
}

func (m *sqlserver) Update(s store, model *Model, cols Columns) error {
	return errors.Wrap(genericUpdate(s, model, cols), "sqlserver update")
}

func (m *sqlserver) Destroy(s store, model *Model) error {
	return errors.Wrap(genericDestroy(s, model), "sqlserver destroy")
}

func (m *sqlserver) SelectOne(s store, model *Model, query Query) error {
	return errors.Wrap(genericSelectOne(s, model, query), "sqlserver select one")
}

func (m *sqlserver) SelectMany(s store, models *Model, query Query) error {
	return errors.Wrap(genericSelectMany(s, models, query), "sqlserver select many")
}

func (m *sqlserver) CreateDB() error {
	// TODO JW -> appropriate call to sqlserver command line
	c := m.ConnectionDetails
	//cmd := exec.Command("sqlserver", "-u", c.User, "-p"+c.Password, "-h", c.Host, "-P", c.Port, "-e", fmt.Sprintf("create database %s", c.Database))
	fmt.Printf("%s -----> \n", c.Password)
	// TODO: remove hardcoded path to sqlcmd
	cmd := exec.Command("/opt/mssql-tools/bin/sqlcmd", "-U", c.User, "-P"+c.Password, "-h", c.Host, "-P", c.Port, "-Q", fmt.Sprintf("create database %s", c.Database))
	err := clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
	return errors.Wrapf(err, "error creating SQLServer database %s", c.Database)
}

func (m *sqlserver) DropDB() error {
	// TODO JW -> appropriate call to sqlserver command line
	c := m.ConnectionDetails
	cmd := exec.Command("sqlserver", "-u", c.User, "-p"+c.Password, "-h", c.Host, "-P", c.Port, "-e", fmt.Sprintf("drop database %s", c.Database))
	err := clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
	return errors.Wrapf(err, "error dropping sqlserver database %s", c.Database)
}

func (m *sqlserver) TranslateSQL(sql string) string {
	return sql
}

func (m *sqlserver) FizzTranslator() fizz.Translator {
	t := translators.NewSQLServer()
	return t
}

func (m *sqlserver) Lock(fn func() error) error {
	return fn()
}

func newSQLServer(deets *ConnectionDetails) dialect {
	cd := &sqlserver{
		ConnectionDetails: deets,
	}

	return cd
}
