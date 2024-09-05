package storage

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	mserver "github.com/acorn-io/mink/pkg/server"
	"github.com/gptscript-ai/otto/pkg/storage/openapi/generated"
	"github.com/gptscript-ai/otto/pkg/storage/registry/apigroups/agent"
	"github.com/gptscript-ai/otto/pkg/storage/scheme"
	"github.com/gptscript-ai/otto/pkg/storage/services"
	"github.com/gptscript-ai/otto/pkg/version"
	k8sversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/server/healthz"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client client.WithWatch

func Start(ctx context.Context, config services.Config) (Client, error) {
	services, err := services.New(config)
	if err != nil {
		return nil, err
	}

	return startMinkServer(ctx, config, services)
}

func startMinkServer(ctx context.Context, config services.Config, services *services.Services) (Client, error) {
	apiGroups, err := mserver.BuildAPIGroups(services, agent.APIGroup)
	if err != nil {
		return nil, err
	}

	var l net.Listener
	if config.StorageListenPort == 0 {
		l, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, err
		}
	}

	minkConfig := &mserver.Config{
		Name:              "Storage Server",
		Version:           version.Get().String(),
		Authenticator:     services.Authn,
		Authorization:     services.Authz,
		HTTPSListenPort:   config.StorageListenPort,
		Listener:          l,
		OpenAPIConfig:     generated.GetOpenAPIDefinitions,
		Scheme:            scheme.Scheme,
		APIGroups:         apiGroups,
		ReadinessCheckers: []healthz.HealthChecker{services.DB},
	}

	//if cfg.AuditLogPolicyFile != "" && cfg.AuditLogPath != "" {
	//	minkConfig.AuditConfig = mserver.NewAuditOptions(cfg.AuditLogPolicyFile, cfg.AuditLogPath)
	//}

	minkConfig.Middleware = []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.URL.Path == "/version" {
					_ = json.NewEncoder(rw).Encode(k8sversion.Info{
						GitVersion: version.Get().String(),
						GitCommit:  version.Get().Commit,
					})
				} else {
					next.ServeHTTP(rw, req)
				}
			})
		},
	}

	minkServer, err := mserver.New(minkConfig)
	if err != nil {
		return nil, err
	}

	_ = minkServer.Handler(ctx)

	cfg := minkServer.Loopback
	return client.NewWithWatch(cfg, client.Options{
		Scheme: scheme.Scheme,
	})
}
