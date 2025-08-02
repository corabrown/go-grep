package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/codecrafters-io/grep-starter-go/grep"
	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
)

var rootCmd = &cobra.Command{
	Use:   "go-grep-cli",
	Short: "go-grep-cli is a tool that mimics grep written in Golang",
	Long:  "go-grep-cli can be used to find pattern matches within files and directories, with a recursive flag option",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("root command for go-grep-cli")
		return nil
	},
}

var goGrep = &cobra.Command{
	Use:   "go-grep",
	Short: "go-grep is a tool that mimics grep written in Golang",
	Long:  "go-grep can be used to find pattern matches within files and directories, with a recursive flag option",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("provide pattern and filename or path")
		}
		pat := args[0]
		if strings.HasSuffix(pat, "\"") && strings.HasPrefix(pat, "\"") && (len(pat) > 1) {
			pat = pat[1 : len(pat)-2]
		}
		filenameOrPath := args[1]

		pattern := grep.Parse(pat)

		info, err := os.Stat(filenameOrPath)
		if err != nil {
			return err
		}

		// if file is not a directory, we search for the pattern in the file
		if !info.IsDir() {
			grep.MatchFile(filenameOrPath, pattern)
			return nil
		}

		// if the file is a directory, we search all the files in the directory (recursively or not depending on flag)
		filesToProcess := make([]string, 0)
		if recursive {
			err := filepath.WalkDir(filenameOrPath, func(path string, entry os.DirEntry, err error) error {
				if !entry.IsDir() {
					filesToProcess = append(filesToProcess, path)
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			entries, err := os.ReadDir(filenameOrPath)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				if !entry.IsDir() {
					filesToProcess = append(filesToProcess, filepath.Join(filenameOrPath, entry.Name()))
				}
			}
		}

		// the sync wait group will allow us to search up to 10 files at once
		sem := semaphore.NewWeighted(10)
		var wg sync.WaitGroup

		for _, file := range filesToProcess {
			wg.Add(1)
			sem.Acquire(context.TODO(), 1)
			go func() {
				defer func() {
					wg.Done()
					sem.Release(1)
				}()
				MatchFile(file, grep)
			}()
		}

		wg.Wait()

		return nil
	},
}

var recursive bool

func init() {
	rootCmd.AddCommand(goGrep)
	goGrep.Flags().BoolVar(&recursive, "r", false, "search files recursively in directory")
}

func Execute() error {
	return rootCmd.Execute()
}
