// Copyright 2017 The Gop Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

// CmdRemove represents
var CmdRemove = cli.Command{
	Name:        "rm",
	Usage:       "remove a dependency",
	Description: `remove a dependency`,
	Action:      runRemove,
}

func runRemove(ctx *cli.Context) error {
	if len(ctx.Args()) <= 0 {
		return errors.New("No package to be removed")
	}

	_, projectRoot, err := analysisDirLevel()
	if err != nil {
		return err
	}

	for _, pkg := range ctx.Args() {
		dstPath := filepath.Join(projectRoot, "src", "vendor", pkg)
		fmt.Println("removing", pkg)
		os.RemoveAll(dstPath)
	}

	return nil
}
