{
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      lastModifiedDate =
        self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;
      supportedSystems =
        [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

    in {

      packages = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};

        in {
          default = pkgs.buildGoModule {
            pname = "lightspeed-dhl-adapter";
            inherit version;
            src = ../.;

            vendorHash = "sha256-23QtnwIv4tivXGsZAYaTjB/M+30WO1tDFjTkN8OMBzw=";
            # vendorHash = pkgs.lib.fakeHash;
          };
        });

      devShells = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools ];
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.default);
    };
}