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
		Add(g.generateMapperInterface()).
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
			Id("tx").Qual(PackageDao, "DBConnection"),
			Id("id").Id("int64"),
		).Params(
			Qual(PackageModel, g.sourceTypeName),
			Error(),
		),
		Id("List").Params(
			Id("ctx").Qual("context", "Context"),
			Id("tx").Qual(PackageDao, "DBConnection"),
		).Params(
			Index().Qual(PackageModel, g.sourceTypeName),
			Error(),
		),
		Id("Create").Params(
			Id("ctx").Qual("context", "Context"),
			Id("tx").Qual(PackageDao, "DBConnection"),
			Id("entity").Qual(PackageModel, g.sourceTypeName),
		).Params(
			Id("int64"),
			Error(),
		),
		Id("Update").Params(
			Id("ctx").Qual("context", "Context"),
			Id("tx").Qual(PackageDao, "DBConnection"),
			Id("entity").Qual(PackageModel, g.sourceTypeName),
		).Error(),
		Id("Delete").Params(
			Id("ctx").Qual("context", "Context"),
			Id("tx").Qual(PackageDao, "DBConnection"),
			Id("id").Id("int64"),
		).Error(),
	)
}

func (g ServiceGenerator) generateMapperInterface() *Statement {
	return Type().Id(g.sourceTypeName+"Mapper").Interface(
		Id("ToListDTO").Params(
			Id("entities").Index().Qual(PackageModel, g.sourceTypeName),
		).Params(
			Index().Qual(PackageModel, g.dtoOutputName),
		),
		Id("ToDTO").Params(
			Id("entity").Qual(PackageModel, g.sourceTypeName),
		).Params(
			Qual(PackageModel, g.dtoOutputName),
		),
		Id("ToEntity").Params(
			Id("id").Id("int64"),
			Id("dto").Qual(PackageModel, g.dtoCreateName),
		).Params(
			Qual(PackageModel, g.sourceTypeName),
		),
	)
}

func (g ServiceGenerator) generateServiceType() *Statement {
	return Type().Id(g.serviceName).Struct(
		Id("dao").Id(g.sourceTypeName + "DAO"),
		Id("txService").Id("TransactionManager"),
		Id("mapper").Id(g.sourceTypeName + "Mapper"),
	)
}

