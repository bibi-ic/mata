package app

import "github.com/bibi-ic/mata/internal/service"

type ServiceServer struct {
	metaService service.MetaService
}

func NewService(metaService service.MetaService) *ServiceServer {
	return &ServiceServer{
		metaService: metaService,
	}
}
