---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dcloud_inventory_telephony Data Source - terraform-provider-dcloud"
subcategory: ""
description: |-
  All the inventory telephony available to be used in a topology
---

# dcloud_inventory_telephony (Data Source)

All the inventory telephony available to be used in a topology



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `topology_uid` (String)

### Read-Only

- `id` (String) The ID of this resource.
- `inventory_telephony` (List of Object) (see [below for nested schema](#nestedatt--inventory_telephony))

<a id="nestedatt--inventory_telephony"></a>
### Nested Schema for `inventory_telephony`

Read-Only:

- `id` (String)
- `name` (String)
- `description` (String)

