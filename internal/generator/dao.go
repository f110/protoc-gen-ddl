package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"vitess.io/vitess/go/vt/sqlparser"

	"go.f110.dev/protoc-ddl/internal/schema"
)

const GoDAOGeneratorVersion = "v0.1"

type GoDAOGenerator struct{}

func (g GoDAOGenerator) Generate(buf *bytes.Buffer, fileOpt *descriptor.FileOptions, messages *schema.Messages) {
	src := new(bytes.Buffer)

	entityPackageName := fileOpt.GetGoPackage()
	if strings.Contains(entityPackageName, ";") {
		s := strings.SplitN(entityPackageName, ";", 2)
		entityPackageName = s[0]
	}
	entityPackageAlias := filepath.Base(entityPackageName)

	src.WriteString("package dao\n")
	src.WriteString("import (\n")
	for _, v := range []string{"context", "database/sql", "fmt", "time", "strings"} {
		src.WriteString("\"" + v + "\"\n")
	}
	src.WriteRune('\n')
	for _, v := range []string{"golang.org/x/xerrors"} {
		src.WriteString("\"" + v + "\"\n")
	}
	src.WriteRune('\n')
	src.WriteString(fmt.Sprintf("\"%s\"\n", entityPackageName))
	src.WriteString(")\n")

	// ListOption
	src.WriteString(`
	type ListOption func(opt *listOpt)

	func Limit(limit int) func(opt *listOpt) {
		return func(opt *listOpt) {
			opt.limit = limit
		}
	}

    func Desc(opt *listOpt) {
        opt.desc = true
    }

	type listOpt struct {
		limit int
        desc  bool
	}

	func newListOpt(opts ...ListOption) *listOpt {
		opt := &listOpt{}
		for _, v := range opts {
			v(opt)
		}
		return opt
	}
`)

	messages.Each(func(m *schema.Message) {
		rels := make(map[string]struct{})
		for f := range m.Relations {
			if f.Virtual {
				continue
			}
			s := strings.Split(f.Type, ".")
			rels[s[len(s)-1]] = struct{}{}
		}
		relations := make([]string, 0)
		for v := range rels {
			relations = append(relations, v)
		}
		sort.Strings(relations)

		src.WriteString(fmt.Sprintf("type %s struct {\nconn *sql.DB\n\n", m.Descriptor.GetName()))
		for _, r := range relations {
			src.WriteString(fmt.Sprintf("%s *%s\n", schema.ToLowerCamel(r), schema.ToCamel(r)))
		}
		src.WriteString("}\n")
		src.WriteRune('\n')
		src.WriteString(fmt.Sprintf("func New%s(conn *sql.DB) *%s {\n", m.Descriptor.GetName(), m.Descriptor.GetName()))
		src.WriteString(fmt.Sprintf("return &%s{\nconn: conn,\n", m.Descriptor.GetName()))
		for _, r := range relations {
			src.WriteString(fmt.Sprintf("%s: New%s(conn),\n", schema.ToLowerCamel(r), schema.ToCamel(r)))
		}
		src.WriteString("}\n")
		src.WriteString("}\n\n")

		g.primaryKeySelect(src, m, entityPackageAlias)

		for _, q := range m.SelectQueries {
			stmt, err := sqlparser.Parse(q.Query)
			if err != nil {
				continue
			}

			switch v := stmt.(type) {
			case *sqlparser.Select:
				g.selectQuery(src, m, q.Query, q.Name, v, entityPackageAlias)
			default:
				log.Printf("%q is not supported: %s", v, q)
			}
		}

		g.create(src, m, entityPackageAlias)
		g.delete(src, m)
		g.update(src, m, entityPackageAlias)
	})

	buf.WriteString("// Generated by protoc-ddl.\n")
	buf.WriteString(fmt.Sprintf("// protoc-gen-dao: %s\n", GoDAOGeneratorVersion))
	b, err := format.Source(src.Bytes())
	if err != nil {
		log.Print(src.String())
		log.Print(err)
		return
	}
	buf.Write(b)
}

