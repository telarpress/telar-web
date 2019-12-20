package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/jinzhu/copier"
	stack "github.com/openfaas/faas-cli/stack"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func MakeStack() *cobra.Command {
	var command = &cobra.Command{
		Use:          "stack",
		Short:        "Conver stack-init.yml to stack.yml",
		Example:      `  tsocial stack`,
		SilenceUsage: false,
	}
	command.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Start creating stack...")
		createStack()
		return nil
	}
	return command
}

func createStack() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// parent := filepath.Dir(wd)
	fmt.Printf("Current address: %s \n", wd)
	stackFile, _ := stack.ParseYAMLFile(path.Join(wd, "stack-init.yml"), "", "", false)

	for name, function := range stackFile.Functions {
		fmt.Println("map:", function.EnvironmentFile, name)
		//read environment variables from the file
		fileEnvironment, err := readFiles(function.EnvironmentFile, wd)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		//combine all environment variables
		allEnvironment, envErr := compileEnvironment([]string{}, function.Environment, fileEnvironment)
		if envErr != nil {
			log.Fatalf("error: %v", envErr)
		}

		// Set environments
		newFuncs := stack.Function{}
		copier.Copy(&newFuncs, stackFile.Functions[name])
		newFuncs.Environment = allEnvironment

		stackFile.Functions[name] = newFuncs
	}

	d, err := yaml.Marshal(&stackFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	errWrite := ioutil.WriteFile(path.Join(wd, "stack.yml"), d, 0644)
	if errWrite != nil {
		log.Fatalf("error: %v", err)
	}
}
