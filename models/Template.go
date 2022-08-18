package models

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Template struct {
	Name    string
	Version string
}

func NewTemplate(name string) *Template {
	return &Template{Name: name}
}

func (t *Template) String() string {
	var name strings.Builder
	name.WriteString("Template: '")
	name.WriteString(t.Name)
	name.WriteString(" ")

	if len(t.Version) != 0 {
		name.WriteString("V")
		name.WriteString(t.Version)
	} else {
		name.WriteString("(not unpacked)")
	}

	name.WriteString("'")
	return name.String()
}

// Unpack the zip file containing all the necessary files for the template
func (t *Template) Unpack(src string, dest string) error {
	var filenames []string
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		fp := filepath.Join(dest, f.Name)
		_ = strings.HasPrefix(fp, filepath.Clean(dest)+string(os.PathSeparator))
		filenames = append(filenames, fp)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fp, os.ModePerm)
			continue
		}
		os.MkdirAll(filepath.Dir(fp), os.ModePerm)
		outFile, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()
		_, filename := filepath.Split(fp)
		if strings.HasSuffix(filename, ".tm") {
			err := t.loadManifest(fp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Parse the manifest file for the template and save the values in the Template structure
// The manifest file should have the same name as the template with a .tm extension
func (t *Template) loadManifest(mf string) error {
	yf, err := os.ReadFile(mf)
	if err != nil {
		return err
	}
	data := make(map[interface{}]interface{})
	err2 := yaml.Unmarshal(yf, &data)
	if err2 != nil {
		return err2
	}
	for k, v := range data {
		fmt.Printf("Key: %s\nVal: %d\n", k, v)
	}
	return nil
}
