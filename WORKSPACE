workspace(name = "jackskj_schematic")

load("//bazel:http_dependencies.bzl", "http_depts")

http_depts()

load("@io_bazel_rules_docker//repositories:repositories.bzl", container_repositories = "repositories")

container_repositories()

load("//bazel:imports.bzl", "dept_imports")

dept_imports()

load("@io_bazel_rules_docker//go:image.bzl", go_image_repos = "repositories")

go_image_repos()

load("//bazel:go_repositories.bzl", "go_repositories")

go_repositories()
