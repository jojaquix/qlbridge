package expr

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	u "github.com/araddon/gou"
	"github.com/araddon/qlbridge/lex"
	"github.com/araddon/qlbridge/value"
)

var (
	_ = u.EMPTY

	// Ensure SqlSelect etc are NodeTypes
	_ Node = (*SqlSelect)(nil)
	_ Node = (*SqlInsert)(nil)
	_ Node = (*SqlUpsert)(nil)
	_ Node = (*SqlUpdate)(nil)
	_ Node = (*SqlInsert)(nil)
)

// The sqlStatement interface, to define the sub-types
//  Select, Insert, Delete etc
type SqlStatement interface {
	Accept(visitor Visitor) (interface{}, error)
	Keyword() lex.TokenType
}

type SqlSelect struct {
	Pos
	Star    bool
	Columns Columns
	From    string
	Where   Node
	Limit   int
}
type SqlInsert struct {
	Pos
	Columns Columns
	Rows    [][]value.Value
	Into    string
}
type SqlUpsert struct {
	Pos
	Columns Columns
	Rows    [][]value.Value
	Into    string
}
type SqlUpdate struct {
	Pos
	kw      lex.TokenType // Update, Upsert
	Columns Columns
	Where   Node
	From    string
}
type SqlDelete struct {
	Pos
	Table string
	Where Node
	Limit int
}
type SqlShow struct {
	Pos
	Identity string
}
type SqlDescribe struct {
	Pos
	Identity string
}

func NewSqlSelect() *SqlSelect {
	req := &SqlSelect{}
	req.Columns = make(Columns, 0)
	return req
}
func NewSqlInsert() *SqlInsert {
	req := &SqlInsert{}
	req.Columns = make(Columns, 0)
	return req
}
func NewSqlUpdate() *SqlUpdate {
	req := &SqlUpdate{kw: lex.TokenUpdate}
	req.Columns = make(Columns, 0)
	return req
}
func NewSqlDelete() *SqlDelete {
	return &SqlDelete{}
}

// Array of Columns
type Columns []*Column

func (m *Columns) AddColumn(col *Column) { *m = append(*m, col) }
func (m *Columns) String() string {
	colCt := len(*m)
	if colCt == 1 {
		return (*m)[0].String()
	} else if colCt == 0 {
		return ""
	}

	s := make([]string, len(*m))
	for i, col := range *m {
		s[i] = col.String()
	}

	return strings.Join(s, ", ")
}
func (m *Columns) FieldNames() []string {
	names := make([]string, len(*m))
	for i, col := range *m {
		names[i] = col.Key()
	}
	return names
}

// Column represents the Column as expressed in a [SELECT]
// expression
type Column struct {
	As      string
	Comment string
	Star    bool
	Tree    *Tree
	Guard   *Tree // If
}

func (m *Column) Key() string    { return m.As }
func (m *Column) String() string { return m.As }

func (m *SqlSelect) Keyword() lex.TokenType { return lex.TokenSelect }
func (m *SqlSelect) Check() error           { return nil }
func (m *SqlSelect) NodeType() NodeType     { return SqlSelectNodeType }
func (m *SqlSelect) Type() reflect.Value    { return nilRv }
func (m *SqlSelect) StringAST() string      { return fmt.Sprintf("%s ", m.Keyword()) }
func (m *SqlSelect) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("SELECT %s FROM %s", m.Columns, m.From))
	if m.Where != nil {
		buf.WriteString(fmt.Sprintf(" WHERE %s ", m.Where.String()))
	}
	return buf.String()
}

func (m *SqlInsert) Keyword() lex.TokenType { return lex.TokenInsert }
func (m *SqlInsert) Check() error           { return nil }
func (m *SqlInsert) Type() reflect.Value    { return nilRv }
func (m *SqlInsert) NodeType() NodeType     { return SqlInsertNodeType }
func (m *SqlInsert) StringAST() string      { return fmt.Sprintf("%s ", m.Keyword()) }
func (m *SqlInsert) String() string         { return fmt.Sprintf("%s ", m.Keyword()) }

func (m *SqlUpsert) Keyword() lex.TokenType { return lex.TokenUpsert }
func (m *SqlUpsert) Check() error           { return nil }
func (m *SqlUpsert) Type() reflect.Value    { return nilRv }
func (m *SqlUpsert) NodeType() NodeType     { return SqlUpsertNodeType }
func (m *SqlUpsert) StringAST() string      { return fmt.Sprintf("%s ", m.Keyword()) }
func (m *SqlUpsert) String() string         { return fmt.Sprintf("%s ", m.Keyword()) }

func (m *SqlUpdate) Keyword() lex.TokenType { return m.kw }
func (m *SqlUpdate) Check() error           { return nil }
func (m *SqlUpdate) Type() reflect.Value    { return nilRv }
func (m *SqlUpdate) NodeType() NodeType     { return SqlUpdateNodeType }
func (m *SqlUpdate) StringAST() string      { return fmt.Sprintf("%s ", m.Keyword()) }
func (m *SqlUpdate) String() string         { return fmt.Sprintf("%s ", m.Keyword()) }

func (m *SqlDelete) Keyword() lex.TokenType { return lex.TokenDelete }
func (m *SqlDelete) Check() error           { return nil }
func (m *SqlDelete) Type() reflect.Value    { return nilRv }
func (m *SqlDelete) NodeType() NodeType     { return SqlDeleteNodeType }
func (m *SqlDelete) StringAST() string      { return fmt.Sprintf("%s ", m.Keyword()) }
func (m *SqlDelete) String() string         { return fmt.Sprintf("%s ", m.Keyword()) }

func (m *SqlDescribe) Keyword() lex.TokenType { return lex.TokenDescribe }
func (m *SqlDescribe) Check() error           { return nil }
func (m *SqlDescribe) Type() reflect.Value    { return nilRv }
func (m *SqlDescribe) NodeType() NodeType     { return SqlDescribeNodeType }
func (m *SqlDescribe) StringAST() string      { return fmt.Sprintf("%s ", m.Keyword()) }
func (m *SqlDescribe) String() string         { return fmt.Sprintf("%s ", m.Keyword()) }

func (m *SqlShow) Keyword() lex.TokenType { return lex.TokenShow }
func (m *SqlShow) Check() error           { return nil }
func (m *SqlShow) Type() reflect.Value    { return nilRv }
func (m *SqlShow) NodeType() NodeType     { return SqlShowNodeType }
func (m *SqlShow) StringAST() string      { return fmt.Sprintf("%s ", m.Keyword()) }
func (m *SqlShow) String() string         { return fmt.Sprintf("%s ", m.Keyword()) }

// Implement Accept() part of SqlStatment interface
func (m *SqlSelect) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitSelect(m)
}

// Implement Accept() part of SqlStatment interface
func (m *SqlInsert) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitInsert(m)
}

// Implement Accept() part of SqlStatment interface
func (m *SqlUpdate) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUpdate(m)
}

// Implement Accept() part of SqlStatment interface
func (m *SqlDelete) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitDelete(m)
}

// Implement Accept() part of SqlStatment interface
func (m *SqlDescribe) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitDescribe(m)
}

func (m *SqlShow) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitShow(m)
}
