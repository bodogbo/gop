// Copyright 2017 The Gop Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Unknwon/com"

	"github.com/urfave/cli"
)

// CmdBuild represents
var CmdBuild = cli.Command{
	Name:            "build",
	Usage:           "Build this project",
	Description:     `Build this project`,
	Action:          runBuild,
	SkipFlagParsing: true,
}

var curTarget *Target

func analysisTarget(level int, targetName, projectRoot string) error {
	if targetName == "" {
		if level == dirLevelTarget {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			dirName := filepath.Base(wd)

			for _, t := range config.Targets {
				if t.Dir == dirName {
					curTarget = &t
					break
				}
			}

			var name = dirName
			if dirName == "main" {
				name = filepath.Base(projectRoot)
			}

			curTarget = &Target{
				Name: name,
				Dir:  dirName,
			}
		}
		if curTarget == nil {
			curTarget = &config.Targets[0]
		}
	} else {
		for _, t := range config.Targets {
			if t.Name == targetName {
				curTarget = &t
				break
			}
			if t.Dir == targetName {
				curTarget = &t
				break
			}
		}

		if curTarget == nil {
			if !com.IsExist(filepath.Join(projectRoot, "src", targetName)) {
				return errors.New("unknow target")
			}

			curTarget = &Target{
				Name: targetName,
				Dir:  targetName,
			}
		}
	}
	return nil
}

func runBuildNoCtx(args []string) error {
	level, projectRoot, err := analysisDirLevel()
	if err != nil {
		return err
	}

	if err = loadConfig(filepath.Join(projectRoot, "gop.yml")); err != nil {
		return err
	}

	var targetName string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		targetName = args[0]
		args = args[1:]
	}

	if err = analysisTarget(level, targetName, projectRoot); err != nil {
		return err
	}

	var find = -1
	for i, arg := range args {
		if arg == "-o" {
			find = i
			break
		}
	}

	var ext string
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	if find > -1 {
		if find < len(args)-1 {
			curTarget.Name = args[find+1]
		} else {
			args = append(args[:find], "-o", curTarget.Name+ext)
		}
	} else {
		args = append(args, "-o", curTarget.Name+ext)
	}

	cmd := NewCommand("build").AddArguments(args...)
	envs := os.Environ()
	var gopathIdx = -1
	for i, env := range envs {
		if strings.HasPrefix(env, "GOPATH=") {
			gopathIdx = i
			break
		}
	}

	newGopath := fmt.Sprintf("GOPATH=%s", projectRoot)
	if gopathIdx > 0 {
		envs[gopathIdx] = newGopath
	} else {
		envs = append(envs, newGopath)
	}
	cmd.Env = envs

	err = cmd.RunInDirPipeline(filepath.Join(projectRoot, "src", curTarget.Dir), os.Stdout, os.Stderr)
	if err != nil {
		return err
	}

	return nil
}

func runBuild(ctx *cli.Context) error {
	return runBuildNoCtx(ctx.Args())
}
