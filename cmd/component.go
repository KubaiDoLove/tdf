/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/KubaiDoLove/tdf/templates"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"text/template"
)

type componentData struct {
	Name string
}

type fileToWrite struct {
	RawTemplate  string
	OutputDir    string
	Name         string
	TemplateData interface{}
}

// componentCmd represents the component command
var componentCmd = &cobra.Command{
	Use:   "component",
	Short: "Generates new frontend React components in TypeScript",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You need to provide a component name.")
			return
		}
		if len(args) > 1 {
			fmt.Println("Command takes only one argument as a component name.")
			os.Exit(1)
		}

		outputDir, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Println("Could not parse output flag: ", err)
			os.Exit(1)
		}

		componentName := args[0]
		componentDir := filepath.Join(outputDir, componentName)
		typesDir := filepath.Join(componentDir, "types")
		if err := os.Mkdir(componentDir, os.ModePerm); err != nil {
			fmt.Println("Could not create directory for component: ", err)
			os.Exit(1)
		}
		if err := os.Mkdir(typesDir, os.ModePerm); err != nil {
			fmt.Println("Could not create directory for component's types: ", err)
			os.Exit(1)
		}

		boilerplate := []fileToWrite{
			{
				RawTemplate: templates.ComponentTemplate,
				OutputDir:   componentDir,
				Name:        componentName + ".tsx",
				TemplateData: componentData{
					Name: componentName,
				},
			},
			{
				RawTemplate: templates.ComponentIndexTemplate,
				OutputDir:   componentDir,
				Name:        "index.ts",
				TemplateData: componentData{
					Name: componentName,
				},
			},
			{
				RawTemplate: templates.ComponentInterfacesTemplate,
				OutputDir:   typesDir,
				Name:        "interfaces.ts",
				TemplateData: componentData{
					Name: componentName,
				},
			},
			{
				RawTemplate: templates.ComponentTypesIndexTemplate,
				OutputDir:   typesDir,
				Name:        "index.ts",
			},
		}

		isScssIncluded, err := cmd.Flags().GetBool("scss")
		if err != nil {
			fmt.Println("Could not parse scss flag: ", err)
			os.Exit(1)
		}
		if isScssIncluded {
			boilerplate = append(boilerplate, fileToWrite{
				OutputDir: componentDir,
				Name:      componentName + ".scss",
			})
		}

		if err := writeComponentFiles(boilerplate...); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(componentCmd)
	componentCmd.Flags().BoolP(
		"scss",
		"s",
		false,
		"Use this flag to include scss file for your component",
	)
	componentCmd.Flags().StringP(
		"output",
		"o",
		"",
		"Set output directory",
	)
}

func writeComponentFiles(files ...fileToWrite) error {
	for _, f := range files {
		file, err := os.Create(filepath.Join(f.OutputDir, f.Name))
		if err != nil {
			errText := fmt.Sprintf(
				"Could not create %s file in %s: %s",
				f.Name,
				f.OutputDir,
				err.Error(),
			)
			return errors.New(errText)
		}
		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()
		componentTemplate := template.Must(template.New(f.Name).Parse(f.RawTemplate))
		if err := componentTemplate.Execute(file, f.TemplateData); err != nil {
			return err
		}
	}

	return nil
}
