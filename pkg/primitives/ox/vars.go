package ox

import (
	"html/template"
)

type Vars struct {
	Map map[string]string
	Cmd string
}

func (v *Vars) Write(data []byte) (n int, err error) {
	v.Cmd += string(data[:])
	return len(data), nil
}

func (v *Vars) Process(cmd string) (string, error) {
	t, err := template.New("ox").Parse(cmd)
	if err != nil {
		return "", err
	}

	err = t.Execute(v, v.Map)
	if err != nil {
		return "", err
	}

	return v.Cmd, nil
}

func GetCmdFromVars(vars map[string]string, cmd string) (string, error) {
	v := Vars{
		Map: vars,
	}

	return v.Process(cmd)
}
