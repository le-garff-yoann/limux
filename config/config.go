package config

import (
	"filemux/processor"
	"filemux/processor/impl/in"
	"filemux/processor/impl/out"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

// Conf covers the entire configuration of fimemux.
type Conf struct {
	Out []*out.Out `yaml:"out"`
	In  []*in.In   `yaml:"in"`
}

// New is the constructor to `filemux.Conf`.
func New(yamlBytes []byte) (*Conf, error) {
	var c Conf
	if err := yaml.Unmarshal(yamlBytes, &c); err != nil {
		return nil, fmt.Errorf("Cannot parse the configuration: %s", err)
	}

	for _, p := range c.Out {
		if err := p.Configure(); err != nil {
			return nil, err
		}
	}

	for _, p := range c.In {
		if err := p.Configure(); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// Processors returns the list of `processor.Processor`.
func (s Conf) Processors() []processor.Processor {
	var processors []processor.Processor

	for _, p := range s.In {
		processors = append(processors, p)
	}

	for _, p := range s.Out {
		processors = append(processors, p)
	}

	return processors
}
