// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/copilot-cli/internal/pkg/aws/identity"
	"github.com/aws/copilot-cli/internal/pkg/aws/sessions"

	"github.com/aws/copilot-cli/internal/pkg/config"
	"github.com/aws/copilot-cli/internal/pkg/deploy"
	"github.com/aws/copilot-cli/internal/pkg/describe"
	"github.com/aws/copilot-cli/internal/pkg/term/color"
	"github.com/aws/copilot-cli/internal/pkg/term/log"
	"github.com/aws/copilot-cli/internal/pkg/term/prompt"
	"github.com/aws/copilot-cli/internal/pkg/term/selector"
	"github.com/spf13/cobra"
)

const (
	envShowAppNamePrompt     = "Which application is the environment in?"
	envShowAppNameHelpPrompt = "An application is a collection of related services."
	envShowNamePrompt        = "Which environment of %s would you like to show?"
	envShowHelpPrompt        = "The detail of an environment will be shown (e.g., region, account ID, services)."
)

type showEnvVars struct {
	appName               string
	name                  string
	shouldOutputJSON      bool
	shouldOutputResources bool
}

type showEnvOpts struct {
	showEnvVars

	w                io.Writer
	store            store
	describer        envDescriber
	sel              configSelector
	initEnvDescriber func() error
}

func newShowEnvOpts(vars showEnvVars) (*showEnvOpts, error) {
	sessProvider := sessions.NewProvider(sessions.UserAgentExtras("env show"))
	defaultSess, err := sessProvider.Default()
	if err != nil {
		return nil, err
	}
	store := config.NewSSMStore(identity.New(defaultSess), ssm.New(defaultSess), aws.StringValue(defaultSess.Config.Region))

	deployStore, err := deploy.NewStore(store)
	if err != nil {
		return nil, fmt.Errorf("connect to copilot deploy store: %w", err)
	}

	opts := &showEnvOpts{
		showEnvVars: vars,
		store:       store,
		w:           log.OutputWriter,
		sel:         selector.NewConfigSelect(prompt.New(), store),
	}
	opts.initEnvDescriber = func() error {
		d, err := describe.NewEnvDescriber(describe.NewEnvDescriberConfig{
			App:             opts.appName,
			Env:             opts.name,
			ConfigStore:     store,
			DeployStore:     deployStore,
			EnableResources: opts.shouldOutputResources,
		})
		if err != nil {
			return fmt.Errorf("creating describer for environment %s in application %s: %w", opts.name, opts.appName, err)
		}
		opts.describer = d
		return nil
	}
	return opts, nil
}

// Validate returns an error if the values provided by the user are invalid.
func (o *showEnvOpts) Validate() error {
	if o.appName == "" {
		return nil
	}
	if _, err := o.store.GetApplication(o.appName); err != nil {
		return err
	}
	if o.name != "" {
		if _, err := o.store.GetEnvironment(o.appName, o.name); err != nil {
			return err
		}
	}
	return nil
}

// Ask asks for fields that are required but not passed in.
func (o *showEnvOpts) Ask() error {
	if err := o.askApp(); err != nil {
		return err
	}
	return o.askEnvName()
}

// Execute shows the environments through the prompt.
func (o *showEnvOpts) Execute() error {
	if err := o.initEnvDescriber(); err != nil {
		return err
	}
	env, err := o.describer.Describe()
	if err != nil {
		return fmt.Errorf("describe environment %s: %w", o.name, err)
	}
	if o.shouldOutputJSON {
		data, err := env.JSONString()
		if err != nil {
			return err
		}
		fmt.Fprint(o.w, data)
	} else {
		fmt.Fprint(o.w, env.HumanString())
	}

	return nil
}

func (o *showEnvOpts) askApp() error {
	if o.appName != "" {
		return nil
	}
	app, err := o.sel.Application(envShowAppNamePrompt, envShowAppNameHelpPrompt)
	if err != nil {
		return fmt.Errorf("select application: %w", err)
	}
	o.appName = app
	return nil
}

func (o *showEnvOpts) askEnvName() error {
	//return if env name is set by flag
	if o.name != "" {
		return nil
	}
	env, err := o.sel.Environment(fmt.Sprintf(envShowNamePrompt, color.HighlightUserInput(o.appName)), envShowHelpPrompt, o.appName)
	if err != nil {
		return fmt.Errorf("select environment for application %s: %w", o.appName, err)
	}
	o.name = env

	return nil
}

// buildEnvShowCmd builds the command for showing environments in an application.
func buildEnvShowCmd() *cobra.Command {
	vars := showEnvVars{}
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Shows info about a deployed environment.",
		Long:  "Shows info about a deployed environment, including region, account ID, and services.",

		Example: `
  Shows info about the environment "test".
  /code $ copilot env show -n test`,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newShowEnvOpts(vars)
			if err != nil {
				return err
			}
			return run(opts)
		}),
	}
	cmd.Flags().StringVarP(&vars.appName, appFlag, appFlagShort, tryReadingAppName(), appFlagDescription)
	cmd.Flags().StringVarP(&vars.name, nameFlag, nameFlagShort, "", envFlagDescription)
	cmd.Flags().BoolVar(&vars.shouldOutputJSON, jsonFlag, false, jsonFlagDescription)
	cmd.Flags().BoolVar(&vars.shouldOutputResources, resourcesFlag, false, envResourcesFlagDescription)
	return cmd
}
