package validation

import "strings"

// /////////////////////////////////////////////////////////////////////////////////////////////////

type Language string

var Languages = languages{
	Go:         "go",
	Python:     "python",
	TypeScript: "typescript",
}

type languages struct {
	Go         Language
	Python     Language
	TypeScript Language
}

type LanguageTypes string

///////////////////////////////////////////////////////////////////////////////////////////////////

type SqlDialect string

var SqlDialects = sqlDialects{
	Mysql:      "mysql",
	Postgresql: "postgresql",
	Sqlite:     "sqlite",
}

type sqlDialects struct {
	Mysql      SqlDialect
	Postgresql SqlDialect
	Sqlite     SqlDialect
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type Entity struct {
	Name       string
	Properties []Property
}

func (this *Entity) Type[T PropertyContraint](name string) *typeProperty[T] {
	typed := Type[T](this, name)
	this.Properties = append(this.Properties, typed)
	return typed
}

func (this *Entity) Reference(name string, entity *Entity) *referenceProperty {
	ref := Refernce(this, name, entity)
	this.Properties = append(this.Properties, ref)
	return ref
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type Property interface {
	Parent() *Entity
	Name() string
}

type PropertyContraint interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

///////////////////////////////////////////////////////////////////////////////////////////////////

func Type[T PropertyContraint](parent *Entity, name string) *typeProperty[T] {
	prop := &typeProperty[T]{
		parent:    parent,
		name:      name,
		validator: validator[T]{},
	}

	return prop
}

type typeProperty[T PropertyContraint] struct {
	parent    *Entity
	name      string
	validator validator[T]
	Array     bool
}

func (this typeProperty[T]) Parent() *Entity {
	return this.parent
}

func (this typeProperty[T]) Name() string {
	return this.name
}

///////////////////////////////////////////////////////////////////////////////////////////////////

func Refernce(parent *Entity, name string, entity *Entity) *referenceProperty {
	return &referenceProperty{
		parent: parent,
		name:   name,
		Entity: *entity,
	}
}

type referenceProperty struct {
	parent *Entity
	name   string
	Entity Entity
}

func (this referenceProperty) Parent() *Entity {
	return this.parent
}

func (this referenceProperty) Name() string {
	return this.name
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type EntityManager struct {
	Entities []Entity
}

func (this EntityManager) FindEntityByName(name string) (Entity, bool) {
	var found Entity
	var ok bool

	for _, entity := range this.Entities {
		if strings.EqualFold(entity.Name, name) {
			found = entity
			ok = true
			break
		}
	}

	return found, ok
}

func (this EntityManager) FindEntityByPropertyName(name string) ([]Entity, bool) {
	var found []Entity

	for _, entity := range this.Entities {
		for _, property := range entity.Properties {
			if strings.EqualFold(property.Name(), name) {
				found = append(found, entity)
				break
			}
		}

	}

	return found, len(found) > 0
}

func (this EntityManager) Types(language Language) {

}

func (this EntityManager) Sql(dialect SqlDialect) {

}

func (this EntityManager) Repository(language Language, dialect SqlDialect) {

}
