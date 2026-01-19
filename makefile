SWAGGER_UI_VERSION:=v5.31.0

generate: generate/openapi

.PHONY: generate/openapi
generate/openapi: api/v2/scriptum.swagger.yaml
	./scripts/generate-openapi-stubs.sh $<

generate/swagger_ui: api/v2/scriptum.swagger.yaml
	 SWAGGER_UI_VERSION=$(SWAGGER_UI_VERSION) ./scripts/generate-swagger-ui.sh $<
