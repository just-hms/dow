package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/just-hms/dow/pkg/osx"
	"github.com/just-hms/dow/pkg/termx"
	"github.com/spf13/cobra"
)

var (
	verboseFlag bool
	yesFlag     bool
)

const maxElapsedBeforeAsking = time.Minute

func getLastFile(path string) (fs.FileInfo, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	lastFile, err := osx.LatestFile(files)
	if err != nil {
		return nil, err
	}
	if lastFile.IsDir() {
		return nil, errors.New("cannot move a directory")
	}
	return lastFile, nil
}

var mvCmd = &cobra.Command{
	Use:          "dow",
	Short:        "move the last downloaded file in the current (or the specified) folder",
	Long:         ``,
	SilenceUsage: false,
	Args:         cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		destPath := "."
		if len(args) == 1 {
			destPath = args[0]
		}

		downloadPath, err := osx.DownloadFolderPath()
		if err != nil {
			return fmt.Errorf("failed to get the dowload path: %v", err)
		}

		var (
			lastFile   fs.FileInfo
			sourcePath string
		)

		// TODO: this is basically a wait for folder
		// - https://superuser.com/questions/860064/how-can-i-find-all-files-open-within-a-given-directory#:~:text=lsof%20has%20switches,open%20files%20recursively)
		// check special names like unconfirmed...

		spinner := termx.NewSpin()
		for {
			lastFile, err = getLastFile(downloadPath)
			if err != nil {
				return fmt.Errorf("failed to get latest file: %v", err)
			}

			sourcePath = filepath.Join(downloadPath, lastFile.Name())

			waited := false
			for osx.IsLocked(sourcePath) {
				lastFile, err = os.Stat(sourcePath)
				if err != nil {
					return fmt.Errorf("failed to get latest file: %v", err)
				}
				waited = true
				spinner.Spin("Downloading " + osx.Size(lastFile))
				time.Sleep(100 * time.Millisecond)
			}
			if !waited {
				break
			}
			time.Sleep(400 * time.Millisecond)
		}
		spinner.Flush()

		if time.Since(lastFile.ModTime()) > maxElapsedBeforeAsking && !yesFlag {
			fmt.Printf(
				"%q is older than %v. Proceed? (y/N) ",
				lastFile.Name(),
				maxElapsedBeforeAsking,
			)

			resp, _ := termx.Read()
			if resp != 'y' {
				fmt.Println("No")
				return nil
			}
			fmt.Println()
		}

		err = osx.Move(sourcePath, destPath)
		if err != nil {
			return fmt.Errorf("failed to %v", err)
		}
		if verboseFlag {
			fmt.Println(filepath.Join(destPath, filepath.Base(sourcePath)))
		}
		return nil
	},
}

func init() {
	mvCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "show the name of the moved file")
	mvCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "force dow to move the latest file even if it's old")
}
