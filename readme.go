package main

// Inspired from [awsls](https://github.com/jckuester/awsls) readme file
// TODO:  should be executed when building for releases
import (
	"os"
	"github.com/rfc2119/aws-tui/common"
	"text/template"
)

var readmeTable = `

# Current Working Services

| Service Name | Implemented | Description |
| :----------: | :---------: | :---------: |
{{ range $serviceName, $desc := . }}|{{$desc.Name}} | {{if $desc.Available}} ✓ {{end}} | {{$desc.Description}}|
{{ end }}
`

func AppendReadMeTable(readmeFileName string) error {
	f, err := os.OpenFile(readmeFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	templ := template.Must(template.New("readmeTable").Parse(readmeTable))
	if err := templ.Execute(f, common.AWServicesDescriptions); err != nil {
		return err
	}
	return nil
}

// {{range pipeline}} T1 {{end}}
// 	The value of the pipeline must be an array, slice, map, or channel.
// 	If the value of the pipeline has length zero, nothing is output;
// 	otherwise, dot is set to the successive elements of the array,
// 	slice, or map and T1 is executed. If the value is a map and the
// 	keys are of basic type with a defined order, the elements will be
// 	visited in sorted key order.

// {{with pipeline}} T1 {{end}}
// 	If the value of the pipeline is empty, no output is generated;
// 	otherwise, dot is set to the value of the pipeline and T1 is
// 	executed.

// range $index, $element := pipeline
// Note that if there is only one variable, it is assigned the element; this is opposite to the convention in Go range clauses
