// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/copilot-cli/internal/pkg/aws/identity"

	"github.com/aws/copilot-cli/internal/pkg/aws/codepipeline"
	"github.com/aws/copilot-cli/internal/pkg/aws/sessions"
	"github.com/aws/copilot-cli/internal/pkg/config"
	"github.com/aws/copilot-cli/internal/pkg/deploy"
	"github.com/aws/copilot-cli/internal/pkg/term/prompt"
	"github.com/aws/copilot-cli/internal/pkg/term/selector"

	"github.com/aws/copilot-cli/internal/pkg/describe"
	"github.com/aws/copilot-cli/internal/pkg/term/log"
	"github.com/spf13/cobra"
)

const (
	appShowNamePrompt     = "Which application would you like to show?"
	appShowNameHelpPrompt = "An application is a collection of related services."
)

type showAppVars struct {
	name             string
	shouldOutputJSON bool
}

type showAppOpts struct {
	showAppVars

	store            store
	w                io.Writer
	sel              appSelector
	pipelineSvc      pipelineGetter
	newVersionGetter func(string) (versionGetter, error)
}

func newShowAppOpts(vars showAppVars) (*showAppOpts, error) {
	defaultSession, err := sessions.NewProvider(sessions.UserAgentExtras("app show")).Default()
	if err != nil {
		return nil, fmt.Errorf("default session: %w", err)
	}
	store := config.NewSSMStore(identity.New(defaultSession), ssm.New(defaultSession), aws.StringValue(defaultSession.Config.Region))
	return &showAppOpts{
		showAppVars: vars,
		store:       store,
		w:           log.OutputWriter,
		sel:         selector.NewSelect(prompt.New(), store),
		pipelineSvc: codepipeline.New(defaultSession),
		newVersionGetter: func(s string) (versionGetter, error) {
			d, err := describe.NewAppDescriber(s)
			if err != nil {
				return d, fmt.Errorf("new app describer for application %s: %v", s, err)
			}
			return d, nil
		},
	}, nil
}

// Validate returns an error if the values provided by the user are invalid.
func (o *showAppOpts) Validate() error {
	if o.name != "" {
		_, err := o.store.GetApplication(o.name)
		if err != nil {
			return fmt.Errorf("get application %s: %w", o.name, err)
		}
	}

	return nil
}

// Ask asks for fields that are required but not passed in.
func (o *showAppOpts) Ask() error {
	if err := o.askName(); err != nil {
		return err
	}

	return nil
}

// Execute writes the application's description.
func (o *showAppOpts) Execute() error {
	description, err := o.description()
	if err != nil {
		return err
	}
	if !o.shouldOutputJSON {
		fmt.Fprint(o.w, description.HumanString())
		return nil
	}
	data, err := description.JSONString()
	if err != nil {
		return fmt.Errorf("get JSON string: %w", err)
	}
	fmt.Fprint(o.w, data)
	return nil
}

func (o *showAppOpts) description() (*describe.App, error) {
	app, err := o.store.GetApplication(o.name)
	if err != nil {
		return nil, fmt.Errorf("get application %s: %w", o.name, err)
	}
	envs, err := o.store.ListEnvironments(o.name)
	if err != nil {
		return nil, fmt.Errorf("list environments in application %s: %w", o.name, err)
	}
	svcs, err := o.store.ListServices(o.name)
	if err != nil {
		return nil, fmt.Errorf("list services in application %s: %w", o.name, err)
	}

	pipelines, err := o.pipelineSvc.GetPipelinesByTags(map[string]string{
		deploy.AppTagKey: o.name,
	})

	if err != nil {
		return nil, fmt.Errorf("list pipelines in application %s: %w", o.name, err)
	}

	var trimmedEnvs []*config.Environment
	for _, env := range envs {
		trimmedEnvs = append(trimmedEnvs, &config.Environment{
			Name:      env.Name,
			AccountID: env.AccountID,
			Region:    env.Region,
			Prod:      env.Prod,
		})
	}
	var trimmedSvcs []*config.Workload
	for _, svc := range svcs {
		trimmedSvcs = append(trimmedSvcs, &config.Workload{
			Name: svc.Name,
			Type: svc.Type,
		})
	}
	versionGetter, err := o.newVersionGetter(o.name)
	if err != nil {
		return nil, err
	}
	version, err := versionGetter.Version()
	if err != nil {
		return nil, fmt.Errorf("get version for application %s: %w", o.name, err)
	}
	return &describe.App{
		Name:      app.Name,
		Version:   version,
		URI:       app.Domain,
		Envs:      trimmedEnvs,
		Services:  trimmedSvcs,
		Pipelines: pipelines,
	}, nil
}

func (o *showAppOpts) askName() error {
	if o.name != "" {
		return nil
	}
	name, err := o.sel.Application(appShowNamePrompt, appShowNameHelpPrompt)
	if err != nil {
		return fmt.Errorf("select application: %w", err)
	}
	o.name = name
	return nil
}

// buildAppShowCmd builds the command for showing details of an application.
func buildAppShowCmd() *cobra.Command {
	vars := showAppVars{}
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Shows info about an application.",
		Long:  "Shows configuration, environments and services for an application.",
		Example: `
  Shows info about the application "my-app"
  /code $ copilot app show -n my-app`,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newShowAppOpts(vars)
			if err != nil {
				return err
			}
			return run(opts)
		}),
	}
	// The flags bound by viper are available to all sub-commands through viper.GetString({flagName})
	cmd.Flags().BoolVar(&vars.shouldOutputJSON, jsonFlag, false, jsonFlagDescription)
	cmd.Flags().StringVarP(&vars.name, nameFlag, nameFlagShort, tryReadingAppName(), appFlagDescription)
	return cmd
}
