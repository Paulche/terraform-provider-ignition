package ignition

import (
	"github.com/coreos/ignition/config/v2_2/types"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetworkdUnit() *schema.Resource {
	return &schema.Resource{
		Exists: resourceNetworkdUnitExists,
		Read:   resourceNetworkdUnitRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dropin": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceNetworkdUnitRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildNetworkdUnit(d, globalCache)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceNetworkdUnitDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourceNetworkdUnitExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildNetworkdUnit(d, globalCache)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildNetworkdUnit(d *schema.ResourceData, c *cache) (string, error) {
	unit := &types.Networkdunit{
		Name:     d.Get("name").(string),
		Contents: d.Get("content").(string),
	}

	for _, raw := range d.Get("dropin").([]interface{}) {
		value := raw.(map[string]interface{})

		dropIn := types.NetworkdDropin{
			Name:     value["name"].(string),
			Contents: value["content"].(string),
		}

		if err := handleReport(dropIn.Validate()); err != nil {
			return "", err
		}

		unit.Dropins = append(unit.Dropins, dropIn)
	}

	return c.addNetworkdUnit(unit), handleReport(unit.Validate())
}
