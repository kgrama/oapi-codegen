package openapiv3

import (
	"errors"
	"fmt"
	"reflect"
)

// Validate checks whether Configuration represent a valid configuration
func (o Configuration) Validate() error {
	if o.PackageName == "" {
		return errors.New("package name must be specified")
	}

	// Only one server type should be specified at a time.
	nServers := 0
	if o.Generate.IrisServer {
		nServers++
	}
	if o.Generate.ChiServer {
		nServers++
	}
	if o.Generate.FiberServer {
		nServers++
	}
	if o.Generate.EchoServer {
		nServers++
	}
	if o.Generate.GorillaServer {
		nServers++
	}
	if o.Generate.StdHTTPServer {
		nServers++
	}
	if o.Generate.GinServer {
		nServers++
	}
	if nServers > 1 {
		return errors.New("only one server type is supported at a time")
	}

	var errs []error
	if problems := o.Generate.Validate(); problems != nil {
		for k, v := range problems {
			errs = append(errs, fmt.Errorf("`generate` configuration for %v was incorrect: %v", k, v))
		}
	}

	if problems := o.Compatibility.Validate(); problems != nil {
		for k, v := range problems {
			errs = append(errs, fmt.Errorf("`compatibility-options` configuration for %v was incorrect: %v", k, v))
		}
	}

	if problems := o.OutputOptions.Validate(); problems != nil {
		for k, v := range problems {
			errs = append(errs, fmt.Errorf("`output-options` configuration for %v was incorrect: %v", k, v))
		}
	}

	err := errors.Join(errs...)
	if err != nil {
		return fmt.Errorf("failed to validate configuration: %w", err)
	}

	return nil
}

// UpdateDefaults sets reasonable default values for unset fields in Configuration
func (o Configuration) UpdateDefaults() Configuration {
	if reflect.ValueOf(o.Generate).IsZero() {
		o.Generate = GenerateOptions{
			EchoServer:   true,
			Models:       true,
			EmbeddedSpec: true,
		}
	}
	return o
}

func (oo GenerateOptions) Validate() map[string]string {
	return nil
}

func (co CompatibilityOptions) Validate() map[string]string {
	return nil
}

func (oo OutputOptions) Validate() map[string]string {
	return nil
}

type OutputOptionsOverlay struct {
	Path string `yaml:"path"`

	// Strict defines whether the Overlay should be applied in a strict way, highlighting any actions that will not take any effect. This can, however, lead to more work when testing new actions in an Overlay, so can be turned off with this setting.
	// Defaults to true.
	Strict *bool `yaml:"strict,omitempty"`
}
