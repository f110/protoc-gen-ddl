package schema

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"

	"go.f110.dev/protoc-ddl"
)

const (
	TimestampType = ".google.protobuf.Timestamp"
)

type Column struct {
	Name     string
	DataType string
	TypeName string
	Size     int
	Null     bool
	Sequence bool
	Default  string
}

type Table struct {
	Name           string
	Columns        []Column
	PrimaryKey     []string
	PrimaryKeyType string
	Indexes        []*ddl.IndexOption
	Engine         string
	WithTimestamp  bool

	packageName string
	descriptor  *descriptor.DescriptorProto
}

type DDLOption struct {
	Dialect    string
	OutputFile string
}

type EntityOption struct {
	Lang       string
	OutputFile string
}

func ParseInput(in io.Reader) (*plugin_go.CodeGeneratorRequest, error) {
	buf, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	var input plugin_go.CodeGeneratorRequest
	err = proto.Unmarshal(buf, &input)
	if err != nil {
		return nil, err
	}

	return &input, nil
}

func ProcessDDL(req *plugin_go.CodeGeneratorRequest) (DDLOption, *Messages) {
	return parseOptionDDL(req.GetParameter()), parseTables(req)
}

type Message struct {
	Descriptor  *descriptor.DescriptorProto
	Package     string
	FullName    string
	TableName   string
	Fields      *Fields
	PrimaryKeys []*Field
	Indexes     []*ddl.IndexOption
	Engine      string
}

func (m *Message) String() string {
	s := make([]string, 0, m.Fields.Len()+2)
	s = append(s, m.FullName)
	m.Fields.Each(func(f *Field) {
		s = append(s, "\t"+f.String())
	})
	s = append(s, fmt.Sprintf("Primary Key: %v", m.PrimaryKeys))

	return strings.Join(s, "\n")
}

type Messages struct {
	messages []*Message
	table    map[string]*Message
}

func NewMessages(messages []*Message) *Messages {
	table := make(map[string]*Message)
	for _, v := range messages {
		table[v.FullName] = v
	}
	return &Messages{messages: messages, table: table}
}

func (m *Messages) Each(fn func(m *Message)) {
	l := make([]*Message, len(m.messages))
	copy(l, m.messages)

	for _, v := range l {
		fn(v)
	}
}

func (m *Messages) Denormalize() {
	m.denormalizePrimaryKey()
	m.denormalizeFields()
}

func (m *Messages) denormalizePrimaryKey() {
	for !m.isPrimaryKeyDenormalized() {
		for _, msg := range m.messages {
			var newPrimaryKey []*Field
			for _, f := range msg.PrimaryKeys {
				if isPrimitiveType(f.Type) {
					newPrimaryKey = append(newPrimaryKey, f)
					continue
				}

				if v, ok := m.table[f.Type]; ok {
					newFields := make([]*Field, 0)
					for _, primaryKey := range v.PrimaryKeys {
						newField := primaryKey.Reference()
						newField.Name = f.Name + "_" + newField.Name
						newFields = append(newFields, newField)

						newPrimaryKey = append(newPrimaryKey, newField)
					}
					msg.Fields.Replace(f.Name, newFields...)
				}
			}
			msg.PrimaryKeys = newPrimaryKey
		}
	}
}

func (m *Messages) denormalizeFields() {
	for !m.isFieldsDenormalized() {
		for _, msg := range m.messages {
			msg.Fields.Each(func(f *Field) {
				if isPrimitiveType(f.Type) {
					return
				}

				if v, ok := m.table[f.Type]; ok {
					newFields := make([]*Field, 0)
					for _, primaryKey := range v.PrimaryKeys {
						newField := primaryKey.Reference()
						newField.Name = f.Name + "_" + newField.Name
						if f.Ext != nil && f.Ext.Null {
							newField.Null = true
						}
						newFields = append(newFields, newField)
					}
					msg.Fields.Replace(f.Name, newFields...)
				}
			})
		}
	}
}

func (m *Messages) isPrimaryKeyDenormalized() bool {
	for _, msg := range m.messages {
		for _, f := range msg.PrimaryKeys {
			if !isPrimitiveType(f.Type) {
				return false
			}
		}
	}

	return true
}

func (m *Messages) isFieldsDenormalized() bool {
	for _, msg := range m.messages {
		ok := true
		msg.Fields.Each(func(f *Field) {
			if !isPrimitiveType(f.Type) {
				ok = false
				return
			}
		})
		if !ok {
			return ok
		}
	}

	return true
}

func isPrimitiveType(typ string) bool {
	switch typ {
	case "TYPE_INT32", "TYPE_INT64", "TYPE_UINT32", "TYPE_UINT64", "TYPE_SINT32", "TYPE_SINT64", "TYPE_FIXED32", "TYPE_FIXED64",
		"TYPE_SFIXED32", "TYPE_SFIXED64":
		fallthrough
	case "TYPE_FLOAT", "TYPE_DOUBLE":
		fallthrough
	case "TYPE_STRING", "TYPE_BYTES", "TYPE_BOOL", TimestampType:
		return true
	default:
		return false
	}
}

