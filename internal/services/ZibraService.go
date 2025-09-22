package services

import (
	"github.com/Talos-hub/ZibraGo/internal/ports"
)

type ZibraService struct {
	walker  ports.Wakler
	arhiver ports.Archiver
	api     ports.ApiCloud
}

func NewZibra(walker ports.Wakler, arhiver ports.Archiver, api ports.ApiCloud) *ZibraService {
	return &ZibraService{
		walker:  walker,
		arhiver: arhiver,
		api:     api,
	}
}
