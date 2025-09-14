{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    crane.url = "github:ipetkov/crane";
    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    filter.url = "github:numtide/nix-filter";
    steiger.url = "github:brainhivenl/steiger/feat/nix-force-host-platform";
    # steiger.url = "path:/Users/julia/projects/2025/steiger";
  };

  outputs = {
    nixpkgs,
    crane,
    rust-overlay,
    steiger,
    filter,
    ...
  }: let
    systems = ["aarch64-darwin" "x86_64-darwin" "x86_64-linux" "aarch64-linux"];
    overlays = [(import rust-overlay) steiger.overlays.ociTools];
    craneLibFor = pkgs: (crane.mkLib pkgs).overrideToolchain (p: p.rust-bin.stable.latest.default);

    eachSystem = nixpkgs.lib.genAttrs systems;
  in {
    packages = eachSystem (system: let
      pkgs = import nixpkgs {inherit system overlays;};
      craneLib = craneLibFor pkgs;
    in {
      default = pkgs.callPackage ./package.nix {inherit craneLib filter;};
    });

    steigerImages = steiger.lib.eachCrossSystem systems (localSystem: crossSystem: let
      pkgs = import nixpkgs {
        system = localSystem;
        inherit overlays;
      };
      pkgsCross = import nixpkgs {
        inherit localSystem crossSystem overlays;
      };

      craneLib = craneLibFor pkgsCross;
      package = pkgsCross.callPackage ./package.nix {inherit craneLib filter;};

      migrations = pkgs.runCommandNoCC "" {} ''
        mkdir -p $out/data
        cp -r ${./migrations} $out/data/migrations
      '';
    in {
      adapter = pkgs.ociTools.buildImage {
        name = package.pname;
        tag = "latest";
        created = "now";

        copyToRoot = pkgsCross.buildEnv {
          name = "${package.pname}-sysroot";
          paths = [
            package
            pkgs.dockerTools.caCertificates
            pkgsCross.diesel-cli
            migrations
          ];
          pathsToLink = [
            "/bin"
            "/etc"
            "/data"
          ];
        };

        config.Cmd = ["/bin/${package.pname}"];
        compressor = "none";
      };
    });

    devShells = eachSystem (system: let
      pkgs = import nixpkgs {inherit system overlays;};
      craneLib = craneLibFor pkgs;
    in {
      default = craneLib.devShell {
        packages = let
          toolchain = pkgs.rust-bin.stable.latest.default.override {
            extensions = ["rust-src" "rustfmt"];
          };
        in
          with pkgs;
          with toolchain; [
            rust-analyzer
            clippy
            valkey
            diesel-cli
            nix-eval-jobs
            steiger.packages.${system}.default
          ];
      };
    });
  };
}
