/*
Provides a resource to create a mariadb backup_time

Example Usage

```hcl
resource "tencentcloud_mariadb_backup_time" "backup_time" {
  instance_id       = "tdsql-9vqvls95"
  start_backup_time = "01:00"
  end_backup_time   = "04:00"
}
```

Import

mariadb backup_time can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_backup_time.backup_time backup_time_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbBackupTime() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbBackupTimeCreate,
		Read:   resourceTencentCloudMariadbBackupTimeRead,
		Update: resourceTencentCloudMariadbBackupTimeUpdate,
		Delete: resourceTencentCloudMariadbBackupTimeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},
			"start_backup_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time of daily backup window in the format of `mm:ss`, such as 22:00.",
			},
			"end_backup_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time of daily backup window in the format of `mm:ss`, such as 23:59.",
			},
		},
	}
}

func resourceTencentCloudMariadbBackupTimeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_backup_time.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbBackupTimeUpdate(d, meta)
}

func resourceTencentCloudMariadbBackupTimeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_backup_time.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	backupTime, err := service.DescribeMariadbBackupTimeById(ctx, instanceId)
	if err != nil {
		return err
	}

	if backupTime == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MariadbBackupTime` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backupTime.InstanceId != nil {
		_ = d.Set("instance_id", backupTime.InstanceId)
	}

	if backupTime.StartBackupTime != nil {
		_ = d.Set("start_backup_time", backupTime.StartBackupTime)
	}

	if backupTime.EndBackupTime != nil {
		_ = d.Set("end_backup_time", backupTime.EndBackupTime)
	}

	return nil
}

func resourceTencentCloudMariadbBackupTimeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_backup_time.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = mariadb.NewModifyBackupTimeRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId
	if v, ok := d.GetOk("start_backup_time"); ok {
		request.StartBackupTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_backup_time"); ok {
		request.EndBackupTime = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyBackupTime(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if *result.Response.Status != MODIFY_BACKUPTIME_SUCCESS {
			return resource.NonRetryableError(fmt.Errorf("update mariadb backupTime status is fail"))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update mariadb backupTime failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMariadbBackupTimeRead(d, meta)
}

func resourceTencentCloudMariadbBackupTimeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_backup_time.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
