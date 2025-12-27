// ABOUTME: Acceptance tests for the zabbix_host data source.
// ABOUTME: Tests looking up hosts by technical name.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostDataSource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHostDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zabbix_host.test", "host", rName),
					resource.TestCheckResourceAttr("data.zabbix_host.test", "name", rName+"-display"),
					resource.TestCheckResourceAttr("data.zabbix_host.test", "status", "0"),
					resource.TestCheckResourceAttrSet("data.zabbix_host.test", "id"),
					resource.TestCheckResourceAttr("data.zabbix_host.test", "interfaces.#", "1"),
					resource.TestCheckResourceAttr("data.zabbix_host.test", "groups.#", "1"),
				),
			},
		},
	})
}

func testAccHostDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zabbix_host_group" "test" {
  name = %[1]q
}

resource "zabbix_host" "test" {
  host   = %[1]q
  name   = "%[1]s-display"
  groups = [zabbix_host_group.test.id]
  status = 0

  interfaces = [{
    type   = "agent"
    ip     = "192.168.1.100"
    dns    = ""
    port   = "10050"
    main   = true
    use_ip = true
  }]
}

data "zabbix_host" "test" {
  host = zabbix_host.test.host
}
`, name)
}
