load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")
load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies")
load("@rules_java//java:repositories.bzl", "rules_java_dependencies", "rules_java_toolchains")

def dept_imports():
    protobuf_deps()
    go_rules_dependencies()
    go_register_toolchains()
    gazelle_dependencies()
    rules_proto_dependencies()
    rules_java_dependencies()
    rules_java_toolchains()
