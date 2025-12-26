// ABOUTME: Acceptance tests for the zabbix_host_group resource.
// ABOUTME: Tests full CRUD lifecycle and import functionality.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostGroupResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHostGroupResourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_host_group.test", "name", rName),
					resource.TestCheckResourceAttrSet("zabbix_host_group.test", "id"),
					resource.TestCheckResourceAttrSet("zabbix_host_group.test", "uuid"),
				),
			},
			{
				ResourceName:      "zabbix_host_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccHostGroupResource_update(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rNameUpdated := acctest.RandomWithPrefix("tf-acc-test-updated")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHostGroupResourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_host_group.test", "name", rName),
				),
			},
			{
				Config: testAccHostGroupResourceConfig(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_host_group.test", "name", rNameUpdated),
				),
			},
		},
	})
}

func testAccHostGroupResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zabbix_host_group" "test" {
  name = %q
}
`, name)
}
