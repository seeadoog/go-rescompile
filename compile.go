package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go-rescompile [dir] [package]")
		return
	}
	CompileResource(os.Args[1], os.Args[2])

}

func CompileResource(dir string, packageName string) {
	result := map[string][]byte{}

	err := compile(dir, result)
	if err != nil {
		panic(err)
	}

	err = writeFile(result, packageName)
	if err != nil {
		panic(err)
	}
}

func compile(dir string, result map[string][]byte) error {
	infos, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, info := range infos {
		path := path.Join(dir, info.Name())

		if info.IsDir() {
			if err := compile(path, result); err != nil {
				return err
			}
			continue
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		result[path] = data
	}

	return nil
}

func writeFile(files map[string][]byte, packag string) error {
	err := os.MkdirAll(packag, 0777)
	if err != nil {
		return err
	}

	fsb := bytes.Buffer{}
	fsb.WriteString("package " + packag + "\n")
	fsb.WriteString(fmt.Sprintf("var Files = map[string][]byte{\n"))
	for name, data := range files {
		n := strings.Index(name, "/")
		if n >= 0 {
			name = name[n+1:]
		}
		if !isValidFileName(name) {
			return fmt.Errorf("invalid filename" + name)
		}
		varName := fmt.Sprintf("%x", md5.Sum([]byte(name)))

		file, err := os.Create(packag + "/" + varName + "_binary.go")
		if err != nil {
			return err
		}
		fmt.Fprintf(file, "package %s \n", packag)
		fmt.Fprintf(file, "var a_%s = []byte{%s}  \n", varName, encode(data))

		fsb.WriteString(fmt.Sprintf("\"%s\":a_%s,\n", name, varName))
	}
	fsb.WriteString("}")

	ioutil.WriteFile(packag+"/index.go", fsb.Bytes(), 0666)
	return nil
}

func encode(b []byte) string {
	bf := strings.Builder{}

	for i, v := range b {

		bf.WriteString(strconv.Itoa(int(v)))
		bf.WriteByte(',')
		if (i+1)%50 == 0 {
			bf.WriteByte('\n')
		}
	}
	return bf.String()
}

var (
	fileNameReg = regexp.MustCompile(`^[0-9a-zA-Z._\-/]+$`)
)

func isValidFileName(name string) bool {
	return true
}
