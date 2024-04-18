/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
)

func Execute() {
	err := mvCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
