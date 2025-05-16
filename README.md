# A devoted friend with ipfs + peerjs

https://codepen.io/naratteu/pen/wvLjOMP 여기서 콘솔로 아래명령으로 접속 테스트

```js
p = new Peer(); await new Promise(r => setTimeout(r, 1000));
c = p.connect("ipfspeerjsfriend"); await new Promise(r => setTimeout(r, 2000));
c.on("data", (data)=>{console.log("https://ipfs.io/ipfs/"+data);})
c.send("ipfs add -Q #똠얌꿍");

# res: https://ipfs.io/ipfs/QmNMtgcdKRwVa8dqBTo2e3oixvbViWJSACc5KogQ78JnuD
```

```js
p = new Peer(); await new Promise(r => setTimeout(r, 1000));
c = p.connect("ipfspeerjsfriend"); await new Promise(r => setTimeout(r, 2000));
c.on("data", (data)=>{console.log(data);})
c.send("ps -o cmd | grep ipfs #");
undefined

# res: ipfs daemon
```
