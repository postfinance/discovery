package server

import (
	discoveryv1connect "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1/discoveryv1connect"
	goauth "gitlab.pnet.ch/linux/go/auth"
)

func rbacConfig() []goauth.Config {
	return []goauth.Config{
		{
			Role: "cop_appl_linux",
			Rules: []goauth.Rule{
				{
					Service: discoveryv1connect.NamespaceAPIName,
					Methods: []string{
						"RegisterNamespace",
						"UnregisterNamespace",
						"ListNamespace",
					},
				},
				{
					Service: discoveryv1connect.ServerAPIName,
					Methods: []string{
						"ListServer",
						"UnregisterServer",
						"RegisterServer",
					},
				},
				{
					Service: discoveryv1connect.ServiceAPIName,
					Methods: []string{
						"RegisterService",
						"UnRegisterService",
						"ListService",
						"ListTargetGroup",
					},
				},
				{
					Service: discoveryv1connect.TokenAPIName,
					Methods: []string{
						"Create",
						"Info",
					},
				},
			},
		},
		{
			Role: "machine",
			Rules: []goauth.Rule{
				{
					Service: discoveryv1connect.NamespaceAPIName,
					Methods: []string{
						"RegisterNamespace",
						"UnregisterNamespace",
						"ListNamespace",
					},
				},
				{
					Service: discoveryv1connect.ServerAPIName,
					Methods: []string{
						"ListServer",
						"RegisterServer",
					},
				},
				{
					Service: discoveryv1connect.ServiceAPIName,
					Methods: []string{
						"RegisterService",
						"UnRegisterService",
						"ListService",
						"ListTargetGroup",
					},
				},
			},
		},
	}
}
