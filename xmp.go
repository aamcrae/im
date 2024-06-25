package im

import (
	"encoding/xml"
	"fmt"
)

type description struct {
	namespace map[string]string
	data      map[string]string
}

type rdf struct {
	Description description `xml:"Description"`
}

type xmpXml struct {
	XMLName xml.Name `xml:"xmpmeta"`
	Rdf     rdf      `xml:"RDF"`
}

func (d *description) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	if d.namespace == nil {
		d.namespace = make(map[string]string)
	}
	if d.data == nil {
		d.data = make(map[string]string)
	}
	for _, attr := range start.Attr {
		// Check for a new namespace declaration.
		if attr.Name.Space == "xmlns" {
			d.namespace[attr.Value] = attr.Name.Local
		}
		n, ok := d.namespace[attr.Name.Space]
		if ok {
			d.data[n+"."+attr.Name.Local] = attr.Value
		}
	}
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		switch tok := token.(type) {
		case xml.StartElement:
			var ele any
			err := decoder.DecodeElement(&ele, &tok)
			if err != nil {
				return err
			}
		case xml.EndElement:
			return nil
		}
	}
	return nil
}

func (im *Imeta) addXmp(b []byte) error {
	var v xmpXml
	err := xml.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	for k, v := range v.Rdf.Description.data {
		fmt.Printf("%s = %s\n", k, v)
	}
	return nil
}
