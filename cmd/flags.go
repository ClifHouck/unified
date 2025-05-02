package cmd

import (
	"github.com/spf13/pflag"
)

var idOnly = false
var listingFlagSet = pflag.NewFlagSet("listing", pflag.ExitOnError)
