package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/3/20
 * @Desc: web.go
**/

var currentDate time.Time = time.Now()

type WebSrvEntry struct {
	http.Server
	webExitChan chan struct{}
}

func RegWebSrv() RegSrvInterface {
	RegWeb := new(RegSrv)
	webSrv := new(WebSrvEntry)
	RegWeb.RegStopSrv(webSrv.StopWeb)
	RegWeb.RegStartSrv(webSrv.StartWeb)
	return RegWeb
}

func (srv *WebSrvEntry) StopWeb(ctx context.Context, wg *sync.WaitGroup) {
	srv.webExitChan <- struct{}{}
	wg.Done()
}

func (srv *WebSrvEntry) stopWebProcess(ctx context.Context) {
	go func() {

		select {
		case <-ctx.Done():
			if err := srv.Shutdown(ctx); nil != err {
				log.Fatalf("server shutdown failed, err: %v\n", err)
			}
			log.Println("server gracefully shutdown")
			return
		case <-srv.webExitChan:
			if err := srv.Shutdown(ctx); nil != err {
				log.Fatalf("server shutdown failed, err: %v\n", err)
			}
			log.Println("server gracefully shutdown")
			return
		}

	}()
}

func (srv *WebSrvEntry) StartWeb(ctx context.Context) {
	// handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(currentDate.String())
		w.Write([]byte(fmt.Sprintf("Hello,%s", currentDate.String())))
	})

	// server
	//srv = http.Server{
	//	Addr:    ctx.Value("port").(string), // :8081
	//	Handler: handler,
	//}
	srv.Addr = ctx.Value("port").(string)
	srv.Handler = handler

	// make sure idle connections returned
	//processed := make(chan struct{})
	srv.webExitChan = make(chan struct{}, 0)
	go srv.stopWebProcess(ctx)

	// serve
	go srv.startWeb()

	// waiting for goroutine above processed
	//<-processed
}

func (srv *WebSrvEntry) startWeb() {
	err := srv.ListenAndServe()
	if http.ErrServerClosed != err {
		log.Fatalf("server not gracefully shutdown, err :%v\n", err)
	}
}
