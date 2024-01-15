{
  description = "twenty-twenty-twenty";

  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
    flake-compat.url = "github:edolstra/flake-compat";
  };

  outputs = { self, nixpkgs, ... }:
    let
      version = "nix-${self.shortRev or self.dirtyShortRev or "unknown-dirty"}";

      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      apps = forAllSystems (system: {
        default = {
          type = "app";
          program = nixpkgs.lib.getExe self.packages.${system}.twenty-twenty-twenty;
        };
      });

      packages = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          default = self.packages.${system}.twenty-twenty-twenty;
          twenty-twenty-twenty = pkgs.callPackage ./twenty-twenty-twenty.nix { inherit version; };
          twenty-twenty-twenty-static = self.packages.${system}.twenty-twenty-twenty.override rec {
            inherit (pkgs.pkgsStatic) alsa-lib stdenv;
            buildGoModule = pkgs.buildGoModule.override { inherit stdenv; };
            extraLdflags = [ "-linkmode external" ''-extldflags "-static"'' ];
          };
          # Also static build because CGO_ENABLED=0
          twenty-twenty-twenty-no-sound = self.packages.${system}.twenty-twenty-twenty.override {
            inherit version;
            withSound = false;
          };
        });

      devShells = forAllSystems (system:
        let pkgs = nixpkgsFor.${system}; in
        {
          default = pkgs.mkShell {
            name = "twenty-twenty-twenty";

            packages = with pkgs; [
              gnumake
              go
              gopls
            ] ++
            lib.optionals stdenv.hostPlatform.isLinux [
              alsa-lib
              gcc
              pkg-config
            ];

            # Keep the current user shell (e.g.: zsh instead of bash)
            shellHook = "exec $SHELL";
          };
        }
      );
    };
}
