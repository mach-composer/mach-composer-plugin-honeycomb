package internal

import (
	"fmt"

	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/mitchellh/mapstructure"
)

type HoneycombPlugin struct {
	provider     string
	environment  string
	globalConfig *HoneycombConfig
	siteConfigs  map[string]*HoneycombConfig
	enabled      bool
}

func NewHoneycombPlugin() schema.MachComposerPlugin {
	state := &HoneycombPlugin{
		provider:    "0.18.1", // Provider version of `honeycombio/honeycombio`
		siteConfigs: map[string]*HoneycombConfig{},
	}
	return plugin.NewPlugin(&schema.PluginSchema{
		Identifier:          "honeycomb",
		Configure:           state.Configure,
		IsEnabled:           func() bool { return state.enabled },
		GetValidationSchema: state.GetValidationSchema,

		SetGlobalConfig: state.SetGlobalConfig,
		SetSiteConfig:   state.SetSiteConfig,

		// Renders
		RenderTerraformProviders: state.RenderTerraformProviders,
		RenderTerraformResources: state.RenderTerraformResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
}

func (p *HoneycombPlugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *HoneycombPlugin) GetValidationSchema() (*schema.ValidationSchema, error) {
	result := getSchema()
	return result, nil
}

func (p *HoneycombPlugin) SetGlobalConfig(data map[string]any) error {
	cfg := HoneycombConfig{}

	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.globalConfig = &cfg
	p.enabled = true

	return nil
}

func (p *HoneycombPlugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := HoneycombConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = &cfg
	p.enabled = true
	return nil
}

func (p *HoneycombPlugin) getSiteConfig(site string) *HoneycombConfig {
	result := &HoneycombConfig{}
	if p.globalConfig != nil {
		result.ApiKey = p.globalConfig.ApiKey
	}

	cfg, ok := p.siteConfigs[site]
	if ok {
		if cfg.ApiKey != "" {
			result.ApiKey = cfg.ApiKey
		}
	}

	return result
}

func (p *HoneycombPlugin) RenderTerraformStateBackend(_ string) (string, error) {
	return "", nil
}

func (p *HoneycombPlugin) RenderTerraformProviders(site string) (string, error) {
	cfg := p.getSiteConfig(site)

	if cfg == nil {
		return "", nil
	}

	result := fmt.Sprintf(`
		honeycombio = {
			source = "honeycombio/honeycombio"
			version = "%s"
		}
	`, helpers.VersionConstraint(p.provider))

	return result, nil
}

func (p *HoneycombPlugin) RenderTerraformResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)

	if cfg == nil {
		return "", nil
	}

	template := `
		provider "honeycombio" {
			{{ renderProperty "api_key" .ApiKey }}
		}
	`
	return helpers.RenderGoTemplate(template, cfg)
}

func (p *HoneycombPlugin) RenderTerraformComponent(site string, _ string) (*schema.ComponentSchema, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return nil, nil
	}

	template := ``

	vars, err := helpers.RenderGoTemplate(template, cfg)
	if err != nil {
		return nil, err
	}

	result := &schema.ComponentSchema{
		Variables: vars,
	}

	return result, nil
}
