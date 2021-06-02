package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/types"
)

type HandlerGenerator struct {
	sourceTypeName string
	serviceName    string
	s              *types.Struct
	dtoOutputName  string
	dtoCreateName  string
}

func NewHandlerGenerator(sourceTypeName string, s *types.Struct) HandlerGenerator {
	return HandlerGenerator{
		s:              s,
		sourceTypeName: sourceTypeName,
		serviceName:    sourceTypeName + "Handler",
		dtoOutputName:  sourceTypeName + "Output",
		dtoCreateName:  sourceTypeName + "Create",
	}
}

func (g HandlerGenerator) GetHandlerFile() *Statement {
	return Add(g.generateServiceInterface()).
		Add(Line()).
		Add(Line()).
		Add(g.generateHandlerType()).
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

func (g HandlerGenerator) generateServiceInterface() *Statement {
	return Type().Id(g.sourceTypeName+"Service").Interface(
		Id("Get").Params(
			Id("ctx").Qual("context", "Context"),
			Id("id").Id("int64"),
		).Params(
			Qual(PackageModel, g.dtoOutputName),
			Error(),
		),
		Id("List").Params(
			Id("ctx").Qual("context", "Context"),
		).Params(
			Index().Qual(PackageModel, g.dtoOutputName),
			Error(),
		),
		Id("Create").Params(
			Id("ctx").Qual("context", "Context"),
			Id("entity").Qual(PackageModel, g.dtoCreateName),
		).Params(
			Qual(PackageModel, g.dtoOutputName),
			Error(),
		),
		Id("Update").Params(
			Id("ctx").Qual("context", "Context"),
			Id("id").Id("int64"),
			Id("entity").Qual(PackageModel, g.dtoCreateName),
		).Params(
			Qual(PackageModel, g.dtoOutputName),
			Error(),
		),
		Id("Delete").Params(
			Id("ctx").Qual("context", "Context"),
			Id("id").Id("int64"),
		).Error(),
	)
}

func (g HandlerGenerator) generateHandlerType() *Statement {
	return Type().Id(g.serviceName).Struct(
		Id("service").Id(g.sourceTypeName + "Service"),
	)
}

func (g HandlerGenerator) generateNewFunction() *Statement {
	return Func().Id(
		fmt.Sprintf("New%s", g.serviceName),
	).Params(
		Id("service").Id(g.sourceTypeName + "Service"),
	).Id(g.serviceName).Block(
		Return(
			Id(g.serviceName).Values(Dict{
				Id("service"): Id("service"),
			}),
		),
	)
}

func (g HandlerGenerator) generateListMethod() *Statement {
	return Func().Params(
		Id("h").Id(g.serviceName),
	).Id("List").Params(
		Id("w").Qual(PackageHttp, "ResponseWriter"),
		Id("r").Op("*").Qual(PackageHttp, "Request"),
	).Call(
		Error(),
	).Block(
		List(Id("dto"), Err()).Op(":=").Id("h").Dot("service").Dot("List").Call(
			Id("r").Dot("Context").Call(),
		),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual(PackageWeb, "NewError").Call(
				Qual(PackageHttp, "StatusInternalServerError"),
				Qual(PackageHttp, "StatusText").Call(
					Qual(PackageHttp, "StatusInternalServerError"),
				),
			)),
		),
		Line(),
		Return(Qual(PackageWeb, "RespondJSON").Call(
			Id("w"),
			Id("dto"),
			Qual(PackageHttp, "StatusOK"))),
	)
}

func (g HandlerGenerator) generateGetMethod() *Statement {
	return Func().Params(
		Id("h").Id(g.serviceName),
	).Id("Get").Params(
		Id("w").Qual(PackageHttp, "ResponseWriter"),
		Id("r").Op("*").Qual(PackageHttp, "Request"),
	).Call(
		Error(),
	).Block(
		List(Id("id"), Err()).Op(":=").Qual(PackageWeb, "Params").Call(Id("r")).Dot("Int").Call(Lit("id")),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual(PackageWeb, "NewError").Call(
				Qual(PackageHttp, "StatusBadRequest"),
				Err().Dot("Error").Call(),
			)),
		),
		Line(),
		List(Id("dto"), Id("err")).Op(":=").Id("h").Dot("service").Dot("Get").Call(
			Id("r").Dot("Context").Call(),
			Id("int64").Call(Id("id")),
		),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual(PackageWeb, "NewError").Call(
				Qual(PackageHttp, "StatusInternalServerError"),
				Qual(PackageHttp, "StatusText").Call(
					Qual(PackageHttp, "StatusInternalServerError"),
				),
			)),
		),
		Line(),
		If(Id("dto").Op("==").Call(Qual(PackageModel, g.dtoOutputName).Values())).Block(
			Return(Qual(PackageWeb, "NewError").Call(
				Qual(PackageHttp, "StatusNotFound"),
				Qual(PackageHttp, "StatusText").Call(
					Qual(PackageHttp, "StatusNotFound"),
				),
			)),
		),
		Line(),
		Return(Qual(PackageWeb, "RespondJSON").Call(
			Id("w"),
			Id("dto"),
			Qual(PackageHttp, "StatusOK"))),
	)
}

