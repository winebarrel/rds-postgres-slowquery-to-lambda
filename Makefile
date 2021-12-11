REPOSITORY_URL=$(shell tfstate-lookup -s ./terraform/terraform.tfstate aws_ecr_repository.postgresql_slowquery.repository_url)

.PHONY: update
update: push
	aws lambda update-function-code --function-name postgresql-slowquery --image-uri $(REPOSITORY_URL):latest

.PHONY: image
image:
	docker build -t $(REPOSITORY_URL):latest .

.PHONY: push
push:
	docker push $(REPOSITORY_URL):latest
