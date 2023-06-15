# Honeycomb Plugin for MACH composer

This repository contains the Honeycomb plugin for Mach Composer. It requires MACH
composer >= 2.5.x

This plugin uses the (Honeycomb Terraform Provider)[https://github.com/honeycombio/terraform-provider-honeycombio]




## Usage

```yaml
mach_composer:
  version: 1
  plugins:
    honeycomb:
      source: mach-composer/honeycomb
      version: 0.1.0

global:
  # ...
  honeycomb:
    api_key: "12345"

sites:
  - identifier: my-site
    # ...
    honeycomb:
      api_key: "34567"
    components:
      - name: my-component
        # ...
```
