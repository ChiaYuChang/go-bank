package gen

type Json struct {
	Out      string `yaml:"out"`
	FileName string `yaml:"filename"`
	Indent   string `yaml:"indent"`
}

func NewJsonGen(out, filename string) *Json {
	return &Json{
		Out:      out,
		FileName: filename,
		Indent:   "    ",
	}
}
