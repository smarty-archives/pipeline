package messaging

import (
	"reflect"
	"strings"
)

type ReflectionDiscovery struct {
	prefix string
}

func NewReflectionDiscovery(prefix string) *ReflectionDiscovery {
	return &ReflectionDiscovery{prefix: prefix}
}

func (this *ReflectionDiscovery) Discover(instance interface{}) (string, error) {
	if instance == nil {
		return "", MessageTypeDiscoveryError
	}

	reflectType := reflect.TypeOf(instance)
	if name := reflectType.Name(); len(name) > 0 {
		return this.prefix + strings.ToLower(name), nil
	}

	name := reflectType.String()
	index := strings.LastIndex(name, ".")

	if index == -1 {
		return "", MessageTypeDiscoveryError
	}

	suffix := strings.ToLower(name[index+1:])
	if strings.HasPrefix(name, "*") {
		return "*" + this.prefix + suffix, nil
	} else {
		return this.prefix + suffix, nil
	}
}
