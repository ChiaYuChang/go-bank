package sqlc

type SQLCEngn string

const (
	PostgreSQL SQLCEngn = "postgresql"
	MySQL      SQLCEngn = "mysql"
)

type SQLCConfig struct {
	Version int    `yaml:"version"`
	SQL     []*SQL `yaml:"sql"`
}

func NewSQLCConfig() *SQLCConfig {
	return &SQLCConfig{
		Version: 2,
		SQL:     make([]*SQL, 0),
	}
}

type SQL struct {
	Schema               string              `yaml:"schema"`
	Queries              string              `yaml:"queries"`
	Engine               SQLCEngn            `yaml:"engine"`
	StrictFunctionChecks bool                `yaml:"strict_function_checks,omitempty"`
	Gen                  map[GenLanguage]any `yaml:"gen"`
}

func NewSQL(schema, queries string, engine SQLCEngn) *SQL {
	return &SQL{
		Schema:  schema,
		Queries: queries,
		Engine:  engine,
		Gen:     make(map[string]any),
	}
}

type GenLanguage = string

const (
	Golang GenLanguage = "go"
	Kotlin GenLanguage = "kotlin"
	Python GenLanguage = "python"
	Json   GenLanguage = "json"
)
