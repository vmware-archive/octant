package plister

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

const docType = `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">` + "\n"

var (
	keyElem = xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "key",
		},
	}
	valueElem = xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "string",
		},
	}
)

type (
	Dict struct {
		XMLName xml.Name `xml:"dict"`
		Items   []*DictItem
	}
	DictItem struct {
		Key   string
		Value interface{}
	}
	Slice struct {
		XMLName xml.Name      `xml:"array"`
		Items   []interface{} `xml:"string"`
	}
	SliceDict struct {
		XMLName xml.Name `xml:"dict"`
		Items   []*Dict
	}
	InfoPlist struct {
		XMLName xml.Name `xml:"plist"`
		Version string   `xml:"version,attr"`
		Dict    *Dict
	}
	Array struct {
		String []string `xml:"string"`
	}
)

func (di DictItem) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if di.Key == "" {
		return nil
	}
	if err := e.EncodeElement(di.Key, keyElem); err != nil {
		return err
	}
	switch di.Value.(type) {
	case bool:
		return e.EncodeElement("", xml.StartElement{
			Name: xml.Name{Space: "", Local: fmt.Sprintf("%t", di.Value)},
		})
	case *Dict, *Slice, *SliceDict:
		return e.Encode(di.Value)
	default:
		return e.EncodeElement(di.Value, valueElem)
	}
}

func (ip *InfoPlist) Get(key string) interface{} {
	for _, i := range ip.Dict.Items {
		if i.Key == key {
			return i.Value
		}
	}
	return nil
}

func (ip *InfoPlist) Set(key string, value interface{}) {
	for _, i := range ip.Dict.Items {
		if i.Key == key {
			i.Value = value
			return
		}
	}
	ip.Dict.Items = append(ip.Dict.Items, &DictItem{key, value})
}

func (ip *InfoPlist) Delete(key string) {
	for i := range ip.Dict.Items {
		if ip.Dict.Items[i].Key == key {
			ip.Dict.Items = append(ip.Dict.Items[:i], ip.Dict.Items[i+1:]...)
			return
		}
	}
}

func MapToInfoPlist(m map[string]interface{}) *InfoPlist {
	return &InfoPlist{
		Version: "1.0",
		Dict:    &Dict{Items: mapToDictItems(m)},
	}
}

func mapToDictItems(m map[string]interface{}) []*DictItem {
	items := make([]*DictItem, 0)
	for k, v := range m {
		if vm, ok := v.(map[string]interface{}); ok {
			v = &Dict{Items: mapToDictItems(vm)}
		} else if vm, ok := v.([]map[string]interface{}); ok {
			v = &SliceDict{Items: arrayToDictSlice(vm)}
		} else if vm, ok := v.([]interface{}); ok {
			v = &Slice{Items: arrayToSlice(vm)}
		}
		items = append(items, &DictItem{k, v})
	}
	return items
}

func arrayToDictSlice(array []map[string]interface{}) []*Dict {
	items := make([]*Dict, 0)
	for _, i := range array {
		items = append(items, &Dict{Items: mapToDictItems(i)})
	}
	return items
}

func arrayToSlice(array []interface{}) []interface{} {
	items := make([]interface{}, 0)
	for _, i := range array {
		if m, ok := i.(map[string]interface{}); ok {
			i = &Dict{Items: mapToDictItems(m)}
		}
		items = append(items, i)
	}
	return items
}

func Fprint(w io.Writer, data *InfoPlist) error {
	body, err := xml.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(append([]byte(xml.Header+docType), body...))
	return err
}

func Generate(path string, data *InfoPlist) error {
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	return Fprint(fp, data)
}

func GenerateFromMap(path string, data map[string]interface{}) error {
	return Generate(path, MapToInfoPlist(data))
}

func Parse(path string) (*InfoPlist, error) {
	result := &InfoPlist{
		XMLName: xml.Name{Local: "plist"},
		Version: "1.0",
		Dict:    &Dict{XMLName: xml.Name{Local: "dict"}},
	}
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	dec := xml.NewDecoder(fp)
	dec.Strict = false
	for {
		token, _ := dec.Token()
		if token == nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			switch start.Name.Local {
			case "dict":
				result.Dict, err = decodeDict(dec)
				if err != nil {
					return result, err
				}
			}
		}
	}
	return result, nil
}

func decodeDict(dec *xml.Decoder) (result *Dict, err error) {
	var workingKey string
	result = &Dict{}
	for {
		token, _ := dec.Token()
		if token == nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			switch start.Name.Local {
			case "key":
				var k string
				if err = dec.DecodeElement(&k, &start); err != nil {
					return
				}
				workingKey = k
			default:
				v, err := decodeValue(dec, token)
				if err != nil {
					return result, err
				}
				result.Items = append(result.Items, &DictItem{
					Key:   workingKey,
					Value: v,
				})
				workingKey = ""
			}
		case xml.EndElement:
			if start.Name.Local == "dict" {
				return
			}
		}
	}
	return
}

func decodeValue(dec *xml.Decoder, token xml.Token) (v interface{}, err error) {
	switch start := token.(type) {
	case xml.StartElement:
		switch start.Name.Local {
		case "string":
			var s string
			return s, dec.DecodeElement(&s, &start)
		case "integer":
			var i int
			return i, dec.DecodeElement(&i, &start)
		case "true", "false":
			return start.Name.Local == "true", nil
		case "dict":
			return decodeDict(dec)
		case "array":
			ai := &Slice{}
			for {
				arrToken, _ := dec.Token()
				if arrToken == nil {
					break
				}
				switch arrToken.(type) {
				case xml.StartElement:
					v, err := decodeValue(dec, arrToken)
					if err != nil {
						return ai, nil
					}
					ai.Items = append(ai.Items, v)
				case xml.EndElement:
					if start.Name.Local == "array" {
						return ai, nil
					}
				}
			}
		}
	}
	return
}
