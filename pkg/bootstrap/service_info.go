package bootstrap

import "os"

type ServiceInfo struct {
	Name     string
	Version  string
	Id       string
	Metadata map[string]string
}

// NewServiceInfo new service info
func NewServiceInfo(name, version, id string) *ServiceInfo {
	if id == "" {
		id, _ = os.Hostname()
	}
	return &ServiceInfo{
		Name:     name,
		Version:  version,
		Id:       id,
		Metadata: map[string]string{},
	}
}

// GetInstanceId returns the service instance id
func (s *ServiceInfo) GetInstanceId() string {
	return s.Id + "." + s.Name
}

// SetMetaData stores the kv
func (s *ServiceInfo) SetMetaData(k, v string) {
	s.Metadata[k] = v
}
