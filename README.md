# A devoted friend with ipfs + peerjs

## How to run

```bash
pkg install git go kubo # termux
git version # git version 2.49.0
go version # go version go1.24.2 android/arm
ipfs version # ipfs version 0.34.1

git clone https://github.com/naratteu/ipfspeerjsfriend
cd ipfspeerjsfriend

chmod 777 ./loop.sh
nohup ./loop.sh &
nohup ipfs daemon &
```

## Usage

- [memo example](https://naratteu.github.io/ipfspeerjsfriend/usage/memo.html)
