package app

import "github.com/bibi-ic/mata/internal/controller"

type Controller struct {
	metaController controller.MetaController
}

func NewController(metaController controller.MetaController) *Controller {
	return &Controller{
		metaController: metaController,
	}
}
