package main

import (
	"strings"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
	"go.mattglei.ch/timber"
)

const (
	appKitPath              = "/System/Library/Frameworks/AppKit.framework/AppKit"
	activationPolicyRegular = 0
	finderBundleID          = "com.apple.finder"
)

var (
	selSharedWorkspace      = objc.RegisterName("sharedWorkspace")
	selFrontmostApplication = objc.RegisterName("frontmostApplication")
	selRunningApplications  = objc.RegisterName("runningApplications")
	selCount                = objc.RegisterName("count")
	selObjectAtIndex        = objc.RegisterName("objectAtIndex:")
	selProcessIdentifier    = objc.RegisterName("processIdentifier")
	selActivationPolicy     = objc.RegisterName("activationPolicy")
	selBundleIdentifier     = objc.RegisterName("bundleIdentifier")
	selLocalizedName        = objc.RegisterName("localizedName")
	selTerminate            = objc.RegisterName("terminate")
	selUTF8String           = objc.RegisterName("UTF8String")
)

func quitApps() {
	_, err := purego.Dlopen(appKitPath, purego.RTLD_GLOBAL|purego.RTLD_LAZY)
	if err != nil {
		timber.Fatal(err, "failed to load AppKit")
	}

	workspace := objc.ID(objc.GetClass("NSWorkspace")).Send(selSharedWorkspace)
	current := objc.Send[int32](workspace.Send(selFrontmostApplication), selProcessIdentifier)
	apps := workspace.Send(selRunningApplications)

	quit := []string{}
	for i := range objc.Send[uint](apps, selCount) {
		app := apps.Send(selObjectAtIndex, i)
		if objc.Send[int](app, selActivationPolicy) != activationPolicyRegular {
			continue
		}
		if objc.Send[int32](app, selProcessIdentifier) == current {
			continue
		}
		if goString(app.Send(selBundleIdentifier)) == finderBundleID {
			continue
		}

		if objc.Send[bool](app, selTerminate) {
			quit = append(quit, goString(app.Send(selLocalizedName)))
		}
	}

	if len(quit) == 0 {
		return
	}
	timber.Donef("quit %s", strings.Join(quit, ", "))
}

func goString(nsString objc.ID) string {
	return objc.Send[string](nsString, selUTF8String)
}
