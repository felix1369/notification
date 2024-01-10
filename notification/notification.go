package notification

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Notification struct {
	Services map[string]NotificationService
}

func (n *Notification) New() {
	n.Services = make(map[string]NotificationService)
}

func (n *Notification) Send(ctx context.Context, request ...interface{}) error {
	var eg errgroup.Group
	for _, service := range n.Services {
		if service == nil {
			continue
		}
		service := service
		eg.Go(func() error {
			return service.Send(ctx, request)
		})
	}
	err := eg.Wait()
	if err != nil {
		err = errors.Wrap(fmt.Errorf(""), err.Error())
	}

	return err
}

func (n *Notification) AddService(service NotificationService) {
	s, _ := n.Services[service.GetName()]
	if s == nil {
		n.Services[service.GetName()] = service
	}
}

func (n *Notification) AddServices(services []NotificationService) {
	for _, service := range services {
		s, _ := n.Services[service.GetName()]
		if s == nil {
			n.Services[service.GetName()] = service
		}
	}
}
