package cmd

import (
	"fmt"
	"os"

	prompt "github.com/ginkida/gptunzip/prompt"

	"github.com/spf13/cobra"
)

var repoPath string
var isParts bool

var rootCmd = &cobra.Command{
	Use:   "gptunzip /path/to/git/repository",
	Short: "gptunzip is a utility to convert a Zip repository to a text file for input into GPT-4",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath = args[0]
		tempDir, err := prompt.CreateTempDirectory()
		if err != nil {
			fmt.Printf("Error creating temp directory: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tempDir)
		if err := prompt.UnzipFile(repoPath, tempDir); err != nil {
			fmt.Printf("Error unzipping file: %v\n", err)
			os.Exit(1)
		}
		subdirPath, err := prompt.FindSingleSubdir(tempDir)
		if err != nil {
			fmt.Printf("Error finding subdir: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Unzipped to:", subdirPath)
		repo, err := prompt.ProcessGitRepo(subdirPath)
		if err != nil {
			fmt.Printf("Error processing git repo: %v\n", err)
			os.Exit(1)
		}
		output, err := prompt.OutputGitRepo(repo)
		if err != nil {
			fmt.Printf("Error generating output: %v\n", err)
			os.Exit(1)
		}
		saveFunc := prompt.SaveTextAsOneFile
		if isParts {
			saveFunc = prompt.SaveTextAsParts
		}
		if err := saveFunc(subdirPath, output); err != nil {
			fmt.Printf("Error saving output: %v\n", err)
		} else {
			fmt.Println("The file has been saved in the current directory.")
		}
		fmt.Printf("Estimated number of tokens: %d\n", prompt.EstimateTokens(output))
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&isParts, "parts", "p", false, "create small parts for chatGPT")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}
