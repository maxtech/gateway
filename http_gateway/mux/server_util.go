package mux_http_gateway

import (
    "context"
    "crypto/tls"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/maxtech/gateway/http_gateway"
    "golang.org/x/crypto/acme/autocert"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
)

type serverUtil struct {
}

var ServerUtil *serverUtil

func (s *serverUtil) InitMux(mode string) *mux.Router {
    return mux.NewRouter()
}

func (s *serverUtil) StartMuxByConfig(route *mux.Router, https bool, certFile, keyFile string, httpConfigFormat http_gateway.HttpConfigFormat) {
    server := &http.Server{
        // Good practice to set timeouts to avoid Slowloris attacks.
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: route, // Pass our instance of gorilla/mux in.
    }

    if https {

        hostPolicy := func(ctx context.Context, host string) error {
            // Note: change to your real domain
            allowedHost := "www.mydomain.com"
            if host == allowedHost {
                return nil
            }
            return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
        }
        dataDir := "."

        m := autocert.Manager{
            Prompt: autocert.AcceptTOS,
            HostPolicy: hostPolicy,
            Cache: autocert.DirCache(dataDir),
        }
        server.Addr = httpConfigFormat.HttpsAddress
        server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

        go startHttpsServer(server, certFile, keyFile)
    } else {
        server.Addr = httpConfigFormat.HttpAddress
        // Run our server in a goroutine so that it doesn't block.
        go startHttpServer(server)
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
    log.Println("shutting down")
    os.Exit(0)
}

func startHttpsServer(server *http.Server, certFile, keyFile string) {
    if err := recover(); err != nil {
        log.Println(err)
        startHttpsServer(server, certFile, keyFile)
    }
    if err := server.ListenAndServe(); err != nil {
        log.Println(err)
    }
}

func startHttpServer(server *http.Server) {
    if err := recover(); err != nil {
        log.Println(err)
        startHttpServer(server)
    }
    if err := server.ListenAndServe(); err != nil {
        log.Println(err)
    }
}
