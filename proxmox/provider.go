package proxmox

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const (
	provider_name           = "terraform-provider-proxmox"
	provider_version string = "v0.1.0"
)

// Provider returns a schema.Provider for Proxmox.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROXMOX_HOST", nil),
				Description: "API host.",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROXMOX_USERNAME", nil),
				Description: "Username for API operations.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROXMOX_PASSWORD", nil),
				Description: "Password for API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"proxmox_resource_vm": resourceVM(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[DEBUG] Configure %s. Version %s", provider_name, provider_version)
	config := Config{
		Host:     d.Get("host").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	return config.Client()
}
