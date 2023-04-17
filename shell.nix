{ pkgs ? import <nixpkgs> {}}:

pkgs.mkShell {
  nativeBuildInputs = [
    pkgs.ctlptl
    pkgs.delve
    pkgs.gdlv
    pkgs.go
    pkgs.gopls
    pkgs.kind
    pkgs.kubernetes-helm
    pkgs.operator-sdk
    pkgs.podman
    pkgs.tilt
    pkgs.vault
  ];
}
