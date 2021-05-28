package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/types"
)

type ServiceGenerator struct {
	sourceTypeName string
	serviceName    string
	s              *types.Struct
	dtoOutputName  string
	dtoCreateName  string
}

func NewServiceGenerator(sourceTypeName string, s *types.Struct) ServiceGenerator {
	return ServiceGenerator{
		s:              s,
		sourceTypeName: sourceTypeName,
		serviceName:    sourceTypeName + "Service",
		dtoOutputName:  sourceTypeName + "Output",
		dtoCreateName:  sourceTypeName + "Create",
	}
}

func (g ServiceGenerator) GetServiceFile() *Statement {
	return Add(g.generateDAOInterface()).
		Add(Line()).
		Add(Line()).
		Add(g.generateServiceType()).
		Add(Line()).
		Add(g.generateNewFunction()).
		Add(Line()).
		Add(Line()).
		Add(g.generateListMethod()).
		Add(Line()).
		Add(Line()).
		Add(g.generateGetMethod()).
		Add(Line()).
		Add(Line()).
		Add(g.generateCreateMethod()).
		Add(Line()).
		Add(Line()).
		Add(g.generateUpdateMethod()).
		Add(Line()).
		Add(Line()).
		Add(g.generateDeleteMethod())
}

func (g ServiceGenerator) generateDAOInterface() *Statement {
	return Type().Id(g.sourceTypeName+"DAO").Interface(
		Id("Get").Params(
			Id("ctx").Qual("context", "Context"),
			Id("id").Id("int64"),
		).Params(
			Qual(PackageModel, g.sourceTypeName),
			Error(),
		),
		Id("List").Params(
			Id("ctx").Qual("context", "Context"),
		).Params(
			Index().Qual(PackageModel, g.sourceTypeName),
			Error(),
		),
		Id("Create").Params(
			Id("ctx").Qual("context", "Context"),
			Id("entity").Qual(PackageModel, g.sourceTypeName),
		).Params(
			Id("int64"),
			Error(),
		),
		Id("Update").Params(
			Id("ctx").Qual("context", "Context"),
			Id("entity").Qual(PackageModel, g.sourceTypeName),
		).Error(),
		Id("Delete").Params(
			Id("ctx").Qual("context", "Context"),
			Id("id").Id("int64"),
		).Error(),
	)
}

func (g ServiceGenerator) generateServiceType() *Statement {
	return Type().Id(g.serviceName).Struct(
		Id("dao").Id(g.sourceTypeName + "DAO"),
	)
}

func (g ServiceGenerator) generateNewFunction() *Statement {
	return Func().Id(
		fmt.Sprintf("New%s", g.serviceName),
	).Params(
		Id("dao").Id(g.sourceTypeName + "DAO"),
	).Id(g.serviceName).Block(
		Return(
			Id(g.serviceName).Values(Dict{
				Id("dao"): Id("dao"),
			}),
		),
	)
}

func (g ServiceGenerator) generateListMethod() *Statement {
	return Func().Params(
		Id("s").Id(g.serviceName),
	).Id("List").Params(
		Id("ctx").Qual("context", "Context"),
	).Call(
		Index().Qual(PackageModel, g.dtoOutputName),
		Error(),
	).Block(
		List(Id("entities"), Id("err")).Op(":=").Id("s").Dot("dao").Dot("List").Call(
			Id("ctx"),
		),
		Line(),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Nil(), Id("err")),
		),
		Line(),
		Return(Id("toListDTO").Call(Id("entities")), Nil()),
	)
}

func (g ServiceGenerator) generateGetMethod() *Statement {
	return Func().Params(
		Id("s").Id(g.serviceName),
	).Id("Get").Params(
		Id("ctx").Qual("context", "Context"),
		Id("id").Id("int64"),
	).Call(
		Qual(PackageModel, g.dtoOutputName),
		Error(),
	).Block(
		List(Id("entity"), Id("err")).Op(":=").Id("s").Dot("dao").Dot("Get").Call(
			Id("ctx"),
			Id("id"),
		),
		Line(),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Nil(), Id("err")),
		),
		Line(),
		Return(Id("toDTO").Call(Id("entity")), Nil()),
	)
}

func (g ServiceGenerator) generateCreateMethod() *Statement {
	return Func().Params(
		Id("s").Id(g.serviceName),
	).Id("Create").Params(
		Id("ctx").Qual("context", "Context"),
		Id("dto").Qual(PackageModel, g.dtoCreateName),
	).Call(
		Qual(PackageModel, g.dtoOutputName),
		Error(),
	).Block(
		Id("entity").Op(":=").Id("toEntity").Call(
			Lit(0),
			Id("dto"),
		),
		Line(),
		List(Id("id"), Id("err")).Op(":=").Id("s").Dot("dao").Dot("Create").Call(
			Id("ctx"),
			Id("entity"),
		),
		Line(),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Id("err")),
		),
		Line(),
		Return(Id("s").Dot("Get").Call(
			Id("ctx"),
			Id("id"),
		)),
	)
}

func (g ServiceGenerator) generateUpdateMethod() *Statement {
	return Func().Params(
		Id("s").Id(g.serviceName),
	).Id("Update").Params(
		Id("ctx").Qual("context", "Context"),
		Id("id").Id("int64"),
		Id("dto").Qual(PackageModel, g.dtoCreateName),
	).Call(
		Qual(PackageModel, g.dtoOutputName),
		Error(),
	).Block(
		If(Id("id").Op("==").Lit(0)).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Lit(fmt.Sprintf(ErrorUpdateNoId, g.sourceTypeName))),
		),
		Line(),
		Id("entity").Op(":=").Id("toEntity").Call(
			Id("id"),
			Id("dto"),
		),
		Line(),
		Id("err").Op(":=").Id("s").Dot("dao").Dot("Update").Call(
			Id("ctx"),
			Id("entity"),
		),
		Line(),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Id("err")),
		),
		Line(),
		Return(Id("s").Dot("Get").Call(
			Id("ctx"),
			Id("id"),
		)),
	)
}

func (g ServiceGenerator) generateDeleteMethod() *Statement {
	return Func().Params(
		Id("s").Id(g.serviceName),
	).Id("Delete").Params(
		Id("ctx").Qual("context", "Context"),
		Id("id").Id("int64"),
	).Error().Block(
		Id("err").Op(":=").Id("s").Dot("dao").Dot("Delete").Call(
			Id("ctx"),
			Id("id"),
		),
		Line(),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Id("err")),
		),
		Line(),
		Return(Nil()),
	)
}
