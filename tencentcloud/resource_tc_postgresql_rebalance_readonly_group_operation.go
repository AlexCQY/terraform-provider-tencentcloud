/*
Provides a resource to create a postgresql rebalance_readonly_group_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_group" "group_rebalance" {
	master_db_instance_id = local.pgsql_id
	name = "test-pg-readonly-group-rebalance"
	project_id = 0
	vpc_id = "vpc-86v957zb"
	subnet_id = "subnet-enm92y0m"
	replay_lag_eliminate = 1
	replay_latency_eliminate =  1
	max_replay_lag = 100
	max_replay_latency = 512
	min_delay_eliminate_reserve = 1
}

resource "tencentcloud_postgresql_rebalance_readonly_group_operation" "rebalance_readonly_group_operation" {
  read_only_group_id = tencentcloud_postgresql_readonly_group.group_rebalance.id
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationCreate,
		Read:   resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationRead,
		Delete: resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationDelete,
		Schema: map[string]*schema.Schema{
			"read_only_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "readonly Group ID.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_rebalance_readonly_group_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = postgresql.NewRebalanceReadOnlyGroupRequest()
		readOnlyGroupId string
	)
	if v, ok := d.GetOk("read_only_group_id"); ok {
		request.ReadOnlyGroupId = helper.String(v.(string))
		readOnlyGroupId = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().RebalanceReadOnlyGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql RebalanceReadonlyGroupOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(readOnlyGroupId)

	return resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_rebalance_readonly_group_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_rebalance_readonly_group_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