func (g GoDAOGenerator) create(src *bytes.Buffer, m *schema.Message, entityName string) {
	src.WriteString(fmt.Sprintf("func (d *%s) Create(ctx context.Context, v *%s.%s) (*%s.%s, error) {\n", m.Descriptor.GetName(), entityName, m.Descriptor.GetName(), entityName, m.Descriptor.GetName()))
	src.WriteString("res, err := d.conn.ExecContext(\nctx,\n")
	cols := make([]string, 0, m.Fields.Len())
	queryArgs := make([]string, 0, m.Fields.Len())
	m.Fields.Each(func(f *schema.Field) {
		if f.Sequence {
			return
		}
		if m.WithTimestamp && (f.Name == "created_at" || f.Name == "updated_at") {
			return
		}
		cols = append(cols, "`"+f.Name+"`")
		queryArgs = append(queryArgs, "v."+schema.ToCamel(f.Name))
	})
	if m.WithTimestamp {
		cols = append(cols, "`created_at`")
		queryArgs = append(queryArgs, "time.Now()")
	}
	args := make([]string, len(cols))
	for i := 0; i < len(args); i++ {
		args[i] = "?"
	}

	src.WriteString(fmt.Sprintf("\"INSERT INTO `%s` (%s) VALUES (%s)\",", m.TableName, strings.Join(cols, ", "), strings.Join(args, ", ")))
	src.WriteString(strings.Join(queryArgs, ",") + ",\n")
	src.WriteString(")\n")
	src.WriteString("if err != nil {\nreturn nil, xerrors.Errorf(\": %w\", err)\n}\n\n")
	src.WriteString("if n, err := res.RowsAffected(); err != nil {\n")
	src.WriteString("return nil, xerrors.Errorf(\": %w\", err)\n")
	src.WriteString("} else if n == 0 {\n")
	src.WriteString("return nil, sql.ErrNoRows\n")
	src.WriteString("}\n\n")
	src.WriteString("v = v.Copy()\n")
	if m.PrimaryKeys[0].Sequence {
		src.WriteString("insertedId, err := res.LastInsertId()\n")
		src.WriteString("if err != nil {\n")
		src.WriteString("return nil, xerrors.Errorf(\": %w\", err)\n")
		src.WriteString("}\n")
		src.WriteString(fmt.Sprintf("v.Id = %s(insertedId)\n", GoDataTypeMap[m.PrimaryKeys[0].Type]))
	}
	src.WriteRune('\n')

	src.WriteString("v.ResetMark()\n")
	src.WriteString("return v, nil\n")
	src.WriteString("}\n\n")
}

func (g GoDAOGenerator) delete(src *bytes.Buffer, m *schema.Message) {
	args := make([]string, 0)
	where := make([]string, 0)
	whereArgs := make([]string, 0)
	for _, v := range m.PrimaryKeys {
		args = append(args, fmt.Sprintf("%s %s", schema.ToLowerCamel(v.Name), GoDataTypeMap[v.Type]))
		where = append(where, fmt.Sprintf("`%s` = ?", v.Name))
		whereArgs = append(whereArgs, schema.ToLowerCamel(v.Name))
	}

	src.WriteString(fmt.Sprintf("func (d *%s) Delete(ctx context.Context,%s) error {\n", m.Descriptor.GetName(), strings.Join(args, ",")))
	src.WriteString(fmt.Sprintf("res, err := d.conn.ExecContext(ctx, \"DELETE FROM `%s` WHERE %s\", %s)\n", m.TableName, strings.Join(where, " AND "), strings.Join(whereArgs, ",")))
	src.WriteString("if err != nil {\n")
	src.WriteString("return xerrors.Errorf(\": %w\", err)\n")
	src.WriteString("}\n\n")
	src.WriteString("if n, err := res.RowsAffected(); err != nil {\n")
	src.WriteString("return xerrors.Errorf(\": %w\", err)\n")
	src.WriteString("} else if n == 0 {\n")
	src.WriteString("return sql.ErrNoRows\n")
	src.WriteString("}\n\n")
	src.WriteString("return nil\n")
	src.WriteString("}\n\n")
}

