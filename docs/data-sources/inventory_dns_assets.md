---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dcloud_inventory_dns_assets Data Source - terraform-provider-dcloud"
subcategory: ""
description: |-
  All the Inventory DNS assets in a given topology
---

# dcloud_inventory_dns_assets (Data Source)

All the Inventory DNS assets in a given topology



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `topology_uid` (String)

### Read-Only

- `id` (String) The ID of this resource.
- `inventory_dns_assets` (List of Object) (see [below for nested schema](#nestedatt--inventory_dns_assets))

<a id="nestedatt--inventory_dns_assets"></a>
### Nested Schema for `inventory_dns_assets`

Read-Only:

- `id` (String)
- `name` (String)


