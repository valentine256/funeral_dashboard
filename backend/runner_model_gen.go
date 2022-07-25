package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/Yamashou/gqlgenc/clientgen"
	gqlgencConfig "github.com/Yamashou/gqlgenc/config"
)

func mutateHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	for _, model := range b.Models {
		for _, field := range model.Fields {
			// this is the logic to add omitempty
			omit := strings.TrimSuffix(field.Tag, `"`)
			field.Tag = fmt.Sprintf(`%v,omitempty"`, omit)
		}
	}
	return b
}

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load gqlgen config", err.Error())
		os.Exit(2)
	}

	clientCfg, err := gqlgencConfig.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load gqlgenc config", err.Error())
		os.Exit(2)
	}

	clientPlugin := clientgen.New(clientCfg.Query, clientCfg.Client, clientCfg.Generate)

	p := modelgen.Plugin{
		MutateHook: mutateHook,
	}

	err = api.Generate(cfg,
		api.NoPlugins(),
		api.AddPlugin(&p),
		api.AddPlugin(clientPlugin),
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
