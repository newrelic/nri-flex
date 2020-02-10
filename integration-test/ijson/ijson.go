// package ijson allows parsing integrations' JSON payloads
package ijson

type Entity struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Metric map[string]interface{}

type Data struct {
	Entity *Entity `json:"entity,omitempty"`
	Metrics []Metric `json:"metrics"`
}

type Payload struct {
	Data []Data `json:"data"`
}