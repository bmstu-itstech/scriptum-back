generate: generate/proto generate/openapi

.PHONY: generate/proto
generate/proto:
	easyp generate

.PHONY: generate/openapi
generate/openapi: api/v2/scriptum.swagger.yaml
	./scripts/generate-openapi-stubs.sh $<
