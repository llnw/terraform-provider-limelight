<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Terraform Provider for Limelight
=============================


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.14 (to build the provider plugin)

Building The Provider
---------------------

Clone repository:

```sh
$ git clone git@github.com:llnw/terraform-provider-limelight.git
```

Enter the provider directory and build the provider:

```sh
$ cd terraform-provider-limelight
$ make build
```

Using the provider
----------------------

See the Limelight Provider for Terraform docs for usage details.

# Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.14+ is *required*).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-limelight
...
```

Before delivering any new code, ensure you run `golint` and address any errors.

### Running tests

The following env variables must be set prior to running the acceptance tests:

```sh
LLNW_API_USERNAME="{MY_LLNW_USERNAME}"
LLNW_API_KEY="{MY_LLNW_APIKEY}"
LLNW_TEST_SHORTNAME="{MY_LLNW_SHORTNAME}"
```

If you are using non-default API endpoints, you can additionally set the following:

```sh
LLNW_EDGEFUNCTIONS_API_URL="{MY_NONDEFAULT_EDGEFUNCTION_API_URL}"
LLNW_CONFIG_API_URL="{MY_NONDEFAULT_CONFIG_API_URL}"
```

To run the acceptance tests, run the `testacc` make target:

```sh
$ make testacc
```

If you want to run against a specific set of tests, run make `testacc` with the `TESTARGS` parameter containing 
the run mask. For example:

```sh
make testacc TESTARGS="-run=TestAccResourceLimelightEdgeFunction"
```