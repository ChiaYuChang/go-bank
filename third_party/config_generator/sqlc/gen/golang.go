package gen

import (
	"gopkg.in/yaml.v3"
)

type SQLPkg string

const (
	Pgxv4 SQLPkg = "pgx/v4"
	DBSQL SQLPkg = "database/sql"
)

type Go struct {
	Package           string `yaml:"package"`
	Out               string `yaml:"out"`
	SQLPackage        SQLPkg `yaml:"sql_package"`
	*GoEmitOpts       `yaml:",inline"`
	JsonTagsCaseStyle JsonTagsCaseStyle `yaml:"json_tags_case_style"`
	*GoOutputOpts     `yaml:",inline"`
	Rename            `yaml:"rename,omitempty"`
	Overrides         []*Override `yaml:"overrides,omitempty"`
}

func NewGoGen(pkg, out string) *Go {
	g := &Go{
		Package:    pkg,
		Out:        out,
		SQLPackage: DBSQL,
		GoEmitOpts: nil,
		GoOutputOpts: &GoOutputOpts{
			DBFileName:      "db.go",
			ModelsFileName:  "models.go",
			QuerierFileName: "querier.go",
		},
		JsonTagsCaseStyle: None,
		Rename:            make(map[string]string),
		Overrides:         []*Override{},
	}
	return g
}

type GoEmitOpts struct {
	DBTags                bool `yaml:"emit_db_tags,omitempty"`
	PreparedQureies       bool `yaml:"emit_prepared_queries,omitempty"`
	Interface             bool `yaml:"emit_interface,omitempty"`
	ExactTableNames       bool `yaml:"emit_exact_table_names,omitempty"`
	EmptySlices           bool `yaml:"emit_empty_slices,omitempty"`
	ExportedQueries       bool `yaml:"emit_exported_queries,omitempty"`
	JsonTags              bool `yaml:"emit_json_tags,omitempty"`
	ResultStructPointers  bool `yaml:"emit_result_struct_pointers,omitempty"`
	ParamsStructPointers  bool `yaml:"emit_params_struct_pointers,omitempty"`
	MethodsWithDBArgument bool `yaml:"emit_methods_with_db_argument,omitempty"`
	EnumValidMethod       bool `yaml:"emit_enum_valid_method,omitempty"`
	AllEnumValues         bool `yaml:"emit_all_enum_values,omitempty"`
}

type JsonTagsCaseStyle string

const (
	Camel  JsonTagsCaseStyle = "camel"
	Pascal JsonTagsCaseStyle = "pascal"
	Snake  JsonTagsCaseStyle = "snake"
	None   JsonTagsCaseStyle = "none"
)

func (j JsonTagsCaseStyle) IsZero() bool {
	return j == None
}

func (j JsonTagsCaseStyle) MarshalYAML() (interface{}, error) {
	return string(j), nil
}

func (j *JsonTagsCaseStyle) UnmarshalYAML(value *yaml.Node) error {
	*j = JsonTagsCaseStyle(value.Value)
	return nil
}

type GoOutputOpts struct {
	DBFileName      string `yaml:"output_db_file_name"`
	ModelsFileName  string `yaml:"output_models_file_name"`
	QuerierFileName string `yaml:"output_querier_file_name"`
	FilesSuffix     string `yaml:"output_files_suffix,omitempty"`
}

type Rename map[string]string

type Override struct {
	DBType   string `yaml:"db_type"`
	GoType   string `yaml:"go_type"`
	Nullable bool   `yaml:"nullable,omitempty"`
}
