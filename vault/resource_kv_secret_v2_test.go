// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/vault/api"

	"github.com/hashicorp/terraform-provider-vault/internal/consts"
	"github.com/hashicorp/terraform-provider-vault/internal/provider"
	"github.com/hashicorp/terraform-provider-vault/testutil"
)

var testKVV2Data = map[string]interface{}{
	"foo": "bar",
	"baz": "qux",
}

func TestAccKVSecretV2(t *testing.T) {
	t.Parallel()
	resourceName := "vault_kv_secret_v2.test"
	mount := acctest.RandomWithPrefix("tf-kvv2")
	name := acctest.RandomWithPrefix("tf-secret")

	updatedMount := acctest.RandomWithPrefix("random-prefix/tf-cloud-metadata")
	updatedName := acctest.RandomWithPrefix("tf-database-creds")

	customMetadata := `{"extra":"cheese","pizza":"please"}`

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKVSecretV2Config_initial(mount, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldMount, mount),
					resource.TestCheckResourceAttr(resourceName, consts.FieldName, name),
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, fmt.Sprintf("%s/data/%s", mount, name)),
					resource.TestCheckResourceAttr(resourceName, "delete_all_versions", "true"),
					resource.TestCheckResourceAttr(resourceName, "data.zip", "zap"),
					resource.TestCheckResourceAttr(resourceName, "data.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "data.flag", "false"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.cas_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.data.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.delete_version_after", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.max_versions", "0"),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "metadata.version", "1"),
					resource.TestCheckResourceAttr(resourceName, "metadata.destroyed", "false"),
					resource.TestCheckResourceAttr(resourceName, "metadata.deletion_time", ""),
					resource.TestCheckResourceAttr(resourceName, "metadata.custom_metadata", "null"),
				),
			},
			{
				Config: testKVSecretV2Config_updated(mount, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldMount, mount),
					resource.TestCheckResourceAttr(resourceName, consts.FieldName, name),
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, fmt.Sprintf("%s/data/%s", mount, name)),
					resource.TestCheckResourceAttr(resourceName, "delete_all_versions", "true"),
					resource.TestCheckResourceAttr(resourceName, "data.zip", "zoop"),
					resource.TestCheckResourceAttr(resourceName, "data.foo", "baz"),
					resource.TestCheckResourceAttr(resourceName, "data.flag", "false"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.cas_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.data.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.data.extra", "cheese"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.data.pizza", "please"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.delete_version_after", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.max_versions", "5"),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "metadata.version", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.destroyed", "false"),
					resource.TestCheckResourceAttr(resourceName, "metadata.deletion_time", ""),
					resource.TestCheckResourceAttr(resourceName, "metadata.custom_metadata", customMetadata),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"data_json", "disable_read",
					"delete_all_versions",
				},
			},
			{
				Config: testKVSecretV2Config_initial(updatedMount, updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldMount, updatedMount),
					resource.TestCheckResourceAttr(resourceName, consts.FieldName, updatedName),
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, fmt.Sprintf("%s/data/%s", updatedMount, updatedName)),
					resource.TestCheckResourceAttr(resourceName, "delete_all_versions", "true"),
					resource.TestCheckResourceAttr(resourceName, "data.zip", "zap"),
					resource.TestCheckResourceAttr(resourceName, "data.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "data.flag", "false"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.cas_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.data.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.delete_version_after", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_metadata.0.max_versions", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"data_json", "disable_read",
					"delete_all_versions",
				},
			},
		},
	})
}

