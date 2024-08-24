package setting

import (
	"encoding/json"
	"fmt"
	"io"
)

type ProgramSetting struct {
	Key string `json:"key"`
}

func (p *ProgramSetting) String() string {
	return fmt.Sprintf("%v", *p)
}

func (p *ProgramSetting) Parse(r io.Reader) error {
	return json.NewDecoder(r).Decode(p)
}

func (p *ProgramSetting) Write(w io.Writer) error {
	return json.NewEncoder(w).Encode(p)
}
