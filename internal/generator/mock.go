package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"vitess.io/vitess/go/vt/sqlparser"

	"go.f110.dev/protoc-ddl/internal/schema"
)

type GoDAOMockGenerator struct{}

func (g GoDAOMockGenerator) Generate(buf *bytes.Buffer, fileOpt *descriptor.FileOptions, messages *schema.Messages, daoPath string) {
	src := newBuffer()

	entityPackageName := fileOpt.GetGoPackage()
	if strings.Contains(entityPackageName, ";") {
		s := strings.SplitN(entityPackageName, ";", 2)
		entityPackageName = s[0]
	}
	entityPackageAlias := filepath.Base(entityPackageName)

	src.Write("package daotest")
	src.Write("import (")
	for _, v := range []string{"context"} {
		src.Write("\"" + v + "\"")
	}
	src.LineBreak()
	for _, v := range []string{"go.f110.dev/protoc-ddl/mock"} {
		src.Write("\"" + v + "\"")
	}
	src.LineBreak()
	src.Writef("%q", entityPackageName)
	src.Writef("%q", daoPath)
	src.Write(")")
	src.LineBreak()

	messages.Each(func(m *schema.Message) {
		s := &GoDAOStruct{m: m, entityPackageName: entityPackageAlias, daoPath: filepath.Base(daoPath)}
		selectFuncs := s.Select(g.selectRowQuery)

		src.Writef("type %s struct{", m.Descriptor.GetName())
		src.Write("*mock.Mock")
		src.Write("}")

		src.Writef("func New%s() *%s {", m.Descriptor.GetName(), m.Descriptor.GetName())
		src.Writef("return &%s{Mock: mock.New()}", m.Descriptor.GetName())
		src.Write("}")
		src.LineBreak()

		src.WriteFunc(s.PrimaryKeySelect(g.primaryKeySelect), g.mockPrimaryKeySelect(entityPackageAlias, m))

		for _, v := range selectFuncs {
			src.WriteString(v.String())
			src.LineBreak()

			f := g.mockSelectQuery(entityPackageAlias, m, v)
			src.WriteString(f.String())
			src.LineBreak()
		}

		src.WriteFunc(s.Create(g.create), s.Delete(g.delete), s.Update(g.update))
	})

	buf.WriteString("// Generated by protoc-ddl.\n")
	buf.WriteString(fmt.Sprintf("// protoc-gen-dao-mock: %s\n", GoDAOGeneratorVersion))
	b, err := src.GoFormat()
	if err != nil {
		r := bufio.NewScanner(strings.NewReader(src.String()))
		line := 1
		for r.Scan() {
			fmt.Fprintf(os.Stderr, "%d: %s\n", line, r.Text())
			line++
		}
		log.Print(err)
		return
	}
	buf.Write(b)
}

func (GoDAOMockGenerator) primaryKeySelect(entityName string, m *schema.Message, args, where, whereArgs []string) string {
	src := newBuffer()
	a := make([]string, 0)
	for _, v := range m.PrimaryKeys {
		a = append(a, "\""+schema.ToLowerCamel(v.Name)+"\":"+schema.ToLowerCamel(v.Name))
	}

	src.Buffer.WriteString("v, err := d.Call(\"Select\",map[string]interface{}{")
	src.Buffer.WriteString(strings.Join(a, ","))
	src.Write("})")
	src.Writef("return v.(*%s.%s), err", entityName, m.Descriptor.GetName())

	return src.String()
}

func (GoDAOMockGenerator) mockPrimaryKeySelect(entityName string, m *schema.Message) *goFunc {
	src := newBuffer()
	funcArgs := make([]*field, 0)
	for _, v := range m.PrimaryKeys {
		funcArgs = append(funcArgs, &field{Name: schema.ToLowerCamel(v.Name), Type: GoDataTypeMap[v.Type]})
	}
	funcArgs = append(funcArgs, &field{Name: "value", Type: entityName + "." + m.Descriptor.GetName(), Pointer: true})

	src.Buffer.WriteString("d.Register(\"Select\", map[string]interface{}{")
	args := make([]string, 0)
	for _, v := range funcArgs[:len(funcArgs)-1] {
		args = append(args, "\""+v.Name+"\":"+v.Name)
	}
	src.Buffer.WriteString(strings.Join(args, ","))
	src.Buffer.WriteString("}, value, nil)")

	return &goFunc{
		Name:     "RegisterSelect",
		Receiver: &field{Name: "d", Type: m.Descriptor.GetName(), Pointer: true},
		Args:     funcArgs,
		Body:     src.String(),
	}
}

