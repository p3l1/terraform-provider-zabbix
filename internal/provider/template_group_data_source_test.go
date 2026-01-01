// ABOUTME: Acceptance tests for the zabbix_template_group data source.
// ABOUTME: Tests looking up existing template groups by name.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTemplateGroupDataSource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateGroupDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zabbix_template_group.test", "name", rName),
					resource.TestCheckResourceAttrSet("data.zabbix_template_group.test", "id"),
					resource.TestCheckResourceAttrSet("data.zabbix_template_group.test", "uuid"),
				),
			},
		},
	})
}

func testAccTemplateGroupDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zabbix_template_group" "test" {
  name = %q
}

data "zabbix_template_group" "test" {
  name = zabbix_template_group.test.name
}
`, name)
}
