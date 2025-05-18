.PHONY: run build exec

run: build exec

build:
	cd cmd && go build -o ipfspeerjsfriend .

exec:
	cd cmd && ./ipfspeerjsfriend