func (g ServiceGenerator) generateNewFunction() *Statement {
	return Func().Id(
		fmt.Sprintf("New%s", g.serviceName),
	).Params(
		Id("dao").Id(g.sourceTypeName + "DAO"),
		Id("txService").Id("TransactionManager"),
		Id("mapper").Id(g.sourceTypeName + "Mapper"),
	).Id(g.serviceName).Block(
		Return(
			Id(g.serviceName).Values(Dict{
				Id("dao"): Id("dao"),
				Id("txService"): Id("txService"),
				Id("mapper"): Id("mapper"),
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
		List(Id("tx"), Err()).Op(":=").Id("s").Dot("txService").Dot("Open").Call(
			Id("ctx"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Nil(), Id("err")),
		),
		Line(),
		List(Id("entities"), Id("err")).Op(":=").Id("s").Dot("dao").Dot("List").Call(
			Id("ctx"),
			Id("tx"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Nil(), Id("err")),
		),
		Line(),
		Return(Id("s").Dot("mapper").Dot("ToListDTO").Call(Id("entities")), Nil()),
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
		List(Id("tx"), Err()).Op(":=").Id("s").Dot("txService").Dot("Open").Call(
			Id("ctx"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Id("err")),
		),
		Line(),
		List(Id("entity"), Id("err")).Op(":=").Id("s").Dot("dao").Dot("Get").Call(
			Id("ctx"),
			Id("tx"),
			Id("id"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Id("err")),
		),
		Line(),
		Return(Id("s").Dot("mapper").Dot("ToDTO").Call(Id("entity")), Nil()),
	)
}

func (g ServiceGenerator) generateCreateMethod() *Statement {
	return Func().Params(
		Id("s").Id(g.serviceName),
	).Id("Create").Params(
		Id("ctx").Qual("context", "Context"),
		Id("dto").Qual(PackageModel, g.dtoCreateName),
	).Call(
		Id("resp").Qual(PackageModel, g.dtoOutputName),
		Err().Error(),
	).Block(
		Id("entity").Op(":=").Id("s").Dot("mapper").Dot("ToEntity").Call(
			Lit(0),
			Id("dto"),
		),
		Line(),
		List(Id("tx"), Err()).Op(":=").Id("s").Dot("txService").Dot("OpenTx").Call(
			Id("ctx"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Id("err")),
		),
		Line(),
		Defer().Func().Call().Block(
			Id("deferErr").Op(":=").Id("s").Dot("txService").Dot("CloseTx").Call(
				Id("ctx"),
				Id("tx"),
				Err(),
			),
			If(Id("deferErr").Op("!=").Nil()).Block(
				Err().Op("=").Qual(PackageUtils,"WrapOrCreateError").Call(
					Err(),
					Id("deferErr"),
				),
				Id("resp").Op("=").Qual(PackageModel, g.dtoOutputName).Values(),
			),
			).Call(),
		Line(),
		List(Id("id"), Id("err")).Op(":=").Id("s").Dot("dao").Dot("Create").Call(
			Id("ctx"),
			Id("tx"),
			Id("entity"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Id("err")),
		),
		Line(),
		List(Id("created"), Id("err")).Op(":=").Id("s").Dot("dao").Dot("Get").Call(
			Id("ctx"),
			Id("tx"),
			Id("id"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Id("err")),
		),
		Line(),
		Return(Id("s").Dot("mapper").Dot("ToDTO").Call(Id("created")), Nil()),
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
		Id("resp").Qual(PackageModel, g.dtoOutputName),
		Err().Error(),
	).Block(
		If(Id("id").Op("==").Lit(0)).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Qual("errors", "New").Call(
				Lit(fmt.Sprintf(ErrorUpdateNoId, g.sourceTypeName)),
			)),
		),
		Line(),
		List(Id("tx"), Err()).Op(":=").Id("s").Dot("txService").Dot("OpenTx").Call(
			Id("ctx"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Id("err")),
		),
		Line(),
		Defer().Func().Call().Block(
			Id("deferErr").Op(":=").Id("s").Dot("txService").Dot("CloseTx").Call(
				Id("ctx"),
				Id("tx"),
				Err(),
			),
			If(Id("deferErr").Op("!=").Nil()).Block(
				Err().Op("=").Qual(PackageUtils,"WrapOrCreateError").Call(
					Err(),
					Id("deferErr"),
				),
				Id("resp").Op("=").Qual(PackageModel, g.dtoOutputName).Values(),
			),
		).Call(),
		Line(),
		Id("entity").Op(":=").Id("s").Dot("mapper").Dot("ToEntity").Call(
			Id("id"),
			Id("dto"),
		),
		Line(),
		Id("err").Op("=").Id("s").Dot("dao").Dot("Update").Call(
			Id("ctx"),
			Id("tx"),
			Id("entity"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Id("err")),
		),
		Line(),
		List(Id("updated"), Id("err")).Op(":=").Id("s").Dot("dao").Dot("Get").Call(
			Id("ctx"),
			Id("tx"),
			Id("id"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Qual(PackageModel, g.dtoOutputName).Values(), Id("err")),
		),
		Line(),
		Return(Id("s").Dot("mapper").Dot("ToDTO").Call(Id("updated")), Nil()),
	)
}

func (g ServiceGenerator) generateDeleteMethod() *Statement {
	return Func().Params(
		Id("s").Id(g.serviceName),
	).Id("Delete").Params(
		Id("ctx").Qual("context", "Context"),
		Id("id").Id("int64"),
	).Error().Block(
		List(Id("tx"), Err()).Op(":=").Id("s").Dot("txService").Dot("OpenTx").Call(
			Id("ctx"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Err()),
		),
		Line(),
		Defer().Func().Call().Block(
			Id("deferErr").Op(":=").Id("s").Dot("txService").Dot("CloseTx").Call(
				Id("ctx"),
				Id("tx"),
				Err(),
			),
			If(Id("deferErr").Op("!=").Nil()).Block(
				Err().Op("=").Qual(PackageUtils,"WrapOrCreateError").Call(
					Err(),
					Id("deferErr"),
				),
			),
		).Call(),
		Id("err").Op("=").Id("s").Dot("dao").Dot("Delete").Call(
			Id("ctx"),
			Id("tx"),
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
