// package integration allows parsing integrations' JSON payloads
package integration

type JSON struct {
	Data []Data `json:"data"`
}

type Data struct {
	Entity  *Entity  `json:"entity,omitempty"`
	Metrics []Metric `json:"metrics"`
}

type Entity struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Metric map[string]interface{}
