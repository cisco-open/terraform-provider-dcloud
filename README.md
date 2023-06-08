# Terraform Provider dCloud Topology Builder

This Terraform provider allows you to manage Topologies, Networks, VMs, and Hardware using the dCloud Topology Builder API. With this provider, you can use Terraform to create, update, and delete Topologies and their associated resources in dCloud.

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.18

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Using the provider

The provider requires an Cisco SSO Token to be set in the environment

`export TB_AUTH_TOKEN=<YOUR_AUTH_TOKEN>`

See the [docs](/docs) and [examples](/examples) directory for examples.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

## Dependencies

The provider makes extensive use of the dCloud Topology Builder Go Client (https://github.com/cisco-open/dcloud-tb-go-client).

Before adding new resource types to the provider the dCloud Topology Builder Go Client should be updated to expose the required resources and released.
