package pl003

import (
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"

	"github.com/giantswarm/awscnfm/v12/pkg/table"
)

const (
	name  = "pl003"
	short = "Execute plan pl003 automatically."
)

type Config struct {
	Logger micrologger.Logger
}

func New(config Config) (*cobra.Command, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	f := &flag{}

	r := &runner{
		flag:   f,
		logger: config.Logger,
	}

	c := &cobra.Command{
		Use:   name,
		Short: short,
		Long:  mustLong(),
		RunE:  r.Run,
	}

	f.Init(c)

	return c, nil
}

func mustLong() string {
	var d time.Duration
	for _, s := range Plan {
		d += s.Backoff.Wait()
	}

	s := "Test plan pl003 launches a basic Tenant Cluster in the previous patch release\n"
	s += "and upgrades the Tenant Cluster to the latest patch release once it is up. Plan\n"
	s += "execution might take up to " + d.String() + ".\n\n"

	t := [][]string{{"ACTION", "RETRY", "WAIT"}}

	for _, s := range Plan {
		t = append(t, []string{s.Action, s.Backoff.Interval().String(), s.Backoff.Wait().String()})
	}

	colourized := table.Colourize(t)
	s += table.Format(colourized)

	return s
}