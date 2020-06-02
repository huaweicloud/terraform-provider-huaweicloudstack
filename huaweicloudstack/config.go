package huaweicloudstack

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/huaweicloud/golangsdk"
	huaweisdk "github.com/huaweicloud/golangsdk/openstack"
)

type Config struct {
	AccessKey        string
	SecretKey        string
	CACertFile       string
	ClientCertFile   string
	ClientKeyFile    string
	DomainID         string
	DomainName       string
	EndpointType     string
	IdentityEndpoint string
	Insecure         bool
	Password         string
	Region           string
	TenantID         string
	TenantName       string
	Token            string
	Username         string
	UserID           string
	terraformVersion string
	endpoints        map[string]string

	HwClient *golangsdk.ProviderClient
}

func (c *Config) LoadAndValidate() error {
	validEndpoint := false
	validEndpoints := []string{
		"internal", "internalURL",
		"admin", "adminURL",
		"public", "publicURL",
		"",
	}

	for _, endpoint := range validEndpoints {
		if c.EndpointType == endpoint {
			validEndpoint = true
		}
	}

	if !validEndpoint {
		return fmt.Errorf("Invalid endpoint type provided")
	}

	if c.Password != "" {
		if c.Username == "" && c.UserID == "" {
			return fmt.Errorf("\"password\": one of `user_name, user_id` must be specified")
		}
		if c.TenantName != "" && c.DomainID == "" && c.DomainName == "" {
			return fmt.Errorf("\"tenant_name\": one of `domain_name, domain_id` must be specified")
		}
	}

	return newhwClient(c)

}

func newhwClient(c *Config) error {

	var ao golangsdk.AuthOptionsProvider

	ao = golangsdk.AuthOptions{
		DomainID:         c.DomainID,
		DomainName:       c.DomainName,
		IdentityEndpoint: c.IdentityEndpoint,
		Password:         c.Password,
		TenantID:         c.TenantID,
		TenantName:       c.TenantName,
		TokenID:          c.Token,
		Username:         c.Username,
		UserID:           c.UserID,
	}

	client, err := huaweisdk.NewClient(ao.GetIdentityEndpoint())
	if err != nil {
		return err
	}

	// Set UserAgent
	client.UserAgent.Prepend(httpclient.TerraformUserAgent(c.terraformVersion))

	config := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := pathorcontents.Read(c.CACertFile)
		if err != nil {
			return fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		config.RootCAs = caCertPool
	}

	if c.Insecure {
		config.InsecureSkipVerify = true
	}

	if c.ClientCertFile != "" && c.ClientKeyFile != "" {
		clientCert, _, err := pathorcontents.Read(c.ClientCertFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Cert: %s", err)
		}
		clientKey, _, err := pathorcontents.Read(c.ClientKeyFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Key: %s", err)
		}

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return err
		}

		config.Certificates = []tls.Certificate{cert}
		config.BuildNameToCertificate()
	}

	// if OS_DEBUG is set, log the requests and responses
	var osDebug bool
	if os.Getenv("OS_DEBUG") != "" {
		osDebug = true
	}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	client.HTTPClient = http.Client{
		Transport: &LogRoundTripper{
			Rt:      transport,
			OsDebug: osDebug,
		},
	}

	err = huaweisdk.Authenticate(client, ao)
	if err != nil {
		return err
	}

	c.HwClient = client
	return nil
}

func (c *Config) determineRegion(region string) string {
	// If a resource-level region was not specified, and a provider-level region was set,
	// use the provider-level region.
	if region == "" && c.Region != "" {
		region = c.Region
	}

	log.Printf("[DEBUG] HuaweiCloudStack Region is: %s", region)
	return region
}

func (c *Config) blockStorageV2Client(region string) (*golangsdk.ServiceClient, error) {
	if sc := c.getCustomerEndpoint("evs"); sc != nil {
		return sc, nil
	}

	return huaweisdk.NewBlockStorageV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) computeV2Client(region string) (*golangsdk.ServiceClient, error) {
	if sc := c.getCustomerEndpoint("ecs"); sc != nil {
		return sc, nil
	}

	return huaweisdk.NewComputeV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) kmsKeyV1Client(region string) (*golangsdk.ServiceClient, error) {
	if sc := c.getCustomerEndpoint("kms"); sc != nil {
		return sc, nil
	}

	// can not use NewKMSV1 as the catalog type "kms" was not registered on HCS
	// NewKmsKeyV1 will get the endpoint from "compute", then replace it with kms
	return huaweisdk.NewKmsKeyV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) identityV3Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewIdentityV3(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) imageV2Client(region string) (*golangsdk.ServiceClient, error) {
	if sc := c.getCustomerEndpoint("ims"); sc != nil {
		return sc, nil
	}

	return huaweisdk.NewImageServiceV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) networkingV2Client(region string) (*golangsdk.ServiceClient, error) {
	if sc := c.getCustomerEndpoint("vpc"); sc != nil {
		sc.ResourceBase = sc.Endpoint + "v2.0/"
		return sc, nil
	}

	return huaweisdk.NewNetworkV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) autoscalingV1Client(region string) (*golangsdk.ServiceClient, error) {
	if sc := c.getCustomerEndpoint("as"); sc != nil {
		return sc, nil
	}

	return huaweisdk.NewAutoScalingService(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) orchestrationV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewOrchestrationV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) getHwEndpointType() golangsdk.Availability {
	if c.EndpointType == "internal" || c.EndpointType == "internalURL" {
		return golangsdk.AvailabilityInternal
	}
	if c.EndpointType == "admin" || c.EndpointType == "adminURL" {
		return golangsdk.AvailabilityAdmin
	}
	return golangsdk.AvailabilityPublic
}

func (c *Config) getCustomerEndpoint(srvType string) *golangsdk.ServiceClient {
	endpoint, ok := c.endpoints[srvType]
	if !ok {
		return nil
	}

	return &golangsdk.ServiceClient{
		ProviderClient: c.HwClient,
		Endpoint:       endpoint,
	}
}
