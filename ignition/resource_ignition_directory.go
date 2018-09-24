package ignition

import (
	"github.com/coreos/ignition/config/v2_2/types"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDirectory() *schema.Resource {
	return &schema.Resource{
		Exists: resourceDirectoryExists,
		Read:   resourceDirectoryRead,
		Schema: map[string]*schema.Schema{
			"filesystem": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"overwrite": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"mode": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"uid": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"gid": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildDirectory(d, globalCache)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceDirectoryExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildDirectory(d, globalCache)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildDirectory(d *schema.ResourceData, c *cache) (string, error) {
	dir := &types.Directory{}
	dir.Filesystem = d.Get("filesystem").(string)
	if err := handleReport(dir.ValidateFilesystem()); err != nil {
		return "", err
	}

	dir.Path = d.Get("path").(string)
	if err := handleReport(dir.ValidatePath()); err != nil {
		return "", err
	}

	overwrite := d.Get("overwrite").(bool)
	dir.Overwrite = &overwrite

	mode := d.Get("mode").(int)
	dir.Mode = &mode
	if err := handleReport(dir.ValidateMode()); err != nil {
		return "", err
	}

	uid := d.Get("uid").(int)
	user := types.NodeUser{ID: &uid}
	if uid != 0 {
		dir.User = &user
	}

	gid := d.Get("gid").(int)
	group := types.NodeGroup{ID: &gid}
	if gid != 0 {
		dir.Group = &group
	}

	return c.addDirectory(dir), nil
}
