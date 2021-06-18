package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/types"
)

type MapperGenerator struct {
	sourceTypeName string
	serviceName    string
	s              *types.Struct
	dtoOutputName  string
	dtoCreateName  string
}


func NewMapperGenerator(sourceTypeName string, s *types.Struct) MapperGenerator {
	return MapperGenerator{
		s:              s,
		sourceTypeName: sourceTypeName,
		serviceName:    sourceTypeName + "Mapper",
		dtoOutputName:  sourceTypeName + "Output",
		dtoCreateName:  sourceTypeName + "Create",
	}
}

func (g MapperGenerator) GetMapperFile() *Statement {
	return Add(g.generateServiceType()).
		Add(Line()).
		Add(g.generateNewFunction()).
		Add(Line()).
		Add(Line()).
		Add(g.generateToListDTO()).
		Add(Line()).
		Add(Line()).
		Add(g.generateToDTO()).
		Add(Line()).
		Add(Line()).
		Add(g.generateToEntity())
}

func (g MapperGenerator) generateServiceType() Code {
	return Type().Id(g.serviceName).Struct()
}

func (g MapperGenerator) generateNewFunction() Code {
	return Func().Id(
		fmt.Sprintf("New%s", g.serviceName),
	).Params().Id(g.serviceName).Block(
		Return(
			Id(g.serviceName).Values(),
		),
	)
}

func (g MapperGenerator) generateToListDTO() Code {
	return Func().Params(
		Id("m").Id(g.serviceName),
	).Id(
		fmt.Sprintf("ToListDTO"),
	).Params(
		Id("entities").Index().Qual(PackageModel, g.sourceTypeName),
	).Index().Qual(PackageModel, g.dtoOutputName).Block(
		Id("dto").Op(":=").Make(
			Index().Qual(PackageModel, g.dtoOutputName),
			Lit(0),
			Len(Id("entities")),
			),
		Line(),
		For(List(Id("_"), Id("entity")).Op(":=").Range().Id("entities")).Block(
			Id("dto").Op("=").Append(Id("dto"), Id("m").Dot("ToDTO").Call(Id("entity"))),
			),
		Return(
			Id("dto"),
		),
	)
}

/*
	return model.PaymentMethodTypeOutput{
		ID:   entity.ID,
		Code: entity.Code,
		Name: entity.Name,
	}
*/

func (g MapperGenerator) generateToDTO() Code {
	return Func().Params(
		Id("m").Id(g.serviceName),
	).Id(
		fmt.Sprintf("ToDTO"),
	).Params(
		Id("entity").Qual(PackageModel, g.sourceTypeName),
	).Qual(PackageModel, g.dtoOutputName).Block(
		Return(Qual(PackageModel, g.dtoOutputName).Values(g.buildMap(g.s, "entity"))),
	)
}

func (g MapperGenerator) generateToEntity() Code {
	return Func().Params(
		Id("m").Id(g.serviceName),
	).Id(
		fmt.Sprintf("ToEntity"),
	).Params(
		Id("id").Id("int64"),
		Id("dto").Qual(PackageModel, g.dtoCreateName),
	).Qual(PackageModel, g.sourceTypeName).Block(
		Return(Qual(PackageModel, g.sourceTypeName).Values(g.buildMapForToEntity(g.s, "dto"))),
	)
}

func (g MapperGenerator) buildMapForToEntity(s *types.Struct, source string) Dict {
	d := map[Code]Code{}
	for i := 0; i < s.NumFields(); i++ {
		_, ok := s.Field(i).Type().Underlying().(*types.Struct)
		if ok {
			continue
		}

		if s.Field(i).Name() == "ID" {
			d[Id(s.Field(i).Name())] = Id("id")
		} else {
			d[Id(s.Field(i).Name())] = Id(source).Dot(s.Field(i).Name())
		}
	}
	return d
}

func (g MapperGenerator) buildMap(s *types.Struct, source string) Dict {
	d := map[Code]Code{}
	for i := 0; i < s.NumFields(); i++ {
		_, ok := s.Field(i).Type().Underlying().(*types.Struct)
		if ok {
			continue
		}
		d[Id(s.Field(i).Name())] = Id(source).Dot(s.Field(i).Name())
	}
	return d
}