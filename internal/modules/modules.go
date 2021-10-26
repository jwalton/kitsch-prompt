// Package modules has modules which can generate parts of the kitsch prompt output.
//
// Each module produces some fragment of output which are assembled together into
// the final shell prompt.  For example the "username" module prints the name of the
// current user.  The "directory" module prints the current working directory.  The
// "block" module combines multiple modules together; it runs each child module
// in parallel, and then assembles up all the results.
//
// Because modules are intended to run in parallel, and because there are certain
// things that many different modules are all going to want to know (e.g. lots
// of programming-language oriented modules will want to know if files with
// a certain name or extension are present in the current folder), each module
// is passed an "env" object, which can be used to access information about the
// environment without duplicating effort (it would be silly if all the various
// programming language modules all read the contents of the current working
// directory - we only need to read it once).
//
package modules

import (
	"fmt"
	"os"
	"os/user"

	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/modtemplate"
	"github.com/jwalton/kitsch-prompt/internal/styling"
)

// ModuleResult represents the output of a module.
type ModuleResult struct {
	// Text contains the rendered output of the module, either the default text
	// generated by the module itself, or the output from the template if one
	// was specified.
	Text string
	// Data contains any template data generated by the module.
	Data interface{}
	// StartStyle contains the foregraound and background colors of the first
	// character in Text.  Note that this is based on the declared style for the
	// module - if the style for the module says the string should be colored
	// blue, but a template is used to change the color of the first character
	// to red, this will still say it is blue.
	StartStyle styling.CharacterColors
	// EndStyle is similar to StartStyle, but contains the colors  of the last
	// character in Text.
	EndStyle styling.CharacterColors
}

// Globals is a collection of "global" values that are passed to all modules.
// These values are available to templates via the ".Globals" property.
type Globals struct {
	// CWD is the current wordking directory.
	CWD string
	// Home is the user's home directory.
	Home string
	// Username is the user's username.
	// TODO: Add the "short" username for MacOS and Windows.
	Username string
	// UserFullName is the user's full name.
	UserFullName string
	// Hostname is the name of the current machine.
	Hostname string
	// Status is the return status of the previous command.
	Status int
	// PreviousCommandDuration is the duration of the previous command, in milliseconds.
	PreviousCommandDuration int64
	// Keymap is the zsh/fish keymap. TODO: What values can this have?
	Keymap string
}

// NewGlobals creates a new Globals object.
func NewGlobals(
	status int,
	previousCommandDuration int64,
	keymap string,
) Globals {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}

	user, err := user.Current()
	username := ""
	name := ""
	if err == nil {
		username = user.Username
		name = user.Name
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	return Globals{
		CWD:                     cwd,
		Home:                    home,
		Username:                username,
		UserFullName:            name,
		Hostname:                hostname,
		Status:                  status,
		PreviousCommandDuration: previousCommandDuration,
	}
}

// Context is a set of common parameters passed to Module.Execute.
type Context struct {
	// Environment is the environment to fetch data from.
	Environment env.Env
	// Styles is the style registry to use to create styles.
	Styles styling.Registry
	// Globals is a collection of "global" values that are passed to all modules.
	// These values are available to templates via the ".Globals" property.
	Globals Globals
}

// Module represnts a module that generates some output to show in the prompt.
type Module interface {
	// Execute will execute this module and return a ModuleResult.
	Execute(context *Context) ModuleResult
}

// CommonConfig is common configuration for all modules.
type CommonConfig struct {
	// Style is the style to apply to this module.
	Style string `yaml:"style"`
	// Template is a golang template to use to render the output of this module.
	Template string `yaml:"template"`
}

// TemplateData is the common data structure passed to a template when it is executed.
type TemplateData struct {
	// Text is the default text produced by this module
	Text string
	// Data is the data for this template.
	Data interface{}
	// Global is the global data.
	Global *Globals
}

// executeModule is called to execute a module.  This handles "common" stuff that
// all modules do, like calling templates.
func executeModule(
	context *Context,
	config CommonConfig,
	data interface{},
	styleStr string,
	defaultText string,
) ModuleResult {
	style, err := context.Styles.Get(styleStr)
	if err != nil {
		style = nil
		context.Environment.Warn(err.Error())
	}

	text := defaultText

	var startStyle styling.CharacterColors
	var endStyle styling.CharacterColors

	if config.Template != "" {
		tmpl, err := modtemplate.CompileTemplate(&context.Styles, "module-template", config.Template)
		if err != nil {
			// FIX: Should add this error to a list of warnings for this module.
			fmt.Printf("Error compiling template: %v", err)
		} else {
			text, err = modtemplate.TemplateToString(tmpl, TemplateData{
				Data:   data,
				Global: &context.Globals,
				Text:   defaultText,
			})
			if err != nil {
				context.Environment.Warn(fmt.Sprintf("Error executing template:\n%s\n%v", config.Template, err))
				text = defaultText
			}
		}
	}

	if style != nil && text != "" {
		text, startStyle, endStyle = style.ApplyGetColors(text)
	}

	return ModuleResult{
		Text:       text,
		Data:       data,
		StartStyle: startStyle,
		EndStyle:   endStyle,
	}
}

// defaultString returns value if it is non-empty, or def otherwise.
func defaultString(value string, def string) string {
	if value != "" {
		return value
	}
	return def
}

func defaultStyle(context *Context, styleString string, defStyle string) *styling.Style {
	style, err := context.Styles.Get(styleString)
	if err != nil {
		context.Environment.Warn(err.Error())
	}
	if styleString == "" || err != nil {
		style, err = context.Styles.Get(defStyle)
		if err != nil {
			panic("Error parsing default style: " + err.Error())
		}
	}

	return style
}
