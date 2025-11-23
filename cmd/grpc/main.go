package main

import (
	"github.com/bmstu-itstech/scriptum-back/internal/infra/docker"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs"
)

func main() {
	l := logs.DefaultLogger()
	_ = docker.MustNewRunner(l)
}
