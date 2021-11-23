package modules

import (
	"os/user"

	"github.com/jwalton/kitsch-prompt/internal/kitsch/log"
	"gopkg.in/yaml.v3"
)

// UsernameModule shows the name of the currently logged in user.  This is,
// by default, hidden unless the user is root or the session is an SSH session.
// The CommonConfig.Style is applied by default, unless the user is Root in which
// case it is overridden by `UsernameConfig.RootStyle`.
//
// The username module provides the following template variables:
//
// • Username - The current user's username.
//
// • IsRoot - True if the user is root, false otherwise.
//
// • IsSSH - True if this is an SSH session, false otherwise.
//
// • Show - True if we should show the username module, false otherwise.
//
type UsernameModule struct {
	CommonConfig `yaml:",inline"`
	// ShowAlways will cause the username to always be shown.  If false (the default),
	// then the username will only be shown if the user is root, or the current
	// session is an SSH session.
	ShowAlways bool `yaml:"showAlways"`
	// RootStyle will be used in place of `Style` if the current user is root.
	// If this style is empty, will fall back to `Style`.
	RootStyle string `yaml:"rootStyle"`
}

type usernameModuleData struct {
	// username is the current user's username.
	username string
	// IsRoot is true if the current user is root.
	IsRoot bool
	// IsSSH is true if the user is in an SSH session.
	IsSSH bool
	// Show is true if the username module should be displayed.
	Show bool
}

// Username is the current user's username.
func (data usernameModuleData) Username() string {
	if data.username != "" {
		return data.username
	}

	// Fetch the user from the OS.  This can be a little slow, eating up around
	// 6ms on MaxOS and Linux style systems, which is why we prefer to get
	// the username from the env.  The good news is that `os/user` caches this
	// value for us, so repeated calls shouldn't be slow.
	user, err := user.Current()
	if err != nil {
		log.Info("Unable to get current user: " + err.Error())
		return ""
	}
	return user.Username
}

// Execute the username module.
func (mod UsernameModule) Execute(context *Context) ModuleResult {
	isRoot := context.Globals.IsRoot
	isSSH := context.Environment.HasSomeEnv("SSH_CLIENT", "SSH_CONNECTION", "SSH_TTY")
	show := isSSH || isRoot || mod.ShowAlways

	data := usernameModuleData{
		username: context.Environment.Getenv("USER"),
		IsRoot:   isRoot,
		IsSSH:    isSSH,
		Show:     show,
	}

	defaultText := ""
	style := mod.Style

	if show {
		defaultText = data.Username()
		if isRoot && mod.RootStyle != "" {
			style = mod.RootStyle
		}
	}

	return executeModule(context, mod.CommonConfig, data, style, defaultText)
}

func init() {
	registerFactory("username", func(node *yaml.Node) (Module, error) {
		var module UsernameModule
		err := node.Decode(&module)
		return &module, err
	})
}