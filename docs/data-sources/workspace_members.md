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
data "bitbucket_workspace" "example" {
  workspace = "gob"
}

data "bitbucket_workspace_members" "example" {
  uuid = data.bitbucket_workspace_members.example.uuid
}
```

## Argument Reference

The following arguments are supported:

* `uuid` - (Required) This is the workspace UUID surrounded by curly-braces

## Attributes Reference

* `members` - A set of string containing the member uuids
* `id` - The workspace's immutable id.
