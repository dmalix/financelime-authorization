package config

import (
	"fmt"
	"strconv"
)

const unset = "unset"

var (
	developmentMode  = unset
	versionNumber    = unset
	versionBuildTime = unset
	versionCommit    = unset
	versionCompiler  = unset
)

func InitVersion() (Version, error) {
	var err error
	var version Version

	if //goland:noinspection GoBoolExpressions
	developmentMode == unset || versionNumber == unset || versionBuildTime == unset || versionCommit == unset || versionCompiler == unset {
		return Version{}, fmt.Errorf("development mode and version variables not set")
	}

	version.DevelopmentMode, err = strconv.ParseBool(developmentMode)
	if err != nil {
		return Version{}, fmt.Errorf("failed to parse the DevelopmentMode variable: %s", err)
	}

	version.Number = versionNumber
	version.BuildTime = versionBuildTime
	version.Commit = versionCommit
	version.Compiler = versionCompiler

	return version, nil
}
