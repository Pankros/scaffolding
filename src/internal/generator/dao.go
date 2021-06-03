package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/types"
	"strings"
)

type DaoGenerator struct {
	sourceTypeName string
	tableName      string
	daoName        string
	s              *types.Struct
}

func NewDaoGenerator(sourceTypeName, tableName string, s *types.Struct) DaoGenerator {
	return DaoGenerator{
		sourceTypeName: sourceTypeName,
		tableName:      tableName,
		s:              s,
		daoName:        sourceTypeName + "DAO",
	}
}

func (g DaoGenerator) GetDaoFile() *Statement {
	return Add(g.generateDAOType()).
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

func (g DaoGenerator) generateNewFunction() *Statement {
	return Func().Id(
		fmt.Sprintf("New%s", g.daoName),
	).Params(
		Id("audit").Id("AuditService"),
	).Id(g.daoName).Block(
		Return(
			Id(g.daoName).Values(Dict{
				Id("audit"): Id("audit"),
			}),
		),
	)
}

func (g DaoGenerator) generateDAOType() *Statement {
	return Type().Id(g.daoName).Struct(
		Id("audit").Id("AuditService"),
	)
}

func (g DaoGenerator) generateListMethod() *Statement {
	query := generateListQuery(g.s, g.tableName)
	return Func().Params(
		Id("p").Id(g.daoName),
	).Id("List").Params(
		Id("ctx").Qual("context", "Context"),
		Id("tx").Id("DBConnection"),
	).Call(
		Index().Qual(PackageModel, g.sourceTypeName),
		Error(),
	).Block(
		Var().Id("query").Op("=").Lit(query),
		Var().Id("rows").Index().Qual(PackageModel, g.sourceTypeName),
		Line(),
		Id("err").Op(":=").Id("tx").Dot("SelectContext").Call(
			Id("ctx"),
			Op("&").Id("rows"),
			Id("query")),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Index().Qual(PackageModel, g.sourceTypeName).Block(), Id("err")),
		),
		Line(),
		Return(Id("rows"), Nil()),
	)
}

func (g DaoGenerator) generateGetMethod() *Statement {
	query := generateGetQuery(g.s, g.tableName)
	return Func().Params(
		Id("p").Id(g.daoName),
	).Id("Get").Params(
		Id("ctx").Qual("context", "Context"),
		Id("tx").Id("DBConnection"),
		Id("id").Id("int64"),
	).Call(
		Qual(PackageModel, g.sourceTypeName),
		Error(),
	).Block(
		Var().Id("query").Op("=").Lit(query),
		Var().Id("row").Qual(PackageModel, g.sourceTypeName),
		Line(),
		List(Id("stmt"), Id("err")).Op(":=").Id("tx").Dot("PrepareContext").Call(
			Id("ctx"),
			Id("query"),
		),
		Line(),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Qual(PackageModel, g.sourceTypeName).Block(), Id("err")),
		),
		Line(),
		Id("err").Op("=").Id("stmt").Dot("Get").Call(
			Op("&").Id("row"),
			Id("id")),
		If(Id("err").Op("!=").Nil()).Block(
			If(Id("err").Op("==").Qual(PackageSQL, "ErrNoRows")).Block(
				Return(Qual(PackageModel, g.sourceTypeName).Block(), Id("ErrNotFound")),
			),
			Return(Qual(PackageModel, g.sourceTypeName).Block(), Id("err")),
		),
		Line(),
		Return(Id("row"), Nil()),
	)
}

func (g DaoGenerator) generateCreateMethod() *Statement {
	query := generateCreateQuery(g.s, g.tableName)
	return Func().Params(
		Id("p").Id(g.daoName),
	).Id("Create").Params(
		Id("ctx").Qual("context", "Context"),
		Id("tx").Id("DBConnection"),
		Id("entity").Qual(PackageModel, g.sourceTypeName),
	).Call(
		Id("int64"),
		Error(),
	).Block(
		Id("entity").Dot("Audit").Op("=").Id("p").Dot("audit").Dot("GetAuditForCreate").Call(
			Id("ctx"),
		),
		Var().Id("query").Op("=").Lit(query),
		Line(),
		List(Id("result"), Id("err")).Op(":=").Id("tx").Dot("NamedExecContext").Call(
			Id("ctx"),
			Id("query"),
			Op("&").Id("entity"),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Lit(0), Id("err")),
		),
		Line(),
		List(Id("id"), Id("_")).Op(":=").Id("result").Dot("LastInsertId").Call(),
		Line(),
		Return(Id("id"), Nil()),
	)
}

