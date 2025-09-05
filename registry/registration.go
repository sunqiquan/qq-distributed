package registry

type ServiceName string

type Registration struct {
	ServiceName      ServiceName   `json:"serviceName"`
	ServiceUrl       string        `json:"serviceUrl"`
	RequiredServices []ServiceName `json:"requiredServices"`
	ServiceUpdateUrl string        `json:"serviceUpdateUrl"`
}

const (
	LogService     = ServiceName("LogService")
	StudentService = ServiceName("StudentService")
)

type patchEntry struct {
	Name ServiceName
	Url  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
