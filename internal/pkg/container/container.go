package container

import (
	"log"

	"github.com/facebookgo/inject"
)

// Service is the service interface
type Service interface {
	Startup() error
	Shutdown() error
}

// Container is the service container interface
type Container interface {
	Ready() error
	GetService(id string) (interface{}, bool)
	RegisterService(id string, svc interface{})
	RegisterServices(services map[string]interface{})
	Shutdown()
}

type container struct {
	graph    inject.Graph
	services map[string]interface{}
	order    []string
	ready    bool
}

// NewContainer creates a new service container
func New() Container {
	return &container{services: make(map[string]interface{}), order: make([]string, 0), ready: false}
}

// GetService fetches a service by its ID
func (c *container) GetService(id string) (interface{}, bool) {
	svc, ok := c.services[id]
	return svc, ok
}

// RegisterService registers a service
func (c *container) RegisterService(id string, svc interface{}) {
	err := c.graph.Provide(&inject.Object{Name: id, Value: svc, Complete: false})
	if err != nil {
		panic(err)
	}
	c.order = append(c.order, id)
	c.services[id] = svc
}

// RegisterServices registers multiple services
func (c *container) RegisterServices(services map[string]interface{}) {
	for id, svc := range services {
		c.RegisterService(id, svc)
	}
}

// Ready starts up the service graph and returns error if it's not ready
func (c *container) Ready() error {
	if c.ready {
		return nil
	}
	if err := c.graph.Populate(); err != nil {
		return err
	}
	for _, key := range c.order {
		obj := c.services[key]
		if s, ok := obj.(Service); ok {
			log.Println("[starting up]", key)
			if err := s.Startup(); err != nil {
				return err
			}
		}
	}
	c.ready = true
	return nil
}

// Shutdown shuts down all services
func (c *container) Shutdown() {
	for _, key := range c.order {
		if service, ok := c.services[key]; ok {
			if s, ok := service.(Service); ok {
				log.Println("[shutting down]", key)
				if err := s.Shutdown(); err != nil {
					log.Println("ERROR: [shutting down]", key, err.Error())
				}
			}
		}
	}
	c.ready = false
}
