package cmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/codecrafters-io/grep-starter-go/grep"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var recursive bool

// goGrep is the base command (and only command) which takes in 2 arguments for the pattern and
// the path to search. It also accepts a --r flag which indicates whether or not
// to search through a directory recursively. This flag is ignored if the given path is a file
// and not a directory.
var goGrep = &cobra.Command{
	Use:   "go-grep",
	Short: "go-grep is a tool that mimics grep written in Golang",
	Long:  "go-grep can be used to find pattern matches within files and directories, with a recursive flag option",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("provide pattern and filename or path")
		}

		// the first argument is the pattern to match against
		pat := args[0]
		// if the pattern is wrapped in double quotes, we remove them
		if strings.HasSuffix(pat, "\"") && strings.HasPrefix(pat, "\"") && (len(pat) > 1) {
			pat = pat[1 : len(pat)-2]
		}

		// the second argument is the path to search. This can be a file or directory path
		path := args[1]

		// parse the pattern once so that it can be checked against each file in path
		pattern := grep.Parse(pat, nil)

		// create a channel to contains the files we want to search 
		fileInputChan := make(chan string, 10)

		// use an error group to fill the channel with files under the given path
		g, _ := errgroup.WithContext(context.Background())
		g.Go(func() error {
			err := getFilesToSearch(path, fileInputChan)
			if err != nil {
				return err
			}
			return nil
		})

		// iterate over the file channel to match the pattern against the contents of the file. 
		// Use a wait group to make sure we finish all the work before exiting  
		var wg sync.WaitGroup
		for file := range fileInputChan {
			wg.Add(1)
			go func() {
				defer func() {
					wg.Done()
				}()
				grep.MatchFile(file, pattern)
			}()
		}

		// check for errors from the getFilesToSearch execution 
		if err := g.Wait(); err != nil {
			return err
		}
		// wait for all files to be processed 
		wg.Wait()

		return nil
	},
}

// getFilesToSearch determines whether the given path is a file or directory. If file, it adds the
// file to the chan. If directory, it searches through the directory by calling searchDirectory.
func getFilesToSearch(path string, input chan<- string) error {
	defer func() {
		close(input)
	}()

	// get information for the filename or path provided
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	// if file is not a directory, we search for the pattern in the file
	if !fileInfo.IsDir() {
		input <- fileInfo.Name()
		return nil
	}

	// if the file is a directory, we search all the files in the directory (recursively or not depending on flag)
	return searchDirectory(path, recursive, input)
}

// searchDirectory takes in a path and returns the file names of all files in that directory (recursive according to
// the recursive bool). If recursive is true, searchDirectory is called recursively
func searchDirectory(path string, recursive bool, input chan<- string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			input <- filepath.Join(path, entry.Name())
		} else if recursive {
			err := searchDirectory(filepath.Join(path, entry.Name()), recursive, input)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	goGrep.Flags().BoolVar(&recursive, "r", false, "search files recursively in directory")
}

func Execute() error {
	return goGrep.Execute()
}
