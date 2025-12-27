// ABOUTME: Acceptance tests for the zabbix_host resource.
// ABOUTME: Tests full CRUD lifecycle, interfaces, templates, tags, and import functionality.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHostResourceConfigBasic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_host.test", "host", rName),
					resource.TestCheckResourceAttr("zabbix_host.test", "name", rName+"-display"),
					resource.TestCheckResourceAttr("zabbix_host.test", "status", "0"),
					resource.TestCheckResourceAttrSet("zabbix_host.test", "id"),
					resource.TestCheckResourceAttr("zabbix_host.test", "interfaces.#", "1"),
					resource.TestCheckResourceAttr("zabbix_host.test", "interfaces.0.type", "agent"),
					resource.TestCheckResourceAttr("zabbix_host.test", "interfaces.0.ip", "192.168.1.100"),
					resource.TestCheckResourceAttr("zabbix_host.test", "interfaces.0.port", "10050"),
					resource.TestCheckResourceAttr("zabbix_host.test", "interfaces.0.main", "true"),
					resource.TestCheckResourceAttr("zabbix_host.test", "interfaces.0.use_ip", "true"),
				),
			},
			{
				ResourceName:      "zabbix_host.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccHostResource_update(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rNameUpdated := acctest.RandomWithPrefix("tf-acc-test-upd")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHostResourceConfigBasic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_host.test", "host", rName),
					resource.TestCheckResourceAttr("zabbix_host.test", "name", rName+"-display"),
					resource.TestCheckResourceAttr("zabbix_host.test", "status", "0"),
				),
			},
			{
				Config: testAccHostResourceConfigUpdated(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_host.test", "host", rNameUpdated),
					resource.TestCheckResourceAttr("zabbix_host.test", "name", rNameUpdated+"-display-updated"),
					resource.TestCheckResourceAttr("zabbix_host.test", "status", "1"),
					resource.TestCheckResourceAttr("zabbix_host.test", "interfaces.0.ip", "192.168.1.200"),
				),
			},
		},
	})
}

func TestAccHostResource_withTags(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHostResourceConfigWithTags(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_host.test", "host", rName),
					resource.TestCheckResourceAttr("zabbix_host.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:      "zabbix_host.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccHostResource_multipleInterfaces(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHostResourceConfigMultipleInterfaces(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_host.test", "host", rName),
					resource.TestCheckResourceAttr("zabbix_host.test", "interfaces.#", "2"),
				),
			},
		},
	})
}

func TestAccHostResource_multipleGroups(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHostResourceConfigMultipleGroups(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zabbix_host.test", "host", rName),
					resource.TestCheckResourceAttr("zabbix_host.test", "groups.#", "2"),
				),
			},
		},
	})
}

func testAccHostResourceConfigBasic(name string) string {
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
`, name)
}

func testAccHostResourceConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "zabbix_host_group" "test" {
  name = %[1]q
}

resource "zabbix_host" "test" {
  host   = %[1]q
  name   = "%[1]s-display-updated"
  groups = [zabbix_host_group.test.id]
  status = 1

  interfaces = [{
    type   = "agent"
    ip     = "192.168.1.200"
    dns    = ""
    port   = "10050"
    main   = true
    use_ip = true
  }]
}
`, name)
}

func testAccHostResourceConfigWithTags(name string) string {
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

func testAccHostResourceConfigMultipleInterfaces(name string) string {
	return fmt.Sprintf(`
resource "zabbix_host_group" "test" {
  name = %[1]q
}

resource "zabbix_host" "test" {
  host   = %[1]q
  name   = "%[1]s-display"
  groups = [zabbix_host_group.test.id]
  status = 0

  interfaces = [
    {
      type   = "agent"
      ip     = "192.168.1.100"
      dns    = ""
      port   = "10050"
      main   = true
      use_ip = true
    },
    {
      type   = "agent"
      ip     = "192.168.1.101"
      dns    = ""
      port   = "10050"
      main   = false
      use_ip = true
    }
  ]
}
`, name)
}

func testAccHostResourceConfigMultipleGroups(name string) string {
	return fmt.Sprintf(`
resource "zabbix_host_group" "test1" {
  name = "%[1]s-group1"
}

resource "zabbix_host_group" "test2" {
  name = "%[1]s-group2"
}

resource "zabbix_host" "test" {
  host   = %[1]q
  name   = "%[1]s-display"
  groups = [zabbix_host_group.test1.id, zabbix_host_group.test2.id]
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
`, name)
}
