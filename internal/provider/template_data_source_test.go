// ABOUTME: Acceptance tests for the zabbix_template data source.
// ABOUTME: Tests looking up templates by technical name and retrieving exported content.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTemplateDataSource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zabbix_template.test", "host", rName),
					resource.TestCheckResourceAttr("data.zabbix_template.test", "name", rName+"-display"),
					resource.TestCheckResourceAttrSet("data.zabbix_template.test", "id"),
					resource.TestCheckResourceAttrSet("data.zabbix_template.test", "uuid"),
					resource.TestCheckResourceAttrSet("data.zabbix_template.test", "exported_content"),
				),
			},
		},
	})
}

func TestAccTemplateDataSource_withOfficialTemplate(t *testing.T) {
	// Fetch the template content at test time
	templateContent := fetchTemplateContent(t, apacheTemplateURL)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateDataSourceConfigWithContent(templateContent),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zabbix_template.test", "host", "Apache by HTTP"),
					resource.TestCheckResourceAttrSet("data.zabbix_template.test", "id"),
					resource.TestCheckResourceAttrSet("data.zabbix_template.test", "uuid"),
					resource.TestCheckResourceAttrSet("data.zabbix_template.test", "exported_content"),
				),
				// The exported_content computed field causes Terraform to show a plan
				// even when nothing has changed. This is expected behavior.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccTemplateDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zabbix_template_group" "test" {
  name = %[1]q
}

resource "zabbix_template" "test" {
  host   = %[1]q
  name   = "%[1]s-display"
  groups = [zabbix_template_group.test.id]
}

data "zabbix_template" "test" {
  host = zabbix_template.test.host
}
`, name)
}

func testAccTemplateDataSourceConfigWithContent(content string) string {
	return fmt.Sprintf(`
resource "zabbix_template" "test" {
  source_format  = "yaml"
  source_content = %q
}

data "zabbix_template" "test" {
  host = zabbix_template.test.host
}
`, content)
}
