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
			},
			"start_at_boot": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"network_devices": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"model": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bridge": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"firewall": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"link_down": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"macaddr": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"queues": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"rate": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"tag": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"trunks": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"virtio_devices": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"file": {
							Type:     schema.TypeString,
							Required: true,
						},
						"format": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"backup": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"iothread": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"snapshot": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceVMCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goproxmox.Client)
	node := d.Get("node").(string)
	vmID := d.Get("vm_id").(int)

	config := new(goproxmox.VMConfig)

	if v, ok := d.GetOk("name"); ok {
		config.Name = goproxmox.String(v.(string))
	}
	if v, ok := d.GetOk("start_at_boot"); ok {
		config.StartAtBoot = goproxmox.Bool(v.(bool))
	}
	if v, ok := d.GetOk("network_devices"); ok {
		networkDevices := v.(*schema.Set)
		for _, element := range networkDevices.List() {
			elem := element.(map[string]interface{})
			log.Printf("[DEBUG] Network device elem %v", elem)
			cardModel, err := goproxmox.NetworkCardModelFromString(elem["model"].(string))
			if err != nil {
				return err
			}
			number := elem["number"].(int)
			networkDevice := &goproxmox.NetworkDevice{
				Model: &cardModel,
			}
			if val, ok := elem["bridge"]; ok {
				networkDevice.Bridge = goproxmox.String(val.(string))
			}
			//if val, ok := elem["firewall"]; ok {
			//	networkDevice.Firewall = goproxmox.Bool(val.(bool))
			//}
			//if val, ok := elem["link_down"]; ok {
			//	networkDevice.LinkDown = goproxmox.Bool(val.(bool))
			//}
			//if val, ok := elem["macaddr"]; ok {
			//	networkDevice.MacAddr = goproxmox.String(val.(string))
			//}
			//if val, ok := elem["queues"]; ok {
			//	networkDevice.Queues = goproxmox.Int(val.(int))
			//}
			//if val, ok := elem["rate"]; ok {
			//	networkDevice.Rate = goproxmox.Float64(val.(float64))
			//}
			//if val, ok := elem["tag"]; ok {
			//	networkDevice.Tag = goproxmox.Int(val.(int))
			//}
			//if val, ok := elem["trunks"]; ok {
			//	networkDevice.Trunks = goproxmox.String(val.(string))
			//}
			log.Printf("[DEBUG] Network device %v", networkDevice.GetQMOptionValue())
			config.AddNetworkDevice(number, networkDevice)
		}
	}
	if v, ok := d.GetOk("virtio_devices"); ok {
		virtIODevices := v.(*schema.Set)
		for _, element := range virtIODevices.List() {
			elem := element.(map[string]interface{})
			number := elem["number"].(int)
			volumeFormat, err := goproxmox.VolumeFormatFromString(elem["format"].(string))
			if err != nil {
				return err
			}
			virtIODevice := &goproxmox.VirtIODevice{
				File:     goproxmox.String(elem["file"].(string)),
				Format:   &volumeFormat,
				//Backup:   goproxmox.Bool(elem["backup"].(bool)),
				//IOThread: goproxmox.Bool(elem["iothread"].(bool)),
				Size:     goproxmox.String(elem["size"].(string)),
				//Snapshot: goproxmox.Bool(elem["snapshot"].(bool)),
			}

			log.Printf("[DEBUG] VirtIO device %v", virtIODevice.GetQMOptionValue())
			config.AddVirtIODevice(number, virtIODevice)
		}
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
	if config.Name != nil {
		d.Set("name", config.Name)
	}

	return nil
}

func resourceVMUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goproxmox.Client)
	node := d.Get("node").(string)
	vmID := d.Get("vm_id").(int)

	config, err := client.VMs.GetVMConfig(node, vmID)
	if err != nil {
		return err
	}
	//if attr, ok := config.GetName(); ok {
	//	d.Set("name", attr)
	//} else {
	//	config.UnsetName()
	//	d.Set("name", nil)
	//}
	//
	//if d.HasChange("name") {
	//	_, n := d.GetChange("name")
	//	config.SetName(n.(string))
	//}

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
