# tap-fitter

Tap Fitter is a utility designed to aid operator developers to generate the configuration files needed to properly run in certain pipeline environments.

This CLI tool will ingest composite templates and generate the appropriate devfiles.

## Getting Started

To build the utility, run `make build` from the root directory of the tap-fitter project.

## Usage

The `tap-fitter` command will ingest composite templates and generate corresponding devfiles.
Note: because templates may use relative paths to related artifacts, `tap-fitter` should be run in the same location as input templates.
```
$ ./tap-fitter -h
tap-fitter reads a composite template and outputs corresponding devfiles to prepare a repository for a catalog production pipeline.
Note: because templates may use relative paths to related artifacts, tap-fitter should be run in the same location as input templates.

Usage:
  tap-fitter [flags]

Flags:
  -c, --catalog-path string     [REQUIRED] the path/URL to the catalog template used for configuration
  -t, --composite-path string   [REQUIRED] the path to the composite template used for configuration
  -h, --help                    help for tap-fitter
  -p, --provider string         the provider of the catalog (default "tap-fitter")
```

### Quick Start Example

This quick start example will use the composite templates found in the [fbc-composite-example](https://github.com/everettraven/fbc-composite-example/tree/main) repository to generate devfiles from.

- Clone the tap-fitter and example composite template repositories
```
# Clone the tap-fitter repo
git clone https://github.com/release-engineering/tap-fitter.git

# Clone the fbc-composite-example repo
git clone https://github.com/everettraven/fbc-composite-example.git
```
- Navigate to the tap-fitter repository to build the CLI tool
```
# Navigate to the tap-fitter repo
cd tap-fitter

# Build the tap-fitter CLI
make build
```
- Navigate to the destination repo and run the tap-fitter command
```
# Navigate to the fbc-composite-example repo
cd ../fbc-composite-example

# Generate devfiles with tap-fitter
../tap-fitter/tap-fitter --catalog-path catalogs.yaml --composite-path contributions.yaml
```
- Verify the devfiles were correctly generated
```
# Using tree (https://github.com/Old-Man-Programmer/tree) to verify
$ tree catalogs
catalogs
├── v4.10
│   ├── devfile.yaml
│   └── my-package
│       └── catalog.yaml
├── v4.11
│   ├── devfile.yaml
│   └── my-package
│       └── catalog.yaml
├── v4.12
│   ├── devfile.yaml
│   └── my-package
│       └── catalog.yaml
└── v4.13
    ├── devfile.yaml
    └── my-package
        └── catalog.yaml

9 directories, 8 files

```

## Future work
- Enable the use of explicit paths
