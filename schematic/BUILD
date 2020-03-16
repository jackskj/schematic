load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "main.go",
        "schematic.go",
    ],
    importpath = "github.com/jackskj/schematic/schematic",
    visibility = ["//visibility:private"],
    deps = [
        "//schematic/pb:go_default_library",
        "@com_github_shopify_sarama//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
    ],
)

go_binary(
    name = "schematic",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
