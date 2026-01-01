// ABOUTME: Acceptance tests for the zabbix_template_group resource.
// ABOUTME: Tests full CRUD lifecycle and import functionality.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTemplateGroupResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateGroupResourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template_group.test", "name", rName),
					resource.TestCheckResourceAttrSet("zabbix_template_group.test", "id"),
					resource.TestCheckResourceAttrSet("zabbix_template_group.test", "uuid"),
				),
			},
			{
				ResourceName:      "zabbix_template_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTemplateGroupResource_update(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rNameUpdated := acctest.RandomWithPrefix("tf-acc-test-updated")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateGroupResourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template_group.test", "name", rName),
				),
			},
			{
				Config: testAccTemplateGroupResourceConfig(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template_group.test", "name", rNameUpdated),
				),
			},
		},
	})
}

func testAccTemplateGroupResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zabbix_template_group" "test" {
  name = %q
}
`, name)
}
