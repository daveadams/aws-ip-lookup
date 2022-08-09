# aws-ip-lookup

A command line tool that uses the
[AWS IP address ranges file](https://docs.aws.amazon.com/general/latest/gr/aws-ip-ranges.html)
to show the region and service of each given IP address.

## Usage

If you provide one or more IP addresses as command line arguments,
`aws-ip-lookup` will print one line for each address giving the matched CIDR
range, the AWS region, and the given service provided in the `ip-ranges.json`
file.

    $ aws-ip-lookup 52.203.233.206
    52.203.233.206 52.200.0.0/13 us-east-1 EC2

    $ aws-ip-lookup 205.251.249.88 52.0.247.15 52.61.40.104 52.83.25.165 13.245.166.133
    205.251.249.88 205.251.249.0/24 GLOBAL CLOUDFRONT
    52.0.247.15 52.0.0.0/15 us-east-1 EC2
    52.83.25.165 52.83.25.160/27 cn-northwest-1 CLOUD9
    13.245.166.133 13.245.166.132/30 af-south-1 AMAZON_APPFLOW

If you do not provide any command line arguments, `aws-ip-lookup` will read
IP addresses in one per line.

    $ echo 52.61.40.104 | aws-ip-lookup
    52.61.40.104 52.61.40.104/29 us-gov-west-1 CODEBUILD

## Installation

    $ go install github.com/daveadams/aws-ip-lookup/cmd/aws-ip-lookup

## Implementation Details

The `ip-ranges.json` file is over a megabyte in size, so this program only
fetches a new copy once every 24 hours. The cached copy is stored in the local
user's cache directory as appropriate on each platform:

* Unix/Linux: `~/.cache/aws-ip-lookup/ip-ranges.json`
* macOS: `~/Library/Caches/aws-ip-lookup/ip-ranges.json`
* Windows: `%AppData%\aws-ip-lookup\ip-ranges.json`

## Possible Future Work

* IPv6 support
* Cache info/management subcommands
* Help messages
* Parsing multiple IP addresses per line when reading input from stdin

## License

This software is public domain. No rights are reserved. See LICENSE for more
information.