func (g GoDAOGenerator) update(src *bytes.Buffer, m *schema.Message, entityName string) {
	src.WriteString(fmt.Sprintf("func (d *%s) Update(ctx context.Context, v *%s.%s) error {\n", m.Descriptor.GetName(), entityName, m.Descriptor.GetName()))
	src.WriteString("if !v.IsChanged() {\n")
	src.WriteString("return nil\n")
	src.WriteString("}\n\n")

	src.WriteString("changedColumn := v.ChangedColumn()\n")
	src.WriteString("cols := make([]string, len(changedColumn)+1)\n")
	src.WriteString("values := make([]interface{}, len(changedColumn)+1)\n")
	src.WriteString("for i := range changedColumn {\n")
	src.WriteString("cols[i] = \"`\" + changedColumn[i].Name + \"` = ?\"\n")
	src.WriteString("values[i] = changedColumn[i].Value\n")
	src.WriteString("}\n")
	if m.WithTimestamp {
		src.WriteString("cols[len(cols)-1] = \"`updated_at` = ?\"\n")
		src.WriteString("values[len(values)-1] = time.Now()\n")
	}
	src.WriteRune('\n')

	where := make([]string, 0)
	whereArgs := make([]string, 0)
	for _, v := range m.PrimaryKeys {
		where = append(where, fmt.Sprintf("`%s` = ?", v.Name))
		whereArgs = append(whereArgs, "v."+schema.ToCamel(v.Name))
	}
	src.WriteString(fmt.Sprintf("query := fmt.Sprintf(\"UPDATE `%s` SET %%s WHERE %s\", strings.Join(cols, \", \"))\n", m.TableName, strings.Join(where, " AND ")))
	src.WriteString("res, err := d.conn.ExecContext(\n")
	src.WriteString("ctx,\n")
	src.WriteString("query,\n")
	src.WriteString(fmt.Sprintf("append(values, %s)...,\n", strings.Join(whereArgs, ",")))
	src.WriteString(")\n")
	src.WriteString("if err != nil {\n")
	src.WriteString("return xerrors.Errorf(\": %w\", err)\n")
	src.WriteString("}\n")

	src.WriteString("if n, err := res.RowsAffected(); err != nil {\n")
	src.WriteString("return xerrors.Errorf(\": %w\", err)\n")
	src.WriteString("} else if n == 0 {\n")
	src.WriteString("return sql.ErrNoRows\n")
	src.WriteString("}\n")
	src.WriteRune('\n')

	src.WriteString("v.ResetMark()\n")
	src.WriteString("return nil\n")
	src.WriteString("}\n\n")
}

func (g GoDAOGenerator) primaryKeySelect(src *bytes.Buffer, m *schema.Message, entityName string) {
	args := make([]string, 0)
	where := make([]string, 0)
	whereArgs := make([]string, 0)
	for _, v := range m.PrimaryKeys {
		args = append(args, fmt.Sprintf("%s %s", schema.ToLowerCamel(v.Name), GoDataTypeMap[v.Type]))
		where = append(where, fmt.Sprintf("`%s` = ?", v.Name))
		whereArgs = append(whereArgs, schema.ToLowerCamel(v.Name))
	}

	src.WriteString(fmt.Sprintf("func (d *%s) Select(ctx context.Context,%s) (*%s.%s, error) {\n", m.Descriptor.GetName(), strings.Join(args, ","), entityName, m.Descriptor.GetName()))
	src.WriteString("row := d.conn.QueryRowContext(ctx,")
	src.WriteString(fmt.Sprintf("\"SELECT * FROM `%s` WHERE %s\", %s)\n", m.TableName, strings.Join(where, " AND "), strings.Join(whereArgs, ",")))
	src.WriteRune('\n')
	src.WriteString(fmt.Sprintf("v := &%s.%s{}\n", entityName, m.Descriptor.GetName()))
	cols := make([]string, 0, m.Fields.Len())
	m.Fields.Each(func(f *schema.Field) {
		cols = append(cols, "&v."+schema.ToCamel(f.Name))
	})
	src.WriteString(fmt.Sprintf("if err := row.Scan(%s); err != nil {\n", strings.Join(cols, ",")))
	src.WriteString("return nil, xerrors.Errorf(\": %w\", err)\n}\n\n")
	if len(m.Relations) > 0 {
		relFields := make([]*schema.Field, 0)
		for f := range m.Relations {
			if f.Virtual {
				continue
			}
			relFields = append(relFields, f)
		}
		sort.Slice(relFields, func(i, j int) bool {
			return relFields[i].Name < relFields[j].Name
		})

		for _, f := range relFields {
			rels := m.Relations[f]

			r := make([]string, 0, len(rels))
			check := make([]string, 0, len(rels))
			for _, v := range rels {
				if f.Null {
					r = append(r, "*v."+schema.ToCamel(v.Name))
				} else {
					r = append(r, "v."+schema.ToCamel(v.Name))
				}
				check = append(check, "v."+schema.ToCamel(v.Name)+" != nil")
			}
			src.WriteString("{\n")
			if f.Null {
				src.WriteString(fmt.Sprintf("if %s {\n", strings.Join(check, " && ")))
			}
			s := strings.Split(f.Type, ".")
			src.WriteString(fmt.Sprintf("rel, err := d.%s.Select(ctx, %s)\n", schema.ToLowerCamel(s[len(s)-1]), strings.Join(r, ",")))
			src.WriteString("if err != nil {\n")
			src.WriteString("return nil, xerrors.Errorf(\": %w\", err)\n")
			src.WriteString("}\n")
			src.WriteString(fmt.Sprintf("v.%s = rel\n", schema.ToCamel(f.Name)))
			if f.Null {
				src.WriteString("}\n")
			}
			src.WriteString("}\n")
		}
		src.WriteRune('\n')
	}
	src.WriteString("v.ResetMark()\n")
	src.WriteString("return v, nil\n")
	src.WriteString("}\n\n")
}

