package proxmox

import (
	"log"

	"errors"
	"fmt"

	"github.com/andrexus/goproxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceVolumeCreate,
		Read:   resourceVolumeRead,
		Delete: resourceVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"node": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vm_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"filename": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goproxmox.Client)
	node := d.Get("node").(string)
	storageName := d.Get("storage_name").(string)
	vmID := d.Get("vm_id").(int)
	filename := d.Get("filename").(string)
	size := d.Get("size").(string)

	if err := client.Storages.CreateVolume(node, storageName, vmID, filename, size, nil); err != nil {
		return err
	}
	volumeId := fmt.Sprintf("%s:%d/%s", storageName, vmID, filename)
	d.SetId(volumeId)

	log.Printf("[INFO] Volume ID: %s", d.Id())

	return resourceVolumeRead(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goproxmox.Client)
	node := d.Get("node").(string)
	storageName := d.Get("storage_name").(string)

	volumes, err := client.Storages.GetStorageVolumes(node, storageName)
	if err != nil {
		return err
	}
	for _, volume := range volumes {
		if volume.VolumeId == d.Id() {
			d.Set("node", node)
			d.Set("storage_name", storageName)
			d.Set("vm_id", volume.VMID)
			d.Set("size", volume.Size)
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Volume with id %s not found", d.Id()))
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goproxmox.Client)
	node := d.Get("node").(string)
	storageName := d.Get("storage_name").(string)

	if err := client.Storages.DeleteVolume(node, storageName, d.Id()); err != nil {
		return err
	}

	return nil
}
