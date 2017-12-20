## Development on dataform has been stopped due to no need for a separate tool




# Dataform

[![Build Status](https://travis-ci.org/MYOB-Technology/dataform.svg?branch=master)](https://travis-ci.org/MYOB-Technology/dataform)

Easily interact with AWS RDS with a unix like cli interface

## Installation

|         |       |
|---------|-------|
| darwin  | amd64 |
| linux   | amd64 |
| windows | amd64 |

You may choose to download it from the [`Github > Release`](https://github.com/MYOB-Technology/dataform/releases) page or from S3.

```bash
# download edge version from the s3 release bucket
aws s3 cp s3://myob-dataform-release/latest/dfm-linux-amd64/dfm-linux-amd64 /usr/local/bin/dfm

# or get a specific version
aws s3 cp s3://myob-dataform-release/0.1.6/dfm-linux-amd64/dfm-linux-amd64 /usr/local/bin/dfm

# Make it executable
chmod +x /usr/local/bin/dfm

# PROFIT!
dfm -h
```
