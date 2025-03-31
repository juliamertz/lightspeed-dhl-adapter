{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
  };

  outputs =
    {
      nixpkgs,
      systems,
      ...
    }:
    let
      forAllSystems =
        function: nixpkgs.lib.genAttrs (import systems) (system: function nixpkgs.legacyPackages.${system});
    in
    {
      packages = forAllSystems (pkgs: {
        default = pkgs.buildGoModule {
          pname = "lightspeed-dhl-adapter";
          version = "0.1.0";
          src = ../.;

          vendorHash = "sha256-o2SNdqIx+YvpKh883rowk9/IlNnpSiutgvc29CAWKj4=";
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