func TestAccKVSecretV2_DisableRead(t *testing.T) {
	t.Parallel()
	resourceName := "vault_kv_secret_v2.test"
	mount := acctest.RandomWithPrefix("tf-kvv2")
	name := acctest.RandomWithPrefix("tf-secret")

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					mountKVEngine(t, mount, name)
				},
				Config: testKVSecretV2Config_DisableRead(mount, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldMount, mount),
					resource.TestCheckResourceAttr(resourceName, consts.FieldName, name),
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, fmt.Sprintf("%s/data/%s", mount, name)),
					resource.TestCheckResourceAttr(resourceName, "disable_read", "true"),
				),
			},
			{
				PreConfig: func() {
					writeKVData(t, mount, name)
				},
				Config: testKVSecretV2Config_DisableRead(mount, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldMount, mount),
					resource.TestCheckResourceAttr(resourceName, consts.FieldName, name),
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, fmt.Sprintf("%s/data/%s", mount, name)),
					resource.TestCheckResourceAttr(resourceName, "disable_read", "true"),
				),
			},
			{
				PreConfig: func() {
					readKVData(t, mount, name)
				},
				Config: testKVSecretV2Config_DisableRead(mount, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldMount, mount),
					resource.TestCheckResourceAttr(resourceName, consts.FieldName, name),
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, fmt.Sprintf("%s/data/%s", mount, name)),
					resource.TestCheckResourceAttr(resourceName, "disable_read", "true"),
				),
			},
		},
	})
}

func readKVData(t *testing.T, mount, name string) {
	t.Helper()
	client := testProvider.Meta().(*provider.ProviderMeta).GetClient()

	// Read data at path
	path := fmt.Sprintf("%s/data/%s", mount, name)
	resp, err := client.Logical().Read(path)
	if err != nil {
		t.Fatalf(fmt.Sprintf("error reading from Vault; err=%s", err))
	}

	if resp == nil {
		t.Fatalf("empty response")
	}
	if len(resp.Data) == 0 {
		t.Fatalf("kvv2 secret data should not be empty")
	}
	if !reflect.DeepEqual(resp.Data["data"], testKVV2Data) {
		t.Fatalf("kvv2 secret data does not match got: %#+v, want: %#+v", resp.Data["data"], testKVV2Data)
	}
}

func writeKVData(t *testing.T, mount, name string) {
	t.Helper()
	client := testProvider.Meta().(*provider.ProviderMeta).GetClient()

	data := map[string]interface{}{
		consts.FieldData: testKVV2Data,
	}
	// Write data at path
	path := fmt.Sprintf("%s/data/%s", mount, name)
	resp, err := client.Logical().Write(path, data)
	if err != nil {
		t.Fatalf(fmt.Sprintf("error writing to Vault; err=%s", err))
	}

	if resp == nil {
		t.Fatalf("empty response")
	}
}

func mountKVEngine(t *testing.T, mount, name string) {
	t.Helper()
	client := testProvider.Meta().(*provider.ProviderMeta).GetClient()

	err := client.Sys().Mount(mount, &api.MountInput{
		Type:        "kv-v2",
		Description: "Mount for testing KV datasource",
	})
	if err != nil {
		t.Fatalf(fmt.Sprintf("error mounting kvv2 engine; err=%s", err))
	}
}

func testKVSecretV2Config_DisableRead(mount, name string) string {
	return fmt.Sprintf(`
resource "vault_kv_secret_v2" "test" {
  mount        = "%s"
  name         = "%s"
  disable_read = true
  data_json    = jsonencode({})
}`, mount, name)
}

func testKVSecretV2Config_initial(mount, name string) string {
	ret := fmt.Sprintf(`
%s

`, kvV2MountConfig(mount))

	ret += fmt.Sprintf(`
resource "vault_kv_secret_v2" "test" {
  mount               = vault_mount.kvv2.path
  name                = "%s"
  delete_all_versions = true
  data_json = jsonencode(
    {
      zip  = "zap",
      foo  = "bar",
      flag = false
    }
  )
}`, name)

	return ret
}

func testKVSecretV2Config_updated(mount, name string) string {
	ret := fmt.Sprintf(`
%s

`, kvV2MountConfig(mount))

	ret += fmt.Sprintf(`
resource "vault_kv_secret_v2" "test" {
  mount               = vault_mount.kvv2.path
  name                = "%s"
  delete_all_versions = true
  data_json = jsonencode(
    {
      zip  = "zoop",
      foo  = "baz",
      flag = false
    }
  )
  custom_metadata {
    max_versions = 5
    data = {
      extra = "cheese",
      pizza = "please"
    }
  }
}`, name)

	return ret
}
