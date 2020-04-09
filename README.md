# Dotscience Terraform Provider

A terraform provider to manage dotscience resources.

## installation

You can install the provider by downloaded the latest binary for your platform from the [github releases page](https://github.com/dotmesh-io/terraform-provider-dotscience/releases).

This file must be placed in the `~/.terraform/plugins` directory with the name: `terraform-provider-dotscience_v<VERSION>`

So - for version `0.0.1` - there should be the following file present on your machine:

```
~/.terraform/plugins/terraform-provider-dotscience_v0.0.1
```

Once you have installed the provider - you will be able to run `terraform init` on a project that uses it.

## usage

The provider requires 2 arguments:

 * `hub_public_url` - the URL of the dotscience hub
 * `hub_admin_password` - the admin password for that user

```
provider "dotscience" {
  hub_public_url = var.hub_public_url
  hub_admin_password = var.hub_admin_password
}
```

### runners

The provider can track managed runners and delete them before the rest of the stack is deleted.

Once you have configured your provider - add a `dotscience_runners` resource to your stack:

```
resource "dotscience_runners" "hub-runners" {
  depends_on = []
}
```

It is important that the `depends_on` lists all other resources needed by the hub to function.

Here is an example of it being used in our [AWS terraform stack](https://github.com/dotmesh-io/dotscience-tf):

`dotscience-aws.tf`:

```
provider "dotscience" {
  hub_public_url = var.hub_public_url
  hub_admin_password = var.hub_admin_password
}

resource "dotscience_runners" "hub-runners" {
  depends_on = [
    aws_instance.ds_hub,
    aws_eip_association.eip_assoc,
    aws_eip.ds_eip,
    module.vpc
  ]
}
```

