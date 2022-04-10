---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_user"
sidebar_current: "docs-bitbucket-data-user"
description: |-
  Provides a data for a Bitbucket user
---

# bitbucket\_user

Provides a way to fetch data on a user.

## Example Usage

```hcl
data "bitbucket_user" "reviewer" {
  account_id = "gob"
}
```

## Argument Reference

The following arguments are supported:

* `uuid` - (Optional) The UUID that bitbucket users to connect a user to various objects
* `account_id` - (Optional) The user's Atlassian account ID.

## Attributes Reference

* `uuid` - the uuid that bitbucket users to connect a user to various objects
* `display_name` - the display name that the user wants to use for GDPR
* `nickname` - Account name defined by the owner. Note that "nickname" cannot be used in place of "username" in URLs and queries, as "nickname" is not guaranteed to be unique.
* `account_status` - The status of the account.
* `account_id` - The user's Atlassian account ID.
* `is_staff` - is staff user.
