// ABOUTME: Acceptance tests for the zabbix_host_group data source.
// ABOUTME: Tests looking up existing host groups by name.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostGroupDataSource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHostGroupDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zabbix_host_group.test", "name", rName),
					resource.TestCheckResourceAttrSet("data.zabbix_host_group.test", "id"),
					resource.TestCheckResourceAttrSet("data.zabbix_host_group.test", "uuid"),
				),
			},
		},
	})
}

func testAccHostGroupDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zabbix_host_group" "test" {
  name = %q
}

data "zabbix_host_group" "test" {
  name = zabbix_host_group.test.name
}
`, name)
}
