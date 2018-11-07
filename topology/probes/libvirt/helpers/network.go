package helpers

import (
	"encoding/xml"
	"fmt"

	"github.com/skydive-project/skydive/logging"
)

// Target describes the XML coding of the target of an interface in libvirt
// Address describes the XML coding of the pci addres of an interface in libvirt
type Address struct {
	Type     string `xml:"type,attr"`
	Domain   string `xml:"domain,attr"`
	Bus      string `xml:"bus,attr"`
	Slot     string `xml:"slot,attr"`
	Function string `xml:"function,attr"`
}

type Interface struct {
	Driver struct {
		Name string `xml:"name,attr"`
	} `xml:"driver"`
	Mac struct {
		Address string `xml:"address,attr"`
	} `xml:"mac"`
	InterfaceType string `xml:"type,attr"`
	Source        struct {
		Bridge string `xml:"bridge,attr"`
	} `xml:"source"`
	Target struct {
		Dev string `xml:"dev,attr"`
	} `xml:"target"`
	Alias struct {
		Name string `xml:"name,attr"`
	} `xml:"alias"`
	Address Address `xml:"address"`
}

type Domain struct {
	Interfaces []Interface `xml:"devices>interface"`
}

func GetVMInterfaces(XMLDesc string) ([]interface{}, error) {
	domain := Domain{}
	if err := xml.Unmarshal([]byte(XMLDesc), &domain); err != nil {
		return nil, fmt.Errorf("XML parsing error: %s", err)
	}

	interfaces := make([]interface{}, len(domain.Interfaces))
	for i, interf := range domain.Interfaces {
		address := interf.Address
		formatted := fmt.Sprintf(
			"%s:%s.%s.%s.%s", address.Type, address.Domain, address.Bus,
			address.Slot, address.Function)
		metadata := map[string]interface{}{
			"MAC":     interf.Mac.Address,
			"Address": formatted,
			"Alias":   interf.Alias.Name,
		}
		if interf.Target.Dev != "" {
			metadata["Target"] = map[string]interface{}{
				"Dev": interf.Target.Dev,
			}
		}
		interfaces[i] = metadata
	}

	logging.GetLogger().Debugf("Found interfaces %s", interfaces)
	return interfaces, nil
}
