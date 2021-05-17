// Copyright (c) 2020 SIGHUP s.r.l All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"os"

	"github.com/sighupio/furyctl/pkg/distribution"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var distributionVersion string

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&distributionVersion, "version", "", "Specify the Kubernetes Fury Distribution version")
	err := initCmd.MarkFlagRequired("version")
	if err != nil {
		logrus.Print(err)
	}
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the minimum distribution configuration",
	Long:  "Initialize the current directory with the minimum distribution configuration",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		err = distribution.Download(distributionVersion, pwd)
		if err != nil {
			return err
		}
		return nil
	},
}
