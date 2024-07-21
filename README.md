# cobra-gen

`cobra-gen` is a boilerplate code generator for the [Cobra](https://github.com/spf13/cobra) CLI framework in Go.

## Dependencies

- Go
- gofmt
- goimports: Install with `go install golang.org/x/tools/cmd/goimports@latest`

## Installation

To install `cobra-gen`, run:

```sh
go install github.com/umutbasal/cobra-gen@latest
```

## Usage

### Example: Designing a Simple AWS CLI

Let's create a basic AWS CLI with `cobra-gen` that handles `s3` and `ec2` commands:

- `aws --profile`
- `aws s3 ls mybucket --page-size 10`
- `aws ec2 create-instance`

1. **Generate Config**

   Convert your commands to the following format and run them in your Go project:

   ```sh
   cobra-gen --profile
   cobra-gen s3 ls +mybucket --page-size
   cobra-gen ec2 create-instance
   ```

   This will generate a YAML configuration file named `.cobra-gen.yaml`. You can modify this file as needed:

   ```yaml
   cmd:
   - --profile
   - s3:
     - ls:
       - +mybucket # `+` denotes arguments
       - --page-size # `--` denotes flags; do not add values for flags
   - ec2:
     - create-instance
   ```

2. **Generate Code**

   Run `cobra-gen` without any arguments to generate the boilerplate code:

   ```sh
   cobra-gen
   ```

   This will create the following files:

   - `./cmd/cmd.go`
   - `./cmd/s3/ls.go`
   - `./cmd/s3/s3.go`
   - `./cmd/ec2/create-instance.go`
   - `./cmd/ec2/ec2.go`

  See the [examples/aws](examples/aws) directory for the generated code.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
