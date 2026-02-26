package commands

import (
	"context"
	"fmt"
)

type Version struct {
	version   string
	commit    string
	buildTime string
}

func NewVersion(version, commit, buildTime string) *Version {
	return &Version{
		version:   version,
		commit:    commit,
		buildTime: buildTime,
	}
}

func (c *Version) Name() string {
	return "version"
}

func (c *Version) Description() string {
	return "Print version information"
}

func (c *Version) Run(ctx context.Context, args []string) error {
	fmt.Printf("version=%s\n", c.version)
	fmt.Printf("commit=%s\n", c.commit)
	fmt.Printf("build_time=%s\n", c.buildTime)
	return nil
}