func (g HandlerGenerator) generateCreateMethod() *Statement {
	return Func().Params(
		Id("h").Id(g.serviceName),
	).Id("Create").Params(
		Id("w").Qual(PackageHttp, "ResponseWriter"),
		Id("r").Op("*").Qual(PackageHttp, "Request"),
	).Call(
		Error(),
	).Block(
		Id("dto").Op(":=").Qual(PackageModel, g.dtoCreateName).Values(),
		Err().Op(":=").Qual(PackageWeb, "Bind").Call(Id("r"), Op("&").Id("dto")),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual(PackageWeb, "NewError").Call(
				Qual(PackageHttp, "StatusBadRequest"),
				Err().Dot("Error").Call(),
			)),
		),
		Line(),
		List(Id("resp"), Err()).Op(":=").Id("h").Dot("service").Dot("Create").Call(
			Id("r").Dot("Context").Call(),
			Id("dto"),
		),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual(PackageWeb, "NewError").Call(
				Qual(PackageHttp, "StatusInternalServerError"),
				Qual(PackageHttp, "StatusText").Call(
					Qual(PackageHttp, "StatusInternalServerError"),
				),
			)),
		),
		Line(),
		Return(Qual(PackageWeb, "RespondJSON").Call(
			Id("w"),
			Id("resp"),
			Qual(PackageHttp, "StatusOK"))),
	)
}

func (g HandlerGenerator) generateUpdateMethod() *Statement {
	return Func().Params(
		Id("h").Id(g.serviceName),
	).Id("Update").Params(
		Id("w").Qual(PackageHttp, "ResponseWriter"),
		Id("r").Op("*").Qual(PackageHttp, "Request"),
	).Call(
		Error(),
	).Block(
		Id("dto").Op(":=").Qual(PackageModel, g.dtoCreateName).Values(),
		Err().Op(":=").Qual(PackageWeb, "Bind").Call(Id("r"), Op("&").Id("dto")),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual(PackageWeb, "NewError").Call(
				Qual(PackageHttp, "StatusBadRequest"),
				Err().Dot("Error").Call(),
			)),
		),
		Line(),
		List(Id("id"), Err()).Op(":=").Qual(PackageWeb, "Params").Call(Id("r")).Dot("Int").Call(Lit("id")),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual(PackageWeb, "NewError").Call(
				Qual(PackageHttp, "StatusBadRequest"),
				Err().Dot("Error").Call(),
			)),
		),
		Line(),
		List(Id("resp"), Err()).Op(":=").Id("h").Dot("service").Dot("Update").Call(
			Id("r").Dot("Context").Call(),
			Id("int64").Call(Id("id")),
			Id("dto"),
		),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual(PackageWeb, "NewError").Call(
				Qual(PackageHttp, "StatusInternalServerError"),
				Qual(PackageHttp, "StatusText").Call(
					Qual(PackageHttp, "StatusInternalServerError"),
				),
			)),
		),
		Line(),
		Return(Qual(PackageWeb, "RespondJSON").Call(
			Id("w"),
			Id("resp"),
			Qual(PackageHttp, "StatusOK"))),
	)
}

func (g HandlerGenerator) generateDeleteMethod() *Statement {
	return Func().Params(
		Id("h").Id(g.serviceName),
	).Id("Delete").Params(
		Id("w").Qual(PackageHttp, "ResponseWriter"),
		Id("r").Op("*").Qual(PackageHttp, "Request"),
	).Call(
		Error(),
	).Block(
		List(Id("id"), Err()).Op(":=").Qual(PackageWeb, "Params").Call(Id("r")).Dot("Int").Call(Lit("id")),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual(PackageWeb, "NewError").Call(
				Qual(PackageHttp, "StatusBadRequest"),
				Err().Dot("Error").Call(),
			)),
		),
		Line(),
		Err().Op("=").Id("h").Dot("service").Dot("Delete").Call(
			Id("r").Dot("Context").Call(),
			Id("int64").Call(Id("id")),
		),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual(PackageWeb, "NewError").Call(
				Qual(PackageHttp, "StatusInternalServerError"),
				Qual(PackageHttp, "StatusText").Call(
					Qual(PackageHttp, "StatusInternalServerError"),
				),
			)),
		),
		Line(),
		Return(Qual(PackageWeb, "RespondJSON").Call(
			Id("w"),
			Nil(),
			Qual(PackageHttp, "StatusNoContent"))),
	)
}
