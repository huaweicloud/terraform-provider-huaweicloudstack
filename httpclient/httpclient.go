package httpclient

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const uaEnvVar = "TF_APPEND_USER_AGENT"

// TerraformUserAgent returns a User-Agent header for a Terraform version string.
//
// Deprecated: This will be removed in v2 without replacement. If you need
// its functionality, you can copy it or reference the v1 package.
func TerraformUserAgent(version string) string {
	ua := fmt.Sprintf("HuaweiCloudStack Terraform Provider/%s (+https://github.com/huaweicloud/terraform-provider-huaweicloudstack)", version)

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}
