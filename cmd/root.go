package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/just-hms/dow/pkg/logx"
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

func waitForDownload(logger logx.Logger, downloadPath string) (fs.FileInfo, error) {
	spinner := termx.NewSpin(logger)
	defer spinner.Flush()

	for {
		lastFile, err := getLastFile(downloadPath)
		if err != nil {
			return nil, err
		}

		sourcePath := filepath.Join(downloadPath, lastFile.Name())

		waited := false
		for osx.IsLocked(sourcePath) {
			lastFile, err = os.Stat(sourcePath)
			if err != nil {
				return nil, err
			}
			waited = true
			spinner.Spin("Downloading " + osx.Size(lastFile))
			time.Sleep(100 * time.Millisecond)
		}

		if waited {
			continue
		}

		time.Sleep(400 * time.Millisecond)

		checkLastFile, err := getLastFile(downloadPath)
		if err != nil {
			return nil, err
		}
		if checkLastFile.Name() != lastFile.Name() {
			continue
		}
		return checkLastFile, nil
	}
}

var rootCmd = &cobra.Command{
	Use:          `dow`,
	Short:        `mv the last downloaded file in the current (or the specified) folder`,
	Hidden:       true,
	SilenceUsage: false,
	Example:      "  dow\n  dow -v | xargs -rd '\n' code",
	Args:         cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		destFolder := "."
		if len(args) == 1 {
			destFolder = args[0]
		}

		downloadPath, err := osx.DownloadFolderPath()
		if err != nil {
			return fmt.Errorf("failed to get the dowload path: %v", err)
		}

		logger := logx.Logger{}

		lastFile, err := waitForDownload(logger, downloadPath)
		if err != nil {
			return fmt.Errorf("failed to get the last downloaded file: %v", err)
		}

		if time.Since(lastFile.ModTime()) > maxElapsedBeforeAsking && !yesFlag {
			logger.Printf(
				"%q is older than %v. Proceed? (y/N) ",
				lastFile.Name(),
				maxElapsedBeforeAsking,
			)

			resp, _ := termx.Read()
			if resp != 'y' {
				logger.Println("No")
				return nil
			}
			logger.Println()
		}

		sourcePath := filepath.Join(downloadPath, lastFile.Name())

		err = osx.Move(sourcePath, destFolder)
		if err != nil {
			return fmt.Errorf("failed to %v", err)
		}
		if verboseFlag {
			fmt.Println(filepath.Join(destFolder, lastFile.Name()))
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	rootCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "show the name of the moved file")
	rootCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "force dow to move the latest file even if it's old")
}
