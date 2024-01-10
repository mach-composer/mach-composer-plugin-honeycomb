package internal

import (
	"embed"
	"fmt"

	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/v2/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/v2/schema"
	"github.com/mitchellh/mapstructure"
)

//go:embed templates/*
var templates embed.FS

type HoneycombPlugin struct {
	provider         string
	environment      string
	globalConfig     *GlobalConfig
	siteConfigs      map[string]*SiteConfig
	componentConfigs map[string]*ComponentConfig
}

func NewHoneycombPlugin() schema.MachComposerPlugin {
	state := &HoneycombPlugin{
		provider:         "0.18.1", // Provider version of `honeycombio/honeycombio`
		siteConfigs:      map[string]*SiteConfig{},
		componentConfigs: map[string]*ComponentConfig{},
	}
	return plugin.NewPlugin(&schema.PluginSchema{
		Identifier:          "honeycomb",
		Configure:           state.Configure,
		GetValidationSchema: state.GetValidationSchema,

		SetGlobalConfig:        state.SetGlobalConfig,
		SetSiteConfig:          state.SetSiteConfig,
		SetComponentConfig:     state.SetComponentConfig,
		SetSiteComponentConfig: state.SetSiteComponentConfig,

		// Renders
		RenderTerraformProviders: state.RenderTerraformProviders,
		RenderTerraformResources: state.RenderTerraformResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
}

func (p *HoneycombPlugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		cfg = &SiteConfig{}
	}
	return cfg.extendGlobalConfig(p.globalConfig)
}

func (p *HoneycombPlugin) getSiteComponentConfig(site, name string) *SiteComponentConfig {
	siteCfg := p.getSiteConfig(site)
	if siteCfg == nil {
		return nil
	}

	cfg, ok := siteCfg.SiteComponents[name]
	if !ok {
		cfg = &SiteComponentConfig{}
	}

	return cfg.extendSiteConfig(siteCfg)
}

func (p *HoneycombPlugin) getComponentConfig(component string) *ComponentConfig {
	cfg, ok := p.componentConfigs[component]
	if !ok {
		cfg = &ComponentConfig{}
	}
	return cfg
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
	cfg := GlobalConfig{}

	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.globalConfig = &cfg

	return nil
}

func (p *HoneycombPlugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := SiteConfig{
		SiteComponents: map[string]*SiteComponentConfig{},
	}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = &cfg
	return nil
}

func (p *HoneycombPlugin) SetSiteComponentConfig(site string, component string, data map[string]any) error {
	cfg := SiteComponentConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}

	siteCfg, ok := p.siteConfigs[site]
	if !ok {
		siteCfg = &SiteConfig{
			SiteComponents: map[string]*SiteComponentConfig{},
		}
		p.siteConfigs[site] = siteCfg
	}
	siteCfg.SiteComponents[component] = &cfg

	return nil
}

func (p *HoneycombPlugin) SetComponentConfig(component, version string, data map[string]any) error {
	cfg := ComponentConfig{
		Version: version,
	}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}

	p.componentConfigs[component] = &cfg

	return nil
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

func (p *HoneycombPlugin) RenderTerraformComponent(site string, component string) (*schema.ComponentSchema, error) {
	cfg := p.getSiteComponentConfig(site, component)
	if cfg == nil {
		return nil, nil
	}

	componentConfig := p.getComponentConfig(component)

	vars, err := terraformRenderVariables(cfg)
	if err != nil {
		return nil, err
	}

	resources, err := terraformRenderComponentResources(component, componentConfig.Version, cfg)
	if err != nil {
		return nil, err
	}

	result := &schema.ComponentSchema{
		Variables: vars,
		Resources: resources,
	}

	return result, nil
}

func terraformRenderVariables(cfg *SiteComponentConfig) (string, error) {
	template := `
		honeycomb = {
			{{ renderProperty "api_key" .ApiKey }}
		}
	`

	vars, err := helpers.RenderGoTemplate(template, cfg)
	if err != nil {
		return "", err
	}

	return vars, nil
}

func terraformRenderComponentResources(component, version string, cfg *SiteComponentConfig) (string, error) {
	templateContext := struct {
		ComponentName string
		Version       string
		Config        *SiteComponentConfig
	}{
		ComponentName: component,
		Version:       version,
		Config:        cfg,
	}

	tpl, err := templates.ReadFile("templates/resources.tmpl")
	if err != nil {
		return "", err
	}

	return helpers.RenderGoTemplate(string(tpl), templateContext)
}
