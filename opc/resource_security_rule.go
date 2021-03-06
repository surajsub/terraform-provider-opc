package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOPCSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCSecurityRuleCreate,
		Read:   resourceOPCSecurityRuleRead,
		Update: resourceOPCSecurityRuleUpdate,
		Delete: resourceOPCSecurityRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flow_direction": {
				Type:     schema.TypeString,
				Required: true,
			},
			"acl": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dst_ip_address_prefixes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"src_ip_address_prefixes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_protocols": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dst_vnic_set": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"src_vnic_set": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsOptionalSchema(),
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOPCSecurityRuleCreate(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.SecurityRules()

	input := compute.CreateSecurityRuleInput{
		Name:          d.Get("name").(string),
		FlowDirection: d.Get("flow_direction").(string),
		Enabled:       d.Get("enabled").(bool),
	}

	if acl, ok := d.GetOk("acl"); ok {
		input.ACL = acl.(string)
	}

	if srcVNicSet, ok := d.GetOk("src_vnic_set"); ok {
		input.SrcVnicSet = srcVNicSet.(string)
	}

	if dstVNicSet, ok := d.GetOk("dst_vnic_set"); ok {
		input.DstVnicSet = dstVNicSet.(string)
	}

	securityProtocols := getStringList(d, "security_protocols")
	if len(securityProtocols) != 0 {
		input.SecProtocols = securityProtocols
	}

	srcIPAddressPrefixes := getStringList(d, "src_ip_address_prefixes")
	if len(srcIPAddressPrefixes) != 0 {
		input.SrcIPAddressPrefixSets = srcIPAddressPrefixes
	}

	dstIPAddressPrefixes := getStringList(d, "dst_ip_address_prefixes")
	if len(dstIPAddressPrefixes) != 0 {
		input.DstIPAddressPrefixSets = dstIPAddressPrefixes
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = description.(string)
	}

	info, err := resClient.CreateSecurityRule(&input)
	if err != nil {
		return fmt.Errorf("Error creating Security Rule: %s", err)
	}

	d.SetId(info.Name)
	return resourceOPCSecurityRuleRead(d, meta)
}

func resourceOPCSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.SecurityRules()

	input := compute.GetSecurityRuleInput{
		Name: d.Id(),
	}

	result, err := resClient.GetSecurityRule(&input)
	if err != nil {
		// SecurityRule does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading security rule %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", result.Name)
	d.Set("flow_direction", result.FlowDirection)
	d.Set("enabled", result.Enabled)
	d.Set("acl", result.ACL)
	d.Set("src_vnic_set", result.SrcVnicSet)
	d.Set("dst_vnic_set", result.DstVnicSet)
	d.Set("description", result.Description)
	d.Set("uri", result.URI)

	if err := setStringList(d, "security_protocols", result.SecProtocols); err != nil {
		return err
	}
	if err := setStringList(d, "dst_ip_address_prefixes", result.DstIPAddressPrefixSets); err != nil {
		return err
	}
	if err := setStringList(d, "src_ip_address_prefixes", result.SrcIPAddressPrefixSets); err != nil {
		return err
	}
	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}
	return nil
}

func resourceOPCSecurityRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.SecurityRules()

	input := compute.UpdateSecurityRuleInput{
		Name:          d.Get("name").(string),
		FlowDirection: d.Get("flow_direction").(string),
		Enabled:       d.Get("enabled").(bool),
	}

	if acl, ok := d.GetOk("acl"); ok {
		input.ACL = acl.(string)
	}

	if srcVNicSet, ok := d.GetOk("src_vnic_set"); ok {
		input.SrcVnicSet = srcVNicSet.(string)
	}

	if dstVNicSet, ok := d.GetOk("dst_vnic_set"); ok {
		input.DstVnicSet = dstVNicSet.(string)
	}

	securityProtocols := getStringList(d, "security_protocols")
	if len(securityProtocols) != 0 {
		input.SecProtocols = securityProtocols
	}

	srcIPAddressPrefixes := getStringList(d, "src_ip_address_prefixes")
	if len(srcIPAddressPrefixes) != 0 {
		input.SrcIPAddressPrefixSets = srcIPAddressPrefixes
	}

	dstIPAddressPrefixes := getStringList(d, "dst_ip_address_prefixes")
	if len(dstIPAddressPrefixes) != 0 {
		input.DstIPAddressPrefixSets = dstIPAddressPrefixes
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		input.Tags = tags
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = description.(string)
	}
	info, err := resClient.UpdateSecurityRule(&input)
	if err != nil {
		return fmt.Errorf("Error updating Security Rule: %s", err)
	}

	d.SetId(info.Name)
	return resourceOPCSecurityRuleRead(d, meta)
}

func resourceOPCSecurityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.SecurityRules()
	name := d.Id()

	input := compute.DeleteSecurityRuleInput{
		Name: name,
	}
	if err := resClient.DeleteSecurityRule(&input); err != nil {
		return fmt.Errorf("Error deleting Security Rule: %s", err)
	}
	return nil
}
