package huaweicloudstack

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
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
	s3sess   *session.Session
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

	if c.AccessKey != "" && c.SecretKey != "" {
		// Setup S3 client/config information for Swift S3 buckets
		log.Println("[INFO] Building Swift S3 auth structure")
		creds, err := GetCredentials(c)
		if err != nil {
			return err
		}
		// Call Get to check for credential provider. If nothing found, we'll get an
		// error, and we can present it nicely to the user
		cp, err := creds.Get()
		if err != nil {
			if sErr, ok := err.(awserr.Error); ok && sErr.Code() == "NoCredentialProviders" {
				return fmt.Errorf("No valid credential sources found for S3 Provider.")
			}

			return fmt.Errorf("Error loading credentials for S3 Provider: %s", err)
		}

		log.Printf("[INFO] S3 Auth provider used: %q", cp.ProviderName)

		sConfig := &aws.Config{
			Credentials: creds,
			Region:      aws.String(c.Region),
			HTTPClient:  cleanhttp.DefaultClient(),
		}

		if osDebug {
			sConfig.LogLevel = aws.LogLevel(aws.LogDebugWithHTTPBody | aws.LogDebugWithRequestRetries | aws.LogDebugWithRequestErrors)
			sConfig.Logger = awsLogger{}
		}

		if c.Insecure {
			transport := sConfig.HTTPClient.Transport.(*http.Transport)
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		// Set up base session for S3
		c.s3sess, err = session.NewSession(sConfig)
		if err != nil {
			return errwrap.Wrapf("Error creating Swift S3 session: {{err}}", err)
		}
	}

	return nil
}

type awsLogger struct{}

func (l awsLogger) Log(args ...interface{}) {
	tokens := make([]string, 0, len(args))
	for _, arg := range args {
		if token, ok := arg.(string); ok {
			tokens = append(tokens, token)
		}
	}
	log.Printf("[DEBUG] [aws-sdk-go] %s", strings.Join(tokens, " "))
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

func (c *Config) computeS3conn(region string) (*s3.S3, error) {
	if c.s3sess == nil {
		return nil, fmt.Errorf("Missing credentials for Swift S3 Provider, need access_key and secret_key values for provider.")
	}

	client, err := huaweisdk.NewImageServiceV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
	endpoint := strings.Replace(client.Endpoint, "//ims", "//oss", 1)

	awsS3Sess := c.s3sess.Copy(&aws.Config{Endpoint: aws.String(endpoint)})
	s3conn := s3.New(awsS3Sess)

	return s3conn, err
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
