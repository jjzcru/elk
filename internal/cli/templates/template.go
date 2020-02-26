package templates

// Config template use for config.yml
var Config = `path: "{{.}}"`

// Elk template use for elk.yml
var Elk = `version: '1'
tasks: {{ range $name, $task := .Tasks }}
  {{ $name }}:
    description: '{{$task.Description}}'
    cmds: {{range $cmd := $task.Cmds}}
      - {{$cmd}}{{end}}
{{end}}
`

// Installation template use when installing elk
var Installation = `
This will create a default elk file

It only covers just a few tasks. 

The installation will include some default events like 'shutdown' 
or 'restart' just to get started but you will be able to add more 
events in the configuration file.

`
