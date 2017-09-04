package proxmox

import (
	"log"

	"strconv"

	"github.com/andrexus/goproxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceVMCreate,
		Read:   resourceVMRead,
		Update: resourceVMUpdate,
		Delete: resourceVMDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"node": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vm_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceVMCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goproxmox.Client)
	node := d.Get("node").(string)
	vmID := d.Get("vm_id").(int)

	config := goproxmox.NewVMConfig()

	if attr, ok := d.GetOk("name"); ok {
		config.SetName(attr.(string))
	}

	if err := client.VMs.CreateVM(node, vmID, config); err != nil {
		return err
	}

	d.SetId(strconv.Itoa(vmID))

	log.Printf("[INFO] VM ID: %s", d.Id())

	return resourceVMRead(d, meta)
}

func resourceVMRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goproxmox.Client)
	node := d.Get("node").(string)
	vmID := d.Get("vm_id").(int)

	config, err := client.VMs.GetVMConfig(node, vmID)
	if err != nil {
		return err
	}
	if value, ok := config.GetName(); ok {
		d.Set("name", value)
	}

	return nil
}

func resourceVMUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goproxmox.Client)
	node := d.Get("node").(string)
	vmID := d.Get("vm_id").(int)

	config := goproxmox.NewVMConfig()

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		config.SetName(n.(string))
	}

	if err := client.VMs.UpdateVM(node, vmID, config, false); err != nil {
		return err
	}

	return resourceVMRead(d, meta)
}

func resourceVMDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goproxmox.Client)
	node := d.Get("node").(string)
	vmID := d.Get("vm_id").(int)

	if err := client.VMs.DeleteVM(node, vmID); err != nil {
		return err
	}

	return nil
}
