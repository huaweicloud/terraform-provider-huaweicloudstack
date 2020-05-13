package huaweicloudstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

// Provider returns a schema.Provider for HuaweiCloudStack.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ACCESS_KEY", ""),
				Description: descriptions["access_key"],
			},

			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SECRET_KEY", ""),
				Description: descriptions["secret_key"],
			},

			"auth_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", nil),
				Description: descriptions["auth_url"],
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["region"],
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", ""),
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: descriptions["user_name"],
			},

			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_ID", ""),
				Description: descriptions["user_name"],
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: descriptions["tenant_id"],
			},

			"tenant_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_NAME",
					"OS_PROJECT_NAME",
				}, ""),
				Description: descriptions["tenant_name"],
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: descriptions["password"],
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_TOKEN", ""),
				Description: descriptions["token"],
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_ID",
					"OS_PROJECT_DOMAIN_ID",
					"OS_DOMAIN_ID",
				}, ""),
				Description: descriptions["domain_id"],
			},

			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_NAME",
					"OS_PROJECT_DOMAIN_NAME",
					"OS_DOMAIN_NAME",
					"OS_DEFAULT_DOMAIN",
				}, ""),
				Description: descriptions["domain_name"],
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_INSECURE", false),
				Description: descriptions["insecure"],
			},

			"endpoint_type": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ENDPOINT_TYPE", ""),
			},

			"endpoints": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: descriptions["cacert_file"],
			},

			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"huaweicloudstack_images_image_v2":        dataSourceImagesImageV2(),
			"huaweicloudstack_kms_key_v1":             dataSourceKmsKeyV1(),
			"huaweicloudstack_kms_data_key_v1":        dataSourceKmsDataKeyV1(),
			"huaweicloudstack_networking_network_v2":  dataSourceNetworkingNetworkV2(),
			"huaweicloudstack_networking_port_v2":     dataSourceNetworkingPortV2(),
			"huaweicloudstack_networking_secgroup_v2": dataSourceNetworkingSecGroupV2(),
			"huaweicloudstack_networking_subnet_v2":   dataSourceNetworkingSubnetV2(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"huaweicloudstack_as_group_v1":                        resourceASGroup(),
			"huaweicloudstack_as_configuration_v1":                resourceASConfiguration(),
			"huaweicloudstack_as_policy_v1":                       resourceASPolicy(),
			"huaweicloudstack_blockstorage_volume_v2":             resourceBlockStorageVolumeV2(),
			"huaweicloudstack_compute_instance_v2":                resourceComputeInstanceV2(),
			"huaweicloudstack_compute_interface_attach_v2":        resourceComputeInterfaceAttachV2(),
			"huaweicloudstack_compute_keypair_v2":                 resourceComputeKeypairV2(),
			"huaweicloudstack_compute_servergroup_v2":             resourceComputeServerGroupV2(),
			"huaweicloudstack_compute_floatingip_associate_v2":    resourceComputeFloatingIPAssociateV2(),
			"huaweicloudstack_compute_volume_attach_v2":           resourceComputeVolumeAttachV2(),
			"huaweicloudstack_kms_key_v1":                         resourceKmsKeyV1(),
			"huaweicloudstack_lb_certificate_v2":                  resourceCertificateV2(),
			"huaweicloudstack_lb_loadbalancer_v2":                 resourceLoadBalancerV2(),
			"huaweicloudstack_lb_listener_v2":                     resourceListenerV2(),
			"huaweicloudstack_lb_pool_v2":                         resourcePoolV2(),
			"huaweicloudstack_lb_member_v2":                       resourceMemberV2(),
			"huaweicloudstack_lb_monitor_v2":                      resourceMonitorV2(),
			"huaweicloudstack_lb_l7policy_v2":                     resourceL7PolicyV2(),
			"huaweicloudstack_lb_l7rule_v2":                       resourceL7RuleV2(),
			"huaweicloudstack_lb_whitelist_v2":                    resourceWhitelistV2(),
			"huaweicloudstack_networking_network_v2":              resourceNetworkingNetworkV2(),
			"huaweicloudstack_networking_subnet_v2":               resourceNetworkingSubnetV2(),
			"huaweicloudstack_networking_floatingip_v2":           resourceNetworkingFloatingIPV2(),
			"huaweicloudstack_networking_floatingip_associate_v2": resourceNetworkingFloatingIPAssociateV2(),
			"huaweicloudstack_networking_port_v2":                 resourceNetworkingPortV2(),
			"huaweicloudstack_networking_router_v2":               resourceNetworkingRouterV2(),
			"huaweicloudstack_networking_router_interface_v2":     resourceNetworkingRouterInterfaceV2(),
			"huaweicloudstack_networking_router_route_v2":         resourceNetworkingRouterRouteV2(),
			"huaweicloudstack_networking_secgroup_v2":             resourceNetworkingSecGroupV2(),
			"huaweicloudstack_networking_secgroup_rule_v2":        resourceNetworkingSecGroupRuleV2(),
			"huaweicloudstack_networking_vip_v2":                  resourceNetworkingVIPV2(),
			"huaweicloudstack_networking_vip_associate_v2":        resourceNetworkingVIPAssociateV2(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return configureProvider(d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "The access key for API operations. You can retrieve this\n" +
			"from the 'My Credential' section of the console.",

		"secret_key": "The secret key for API operations. You can retrieve this\n" +
			"from the 'My Credential' section of the console.",

		"auth_url": "The Identity authentication URL.",

		"region": "The HuaweiCloudStack region to connect to.",

		"user_name": "Username to login with.",

		"user_id": "User ID to login with.",

		"tenant_id": "The ID of the Tenant (Identity v2) or Project (Identity v3)\n" +
			"to login with.",

		"tenant_name": "The name of the Tenant (Identity v2) or Project (Identity v3)\n" +
			"to login with.",

		"password": "Password to login with.",

		"token": "Authentication token to use as an alternative to username/password.",

		"domain_id": "The ID of the Domain to scope to (Identity v3).",

		"domain_name": "The name of the Domain to scope to (Identity v3).",

		"insecure": "Trust self-signed certificates.",

		"cacert_file": "A Custom CA certificate.",

		"endpoint_type": "The catalog endpoint type to use.",

		"endpoints": "The custom endpoints used to override the default endpoint URL.",

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",
	}
}

func configureProvider(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		AccessKey:        d.Get("access_key").(string),
		SecretKey:        d.Get("secret_key").(string),
		CACertFile:       d.Get("cacert_file").(string),
		ClientCertFile:   d.Get("cert").(string),
		ClientKeyFile:    d.Get("key").(string),
		DomainID:         d.Get("domain_id").(string),
		DomainName:       d.Get("domain_name").(string),
		EndpointType:     d.Get("endpoint_type").(string),
		IdentityEndpoint: d.Get("auth_url").(string),
		Insecure:         d.Get("insecure").(bool),
		Password:         d.Get("password").(string),
		Region:           d.Get("region").(string),
		Token:            d.Get("token").(string),
		TenantID:         d.Get("tenant_id").(string),
		TenantName:       d.Get("tenant_name").(string),
		Username:         d.Get("user_name").(string),
		UserID:           d.Get("user_id").(string),
		terraformVersion: terraformVersion,
	}

	config.endpoints = configureProviderEndpoints(d)

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func configureProviderEndpoints(d *schema.ResourceData) map[string]string {
	var availbaleServiceTypes = []string{"as", "ecs", "evs", "ims", "kms", "vpc"}

	endpoints := d.Get("endpoints").(map[string]interface{})
	epMap := make(map[string]string)

	for key, val := range endpoints {
		// ignore unsupportted type
		if !validateEndpointType(key, availbaleServiceTypes) {
			log.Printf("[WARN] the endpoint type %s is unsupportted", key)
			continue
		}

		endpoint := strings.TrimSpace(val.(string))
		// ignore empty string
		if endpoint == "" {
			log.Printf("[WARN] the value of endpoint %s is empty", key)
			continue
		}
		// add prefix "https://" and suffix "/"
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", endpoint)
		}
		if !strings.HasSuffix(endpoint, "/") {
			endpoint = fmt.Sprintf("%s/", endpoint)
		}

		log.Printf("[DEBUG] set endpoint of service %s: %s", key, endpoint)
		epMap[key] = endpoint
	}
	return epMap
}

func validateEndpointType(value string, valid []string) bool {
	for _, str := range valid {
		if value == str {
			return true
		}
	}
	return false
}