func (g GoDAOGenerator) selectQuery(src *bytes.Buffer, m *schema.Message, rawQuery, name string, stmt *sqlparser.Select, entityName string) {
	if len(stmt.From) != 1 {
		log.Printf("Multiple tables is not supported")
		return
	}

	allColumn := false
	var cols []string
	for _, c := range stmt.SelectExprs {
		if _, ok := c.(*sqlparser.StarExpr); ok {
			allColumn = true
			cols = nil
			break
		}

		co, ok := c.(*sqlparser.AliasedExpr)
		if !ok {
			log.Printf("%v is %v", c, reflect.TypeOf(c))
			continue
		}
		col, ok := co.Expr.(*sqlparser.ColName)
		if !ok {
			log.Printf("%v is not ColName", c)
			continue
		}
		cols = append(cols, col.Name.String())
	}
	if allColumn {
		cols = make([]string, 0)
		stmt.SelectExprs = make([]sqlparser.SelectExpr, 0)
		for _, v := range m.Fields.List() {
			cols = append(cols, v.Name)
			stmt.SelectExprs = append(stmt.SelectExprs, &sqlparser.AliasedExpr{
				Expr: &sqlparser.ColName{
					Name: sqlparser.NewColIdent(v.Name),
				},
			})
		}
	}

	var comp []*schema.Field
	if stmt.Where != nil {
		comp = g.findArgs(m.Fields.List(), stmt.Where)
	}

	args := make([]string, len(comp))
	for i := range comp {
		args[i] = fmt.Sprintf("%s %s", schema.ToLowerCamel(comp[i].Name), GoDataTypeMap[comp[i].Type])
	}
	args = append(args, "opt ...ListOption")
	// Query execution
	src.WriteString(
		fmt.Sprintf(
			"func (d *%s) List%s(ctx context.Context, %s) ([]*%s.%s, error) {\n",
			m.Descriptor.GetName(), name,
			strings.Join(args, ","),
			entityName, m.Descriptor.GetName(),
		))
	primaryKeys := make([]string, 0)
	for _, v := range m.PrimaryKeys {
		primaryKeys = append(primaryKeys, "`"+v.Name+"`")
	}
	src.WriteString("listOpts := newListOpt(opt...)\n")
	src.WriteString(fmt.Sprintf("query := %q\n", printSelectQueryAST(stmt)))
	src.WriteString("if listOpts.limit > 0 {\n")
	src.WriteString("order := \"ASC\"\n")
	src.WriteString("if listOpts.desc {\n")
	src.WriteString("order = \"DESC\"\n")
	src.WriteString("}\n")
	src.WriteString(fmt.Sprintf("query = query + fmt.Sprintf(\" ORDER BY %s %%s LIMIT %%d\",order, listOpts.limit)\n", strings.Join(primaryKeys, ", ")))
	src.WriteString("}\n")
	src.WriteString("rows, err := d.conn.QueryContext(\nctx,\nquery,\n")
	for _, a := range comp {
		src.WriteString(schema.ToLowerCamel(a.Name) + ",\n")
	}
	src.WriteString(")\n")
	src.WriteString("if err != nil {\nreturn nil, xerrors.Errorf(\": %w\", err)\n}\n")
	src.WriteRune('\n')

	// Object mapping
	src.WriteString(fmt.Sprintf("res := make([]*%s.%s, 0)\n", entityName, m.Descriptor.GetName()))
	src.WriteString("for rows.Next() {\n")
	src.WriteString(fmt.Sprintf("r := &%s.%s{}\n", entityName, m.Descriptor.GetName()))
	scanCols := make([]string, len(cols))
	for i := range cols {
		scanCols[i] = "&r." + schema.ToCamel(cols[i])
	}
	src.WriteString(fmt.Sprintf("if err := rows.Scan(%s); err != nil {\n", strings.Join(scanCols, ",")))
	src.WriteString("return nil, xerrors.Errorf(\": %w\", err)\n")
	src.WriteString("}\n")
	src.WriteString("r.ResetMark()\n")
	src.WriteString("res = append(res, r)\n")
	src.WriteString("}\n")

	if len(m.Relations) > 0 {
		relFields := make([]*schema.Field, 0)
		for f := range m.Relations {
			if f.Virtual {
				continue
			}
			relFields = append(relFields, f)
		}
		sort.Slice(relFields, func(i, j int) bool {
			return relFields[i].Name < relFields[j].Name
		})

		src.WriteString("if len(res) > 0 {\n")
		src.WriteString("for _, v := range res {\n")
		for _, f := range relFields {
			rels := m.Relations[f]

			r := make([]string, 0, len(rels))
			check := make([]string, 0, len(rels))
			for _, v := range rels {
				if f.Null {
					r = append(r, "*v."+schema.ToCamel(v.Name))
				} else {
					r = append(r, "v."+schema.ToCamel(v.Name))
				}
				check = append(check, "v."+schema.ToCamel(v.Name)+" != nil")
			}
			src.WriteString("{\n")
			if f.Null {
				src.WriteString(fmt.Sprintf("if %s {\n", strings.Join(check, " && ")))
			}
			s := strings.Split(f.Type, ".")
			src.WriteString(fmt.Sprintf("rel, err := d.%s.Select(ctx, %s)\n", schema.ToLowerCamel(s[len(s)-1]), strings.Join(r, ",")))
			src.WriteString("if err != nil {\n")
			src.WriteString("return nil, xerrors.Errorf(\": %w\", err)\n")
			src.WriteString("}\n")
			src.WriteString(fmt.Sprintf("v.%s = rel\n", schema.ToCamel(f.Name)))
			if f.Null {
				src.WriteString("}\n")
			}
			src.WriteString("}\n")
		}
		src.WriteString("}\n")
		src.WriteString("}\n")
	}

	src.WriteRune('\n')
	src.WriteString("return res, nil\n")

	src.WriteString("}\n\n")
}

