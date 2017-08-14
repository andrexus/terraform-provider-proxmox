terraform-provider-proxmox
==========================

Terraform provider for Proxmox VE

## Description

With this custom terraform provider plugin you can manage your Proxmox resources.

## Usage

Add plugin binary to your ~/.terraformrc file
```
providers {
    proxmox = "/path/to/your/bin/terraform-provider-proxmox"
}
```

### Provider Configuration

```
provider "proxmox" {
    host = "${var.proxmox_host}"
    username  = "${var.proxmox_username}"
    password  = "${var.proxmox_password}"
}
```

##### Argument Reference

The following arguments are required.

* `host` - API host
* `username` - username for accessing Proxmox Control Panel (like root@pam).
* `password` - password for accessing Proxmox Control Panel.

### Resource Configuration

work in progress

## Contribution
This project is based on the [goproxmox] library (https://github.com/andrexus/goproxmox) which is under active development.
So if you want a new feature feel free to send a pull request for the library.


## Licence

[MIT License](https://raw.githubusercontent.com/andrexus/terraform-provider-goproxmox/master/LICENSE.txt)

## Author

[andrexus](https://github.com/andrexus)
