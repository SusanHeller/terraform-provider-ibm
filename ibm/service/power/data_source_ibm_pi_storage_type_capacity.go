// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"

	"log"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	PITypeName = "pi_storage_type"
)

func DataSourceIBMPIStorageTypeCapacity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIStorageTypeCapacityRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			PITypeName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Storage type name",
			},
			// Computed Attributes
			Attr_MaximumStorageAllocation: {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Maximum storage allocation",
			},
			Attr_StoragePoolsCapacity: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Storage pools capacity",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_MaxAllocationSize: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum allocation storage size (GB)",
						},
						Attr_PoolName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pool name",
						},
						Attr_StorageType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage type of the storage pool",
						},
						Attr_TotalCapacity: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total pool capacity (GB)",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIStorageTypeCapacityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	storageType := d.Get(PITypeName).(string)

	client := st.NewIBMPIStorageCapacityClient(ctx, sess, cloudInstanceID)
	stc, err := client.GetStorageTypeCapacity(storageType)
	if err != nil {
		log.Printf("[ERROR] get storage type capacity failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, storageType))

	if stc.MaximumStorageAllocation != nil {
		msa := stc.MaximumStorageAllocation
		data := map[string]interface{}{
			Attr_MaxAllocationSize: *msa.MaxAllocationSize,
			Attr_StoragePool:       *msa.StoragePool,
			Attr_StorageType:       *msa.StorageType,
		}
		d.Set(Attr_MaximumStorageAllocation, flex.Flatten(data))
	}

	result := make([]map[string]interface{}, 0, len(stc.StoragePoolsCapacity))
	for _, sp := range stc.StoragePoolsCapacity {
		data := map[string]interface{}{
			Attr_MaxAllocationSize: *sp.MaxAllocationSize,
			Attr_PoolName:          sp.PoolName,
			Attr_StorageType:       sp.StorageType,
			Attr_TotalCapacity:     sp.TotalCapacity,
		}
		result = append(result, data)
	}
	d.Set(Attr_StoragePoolsCapacity, result)

	return nil
}
