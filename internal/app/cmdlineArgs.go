package application

import (
	"flag"
	"fmt"
	"os"

	"github.com/coreos/pkg/flagutil"
	yaml "gopkg.in/yaml.v2"
)

type cmdlineArgs struct {
	InsecurePort        int    `yaml:"InsecurePort"`
	InsecureBindAddress string `yaml:"InsecureBindAddress"`
	Port                int    `yaml:"Port"`
	BindAddress         string `yaml:"BindAddress"`
	RootPath            string `yaml:"RootPath"`
	CertDir             string `yaml:"CertDir"`
	TLSCertFile         string `yaml:"TLSCertFile"`
	TLSKeyFile          string `yaml:"TLSKeyFile"`
	MTLSCACertFile      string `yaml:"MTLSCACertFile"`
	WriteDefaultConfig  bool   `yaml:"-"`
	ConfigFile          string `yaml:"-"`
}

func (app *Application) ParseCommandlineAndEnvironment(args []string) {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.IntVar(&app.args.InsecurePort, "insecure-port", 0, "port to listen to for HTTP requests. (set to \"0\" to disable insecure communication)")
	flags.StringVar(&app.args.InsecureBindAddress, "insecure-bind-address", "127.0.0.1", "The IP address on which to serve the --insecure-port (set to 127.0.0.1 for loopback only).")
	flags.IntVar(&app.args.Port, "port", 443, "The secure port to listen to for incoming HTTPS requests.")
	flags.StringVar(&app.args.BindAddress, "bind-address", "0.0.0.0", "The IP address on which to serve the --port (set to 0.0.0.0 for all interfaces).")
	flags.StringVar(&app.args.RootPath, "root-path", "/", "The root path to serve.")
	flags.StringVar(&app.args.CertDir, "default-cert-dir", "/certs", "Directory path containing --tls-cert-file and --tls-key-file files. Relative to the container, not the host.")
	flags.StringVar(&app.args.TLSCertFile, "tls-cert-file", "tls.crt", "File containing the default x509 Certificate for HTTPS.")
	flags.StringVar(&app.args.TLSKeyFile, "tls-key-file", "tls.key", "File containing the default x509 private key matching --tls-cert-file.")
	flags.StringVar(&app.args.MTLSCACertFile, "mtls-cacert-file", "", "File containing the CA Certificate matching the --tls-key-file and --tls-cert-file. If provided the server switches to mTLS communication")
	flags.BoolVar(&app.args.WriteDefaultConfig, "write-default-config", false, "write a default config file")
	flags.StringVar(&app.args.ConfigFile, "config-file", "", "read parameters from a yaml config file")

	flags.Parse(args)
	flagutil.SetFlagsFromEnv(flags, "SERVER")

	if app.args.ConfigFile != "" {
		app.Logger.Printf("read configuration from %s", app.args.ConfigFile)

		yamlFile, err := os.ReadFile(app.args.ConfigFile)
		if err != nil {
			app.Logger.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &app.args)

		if err != nil {
			app.Logger.Fatalf("Unmarshal: %v", err)
		}
	}

	if app.args.WriteDefaultConfig {
		app.Logger.Println("writing config defaults")
		yamlData, err := yaml.Marshal(app.args)

		if err != nil {
			app.Logger.Printf("Error while Marshaling. %v", err)
		}

		fmt.Println(string(yamlData))
		os.Exit(0)
	}
}

func (app *Application) GetInsecureAddrAsString() string {
	app.Logger.Printf("GetInsecureAddrAsString: %s:%d", app.args.BindAddress, app.args.Port)
	return fmt.Sprintf("%s:%d", app.args.InsecureBindAddress, app.args.InsecurePort)
}

func (app *Application) GetAddrAsSring() string {
	app.Logger.Printf("GetAddrAsString: %s:%d", app.args.BindAddress, app.args.Port)
	return fmt.Sprintf("%s:%d", app.args.BindAddress, app.args.Port)
}

func (app *Application) AllowInsecureConnections() bool {
	app.Logger.Println("AllowInsecureConnection called")
	return app.args.InsecurePort > 0 && app.args.InsecurePort < 65536
}
