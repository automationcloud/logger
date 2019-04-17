package logging

type ServiceContext struct {
	Service string `json:"service"`
	Version string `json:"version"`
}
