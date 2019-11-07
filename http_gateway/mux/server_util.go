package mux_http_gateway

import (
    "context"
    "crypto/tls"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/maxtech/gateway/http_gateway"
    "golang.org/x/crypto/acme/autocert"
    "net/http"
    "os"
    "os/signal"
    "time"
)

type serverUtil struct {
}

var (
    ServerUtil *serverUtil
    AllowHosts []string
)

func (su *serverUtil) InitMux() *mux.Router {
    return mux.NewRouter()
}

func (su *serverUtil) StartMuxByConfig(_route *mux.Router, _httpConfigFormat http_gateway.HttpConfigFormat, _hostPolicy func(ctx context.Context, host string) error) {
    server := &http.Server{
        // Good practice to set timeouts to avoid Slowloris attacks.
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: _route, // Pass our instance of gorilla/mux in.
    }

    if _httpConfigFormat.Https {

        dataDir := "."

        m := autocert.Manager{
            Prompt: autocert.AcceptTOS,
            HostPolicy: _hostPolicy,
            Cache: autocert.DirCache(dataDir),
        }
        server.Addr = _httpConfigFormat.HttpsAddress
        server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

        go startHttpsServer(server, _httpConfigFormat.CertFile, _httpConfigFormat.KeyFile)

        _, _ = fmt.Fprintln(os.Stdout, fmt.Sprintf("https server started: %v", _httpConfigFormat.HttpsAddress))
    } else {
        server.Addr = _httpConfigFormat.HttpAddress
        // Run our server in a goroutine so that it doesn't block.
        go startHttpServer(server)

        _, _ = fmt.Fprintln(os.Stdout, fmt.Sprintf("http server started: %v", _httpConfigFormat.HttpAddress))
    }

    c := make(chan os.Signal, 1)
    // We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
    // SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
    signal.Notify(c, os.Interrupt)

    // Block until we receive our signal.
    <-c

    // Create a deadline to wait for.
    ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
    defer cancel()
    // Doesn't block if no connections, but will otherwise wait
    // until the timeout deadline.
    _ = server.Shutdown(ctx)
    // Optionally, you could run srv.Shutdown in a goroutine and block on
    // <-ctx.Done() if your application should wait for other services
    // to finalize based on context cancellation.

    _, _ = fmt.Fprintln(os.Stdout, "shutting down")
    os.Exit(0)
}

func startHttpsServer(_server *http.Server, certFile, keyFile string) {
    if err := recover(); err != nil {
        _, _ = fmt.Fprintln(os.Stderr, err)
        startHttpsServer(_server, certFile, keyFile)
    }
    if err := _server.ListenAndServe(); err != nil {
        _, _ = fmt.Fprintln(os.Stderr, err.Error())
    }
}

func startHttpServer(_server *http.Server) {
    if err := recover(); err != nil {
        _, _ = fmt.Fprintln(os.Stderr, err)
        startHttpServer(_server)
    }
    if err := _server.ListenAndServe(); err != nil {
        _, _ = fmt.Fprintln(os.Stderr, err.Error())
    }
}

func (su *serverUtil) HostPolicy(_ctx context.Context, _host string) error {
    // Note: change to your real domain
    for _, allowedHost := range AllowHosts{
        if _host == allowedHost {
            return nil
        }
    }

    return fmt.Errorf("acme/autocert: only %v host is allowed", AllowHosts)
}
