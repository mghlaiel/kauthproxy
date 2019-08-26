// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"github.com/int128/kauthproxy/pkg/cmd"
	"github.com/int128/kauthproxy/pkg/logger"
	"github.com/int128/kauthproxy/pkg/network"
	"github.com/int128/kauthproxy/pkg/portforwarder"
	"github.com/int128/kauthproxy/pkg/resolver"
	"github.com/int128/kauthproxy/pkg/reverseproxy"
	"github.com/int128/kauthproxy/pkg/usecases"
)

// Injectors from di.go:

func NewCmd() cmd.Interface {
	loggerLogger := &logger.Logger{}
	reverseProxy := &reverseproxy.ReverseProxy{
		Logger: loggerLogger,
	}
	portForwarder := &portforwarder.PortForwarder{
		Logger: loggerLogger,
	}
	factory := &resolver.Factory{
		Logger: loggerLogger,
	}
	networkNetwork := &network.Network{}
	authProxy := &usecases.AuthProxy{
		ReverseProxy:    reverseProxy,
		PortForwarder:   portForwarder,
		ResolverFactory: factory,
		Network:         networkNetwork,
		Logger:          loggerLogger,
	}
	cmdCmd := &cmd.Cmd{
		AuthProxy: authProxy,
		Logger:    loggerLogger,
	}
	return cmdCmd
}