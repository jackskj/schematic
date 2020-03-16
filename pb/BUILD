load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "pb_proto",
    srcs = [
        "biometrics.proto",
        "collector.proto",
        "schematic.proto",
    ],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "pb_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/jackskj/schematic/pb",
    proto = ":pb_proto",
    visibility = ["//visibility:public"],
)
