package internal

// BaseConfig is the base honeycomb config.
type BaseConfig struct {
	ApiKey  string `mapstructure:"api_key"`
	DataSet string `mapstructure:"dataset"`
	Type    string `mapstructure:"type"`
	Url     string `mapstructure:"url"`
}

type GlobalConfig struct {
	BaseConfig `mapstructure:",squash"`
}

type SiteConfig struct {
	BaseConfig `mapstructure:",squash"`
	Components map[string]*SiteComponentConfig `mapstructure:"-"`
}

func (c *SiteConfig) extendGlobalConfig(g *GlobalConfig) *SiteConfig {
	cfg := &SiteConfig{
		BaseConfig: g.BaseConfig,
		Components: c.Components,
	}

	if c.ApiKey != "" {
		cfg.ApiKey = c.ApiKey
	}
	if c.DataSet != "" {
		cfg.DataSet = c.DataSet
	}
	if c.Type != "" {
		cfg.Type = c.Type
	}
	if c.Url != "" {
		cfg.Url = c.Url
	}

	return cfg
}

// SiteComponentConfig is for component specific sentry DSN settings
type SiteComponentConfig struct {
	BaseConfig `mapstructure:",squash"`
}

func (c *SiteComponentConfig) extendSiteConfig(s *SiteConfig) *SiteComponentConfig {
	cfg := &SiteComponentConfig{
		BaseConfig: s.BaseConfig,
	}

	if c.ApiKey != "" {
		cfg.ApiKey = c.ApiKey
	}
	if c.DataSet != "" {
		cfg.DataSet = c.DataSet
	}
	if c.Type != "" {
		cfg.Type = c.Type
	}
	if c.Url != "" {
		cfg.Url = c.Url
	}

	return cfg
}

type ComponentConfig struct {
	Version string `mapstructure:"-"`
}
