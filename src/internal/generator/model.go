package generator

import (
	"github.com/Pankros/scaffolding/src/internal/utils"
	. "github.com/dave/jennifer/jen"
	"go/types"
)

type ModelGenerator struct {
	sourceTypeName string
	s              *types.Struct
	dtoOutputName  string
	dtoCreateName  string
}

func NewModelGenerator(sourceTypeName string, s *types.Struct) ModelGenerator {
	return ModelGenerator{
		s:              s,
		sourceTypeName: sourceTypeName,
		dtoOutputName:  sourceTypeName + "Output",
		dtoCreateName:  sourceTypeName + "Create",
	}
}

func (g ModelGenerator) GetModelFile() *Statement {
	return Add(g.generateDTOOutputType()).
		Add(Line()).
		Add(Line()).
		Add(g.generateDTOCreateType())
}

func (g ModelGenerator) generateDTOOutputType() Code {
	return Type().Id(g.dtoOutputName).Struct(g.buildStructOutput()...)
}

func (g ModelGenerator) generateDTOCreateType() Code {
	return Type().Id(g.dtoCreateName).Struct(g.buildStructCreate()...)
}

func (g ModelGenerator) buildStructOutput() []Code {
	var d []Code
	for i := 0; i < g.s.NumFields(); i++ {
		field := g.s.Field(i)
		_, ok := field.Type().Underlying().(*types.Struct)
		if ok {
			continue
		}
		d = append(d,
			Id(field.Name()).
				Id(field.Type().String()).
				Tag(map[string]string{"json": g.toSnakeCase(field.Name())}))
	}
	return d
}

func (g ModelGenerator) buildStructCreate() []Code {
	var d []Code
	for i := 0; i < g.s.NumFields(); i++ {
		field := g.s.Field(i)
		_, ok := field.Type().Underlying().(*types.Struct)
		if ok {
			continue
		}
		if field.Name() == "ID" {
			continue
		}
		d = append(d,
			Id(field.Name()).
			Id(field.Type().String()).
			Tag(map[string]string{"json": g.toSnakeCase(field.Name())}))
	}
	return d
}

func (g ModelGenerator) toSnakeCase(s string) string {
	return utils.ToSnakeCase(s)
}
