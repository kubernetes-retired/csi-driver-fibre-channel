REGISTRY_NAME=quay.io/k8scsi
IMAGE_NAME=fcplugin
IMAGE_VERSION=canary
IMAGE_TAG=$(REGISTRY_NAME)/$(IMAGE_NAME):$(IMAGE_VERSION)

#.PHONY: all clean fibrechannel fc-container

all: fibrechannel

fibrechannel:
	if [ ! -d ./vendor ]; then dep ensure -vendor-only; fi
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o _output/fibrechannel ./app
fc-container: fibrechannel
	docker build -t $(IMAGE_TAG) -f ./app/Dockerfile .
push: fc-container
	docker push $(IMAGE_TAG)
clean:
	go clean -r -x
	-rm -rf _output
