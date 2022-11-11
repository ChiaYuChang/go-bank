package gen

type Python struct {
	Package         string `yaml:"package"`
	Out             string `yaml:"out"`
	*PythonEmitOpts `yaml:",inline"`
}

type PythonEmitOpts struct {
	ExactTableNames bool `yaml:"emit_exact_table_names,omitempty"`
	SyncQuerier     bool `yaml:"emit_sync_querier,omitempty"`
	AsyncQuerier    bool `yaml:"emit_async_querier,omitempty"`
	PydanticModels  bool `yaml:"emit_pydantic_models,omitempty"`
}
