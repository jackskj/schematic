load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "consumer.go",
        "dynamic_record.go",
        "main.go",
    ],
    importpath = "github.com/jackskj/schematic/consumerGo",
    visibility = ["//visibility:private"],
    deps = [
        "//consumerGo/pb:go_default_library",
        "@com_github_golang_protobuf//protoc-gen-go/descriptor:go_default_library",
        "@com_github_shopify_sarama//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_google_protobuf//proto:go_default_library",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
        "@org_golang_google_protobuf//types/dynamicpb:go_default_library",
    ],
)

go_binary(
    name = "consumerGo",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
