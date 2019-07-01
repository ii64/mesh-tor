package main
// Author: Anysz
import (
    "fmt"
    "os"
    "io"
    "flag"
    "net/http"
    "net/url"
    _ "net/http/httputil"
)
func main() {
    proxyPtr := flag.String("proxy", os.Getenv("PROXY_SERVICE"), "Proxy for host")
    destPtr  := flag.String("dst", os.Getenv("DST_SERVICE"), "Destinaton url")

    addrPtr  := flag.String("addr", os.Getenv("SERVICE_ADDR"), "Server bind address")

    flag.Parse()

    proxyUrl, err := url.Parse(*proxyPtr)
    if err != nil { panic(err) }
    destUrl, err := url.Parse(*destPtr)
    if err != nil { panic(err) }
    http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

    //fmt.Printf("%#+v\n" , *destUrl)

    _ = destUrl

    http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        rr, _ := http.NewRequest(r.Method, *destPtr + r.RequestURI, r.Body)
        res, err := http.DefaultClient.Do(rr)
        if err != nil {
            fmt.Fprintf(w, "Err %#+v\n", err.Error())
        }else{
            w.WriteHeader(res.StatusCode)
            for k, v := range res.Header {
                w.Header().Set(k, v[0])
            }
            defer res.Body.Close()
            io.Copy(w, res.Body)
        }
    }))
    http.ListenAndServe(*addrPtr, nil)
}

