package application

import (
	"flag"
	"fmt"
	"os"

	"github.com/coreos/pkg/flagutil"
)

type cmdlineArgs struct {
	argInsecurePort        int
	argInsecureBindAddress string
	argPort                int
	argBindAddress         string
	argRootPath            string
	argCertDir             string
	argTLSCertFile         string
	argTLSKeyFile          string
	argMTLSCACertFile      string
	argEnableProfiling     bool
}

func (app *Application) ParseCommandlineAndEnvironment(args []string) {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.IntVar(&app.args.argInsecurePort, "insecure-port", 0, "port to listen to for HTTP requests. (set to \"0\" to disable insecure communication)")
	flags.StringVar(&app.args.argInsecureBindAddress, "insecure-bind-address", "127.0.0.1", "The IP address on which to serve the --insecure-port (set to 127.0.0.1 for loopback only).")
	flags.IntVar(&app.args.argPort, "port", 443, "The secure port to listen to for incoming HTTPS requests.")
	flags.StringVar(&app.args.argBindAddress, "bind-address", "0.0.0.0", "The IP address on which to serve the --port (set to 0.0.0.0 for all interfaces).")
	flags.StringVar(&app.args.argRootPath, "root-path", "/", "The root path to serve.")
	flags.StringVar(&app.args.argCertDir, "default-cert-dir", "/certs", "Directory path containing --tls-cert-file and --tls-key-file files. Relative to the container, not the host.")
	flags.StringVar(&app.args.argTLSCertFile, "tls-cert-file", "tls.crt", "File containing the default x509 Certificate for HTTPS.")
	flags.StringVar(&app.args.argTLSKeyFile, "tls-key-file", "tls.key", "File containing the default x509 private key matching --tls-cert-file.")
	flags.StringVar(&app.args.argMTLSCACertFile, "mtls-cacert-file", "", "File containing the CA Certificate matching the --tls-key-file and --tls-cert-file. If provided the server switches to mTLS communication")
	flags.BoolVar(&app.args.argEnableProfiling, "enable-profiling", false, "serve profiling on /debug endpoint.")

	flags.Parse(args)
	flagutil.SetFlagsFromEnv(flags, "SERVER")
}

func (app *Application) GetInsecureAddrAsString() string {
	app.Logger.Printf("GetInsecureAddrAsString: %s:%d", app.args.argBindAddress, app.args.argPort)
	return fmt.Sprintf("%s:%d", app.args.argInsecureBindAddress, app.args.argInsecurePort)
}

func (app *Application) GetAddrAsSring() string {
	app.Logger.Printf("GetAddrAsString: %s:%d", app.args.argBindAddress, app.args.argPort)
	return fmt.Sprintf("%s:%d", app.args.argBindAddress, app.args.argPort)
}

func (app *Application) AllowInsecureConnections() bool {
	app.Logger.Println("AllowInsecureConnection called")
	return app.args.argInsecurePort > 0 && app.args.argInsecurePort < 65536
}
