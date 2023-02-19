package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "net"
    "fmt"
    "time"
)

type URL struct {
    IP string
    Addr string
    Live bool
}

type RR struct {
    Index int
    Addrs []string
}

func (rr *RR) Add(url string) {
    if url == "" {
        fmt.Println("Add error: empty url")
        return
    }
    rr.Addrs = append(rr.Addrs, url)
}

func (rr *RR) Del(url string) {
    if url == "" {
        fmt.Println("Add error: empty url")
        return
    }
    for i := 0; i < len(rr.Addrs); i++ {
        if rr.Addrs[i] == url {
            rr.Addrs = append(rr.Addrs[:i], rr.Addrs[i+1:]...)
            return
        }
    }
}

func (rr *RR) RoundRobin() string {
    if rr.Index >= len(rr.Addrs) {
        rr.Index = 0
    }
    res := rr.Addrs[rr.Index]
    rr.Index++
    return res
}

func (url *URL) LivenessCheck(rr *RR)  {
    for {
        conn, err := net.Dial("tcp", url.IP)
        if err != nil {
            fmt.Println(url.IP, "已断开连接")
            rr.Del(url.Addr)
            url.Live = false
            continue
        } else {
            if url.Live == false {
                url.Live = true
                rr.Add(url.Addr)
            }
        }
        time.Sleep(time.Second * 3)
        defer conn.Close()
    }
}

func main() {
    var rr = &RR{}
    u1 := URL{IP : "192.168.108.164:8080", Addr : "http://192.168.108.164:8080/home", Live : true,}
    u2 := URL{IP : "192.168.108.167:8080", Addr : "http://192.168.108.167:8080/home", Live : true,}
    u3 := URL{IP : "192.168.108.168:8080", Addr : "http://192.168.108.168:8080/home", Live : true,}
    rr.Add(u1.Addr)
    rr.Add(u2.Addr)
    rr.Add(u3.Addr)
    r := gin.Default()
    go u1.LivenessCheck(rr)
    go u2.LivenessCheck(rr)
    go u3.LivenessCheck(rr)
    r.GET("/home", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, rr.RoundRobin())
    })
    r.Run(":8081")
}