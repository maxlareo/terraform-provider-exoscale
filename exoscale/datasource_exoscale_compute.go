package exoscale

import (
	"context"
	"errors"
	"fmt"

	"github.com/exoscale/egoscale"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCompute() *schema.Resource {
	return &schema.Resource{
		Description: "Fetch Exoscale Compute Instances data.",
		Schema: map[string]*schema.Schema{
			"id": {
				Description:   "The compute instance ID to match.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"hostname", "tags"},
			},
			"hostname": {
				Description:   "The instance hostname to match.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "tags"},
			},
			"tags": {
				Description: "The instance tags to match (map of key/value).",
				Type:        schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:      true,
				ConflictsWith: []string{"id", "hostname"},
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the Compute instance was created",
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the availability zone for the Compute instance",
			},
			"template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the template for the Compute instance",
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current size of the Compute instance",
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the Compute instance disk",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of cpu the Compute instance is running with",
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Memory allocated for the Compute instance",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the Compute instance",
			},

			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Compute instance public ipv4 address",
			},
			"ip6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Compute instance public ipv6 address (if ipv6 is enabled)",
			},
			"private_network_ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Compute instance private IP addresses (in managed Private Networks only)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},

		Read: dataSourceComputeRead,
	}
}

func dataSourceComputeRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()

	client := GetComputeClient(meta)

	req := egoscale.VirtualMachine{}

	computeName, byName := d.GetOk("hostname")
	computeID, byID := d.GetOk("id")
	computeTag, byTag := d.GetOk("tags")

	switch {
	case byName:
		req.Name = computeName.(string)

	case byID:
		var err error
		if req.ID, err = egoscale.ParseUUID(computeID.(string)); err != nil {
			return fmt.Errorf("invalid value for id: %s", err)
		}

	case byTag:
		for key, value := range computeTag.(map[string]interface{}) {
			req.Tags = append(req.Tags, egoscale.ResourceTag{
				Key:   key,
				Value: value.(string),
			})
		}

	default:
		return errors.New("either hostname, id or tags must be specified")
	}

	resp, err := client.GetWithContext(ctx, &req)
	if err != nil {
		return err
	}
	instance := resp.(*egoscale.VirtualMachine)

	// Querying VM NICs separately because the non-default NICs IP addresses
	// are not returned in the CS listVirtualMachines operation results.
	resp, err = client.RequestWithContext(ctx, &egoscale.ListNics{VirtualMachineID: instance.ID})
	if err != nil {
		return err
	}
	instance.Nic = resp.(*egoscale.ListNicsResponse).Nic

	resp, err = client.GetWithContext(ctx, &egoscale.Volume{
		VirtualMachineID: instance.ID,
		Type:             "ROOT",
	})
	if err != nil {
		return err
	}
	diskSize := resp.(*egoscale.Volume).Size >> 30

	return dataSourceComputeApply(d, instance, diskSize)
}

func dataSourceComputeApply(d *schema.ResourceData, instance *egoscale.VirtualMachine, diskSize uint64) error {
	d.SetId(instance.ID.String())

	if err := d.Set("id", d.Id()); err != nil {
		return err
	}
	if err := d.Set("hostname", instance.Name); err != nil {
		return err
	}
	if err := d.Set("created", instance.Created); err != nil {
		return err
	}
	if err := d.Set("zone", instance.ZoneName); err != nil {
		return err
	}
	if err := d.Set("template", instance.TemplateName); err != nil {
		return err
	}
	if err := d.Set("size", instance.ServiceOfferingName); err != nil {
		return err
	}
	if err := d.Set("disk_size", diskSize); err != nil {
		return err
	}
	if err := d.Set("cpu", instance.CPUNumber); err != nil {
		return err
	}
	if err := d.Set("memory", instance.Memory); err != nil {
		return err
	}
	if err := d.Set("state", instance.State); err != nil {
		return err
	}
	if err := d.Set("ip_address", instance.DefaultNic().IPAddress.String()); err != nil {
		return err
	}
	if err := d.Set("ip6_address", instance.DefaultNic().IP6Address.String()); err != nil {
		return err
	}

	privateNetworkIPs := make([]string, 0)
	for _, nic := range instance.Nic {
		if nic.IsDefault {
			continue
		}
		privateNetworkIPs = append(privateNetworkIPs, nic.IPAddress.String())
	}
	if err := d.Set("private_network_ip_addresses", privateNetworkIPs); err != nil {
		return err
	}

	tags := make(map[string]interface{})
	for _, tag := range instance.Tags {
		tags[tag.Key] = tag.Value
	}
	if err := d.Set("tags", tags); err != nil {
		return err
	}

	return nil
}
