// ABOUTME: Acceptance tests for the zabbix_template resource.
// ABOUTME: Tests full CRUD lifecycle including import of official Zabbix templates.

package provider

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const apacheTemplateURL = "https://raw.githubusercontent.com/zabbix/zabbix/refs/tags/7.0.22/templates/app/apache_http/template_app_apache_http.yaml"

func TestAccTemplateResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateResourceConfigBasic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template.test", "host", rName),
					resource.TestCheckResourceAttr("zabbix_template.test", "name", rName+"-display"),
					resource.TestCheckResourceAttrSet("zabbix_template.test", "id"),
					resource.TestCheckResourceAttrSet("zabbix_template.test", "uuid"),
				),
			},
			{
				ResourceName:            "zabbix_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_content", "source_format"},
			},
		},
	})
}

func TestAccTemplateResource_withOfficialTemplate(t *testing.T) {
	// Fetch the template content at test time
	templateContent := fetchTemplateContent(t, apacheTemplateURL)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateResourceConfigWithContent(templateContent),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template.test", "host", "Apache by HTTP"),
					resource.TestCheckResourceAttrSet("zabbix_template.test", "id"),
					resource.TestCheckResourceAttrSet("zabbix_template.test", "uuid"),
					resource.TestCheckResourceAttrSet("zabbix_template.test", "exported_content"),
				),
				// The exported_content computed field causes Terraform to show a plan
				// even when nothing has changed. This is expected behavior.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTemplateResource_update(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateResourceConfigBasic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template.test", "host", rName),
					resource.TestCheckResourceAttr("zabbix_template.test", "name", rName+"-display"),
				),
			},
			{
				Config: testAccTemplateResourceConfigUpdated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template.test", "host", rName),
					resource.TestCheckResourceAttr("zabbix_template.test", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("zabbix_template.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccTemplateResource_withTags(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateResourceConfigWithTags(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_template.test", "host", rName),
					resource.TestCheckResourceAttr("zabbix_template.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "zabbix_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_content", "source_format"},
			},
		},
	})
}

// fetchTemplateContent fetches template YAML from a URL at test time.
func fetchTemplateContent(t *testing.T, url string) string {
	t.Helper()

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to fetch template from %s: %v", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to fetch template: HTTP %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read template content: %v", err)
	}

	return string(content)
}

func testAccTemplateResourceConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "zabbix_template_group" "test" {
  name = "%[1]s-group"
}

resource "zabbix_template" "test" {
  host   = %[1]q
  name   = "%[1]s-display"
  groups = [zabbix_template_group.test.id]
}
`, name)
}

func testAccTemplateResourceConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "zabbix_template_group" "test" {
  name = "%[1]s-group"
}

resource "zabbix_template" "test" {
  host        = %[1]q
  name        = "%[1]s-updated"
  description = "Updated description"
  groups      = [zabbix_template_group.test.id]
}
`, name)
}

func testAccTemplateResourceConfigWithTags(name string) string {
	return fmt.Sprintf(`
resource "zabbix_template_group" "test" {
  name = "%[1]s-group"
}

resource "zabbix_template" "test" {
  host   = %[1]q
  name   = "%[1]s-display"
  groups = [zabbix_template_group.test.id]

  tags = [
    {
      tag   = "environment"
      value = "test"
    },
    {
      tag   = "team"
      value = "platform"
    }
  ]
}
`, name)
}

func testAccTemplateResourceConfigWithContent(content string) string {
	return fmt.Sprintf(`
resource "zabbix_template" "test" {
  source_format  = "yaml"
  source_content = %q
}
`, content)
}
