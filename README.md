# devcontainer-features-update-presets

This repository contains preset files for keeping versions of tools installed via our devcontainer-features up-to-date.

Presets for the following tools are available:
* [Renovate](https://docs.renovatebot.com)
* [gonovate](https://github.com/Roemer/gonovate)

You can either use the main version of the presets or one of the git-tags in case you want a fixed version of the presets.

## Currently unsupported features

The following features currently do not support auto-updating:

* instant-client
* nvidia-cuda

## Usage

### Renovate

Extend the JSON preset from your Renovate configuration:

```json
{
	"extends": [
		"github>postfinance/devcontainer-features-update-presets:renovate-preset"
	]
}
```

or with a fixed tag:

```json
{
	"extends": [
		"github>postfinance/devcontainer-features-update-presets:renovate-preset#v1.0.0"
	]
}
```

Also add the custom regex manager:

```json
{
    "enabledManagers": [
        "custom.regex"
    ],
}
```

If you also want to update the devcontainer features themselves, also add the official `devcontainer` manager:

```json
{
    "enabledManagers": [
        "devcontainer",
        "custom.regex"
    ],
}
```

### Gonovate

Add the preset to the extends of you gonovate config:

```yaml
extends:
  - defaults
  - https://raw.githubusercontent.com/postfinance/devcontainer-features-update-presets/refs/heads/main/gonovate-preset.yaml
```

Enable the devcontainer manager:

```yaml
managers:
  - id: devcontainer
    type: devcontainer
```
