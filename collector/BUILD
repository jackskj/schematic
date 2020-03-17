load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_push", "container_image")

go_library(
    name = "go_default_library",
    srcs = [
        "collector.go",
        "main.go",
        "producer.go",
        "schema.go",
    ],
    importpath = "github.com/jackskj/schematic/collector",
    visibility = ["//visibility:private"],
    deps = [
        "//collector/pb:go_default_library",
        "//collector/schemas/pb:go_default_library",
        "//stresstest/stress:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//protoc-gen-go/descriptor:go_default_library",
        "@com_github_shopify_sarama//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_image(
    name = "collector",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

container_image(
    name = "collector_container",
    base = ":collector",
    env = {
        "MODE": "DOCKER",
    },
)

container_push(
    name = "collector_image",
    image = ":collector_container",
    format = "Docker",
    registry = "index.docker.io",
    repository = "jackskj/schematic/collector",
    tag = "latest",
)
