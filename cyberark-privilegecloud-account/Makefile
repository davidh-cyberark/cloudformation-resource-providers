.PHONY: build test clean Makefile

build:
	go env -w CGO_ENABLED=0
	make -f makebuild  # this runs build steps required by the cfn cli
	go env -w CGO_ENABLED=1

test:
	cfn generate
	go env -w CGO_ENABLED=0
	env GOOS=linux go build -ldflags="-s -w" -tags="lambda.norpc,$(TAGS)" -o bin/bootstrap cmd/main.go
	go env -w CGO_ENABLED=1

clean::
	rm -rf bin

AWS_REGION ?= $(error ERROR: set AWS_REGION variable in Makefile.local)
AWS_LAMBDA_ARN ?= $(error ERROR: set AWS_LAMBDA_ARN variable in Makefile.local)

register: build
	cfn submit --set-default --role-arn "$(AWS_LAMBDA_ARN)" --region "$(AWS_REGION)"

-include Makefile.local