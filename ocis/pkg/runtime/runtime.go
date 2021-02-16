package runtime

import (
	"context"
	"fmt"
	golog "log"
	"net/rpc"
	"os"
	"os/signal"
	"time"

	storage "github.com/owncloud/ocis/storage/pkg/command"
	storageConfig "github.com/owncloud/ocis/storage/pkg/config"

	"github.com/thejerf/suture"

	accounts "github.com/owncloud/ocis/accounts/pkg/command"
	settings "github.com/owncloud/ocis/settings/pkg/command"

	"github.com/micro/go-micro/v2"

	"github.com/owncloud/ocis/ocis/pkg/config"

	cli "github.com/micro/cli/v2"
	"github.com/micro/micro/v2/client/api"
	"github.com/micro/micro/v2/service/registry"

	"github.com/owncloud/ocis/ocis/pkg/runtime/process"
)

var (
	// OwncloudNamespace is the base path for micro' services to use
	OwncloudNamespace = "com.owncloud."

	// MicroServices to start as part of the fullstack option
	MicroServices = []string{
		"api", // :8080
		//"web",      // :8082
		"registry", // :8000
	}

	// Extensions are oCIS extension services
	Extensions = []string{
		"glauth",
		"idp",
		"ocs",
		"onlyoffice",
		"proxy",
		"settings",
		"store",
		"storage-frontend",
		"storage-gateway",
		"storage-userprovider",
		"storage-auth-basic",
		"storage-auth-bearer",
		"storage-home",
		"storage-users",
		"storage-metadata",
		"storage-public-link",
		"thumbnails",
		"web",
		"webdav",
		//"graph",
		//"graph-explorer",
	}

	// There seem to be a race condition when reva-sharing needs to read the sharing.json file and the parent folder is not present.
	dependants = []string{
		"accounts",
		"storage-sharing",
	}

	// Maximum number of retries until getting a connection to the rpc runtime service.
	maxRetries int = 10
)

// Runtime represents an oCIS runtime environment.
type Runtime struct {
	c *config.Config
}

// New creates a new oCIS + micro runtime
func New(cfg *config.Config) Runtime {
	return Runtime{
		c: cfg,
	}
}

// serviceTokens keeps in memory a set of [service name] = []suture.ServiceToken that is used to shutdown services.
// Shutting down a service implies removing it from the supervisor AND cancelling its context, this should be done
// within the service Stop() method. Services should cancel their context.
type serviceTokens map[string][]suture.ServiceToken

// Start rpc runtime
func (r *Runtime) Start() error {
	halt := make(chan os.Signal, 1)
	signal.Notify(halt, os.Interrupt)

	// tokens is a global state that keeps service tokens in memory
	tokens := serviceTokens{}

	supervisor := suture.NewSimple("ocis")
	globalCtx, _ := context.WithCancel(context.Background())

	settingsCtx, settingsCancel := context.WithCancel(globalCtx)
	r.c.Settings.Context = settingsCtx
	tokens["settings"] = append(tokens["settings"], supervisor.Add(settings.NewSutureService(settingsCtx, settingsCancel, r.c.Settings)))

	scfg := storageConfig.New()
	scfg.Log.Color = r.c.Log.Color
	scfg.Log.Pretty = r.c.Log.Pretty
	scfg.Log.Level = r.c.Log.Level
	storageMetadataCtx, storageMetadataCancel := context.WithCancel(globalCtx)
	scfg.Context = storageMetadataCtx
	tokens["storage-metadata"] = append(tokens["storage-metadata"], supervisor.Add(storage.NewStorageMetadataSutureService(storageMetadataCtx, storageMetadataCancel, scfg)))

	accountsCtx, accountsCancel := context.WithCancel(globalCtx)
	r.c.Accounts.Context = accountsCtx
	tokens["accounts"] = append(tokens["accounts"], supervisor.Add(accounts.NewAccountSutureService(accountsCtx, accountsCancel, r.c.Accounts)))

	go supervisor.ServeBackground()

	<-halt
	close(halt)
	//globalCancel()
	return nil
}

// Launch oCIS default oCIS extensions.
func (r *Runtime) Launch() {
	var client *rpc.Client
	var err error
	var try int

	for {
		if try >= maxRetries {
			golog.Fatal("could not get a connection to rpc runtime on localhost:10666")
		}
		client, err = rpc.DialHTTP("tcp", "localhost:10666")
		if err != nil {
			try++
			fmt.Println("runtime not available, retrying in 1 second...")
			time.Sleep(1 * time.Second)
		} else {
			goto OUT
		}
	}

OUT:
	for _, v := range MicroServices {
		RunService(client, v)
	}

	for _, v := range Extensions {
		RunService(client, v)
	}

	if len(dependants) > 0 {
		// TODO(refs) this should disappear and tackled at the runtime (pman) level.
		// see https://github.com/cs3org/reva/issues/795 for race condition.
		// dependants might not be needed on a ocis_simple build, therefore
		// it should not be started under these circumstances.
		time.Sleep(2 * time.Second)
		for _, v := range dependants {
			RunService(client, v)
		}
	}
}

// RunService sends a Service.Start command with the given service name  to pman
func RunService(client *rpc.Client, service string) int {
	args := process.NewProcEntry(service, os.Environ(), []string{service}...)

	all := append(Extensions, append(dependants, MicroServices...)...)
	if !contains(all, service) {
		return 1
	}

	var reply int
	if err := client.Call("Service.Start", args, &reply); err != nil {
		golog.Fatal(err)
	}
	return reply
}

// AddMicroPlatform adds the micro subcommands to the cli app
func AddMicroPlatform(app *cli.App, opts micro.Options) {
	setDefaults()

	app.Commands = append(app.Commands, api.Commands(micro.Registry(opts.Registry))...)
	//app.Commands = append(app.Commands, web.Commands(micro.Registry(opts.Registry))...)
	app.Commands = append(app.Commands, registry.Commands(micro.Registry(opts.Registry))...)
}

// provide a config.Config with default values?
func setDefaults() {
	// api
	api.Name = OwncloudNamespace + "api"
	api.Namespace = OwncloudNamespace + "api"
	api.HeaderPrefix = "X-Micro-Owncloud-"

	// web
	//web.Name = OwncloudNamespace + "web"
	//web.Namespace = OwncloudNamespace + "web"

	// registry
	registry.Name = OwncloudNamespace + "registry"
}

func contains(a []string, b string) bool {
	for i := range a {
		if a[i] == b {
			return true
		}
	}
	return false
}
