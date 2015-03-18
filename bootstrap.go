package main

import (
	"crypto/tls"
	"flag"
	"github.com/mikechack/bootstrap/certs"
	"github.com/mikechack/bootstrap/controllers"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type TcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln TcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func httpTlsListener(port string, c chan bool) {
	log.Printf("About to listen securely on port %s\n", port)
	/*
		err := http.ListenAndServeTLS(":443", "certs/server.pem", "certs/server.key", nil)
		if err != nil {
			log.Fatal(err)
		}
	*/
	srv := &http.Server{Addr: ":" + port, Handler: nil}
	addr := srv.Addr

	//dummy comment

	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	//config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)

	clientCert, _ := certs.Asset("certs/client0.crt")
	clientKey, _ := certs.Asset("certs/client0.key")
	config.Certificates[0], err = tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		//return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		//return err
	}

	tlsListener := tls.NewListener(TcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	srv.Serve(tlsListener)

}

func httpListener(port string, c chan bool) {
	log.Printf("About to listen insecurely on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}

func main() {

	port := flag.String("http-port", "80", "http port number")
	sport := flag.String("https-port", "443", "https port number")

	flag.Parse()

	rand.Seed(time.Now().Unix())

	http.Handle("/api/v1/status", http.HandlerFunc(controllers.Status_controller))
	http.Handle("/api/v1/status/", http.HandlerFunc(controllers.Status_controller))
	http.Handle("/api/v1/redirectUri", http.HandlerFunc(controllers.RedirectUri_controller))
	http.Handle("/api/v1/getMachineAccount", http.HandlerFunc(controllers.GetMachineAccount_controller))
	http.Handle("/api/v1/logonRegister", http.HandlerFunc(controllers.LogonRegister_controller))

	// test methods
	http.Handle("/api/v1/logmeon", http.HandlerFunc(controllers.LogMeOn))
	http.Handle("/api/v1/token", http.HandlerFunc(controllers.Token))
	http.Handle("/api/v1/cachedAccount", http.HandlerFunc(controllers.UseCachedAccount))

	chan1 := make(chan bool)
	go httpListener(*port, chan1)
	go httpTlsListener(*sport, chan1)

	select {
	case <-chan1:
		log.Printf("\n\nDetected context done\n\n")
	}
	log.Printf("goodbye for now - cya later\n")
}