func (m *Messages) String() string {
	s := make([]string, len(m.messages))
	for i := range m.messages {
		s[i] = m.messages[i].String()
	}
	return strings.Join(s, "\n")
}

type Field struct {
	Descriptor   *descriptor.FieldDescriptorProto
	Ext          *ddl.ColumnOptions
	Name         string
	Type         string
	OptionalType string
	Size         int
	Null         bool
	Sequence     bool
	Default      string
}

func (f *Field) Copy() *Field {
	n := &Field{}
	*n = *f
	return n
}

func (f *Field) Reference() *Field {
	n := &Field{}
	*n = *f
	n.Sequence = false
	return n
}

func (f *Field) String() string {
	return fmt.Sprintf("%s %s", f.Name, f.Type)
}

type Fields struct {
	list  []*Field
	table map[string]*Field
}

func NewFields(fields []*Field) *Fields {
	table := make(map[string]*Field)
	for _, v := range fields {
		table[v.Name] = v
	}

	return &Fields{list: fields, table: table}
}

func (f *Fields) Len() int {
	return len(f.list)
}

func (f *Fields) Each(fn func(f *Field)) {
	l := make([]*Field, len(f.list))
	copy(l, f.list)

	for _, v := range l {
		fn(v)
	}
}

func (f *Fields) Replace(oldName string, newField ...*Field) {
	newList := make([]*Field, 0, len(f.list))
	for _, v := range f.list {
		if v.Name == oldName {
			newList = append(newList, newField...)
			continue
		}

		newList = append(newList, v)
	}
	f.list = newList

	delete(f.table, oldName)
}

func (f *Fields) Get(name string) *Field {
	return f.table[name]
}

func (f *Fields) String() string {
	s := make([]string, f.Len())
	i := 0
	f.Each(func(f *Field) {
		s[i] = f.String()
		i++
	})
	return strings.Join(s, "\n")
}

func ProcessEntity(req *plugin_go.CodeGeneratorRequest) (EntityOption, *descriptor.FileOptions, *Messages) {
	files := make(map[string]*descriptor.FileDescriptorProto)
	for _, f := range req.ProtoFile {
		files[f.GetName()] = f
	}

	var opt *descriptor.FileOptions
	for _, filename := range req.FileToGenerate {
		opt = files[filename].GetOptions()
	}

	targetMessages := make([]*Message, 0)
	for _, filename := range req.FileToGenerate {
		f := files[filename]
		for _, m := range f.GetMessageType() {
			opt := m.GetOptions()
			e, err := proto.GetExtension(opt, ddl.E_Table)
			if err == proto.ErrMissingExtension {
				continue
			}
			ext := e.(*ddl.TableOptions)

			foundFields := make([]*Field, 0)
			for _, v := range m.Field {
				switch v.GetType() {
				case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
					foundFields = append(foundFields, &Field{
						Name: v.GetName(),
						Type: v.GetTypeName(),
					})
				default:
					foundFields = append(foundFields, &Field{
						Name: v.GetName(),
						Type: v.GetType().String(),
					})
				}
			}
			if ext.WithTimestamp {
				foundFields = append(foundFields,
					&Field{Name: "created_at", Type: TimestampType},
					&Field{Name: "updated_at", Type: TimestampType},
				)
			}

			fields := NewFields(foundFields)
			primaryKey := make([]*Field, 0)
			for _, v := range ext.PrimaryKey {
				primaryKey = append(primaryKey, fields.Get(v))
			}

			targetMessages = append(targetMessages, &Message{
				Descriptor:  m,
				Package:     "." + f.GetPackage(),
				FullName:    "." + f.GetPackage() + "." + m.GetName(),
				Fields:      fields,
				PrimaryKeys: primaryKey,
			})
		}
	}

	msgs := NewMessages(targetMessages)
	msgs.Denormalize()
	return parseOptionEntity(req.GetParameter()), opt, msgs
}

