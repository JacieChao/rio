package populate

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/rancher/norman/pkg/objectset"
	v1alpha3client "github.com/rancher/rio/types/apis/networking.istio.io/v1alpha3"
	v1 "github.com/rancher/rio/types/apis/rio.cattle.io/v1"
	"istio.io/api/networking/v1alpha3"
)

func ServiceEntry(svc *v1.ExternalService, stack *v1.Stack, os *objectset.ObjectSet) error {
	if svc.Spec.FQDN != "" {
		se, err := populateServiceEntryForFqdn(svc.Spec.FQDN, svc)
		if err != nil {
			return err
		}
		os.Add(se)
	}
	return nil
}

func populateServiceEntryForFqdn(fqdn string, svc *v1.ExternalService) (*v1alpha3client.ServiceEntry, error) {
	u, err := ParseTargetUrl(fqdn)
	if err != nil {
		return nil, err
	}

	scheme := u.Scheme
	if scheme == "" {
		scheme = "http"
	}

	port, _ := strconv.ParseUint(u.Port(), 10, 64)
	if port == 0 {
		if scheme == "http" {
			port = 80
		} else if scheme == "https" {
			port = 443
		}
	}

	se := v1alpha3client.NewServiceEntry(svc.Namespace, svc.Name, v1alpha3client.ServiceEntry{
		Spec: v1alpha3client.ServiceEntrySpec{
			Hosts:      []string{u.Host},
			Location:   v1alpha3.ServiceEntry_MESH_EXTERNAL,
			Resolution: v1alpha3.ServiceEntry_DNS,
			Ports: []v1alpha3client.Port{
				{
					Protocol: strings.ToUpper(scheme),
					Number:   uint32(port),
					Name:     fmt.Sprintf("%s-%v", scheme, port),
				},
			},
			Endpoints: []v1alpha3client.ServiceEntry_Endpoint{
				{
					Address: u.Host,
					Ports: map[string]uint32{
						scheme: uint32(port),
					},
				},
			},
		},
	})
	return se, nil
}

func ParseTargetUrl(target string) (*url.URL, error) {
	if !strings.HasPrefix(target, "https://") && !strings.HasPrefix(target, "http://") {
		target = "http://" + target
	}
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return u, nil
}
