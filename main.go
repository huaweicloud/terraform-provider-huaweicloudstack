package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/huaweicloud/terraform-provider-huaweicloudstack/huaweicloudstack"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: huaweicloudstack.Provider})
}
