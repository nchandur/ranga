package uci

import (
	"fmt"
	"strings"
)

// defines metadata and update action for UCI options
type Option struct {
	Name       string
	Type       string // "spin", "check", "string", "button", "combo"
	DefaultVal string
	Min        int
	Max        int
	Apply      func(value string) error
}

// holds registered options
type OptionRegistry struct {
	options map[string]Option
	order   []string
}

func NewOptionRegistry() *OptionRegistry {
	return &OptionRegistry{
		options: make(map[string]Option),
		order:   make([]string, 0),
	}
}

// adds option to engine
func (r *OptionRegistry) Register(opt Option) {
	key := strings.ToLower(opt.Name)
	if _, exists := r.options[key]; !exists {
		r.order = append(r.order, key)
	}
	r.options[key] = opt
}

// outputs the 'option name ...'.
func (r *OptionRegistry) Print() {
	for _, key := range r.order {
		opt := r.options[key]
		switch opt.Type {
		case "spin":
			fmt.Printf("option name %s type spin default %s min %d max %d\n",
				opt.Name, opt.DefaultVal, opt.Min, opt.Max)
		case "check":
			fmt.Printf("option name %s type check default %s\n",
				opt.Name, opt.DefaultVal)
		case "string":
			fmt.Printf("option name %s type string default %s\n",
				opt.Name, opt.DefaultVal)
		case "button":
			fmt.Printf("option name %s type button\n", opt.Name)
		}
	}
}

// parses payload and invokes registered handler
func (r *OptionRegistry) Set(payload string) {
	var optName, optValue string

	if strings.Contains(payload, " value ") {
		parts := strings.SplitN(payload, " value ", 2)
		optName = strings.TrimSpace(strings.TrimPrefix(parts[0], "name "))
		optValue = strings.TrimSpace(parts[1])
	} else {
		optName = strings.TrimSpace(strings.TrimPrefix(payload, "name "))
	}

	key := strings.ToLower(optName)
	opt, found := r.options[key]
	if !found {
		return
	}

	if opt.Apply != nil {
		opt.Apply(optValue)
	}
}
