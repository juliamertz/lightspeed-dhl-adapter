{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
  };

  outputs = {
    nixpkgs,
    systems,
    ...
  }: let
    forAllSystems = function: nixpkgs.lib.genAttrs (import systems) (system: function nixpkgs.legacyPackages.${system});
  in {
    packages = forAllSystems (pkgs: {
      default = pkgs.buildGoModule {
        pname = "lightspeed-dhl-adapter";
        version = "0.1.1";
        src = ./src;

        vendorHash = "sha256-B/NtKPAOHRcVq1VBK/L/kCQ04Fvyitfpt5c3M273I8M=";
        meta.mainProgram = "lightspeed-dhl";

        GO_TEST = "none";
      };
    });

    devShells = forAllSystems (pkgs: {
      default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          gotools
          go-tools
          golangci-lint
        ];
      };
    });
  };
}
