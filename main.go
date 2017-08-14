package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/andrexus/terraform-provider-proxmox/proxmox"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: proxmox.Provider,
	})
}
