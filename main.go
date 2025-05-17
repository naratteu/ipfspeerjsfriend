package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	exts "ipfspeerjsfriend/peerjs_extensions"
	bp "ipfspeerjsfriend/peerjs_js_binarypack"

	"github.com/bitfield/script"
	peer "github.com/muka/peerjs-go"
)

func main() {
	for {
		func() {
			newPeer := func(id string, opts *peer.Options) (*peer.Peer, *exts.NamedEvents) {
				log.Printf("%s 로 시도해볼게용\n", id)
				p, err := peer.NewPeer(id, *opts)
				if err != nil {
					panic(err)
				}
				ok, end := taskCompletionSource[bool](2 * time.Second)
				nevs := &exts.NamedEvents{}
				nevs.Open = func(id string) { end(true) }
				nevs.Error = func(err peer.PeerError) { end(false) }
				defer func() {
					nevs.Open = nil
					nevs.Error = nil
				}()
				nevs.Join(p, evLog)
				if <-ok {
					return p, nevs
				} else {
					p.Close()
					return nil, nil
				}
			}
			opts := peer.NewOptions()
			p, nevs := newPeer("ipfspeerjsfriend", &opts)
			for i := 0; p == nil; i++ {
				p, nevs = newPeer(fmt.Sprintf("ipfspeerjsfriend%d", i), &opts)
			}
			defer p.Close()
			defer log.Println("종료댐.", p.ID)
			log.Printf("%s 로 됫어용\n", p.ID)
			reset, end := taskCompletionSource[string](time.Hour) // 1시간마다 아무일 없어도 리셋
			nevs.Connection = connection
			nevs.Disconnected = end
			<-reset
		}()
		log.Println("1분 뒤에 재시작 할게요..")
		time.Sleep(time.Minute) // 네트워크나 몬가 문제가 있었겠거니 하고 1분후에 다시시작
	}
}

func taskCompletionSource[T any](timeout time.Duration) (<-chan T, func(T)) {
	ch, once := make(chan T), new(sync.Once)
	go func() { time.Sleep(timeout); once.Do(func() { close(ch) }) }()
	return ch, func(t T) { once.Do(func() { ch <- t; close(ch) }) }
}

func evLog(ev string, arg any) {
	log.Printf("\033[0;90m[peerjs] on %s[%T]:%s\033[0m\n", ev, arg, arg)
}

func connection(conn *peer.DataConnection) {
	id := conn.GetPeerID()
	log.Println("누군가접속함!:", id)
	conn.On("data", func(data any) {
		bin := data.([]byte)
		log.Println("Received", id, len(bin))

		str := bp.UnpackStr(bin)

		cmd_add := "ipfs add -Q #"
		cmd_ps := "ps -Ao cmd #ipfs"

		res, err := "unknown command", error(nil)
		switch {
		case strings.HasPrefix(str, cmd_add):
			document := strings.TrimPrefix(str, cmd_add)
			res, err = script.Echo(document).Exec(cmd_add).String()
			if err != nil {
				res = err.Error()
			}
		case strings.HasPrefix(str, cmd_ps):
			res, err = script.Exec(cmd_ps).Match("ipfs").String()
			if err != nil {
				res = err.Error()
			}
		}
		conn.Send(bp.PackStr(res), false)
	})
}
