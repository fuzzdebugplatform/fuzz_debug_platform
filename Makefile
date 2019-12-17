

bin:
	go build -o bin/platform cmd/platform/*.go


docker: bin
	cd web && yarn run build && cd ..
	docker build -t ddfddf/fuzzdebugplatform:v0.1 .

.PHONY: bin docker