//
// Slack File Upload Gateway
// Allows to anonymously upload files to Slack channels specified in the url
// as@portworx.com v1.0
// 

package main

import (
    "fmt"
    "log"
    "io"
    "os"
    "time"
    "strconv"
    "strings"
    "net/http"
    "crypto/tls"
    "context"

    "github.com/nlopes/slack"
    "github.com/dchest/uniuri"
    "golang.org/x/crypto/acme/autocert"
)

var apitok string
var maxbytes int64
var timeout int
var datadir string
var hostn string

func msg(w http.ResponseWriter, r *http.Request, format string, args ...interface{}) {
    nformat := fmt.Sprintf("%s %s", r.RemoteAddr, format)
    fmt.Fprintf(w, nformat, args...)
    log.Printf(nformat, args...)
}

func slackit(w http.ResponseWriter, r *http.Request) {
    channel_name := strings.Replace(r.URL.Path, "/", "", -1)

    if r.ContentLength == 0 || len(channel_name) == 0 {
        msg(w, r, "ERROR: ContentLength: %d Path: %s\r\n\r\n", r.ContentLength, r.URL.Path)
        fmt.Fprintf(w, "Usage: curl -F \"file=@myfile.ext\" https://" + hostn + "/channel-name\r\n")
        return
    }

    r.Body = http.MaxBytesReader(w, r.Body, maxbytes)

    u, h, err := r.FormFile("file")
    if err != nil {
        msg(w, r, "ERROR: Unable to receive file upload, %s\n", err)
        return
    }

    upload_file_dir := datadir + "/slackit-" + uniuri.New()
    upload_file_path := upload_file_dir + "/" + h.Filename

    os.MkdirAll(upload_file_dir, 0755)
    defer os.RemoveAll(upload_file_dir)

    file, err := os.OpenFile(upload_file_path, os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        msg(w, r, "ERROR: Unable to create file %s, %s\n", h.Filename, err)
        return
    }
    io.Copy(file, u)

    u.Close()
    file.Close()

    msg(w, r, "Received file %s (%d) for channel %s\n", h.Filename, h.Size, channel_name)

    sc := slack.New(apitok)

    channel := make([]string, 1)

    conv_params := slack.GetConversationsParameters{
        Limit: 1000,
        Types: []string{"public_channel", "private_channel"},
    }

    conv, _, err := sc.GetConversations(&conv_params)
    if err != nil {
        msg(w, r, "ERROR, Unable to get Slack channels, %s\n", err)
        return
    }
    for _, c := range conv {
        if c.Name == channel_name {
            channel[0]=c.ID
        }
    }

    if len(channel[0]) != 9 {
        msg(w, r, "Error: Unable to find channel ID for %s, [%s]\r\n", channel_name, channel[0])
        return
    }

    params := slack.FileUploadParameters{
        File: upload_file_path,
        Channels:  channel,
    }

    s, err := sc.UploadFile(params)
    if err != nil {
        msg(w, r, "ERROR: Unable to upload file %s to Slack, %s\r\n", h.Filename, err)
        return
    }
    msg(w, r, "Uploaded file %s [%s] to channel %s [%s]\r\n", s.Name, s.ID, channel_name, channel[0])

}

func mkmux() *http.Server {
    mux := &http.ServeMux{}
    mux.HandleFunc("/", slackit)

    return &http.Server{
        ReadTimeout:  time.Duration(timeout) * time.Second,
        WriteTimeout: time.Duration(timeout) * time.Second,
        IdleTimeout:  time.Duration(timeout) * time.Second,
        Handler:      mux,
    }
}

func main() {
    var acm *autocert.Manager
    var ssl *http.Server
    var http *http.Server

    apitok = os.Getenv("APITOK")
    if len(apitok) == 0 {
        log.Fatal("APITOK env var is not defined")
    }

    hostn = os.Getenv("HOSTN")
    if len(hostn) == 0 {
        log.Fatal("HOSTN env var is not defined")
    }

    maxbytes, _ = strconv.ParseInt(os.Getenv("MAXBYTES"), 10, 64)
    if maxbytes == 0 {
        maxbytes = 1000 * 1024 * 1024
    }

    timeout, _ = strconv.Atoi(os.Getenv("TIMEOUT"))
    if timeout == 0 {
        timeout = 3600
    }

    datadir = os.TempDir()

    hostPolicy := func(ctx context.Context, host string) error {
        allowedHost := hostn
        if host == allowedHost {
            return nil
        }
        return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
    }

    acm = &autocert.Manager{
        Prompt:     autocert.AcceptTOS,
        HostPolicy: hostPolicy,
        Cache:      autocert.DirCache(datadir),
    }

    ssl = mkmux()
    ssl.Addr = ":8443"
    ssl.TLSConfig = &tls.Config{GetCertificate: acm.GetCertificate}

    go func() {
        log.Print("Starting HTTPS server on port", ssl.Addr)
        err := ssl.ListenAndServeTLS("", "")
        if err != nil {
            log.Fatalf("ListendAndServeTLS() failed with %s", err)
        }
    }()

    http = mkmux()

    if acm != nil {
        http.Handler = acm.HTTPHandler(http.Handler)
    }

    http.Addr = ":8080"

    log.Print("Starting HTTP server on port", http.Addr)
    log.Fatal(http.ListenAndServe())

}
