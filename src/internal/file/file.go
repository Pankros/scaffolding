package file

import (
	"github.com/dave/jennifer/jen"
	"os"
)

func generateJenFile(pkg string) *jen.File {
	f := jen.NewFile(pkg)
	return f
}

func generateFolder(path string) error {
	return os.Mkdir(path, 0755)
}

func saveFile(file *jen.File, path, fileName string) error {
	_ = os.Remove(path + "/" + fileName)
	return file.Save(path + "/" + fileName + ".go")
}

func SaveFile(pkg, path, fileName string, content *jen.Statement) error {
	f := generateJenFile(pkg)
	f.Add(content)
	_ = generateFolder(path)
	return saveFile(f, path, fileName)
}
