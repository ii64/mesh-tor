package main
// Author: Anysz
import (
    "fmt"
    "log"
    "os"
    "io"
    "flag"
    "net"
    "net/http"
    "net/url"
    "context"
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

    makeClient := func() (cl http.Client, resolv_ip string) {
        ips, err := net.LookupIP(proxyUrl.Hostname())
        if err != nil { panic(err) }
	proxyStruct := &url.URL{}
	if len(ips) > 0 {
          resolv_ip = fmt.Sprintf("%s", ips[0])
          proxyStruct = &url.URL{
            Scheme: proxyUrl.Scheme,
            Opaque : proxyUrl.Opaque,
            User: proxyUrl.User,
            Host: fmt.Sprintf("%s:%s", resolv_ip, proxyUrl.Port()),
            Path: proxyUrl.Path,
            RawPath: proxyUrl.RawPath,
            ForceQuery: proxyUrl.ForceQuery,
            RawQuery: proxyUrl.RawQuery,
            Fragment: proxyUrl.Fragment,
          }
        }else{
          resolv_ip = "no_scope"
          proxyStruct = proxyUrl
        }
        cl = http.Client{
            Transport: &http.Transport{
                Proxy: http.ProxyURL(proxyStruct),
                DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
                    conn, err := net.Dial(network, addr)
                    if err == nil {
                       log.Println("Addr", addr, "Ip:", conn.RemoteAddr())
                    }else{ log.Println(err) }
                       return conn, err
                    },
            },
        }
	return
    }

    _=destUrl

    http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        rr, err  := http.NewRequest(r.Method, *destPtr + r.RequestURI, r.Body)
        if err != nil { fmt.Fprintf(w, "Err %+v\n", err.Error()) }
        cl, proxy_node := makeClient()
        res, err := cl.Do(rr)
        if err != nil {
            fmt.Fprintf(w, "Err %#+v\n", err.Error())
        }else{
            w.WriteHeader(res.StatusCode)
            for k, v := range res.Header {
                w.Header().Set(k, v[0])
            }
            w.Header().Set("Node-IP", proxy_node)
            defer res.Body.Close()
            io.Copy(w, res.Body)
        }
    }))
    http.ListenAndServe(*addrPtr, nil)
}