func (g DaoGenerator) generateUpdateMethod() *Statement {
	query := generateUpdateQuery(g.s, g.tableName)
	return Func().Params(
		Id("p").Id(g.daoName),
	).Id("Update").Params(
		Id("ctx").Qual("context", "Context"),
		Id("tx").Id("DBConnection"),
		Id("entity").Qual(PackageModel, g.sourceTypeName),
	).Error().Block(
		Id("entity").Dot("Audit").Op("=").Id("p").Dot("audit").Dot("GetAuditForUpdate").Call(
			Id("ctx"),
		),
		Var().Id("query").Op("=").Lit(query),
		Line(),
		If(Id("entity").Dot("ID").Op("==").Lit(0)).Block(
			Return(Qual("errors", "New").Call(Lit("can't update an entity without ID"))),
		),
		Line(),
		List(Id("_"), Id("err")).Op(":=").Id("tx").Dot("NamedExecContext").Call(
			Id("ctx"),
			Id("query"),
			Op("&").Id("entity"),
		),
		Line(),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Id("err")),
		),
		Line(),
		Return(Nil()),
	)
}

func (g DaoGenerator) generateDeleteMethod() *Statement {
	query := generateDeleteQuery(g.tableName)
	return Func().Params(
		Id("p").Id(g.daoName),
	).Id("Delete").Params(
		Id("ctx").Qual("context", "Context"),
		Id("tx").Id("DBConnection"),
		Id("id").Id("int64"),
	).Error().Block(
		Var().Id("query").Op("=").Lit(query),
		Line(),
		If(Id("id").Op("==").Lit(0)).Block(
			Return(Qual("errors", "New").Call(Lit("can't delete an entity without ID"))),
		),
		Line(),
		List(Id("_"), Id("err")).Op(":=").Id("tx").Dot("NamedExecContext").Call(
			Id("ctx"),
			Id("query"),
			Map(String()).Interface().Values(Dict{Lit("id"): Id("id")}),
		),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Id("err")),
		),
		Line(),
		Return(Nil()),
	)
}

func generateDeleteQuery(table string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE id = :id", table)
}

func generateUpdateQuery(s *types.Struct, table string) string {
	columns := getColumnNames(s)
	columns = removeColumn(columns, "id")
	columns = removeColumn(columns, "created_at")
	columns = removeColumn(columns, "created_by")
	for i, c := range columns {
		columns[i] = fmt.Sprintf("%s = :%s", c, c)
	}
	return fmt.Sprintf("UPDATE %s SET %s WHERE id = :id", table, strings.Join(columns, ", "))
}

func generateCreateQuery(s *types.Struct, table string) string {
	columns := getColumnNames(s)
	columns = removeColumn(columns, "id")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (:%s)", table, strings.Join(columns, ", "), strings.Join(columns, ", :"))
}

func generateListQuery(s *types.Struct, table string) string {
	return fmt.Sprintf("SELECT %s FROM %s", strings.Join(getColumnNames(s), ", "), table)
}

func generateGetQuery(s *types.Struct, table string) string {
	return fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", strings.Join(getColumnNames(s), ", "), table)
}

func removeColumn(ss []string, id string) []string {
	var resp []string
	for _, s := range ss {
		if s != id {
			resp = append(resp, s)
		}
	}
	return resp
}

func getColumnNames(s *types.Struct) []string {
	var l []string
	for i := 0; i < s.NumFields(); i++ {
		t := s.Field(i).Type()
		if strings.Contains(t.String(), ".Audit") {
			l = append(l, getColumnNames(t.Underlying().(*types.Struct))...)
		} else {
			tagValue := s.Tag(i)
			a := strings.Split(tagValue, ":")[1]
			a = strings.ReplaceAll(a, "\"", "")
			l = append(l, a)
		}
	}
	return l
}
