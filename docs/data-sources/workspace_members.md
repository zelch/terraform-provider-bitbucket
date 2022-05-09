---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_workspace_members"
sidebar_current: "docs-bitbucket-data-workspace-members"
description: |-
  Provides a data for a Bitbucket workspace members
---

# bitbucket\_workspace\_members

Provides a way to fetch data on a the members of a workspace.

## Example Usage

```hcl
data "bitbucket_workspace_members" "example" {
  workspace = "gob"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) This can either be the workspace ID (slug) or the workspace UUID surrounded by curly-braces.

## Attributes Reference

* `members` - A set of string containing the member UUIDs.
* `id` - The workspace's immutable id.