func parseTables(req *plugin_go.CodeGeneratorRequest) *Messages {
	tables := make([]*Table, 0)

	files := make(map[string]*descriptor.FileDescriptorProto)
	for _, f := range req.ProtoFile {
		files[f.GetName()] = f
	}

	targetMessages := make([]*Message, 0)
	for _, fileName := range req.FileToGenerate {
		f := files[fileName]
		for _, m := range f.GetMessageType() {
			opt := m.GetOptions()
			_, err := proto.GetExtension(opt, ddl.E_Table)
			if err == proto.ErrMissingExtension {
				continue
			}

			targetMessages = append(targetMessages, &Message{
				Descriptor: m,
				Package:    "." + f.GetPackage(),
				FullName:   "." + f.GetPackage() + "." + m.GetName(),
			})
		}
	}

	msgs := NewMessages(targetMessages)
	msgs.Each(func(m *Message) {
		e, err := proto.GetExtension(m.Descriptor.GetOptions(), ddl.E_Table)
		if err != nil {
			return
		}
		ext := e.(*ddl.TableOptions)

		foundFields := make([]*Field, 0)
		for _, v := range m.Descriptor.Field {
			var f *Field
			switch v.GetType() {
			case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
				f = &Field{
					Descriptor: v,
					Name:       v.GetName(),
					Type:       v.GetTypeName(),
				}
			default:
				f = &Field{
					Descriptor: v,
					Name:       v.GetName(),
					Type:       v.GetType().String(),
				}
			}
			foundFields = append(foundFields, f)

			e, err := proto.GetExtension(v.GetOptions(), ddl.E_Column)
			if err != nil {
				continue
			}
			f.Ext = e.(*ddl.ColumnOptions)
		}
		if ext.WithTimestamp {
			foundFields = append(foundFields,
				&Field{Name: "created_at", Type: TimestampType},
				&Field{Name: "updated_at", Type: TimestampType, Null: true},
			)
		}

		fields := NewFields(foundFields)
		primaryKey := make([]*Field, 0)
		for _, v := range ext.PrimaryKey {
			primaryKey = append(primaryKey, fields.Get(v))
		}

		m.TableName = ToSnake(m.Descriptor.GetName())
		if ext.TableName != "" {
			m.TableName = ext.TableName
		}
		m.Fields = fields
		m.PrimaryKeys = primaryKey
		m.Indexes = ext.GetIndexes()
		m.Engine = ext.Engine

		m.Fields.Each(func(f *Field) {
			e, err := proto.GetExtension(f.Descriptor.GetOptions(), ddl.E_Column)
			if err != nil {
				return
			}

			ext := e.(*ddl.ColumnOptions)
			f.Sequence = ext.Sequence
			f.Null = ext.Null
			f.Default = ext.Default
			f.Size = int(ext.Size)
			f.OptionalType = ext.Type
		})
	})

	for _, t := range tables {
		for _, f := range t.descriptor.GetField() {
			t.Columns = append(t.Columns, convertToColumn(f, tables))
		}

		if t.WithTimestamp {
			t.Columns = append(t.Columns,
				Column{Name: "created_at", DataType: TimestampType},
				Column{Name: "updated_at", DataType: TimestampType, Null: true},
			)
		}
	}

	msgs.Denormalize()
	return msgs
}

func convertToTable(packageName string, msg *descriptor.DescriptorProto, opt *ddl.TableOptions) *Table {
	t := &Table{descriptor: msg, packageName: packageName, WithTimestamp: opt.WithTimestamp}

	if opt.TableName != "" {
		t.Name = opt.TableName
	} else {
		t.Name = ToSnake(msg.GetName())
	}

	if len(opt.PrimaryKey) > 0 {
		t.PrimaryKey = opt.PrimaryKey
		for _, f := range msg.GetField() {
			if f.GetName() == t.PrimaryKey[0] {
				t.PrimaryKeyType = f.GetType().String()
				break
			}
		}
	}

	t.Indexes = opt.GetIndexes()
	t.Engine = opt.Engine

	return t
}

func convertToColumn(field *descriptor.FieldDescriptorProto, tables []*Table) Column {
	f := Column{}

	f.Name = field.GetName()
	if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		foreignTable := false
		for _, t := range tables {
			if "."+t.packageName+"."+t.descriptor.GetName() == field.GetTypeName() {
				foreignTable = true
				f.Name += "_" + t.PrimaryKey[0]
				f.DataType = t.PrimaryKeyType
			}
		}
		if !foreignTable {
			f.DataType = field.GetTypeName()
		}
	} else {
		f.DataType = field.GetType().String()
	}

	opt := field.GetOptions()
	if opt == nil {
		return f
	}
	e, err := proto.GetExtension(opt, ddl.E_Column)
	if err == proto.ErrMissingExtension {
		return f
	}
	ext := e.(*ddl.ColumnOptions)
	f.Sequence = ext.Sequence
	f.Null = ext.Null
	f.Default = ext.Default
	f.Size = int(ext.Size)
	f.TypeName = ext.Type

	return f
}

func parseOptionDDL(p string) DDLOption {
	opt := DDLOption{OutputFile: "sql/schema.sql"}
	params := strings.Split(p, ",")
	for _, param := range params {
		s := strings.SplitN(param, "=", 2)
		if len(s) == 1 {
			opt.OutputFile = s[0]
			continue
		}
		key := s[0]
		value := s[1]

		switch key {
		case "dialect":
			opt.Dialect = value
		}
	}
	return opt
}

func parseOptionEntity(p string) EntityOption {
	opt := EntityOption{}
	params := strings.Split(p, ",")
	for _, param := range params {
		s := strings.SplitN(param, "=", 2)
		if len(s) == 1 {
			opt.OutputFile = s[0]
			continue
		}
		key := s[0]
		value := s[1]

		switch key {
		case "lang":
			opt.Lang = value
		}
	}

	return opt
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnake(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

var link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")

func ToCamel(str string) string {
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}
