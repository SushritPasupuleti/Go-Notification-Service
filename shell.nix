{ pkgs ? import <nixpkgs> { } }:
with pkgs;
mkShell {
  nativeBuildInputs = [
    # rustc
    # cargo
    # openssl.dev
    # sqlx-cli
    pkg-config
    go
    golangci-lint
    cloc
    wrk2
    # haskellPackages.postgrest
  ];
}
