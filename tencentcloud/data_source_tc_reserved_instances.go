/*
Use this data source to query reserved instances.

Example Usage

```hcl
data "tencentcloud_reserved_instances" "instances" {
  availability_zone = "na-siliconvalley-1"
  instance_type     = "S2.MEDIUM8"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudReservedInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudReservedInstancesRead,

		Schema: map[string]*schema.Schema{
			"reserved_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the reserved instance to be query.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone that the reserved instance locates at.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of reserved instance.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"reserved_instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of reserved instance. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reserved_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the reserved instance.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of reserved instance.",
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of reserved instance.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone of the reserved instance.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time of the reserved instance.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiry time of the reserved instance.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the reserved instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudReservedInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_reserved_instances.read")
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	filter := make(map[string]string)
	if v, ok := d.GetOk("reserved_instance_id"); ok {
		filter["reserved-instances-id"] = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		filter["zone"] = v.(string)
	}
	if v, ok := d.GetOk("instance_type"); ok {
		filter["instance-type"] = v.(string)
	}

	var instances []*cvm.ReservedInstances
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instances, errRet = cvmService.DescribeReservedInstanceByFilter(ctx, filter)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	instanceList := make([]map[string]interface{}, 0, len(instances))
	ids := make([]string, 0, len(instances))
	for _, instance := range instances {
		mapping := map[string]interface{}{
			"reserved_instance_id": instance.ReservedInstancesId,
			"instance_type":        instance.InstanceType,
			"instance_count":       instance.InstanceCount,
			"availability_zone":    instance.Zone,
			"start_time":           instance.StartTime,
			"end_time":             instance.EndTime,
			"status":               instance.State,
		}
		instanceList = append(instanceList, mapping)
		ids = append(ids, *instance.ReservedInstancesId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("reserved_instance_list", instanceList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set reserved instance list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), instanceList); err != nil {
			return err
		}
	}
	return nil
}
