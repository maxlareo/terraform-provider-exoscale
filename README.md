Terraform Provider
==================

- Website: https://www.terraform.io
- [![Actions Status](https://github.com/exoscale/terraform-provider-exoscale/workflows/CI/badge.svg)](https://github.com/exoscale/terraform-provider-exoscale/actions?query=workflow%3ACI)
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.13+ (to build the provider plugin)

Building The Provider
---------------------

Clone repository:

```sh
$ git clone https://github.com/exoscale/terraform-provider-exoscale
```

Enter the provider directory and build the provider

```sh
$ cd /path/to/terraform-provider-exoscale
$ make build
$ make install
```

Using the provider
----------------------
If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) After placing it into your plugins directory,  run `terraform init` to initialize it.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-exoscale
...
```

In order to test the provider, you can simply run `make test`.

*Note:* Make sure no `CLOUDSTACK_KEY` or `CLOUDSTACK_SECRET` variables are set, and there's no `[cloudstack]` section in the CloudStack credentials file `~/.cloudstack.ini`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make test-acc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make test-acc
```

In order to test a specific part of the acceptance test suite, you may run:

``` sh
make GO_TEST_EXTRA_ARGS="-v -run ^TestAcc..." test-acc
```