func (g GoDAOGenerator) findArgs(fields []*schema.Field, stmt *sqlparser.Where) []*schema.Field {
	fieldMap := make(map[string]*schema.Field)
	for _, v := range fields {
		fieldMap[v.Name] = v
	}

	res := g.findArgFieldFromExprIfExist(fieldMap, stmt.Expr)
	if len(res) == 0 {
		return nil
	}

	return res
}

func (g GoDAOGenerator) findArgFieldFromExprIfExist(fields map[string]*schema.Field, stmt sqlparser.Expr) []*schema.Field {
	res := make([]*schema.Field, 0)
	switch v := stmt.(type) {
	case *sqlparser.ComparisonExpr:
		f := g.findArgFieldFromComparisonExprIfExist(fields, v)
		if len(f) > 0 {
			res = append(res, f...)
		}
	case *sqlparser.AndExpr:
		f := g.findArgFieldFromAndExprIfExist(fields, v)
		if len(f) > 0 {
			res = append(res, f...)
		}
	case *sqlparser.OrExpr:
		f := g.findArgFieldFromOrExprIfExist(fields, v)
		if len(f) > 0 {
			res = append(res, f...)
		}
	default:
		log.Printf("%T", v)
	}

	return res
}

func (GoDAOGenerator) findArgFieldFromComparisonExprIfExist(fields map[string]*schema.Field, stmt *sqlparser.ComparisonExpr) []*schema.Field {
	res := make([]*schema.Field, 0)
	if left, ok := stmt.Left.(*sqlparser.ColName); ok {
		switch r := stmt.Right.(type) {
		case *sqlparser.SQLVal:
			if r.Type == sqlparser.ValArg {
				if f, ok := fields[left.Name.String()]; ok {
					res = append(res, f)
				}
			}
		}
	}

	return res
}

func (g GoDAOGenerator) findArgFieldFromAndExprIfExist(fields map[string]*schema.Field, stmt *sqlparser.AndExpr) []*schema.Field {
	res := make([]*schema.Field, 0)

	for _, s := range []sqlparser.Expr{stmt.Left, stmt.Right} {
		f := g.findArgFieldFromExprIfExist(fields, s)
		if len(f) > 0 {
			res = append(res, f...)
		}
	}

	return res
}

func (g GoDAOGenerator) findArgFieldFromOrExprIfExist(fields map[string]*schema.Field, stmt *sqlparser.OrExpr) []*schema.Field {
	res := make([]*schema.Field, 0)

	for _, s := range []sqlparser.Expr{stmt.Left, stmt.Right} {
		f := g.findArgFieldFromExprIfExist(fields, s)
		if len(f) > 0 {
			res = append(res, f...)
		}
	}

	return res
}

func printSelectQueryAST(stmt *sqlparser.Select) string {
	buf := sqlparser.NewTrackedBuffer(func(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
		switch v := node.(type) {
		case *sqlparser.SQLVal:
			if v.Type == sqlparser.ValArg {
				buf.WriteString("?")
				return
			}
		}

		node.Format(buf)
	})
	stmt.Format(buf)

	return buf.String()
}
