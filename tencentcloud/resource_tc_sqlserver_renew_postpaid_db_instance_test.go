package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverRenewPostpaidDBInstanceResource_basic -v
func TestAccTencentCloudSqlserverRenewPostpaidDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRenewPostpaidDBInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_renew_postpaid_db_instance.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_renew_postpaid_db_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverRenewPostpaidDBInstance = defaultVpcSubnets + defaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = local.vpc_id
  subnet_id              = local.subnet_id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [local.sg_id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_config_terminate_db_instance" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
}

resource "tencentcloud_sqlserver_renew_postpaid_db_instance" "example" {
  instance_id = tencentcloud_sqlserver_config_terminate_db_instance.example.id
}
`
