# cobra-gen

cobra-gen is boilerplate code generator for [Cobra](https://github.com/spf13/cobra)

## Deps

- go
- gofmt
- goimports `go install golang.org/x/tools/cmd/goimports@latest`

## Install

```sh
go install github.com/umutbasal/cobra-gen@latest
```

## Usage

Lets design simple aws cli with cobra-gen
we will handle below s3 and ec2 commands

- `aws --profile`
- `aws s3 ls mybucket --page-size 10`
- `aws ec2 create-instance`

Generate config by simply converting commands to below format and run them in your fresh go project.

```sh
cobra-gen --profile
cobra-gen s3 ls +mybucket --page-size
cobra-gen ec2 create-instance
```

This will generate a yaml config in `.cobra-gen.yaml`, you can edit it as you like with same format.

```yaml
cmd:
- --profile
- s3:
  - ls:
    - +mybucket # + is for arguments
    - --page-size # -- is for flags, don't add value for them
- ec2:
  - create-instance
```

Then generate code by running `cobra-gen` without any arguments.

This will generate boilerplate code for you with this files.

- ./cmd/cmd.go
- ./cmd/s3/ls.go
- ./cmd/s3/s3.go
- ./cmd/ec2/create-instance.go
- ./cmd/ec2/ec2.go
