package compose

import (
	"context"
	"fmt"

	"github.com/leandrodaf/harborctl/internal/config"
)

// traefikBuilder implementa TraefikBuilder
type traefikBuilder struct{}

func NewTraefikBuilder() TraefikBuilder {
	return &traefikBuilder{}
}

func (b *traefikBuilder) Build(ctx context.Context, stack *config.Stack) map[string]any {
	args := []string{
		"--providers.docker=true",
		"--providers.docker.exposedbydefault=false",
		"--entrypoints.websecure.address=:443",
		"--api=false",

		// Força HTTPS apenas - remove HTTP
		"--entrypoints.websecure.http.tls=true",

		// Configurações de Segurança
		"--global.checknewversion=false",
		"--global.sendanonymoususage=false",

		// Headers de Segurança
		"--entrypoints.websecure.http.middlewares=security-headers@docker",

		// Timeouts e Limites
		"--entrypoints.websecure.transport.respondingtimeouts.readtimeout=60s",
		"--entrypoints.websecure.transport.respondingtimeouts.writetimeout=60s",
		"--entrypoints.websecure.transport.respondingtimeouts.idletimeout=180s",

		// Rate Limiting Global
		"--entrypoints.websecure.http.middlewares=rate-limit@docker",

		// Tamanho máximo de request
		"--entrypoints.websecure.http.middlewares=request-size@docker",
	}

	if stack.TLS.Mode == "acme" {
		res := stack.TLS.Resolver
		args = append(args,
			fmt.Sprintf("--certificatesresolvers.%s.acme.email=%s", res, stack.TLS.Email),
			fmt.Sprintf("--certificatesresolvers.%s.acme.storage=/letsencrypt/acme.json", res),
			fmt.Sprintf("--certificatesresolvers.%s.acme.httpchallenge=true", res),
			fmt.Sprintf("--certificatesresolvers.%s.acme.httpchallenge.entrypoint=websecure", res),
		)

		if stack.TLS.DNS != nil && stack.TLS.DNS.Provider != "" {
			args = append(args,
				fmt.Sprintf("--certificatesresolvers.%s.acme.dnschallenge=true", res),
				fmt.Sprintf("--certificatesresolvers.%s.acme.dnschallenge.provider=%s", res, stack.TLS.DNS.Provider),
			)
		}
	}

	return map[string]any{
		"image":   "traefik:v3.5",
		"command": args,
		"ports":   []string{"443:443"}, // Apenas HTTPS
		"volumes": []string{
			"traefik_acme:/letsencrypt",
			"/var/run/docker.sock:/var/run/docker.sock:ro",
		},
		"labels": map[string]string{
			// Middlewares de Segurança Globais
			"traefik.http.middlewares.security-headers.headers.framedeny":                              "true",
			"traefik.http.middlewares.security-headers.headers.sslredirect":                            "true",
			"traefik.http.middlewares.security-headers.headers.stsincludesubdomains":                   "true",
			"traefik.http.middlewares.security-headers.headers.stspreload":                             "true",
			"traefik.http.middlewares.security-headers.headers.stsseconds":                             "63072000",
			"traefik.http.middlewares.security-headers.headers.contenttypenosniff":                     "true",
			"traefik.http.middlewares.security-headers.headers.browserxssfilter":                       "true",
			"traefik.http.middlewares.security-headers.headers.referrerpolicy":                         "strict-origin-when-cross-origin",
			"traefik.http.middlewares.security-headers.headers.permissionspolicy":                      "camera=(), microphone=(), payment=(), usb=()",
			"traefik.http.middlewares.security-headers.headers.customrequestheaders.X-Forwarded-Proto": "https",

			// Rate Limiting - 100 req/min por IP
			"traefik.http.middlewares.rate-limit.ratelimit.burst":                            "20",
			"traefik.http.middlewares.rate-limit.ratelimit.average":                          "100",
			"traefik.http.middlewares.rate-limit.ratelimit.period":                           "1m",
			"traefik.http.middlewares.rate-limit.ratelimit.sourcecriterion.ipstrategy.depth": "1",

			// Limite de tamanho de request - 10MB
			"traefik.http.middlewares.request-size.buffering.maxrequestbodybytes": "10485760",

			// Timeout personalizado
			"traefik.http.middlewares.timeout.circuitbreaker.expression": "ResponseCodeRatio(500, 600, 0, 600) > 0.25",

			// Remove headers sensíveis
			"traefik.http.middlewares.secure-headers.headers.customresponseheaders.Server":       "",
			"traefik.http.middlewares.secure-headers.headers.customresponseheaders.X-Powered-By": "",

			// Traefik não expõe a si mesmo
			"traefik.enable": "false",
		},
		"networks": []string{"public", "private"},
		"restart":  "always",

		// Configurações de segurança do container
		"security_opt": []string{
			"no-new-privileges:true",
		},
		"read_only": true,
		"tmpfs": []string{
			"/tmp:rw,noexec,nosuid,size=100m",
		},

		// Limites de recursos para evitar DoS
		"deploy": map[string]any{
			"resources": map[string]any{
				"limits": map[string]string{
					"cpus":   "1.0",
					"memory": "512M",
				},
				"reservations": map[string]string{
					"cpus":   "0.25",
					"memory": "128M",
				},
			},
		},
	}
}
