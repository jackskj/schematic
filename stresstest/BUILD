load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_push", "container_image")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/jackskj/schematic/stresstest",
    visibility = ["//visibility:private"],
    deps = [
        "//stresstest/pb:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_image(
    name = "stresstest",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

container_image(
    name = "stresstest_container",
    base = ":stresstest",
    env = {
        "MODE": "DOCKER",
    },
)

container_push(
    name = "stresstest_image",
    image = ":stresstest_container",
    format = "Docker",
    registry = "index.docker.io",
    repository = "jackskj/schematic/stresstest",
    tag = "latest",
)
