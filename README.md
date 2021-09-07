# StaticSync

A [blox](https://github.com/cueblox/blox) postbuild plugin to sync static files to your CDN. 

StaticSync exposes two different plugins: 
* 'staticsync' synchronizes your blox *static_dir* to the bucket you specify after blox completes a build.
* 'datasync' synchronizes your blox *build_dir* to the bucket you specify after blox completes a build.

## Getting Started

StaticSync is written in Go. You'll need a local Go installation to use it, so follow the [instructions](https://go.dev/learn/) at [go.dev](https://go.dev/) to install Go on your computer.

You'll also need an account at a CDN provider supported by the [GoCloud SDK](https://gocloud.dev/). See the [documentation](https://gocloud.dev/howto/blob/#services) for specific information on how to choose the correct URL to export as an environment variable. See usage below for more details.

### Prerequisites

The things you need before installing the software.

* Go version 1.17 or later.
* CDN provider supported by the GoCloud SDK.


### Installation

A step by step guide that will tell you how to get the development environment up and running.

```
$ go get github.com/cueblox/staticsync
```

## Usage

Add the plugins you want to use to your `blox.cue` configuration file as a postbuild plugin.

For staticsync, add the following to your `blox.cue` configuration file:
```
{
  data_dir: "data"
  schemata_dir: "data/schemata"
  build_dir: ".build"
  template_dir: "data/tpl"
  static_dir: "public/static"
  postbuild: [ {
    name: "staticsync"
    executable: "staticsync"
  }]
}
```

Or you can use both plugins:
```
{
  data_dir: "data"
  schemata_dir: "data/schemata"
  build_dir: ".build"
  template_dir: "data/tpl"
  static_dir: "public/static"
  prebuild: [ {
    name: "bloximages"
    executable: "bloximages"
  }]
  postbuild: [ {
    name: "staticsync"
    executable: "staticsync"
  },
  {
    name: "datasync"
    executable: "staticsync"
  }]
}
```

## Configuration

In addition to the environment variables required by GoCloud, you'll need to export the following environment variables:

* `DATA_BUCKET`: The bucket to sync your built data to when using the datasync plugin
* `STATIC_BUCKET`: The bucket to sync your static files to when using the staticsync plugin

An example .envrc file:

```
export AZURE_STORAGE_ACCOUNT=yourstorageaccount
export AZURE_STORAGE_KEY=SomeBigLongKeyInBase64==
export STATIC_BUCKET=azblob://images
export DATA_BUCKET=azblob://data
```

## Deployment

StaticSync uploads your static content at the end of your `blox build`. If you run `blox build` locally, it will synchronize to your CDN. 

## Contributing

See [contributing guide](CONTRIBUTING.md) for more information.