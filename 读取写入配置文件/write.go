package main

import (
	"fmt"

	"github.com/ghodss/yaml"
)

type RuleGroup struct {
	Groups []SubGroup `json:"groups"`
}
type SubGroup struct {
	Name  string   `json:"name"`
	Rules []Metric `json:"rules"`
}

type Metric struct {
	Expr        string                 `json:"expr"`
	Alert       string                 `json:"alert"`
	For         string                 `json:"for"`
	Labels      map[string]interface{} `json:"labels"`
	Annotations Annotations            `json:"annotations"`
}
type Annotations struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

func main() {
	// Marshal a Person struct to YAML.
	m := Metric{
		Expr:   "cpu >%d",
		Alert:  "cpu",
		For:    "2m",
		Labels: map[string]interface{}{"term": "node"},
		Annotations: Annotations{
			Summary:     "IP:{{$labels.instance}} HOSTID:{{$labels.hostid}}: cpu high eq 10",
			Description: "{{$labels.instance}}: {{$labels.job}} usage is above 10% (current value is:{{ $value }})",
		},
	}
	p := RuleGroup{
		Groups: []SubGroup{
			{

				Name:  "03000200-0400-0500-0006-000700080009",
				Rules: []Metric{m},
			},
		},
	}
	y, err := yaml.Marshal(p)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(y))
	/* Output:
	age: 30
	name: John
	*/

	// Unmarshal the YAML back into a Person struct.
	var p2 RuleGroup
	err = yaml.Unmarshal(y, &p2)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(p2)
	/* Output:
	{John 30}
	*/
}
