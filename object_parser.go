package object_parser

import (
	"reflect"
	"strings"
)

type ObjectParser struct {
	object          interface{}
	objectType      reflect.Type
	objectValue     reflect.Value
	objectFieldTags map[FieldName]*FieldTag
}

type FieldName string

type TagName string

type FieldTag struct {
	Tags map[TagName]*FieldSpecifyTag
}

type FieldSpecifyTag struct {
	Values []string
}

func NewObjectParser(object interface{}) *ObjectParser {
	t, v := reflect.TypeOf(object), reflect.ValueOf(object)
	tags := getFieldsTag(t)

	return &ObjectParser{
		object:          object,
		objectType:      t,
		objectValue:     v,
		objectFieldTags: tags,
	}
}

// Ignore field
// 1. field wasn't granted target tag
// 2. field is ptr & nil
// 3. target tag contains omitempty & field is zero value
// if field type implements interface `Convert() interface{}` by value receiver, Auto convert field type to other type
func (objectParser *ObjectParser) TagValueMap(targetTag string) map[string]interface{} {
	namedParam := make(map[string]interface{})
	t, v := objectParser.getTypeAndValue()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		iValue := objectParser.getInterfaceValue(field, v.Field(i), targetTag)
		if iValue == nil {
			continue
		}

		namedParam[objectParser.objectFieldTags[FieldName(field.Name)].Tags[TagName(targetTag)].Values[0]] = iValue
	}

	return namedParam
}

func (objectParser *ObjectParser) getInterfaceValue(field reflect.StructField, value reflect.Value, targetTag string) interface{} {
	vValue := getValue(value)

	if field.Type.Kind() == reflect.Ptr && value.IsNil() {
		return nil
	}

	ofTags := objectParser.objectFieldTags[FieldName(field.Name)]
	if ofTags == nil {
		return nil
	}

	if ofTags.Tags == nil {
		return nil
	}

	if ofTags.Tags[TagName(targetTag)] == nil {
		return nil
	}

	if len(ofTags.Tags[TagName(targetTag)].Values) == 0 {
		return nil
	}

	if tagContains(ofTags.Tags[TagName(targetTag)].Values, "omitempty") && vValue.IsZero() {
		return nil
	}

	if cast, ok := vValue.Interface().(interface{ Convert() interface{} }); ok {
		return cast.Convert()
	}

	return vValue.Interface()
}

func (objectParser *ObjectParser) getTypeAndValue() (reflect.Type, reflect.Value) {
	if objectParser.objectType.Kind() == reflect.Ptr {
		return objectParser.objectType.Elem(), objectParser.objectValue.Elem()
	}
	return objectParser.objectType, objectParser.objectValue
}

func getFieldsTag(t reflect.Type) map[FieldName]*FieldTag {
	fieldsTagMap := make(map[FieldName]*FieldTag)
	tt := getType(t)
	for i := 0; i < tt.NumField(); i++ {
		field := tt.Field(i)
		tag := getFieldTag(field)
		fieldsTagMap[FieldName(field.Name)] = tag
	}
	return fieldsTagMap
}

func getFieldTag(field reflect.StructField) *FieldTag {
	tags := make([]string, 0, 0)
	separatedTag := strings.Split(string(field.Tag), " ")
	for _, tag := range separatedTag {
		tags = append(tags, tag)
	}

	ft := new(FieldTag)
	ft.Tags = make(map[TagName]*FieldSpecifyTag)
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		tagNameAndValues := strings.Split(tag, ":")
		name, values := tagNameAndValues[0], strings.Trim(tagNameAndValues[1], "\"")
		fst := new(FieldSpecifyTag)
		fst.Values = append(fst.Values, strings.Split(values, ",")...)
		ft.Tags[TagName(name)] = fst
	}

	return ft
}

func getType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func getValue(value reflect.Value) reflect.Value {
	if value.Type().Kind() == reflect.Ptr {
		return value.Elem()
	}
	return value
}

func tagContains(tagValues []string, element string) bool {
	for _, tag := range tagValues {
		if tag == element {
			return true
		}
	}
	return false
}