func (GoDAOMockGenerator) mockMember(f *goFunc) string {
	v := f.Returns[0].Copy()
	v.Name = "mock" + f.Name
	return v.String()
}

func (GoDAOMockGenerator) selectRowQuery(m *schema.Message, name string, stmt *sqlparser.Select, comp []*schema.Field, cols, args []string, entityName string, single bool) string {
	src := newBuffer()
	a := make([]string, 0)
	for _, v := range args[:len(args)-1] {
		s := strings.Split(v, " ")
		a = append(a, fmt.Sprintf("\"%s\":%s", s[0], s[0]))
	}
	src.Writef("v, err := d.Call(\"%s\", map[string]interface{}{%s})", name, strings.Join(a, ","))
	if single {
		src.Writef("return v.(*%s.%s), err", entityName, m.Descriptor.GetName())
	} else {
		src.Writef("return v.([]*%s.%s), err", entityName, m.Descriptor.GetName())
	}

	return src.String()
}

func (GoDAOMockGenerator) mockSelectQuery(entityName string, m *schema.Message, selectFunc *goFunc) *goFunc {
	args := selectFunc.Args[1 : len(selectFunc.Args)-1]
	if selectFunc.Returns[0].Slice {
		args = append(args, &field{Name: "value", Type: entityName + "." + m.Descriptor.GetName(), Pointer: true, Slice: true})
	} else {
		args = append(args, &field{Name: "value", Type: entityName + "." + m.Descriptor.GetName(), Pointer: true})
	}

	src := newBuffer()
	src.Buffer.WriteString(fmt.Sprintf("d.Register(\"%s\", map[string]interface{}{", selectFunc.Name))
	a := make([]string, 0)
	for _, v := range args[:len(args)-1] {
		a = append(a, fmt.Sprintf("\"%s\":%s", v.Name, v.Name))
	}
	src.Buffer.WriteString(strings.Join(a, ","))
	src.Buffer.WriteString("}, value, nil)")

	return &goFunc{
		Name:     "Register" + selectFunc.Name,
		Body:     src.String(),
		Args:     args,
		Receiver: &field{Name: "d", Type: m.Descriptor.GetName(), Pointer: true},
	}
}

func (GoDAOMockGenerator) create(m *schema.Message, f *goFunc) string {
	src := newBuffer()

	a := make([]string, 0)
	for _, v := range f.Args[1 : len(f.Args)-1] {
		a = append(a, fmt.Sprintf("\"%s\":%s", v.Name, v.Name))
	}

	src.Writef("_, _ = d.Call(\"%s\", map[string]interface{}{%s})", f.Name, strings.Join(a, ","))
	src.Writef("return %s, nil", schema.ToLowerCamel(m.Descriptor.GetName()))

	return src.String()
}

func (GoDAOMockGenerator) delete(m *schema.Message, f *goFunc, _, _ []string) string {
	src := newBuffer()

	a := make([]string, 0)
	for _, v := range f.Args[1 : len(f.Args)-1] {
		a = append(a, fmt.Sprintf("\"%s\":%s", v.Name, v.Name))
	}

	src.Writef("_, _ = d.Call(\"%s\", map[string]interface{}{%s})", f.Name, strings.Join(a, ","))
	src.Writef("return nil")

	return src.String()
}

func (GoDAOMockGenerator) update(m *schema.Message, f *goFunc) string {
	src := newBuffer()

	a := make([]string, 0)
	for _, v := range f.Args[1 : len(f.Args)-1] {
		a = append(a, fmt.Sprintf("\"%s\":%s", v.Name, v.Name))
	}

	src.Writef("_, _ = d.Call(\"%s\", map[string]interface{}{%s})", f.Name, strings.Join(a, ","))
	src.Writef("return nil")

	return src.String()
}
