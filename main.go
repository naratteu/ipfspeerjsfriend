package main

import (
	"log"
	"strings"

	bp "ipfspeerjsfriend/peerjs_js_binarypack"

	"github.com/bitfield/script"
	peer "github.com/muka/peerjs-go"
)

func main() {
	p, err := peer.NewPeer("ipfspeerjsfriend", peer.NewOptions())
	if err != nil {
		panic(err)
	}
	defer p.Close()
	p.On("open", func(data any) {
		log.Println("Peer Opened. ID:", p.ID)
	})
	p.On("connection", func(connection any) {
		conn := connection.(*peer.DataConnection)
		id := conn.GetPeerID()
		log.Println("누군가접속함!:", id)
		conn.On("data", func(data any) {
			bin := data.([]byte)
			log.Println("Received", id, len(bin))

			str := bp.UnpackStr(bin)

			cmd_add := "ipfs add -Q #"
			cmd_ps := "ps -o cmd | grep ipfs #"

			res := "unknown command"
			switch {
			case strings.HasPrefix(str, cmd_add):
				document := strings.TrimPrefix(str, cmd_add)
				res, err = script.Echo(document).Exec(cmd_add).String()
				if err != nil {
					res = err.Error()
				}
			case strings.HasPrefix(str, cmd_ps):
				res, err = script.Exec("ps -ao cmd").Match("ipfs").String()
				if err != nil {
					res = err.Error()
				}
			}
			conn.Send(bp.PackStr(res), false)
		})
	})
	select {}
}
