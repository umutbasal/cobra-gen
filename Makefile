DIR := examples/aws
MODULE := awsexample

all: setup generate

setup:
	rm -rf $(DIR)
	mkdir -p $(DIR)
	cd $(DIR) && go mod init $(MODULE)

generate:
	cd $(DIR) && cobra-gen --profile
	cd $(DIR) && cobra-gen s3 ls +mybucket --page-size
	cd $(DIR) && cobra-gen ec2 create-instance
	cd $(DIR) && cobra-gen
	cd $(DIR) && go mod tidy
	cp $(DIR)/examples/cobra-gen/main.go $(DIR)/main.go