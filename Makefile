gazelle:
	# generates build files
	bazel run //:gazelle_update

repos:
	# generates go repos
	bazel run //:gazelle_update -- update-repos -from_file=collector/go.mod -from_file=stresstest/go.mod -to_macro=bazel/go_repositories.bzl%go_repositories

fix:
	# fixes deprecated usage of rules
	bazel run //:gazelle_update -- fix

build:
	bazel build //collector:collector
	bazel build //stresstest:stresstest

docker:
	# build collector and stresstest images for testing
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64  //stresstest:stresstest_container
	bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //collector:collector_container
	# running container push with docker.io is giving me issues
	# running tag and push manually
	docker tag bazel/stresstest:stresstest_container jackskj/stresstest:stresstest_container
	docker tag bazel/collector:collector_container jackskj/collector:collector_container
	docker push jackskj/stresstest:stresstest_container
	docker push jackskj/collector:collector_container

helm:
	helm install zookeeper incubator/zookeeper -f k8s/helm_values_zookeeper.yaml
	helm install kafka incubator/kafka -f k8s/helm_values_kafka.yaml

unhelm:
	helm delete zookeeper kafka

manifest:
	kubectl apply -f k8s/manifest.yaml

unmanifest:
	kubectl delete deployment collector stresstest
	kubectl delete service collector

pb:
	# TODO, change protob-gen-go to second major version 
	# BIOMETRICS
	protoc  --go_out=collector/schemas \
		--swift_out=watchExporter \
		./pb/biometrics.proto
	python3 -m grpc_tools.protoc -I=. \
		--python_out=consumerPy \
		--proto_path=consumerPy \
		./pb/biometrics.proto
	\
	# COLLECTOR
	protoc  --go_out=plugins=grpc:collector  \
		--go_out=plugins=grpc:stresstest  \
		--swift_out=watchExporter \
		--grpc-swift_out=watchExporter \
		./pb/collector.proto
	\
	# SCHEMATIC
	protoc  --go_out=plugins=grpc:collector \
		--go_out=plugins=grpc:consumerGo \
		--go_out=plugins=grpc:schematic \
		./pb/schematic.proto
	python3 -m grpc_tools.protoc -I=. \
		--python_out=consumerPy \
		--grpc_python_out=consumerPy \
		--proto_path=consumerPy \
		./pb/schematic.proto

up:
	docker-compose  up -d --build
down:
	docker-compose down

.PHONY: pb up down gazelle fix repos build docker helm unhelm manifest unmanifest
