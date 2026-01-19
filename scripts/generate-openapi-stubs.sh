#!/bin/bash
set -e

readonly specification="$1"

oapi-codegen -generate types      -o "internal/api/v2/openapi_types.gen.go" -package apiv2 "$specification"
oapi-codegen -generate chi-server -o "internal/api/v2/openapi_api.gen.go"   -package apiv2 "$specification"

mkdir -p "pkg/clients"
oapi-codegen -generate types  -o "pkg/clients/api/v2/openapi_types.gen.go"  -package apiv2 "$specification"
oapi-codegen -generate client -o "pkg/clients/api/v2/openapi_client_gen.go" -package apiv2 "$specification"
