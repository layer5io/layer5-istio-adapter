gen:
	cd ../../..; BUILD_WITH_CONTAINER=1 make gen; cd mixer/adapter/layer5

docker:
	docker system prune -f
	docker build -t layer5/layer5-istio-adapter .

build:
	go build ./...
	cd cmd; go build ./...; cd ..

run:
	go run cmd/main.go 49494

mixs:
	cd ../../..; make mixs; cd mixer/adapter/layer5
	$(GOPATH)/out/darwin_amd64/release/mixs server --configStoreURL=fs://`pwd`/testdata

mixc:
	cd ../../..; make mixc; cd mixer/adapter/layer5
	$(GOPATH)/out/darwin_amd64/release/mixc report -s destination.service="svc.cluster.local" -i request.size=1235

