images += alpine

default:
	@echo Possible targets: $(images)

$(images):
	$(eval tag := $(shell ./version.sh $@))
	$(eval image := prince1809/$@)
	docker build -t $(image):$(tag) $@
	docker tag $(image):$(tag) $(image):latest
	docker push $(image):$(tag)
	docker push $(image):latest

.PHONY: $(images)
