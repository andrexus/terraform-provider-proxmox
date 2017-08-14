package proxmox

import (
	"log"

	"github.com/andrexus/goproxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVMQemu() *schema.Resource {
	return &schema.Resource{
		Create: resourceVMQemuCreate,
		Read:   resourceVMQemuRead,
		Delete: resourceVMQemuDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"node": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vmId": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"server_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceVMQemuCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goproxmox.Client)

	err := client.VMs.CreateVM(d.Get("node").(string), d.Get("vmId").(string), nil)

	if err != nil {
		return err
	}

	d.SetId(d.Get("vmId").(string))

	log.Printf("[INFO] VM ID: %s", d.Id())

	return resourceVMQemuRead(d, meta)
}

func resourceVMQemuRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*goproxmox.Client)

	// not implemented yet

	return nil
}

func resourceVMQemuDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goproxmox.Client)

	err := client.VMs.DeleteVM(d.Get("node").(string), d.Id())

	if err != nil {
		return err
	}

	return nil
}
