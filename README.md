# Intercity Server

Intercity Server is a CLI tool for managing your Intercity instance.
Currently the only thing that it is doing it installing your instance.

## Building the packages
Building the packages is very simple, just run `rake package:linux:x86_64` and it will create
a package for Linux-x86_64 as a tar ball. You can then distribute this tarball.

## Using the package
When you extracted the tarball on an empty server, you go to the folder where
you installed the CLI and run: `./intercity-server install`; This will ask you
for the Hostname and From Email and then starts the installation process for
Intercity.

## License
This project is licensed under the MIT license.
