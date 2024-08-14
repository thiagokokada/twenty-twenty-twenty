{
  description = "twenty-twenty-twenty";

  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
    flake-compat.url = "github:edolstra/flake-compat";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      version = "nix-${self.shortRev or self.dirtyShortRev or "unknown-dirty"}";

      supportedSystems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-linux"
        "aarch64-darwin"
      ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      apps = forAllSystems (system: {
        default = {
          type = "app";
          program =
            let
              inherit (self.packages.${system}) twenty-twenty-twenty;
            in
            if (system == "aarch64-darwin" || system == "x86_64-darwin") then
              "${twenty-twenty-twenty}/Applications/TwentyTwentyTwenty.app/Contents/MacOS/TwentyTwentyTwenty"
            else
              nixpkgs.lib.getExe twenty-twenty-twenty;
        };
      });

      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = self.packages.${system}.twenty-twenty-twenty;
          twenty-twenty-twenty = pkgs.callPackage ./twenty-twenty-twenty.nix { inherit version; };
          twenty-twenty-twenty-no-sound = self.packages.${system}.twenty-twenty-twenty.override {
            withSound = false;
          };
          twenty-twenty-twenty-no-systray = self.packages.${system}.twenty-twenty-twenty.override {
            withSystray = false;
          };
          twenty-twenty-twenty-minimal = self.packages.${system}.twenty-twenty-twenty.override {
            withSound = false;
            withSystray = false;
          };
          twenty-twenty-twenty-static = pkgs.pkgsStatic.callPackage ./twenty-twenty-twenty.nix {
            inherit version;
            withStatic = true;
          };
          twenty-twenty-twenty-aarch64-linux-static =
            pkgs.pkgsCross.aarch64-multiplatform.pkgsStatic.callPackage ./twenty-twenty-twenty.nix
              {
                inherit version;
                withStatic = true;
              };
        }
      );

      formatter = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        pkgs.nixfmt-rfc-style
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            name = "twenty-twenty-twenty";

            packages =
              with pkgs;
              [
                gnumake
                go
                gopls
              ]
              ++ lib.optionals stdenv.hostPlatform.isLinux [
                alsa-lib
                gcc
                pkg-config
              ]
              ++ lib.optionals stdenv.hostPlatform.isDarwin [
                darwin.apple_sdk_11_0.frameworks.Cocoa
                darwin.apple_sdk_11_0.frameworks.MetalKit
                darwin.apple_sdk_11_0.frameworks.UserNotifications
              ];

            # Keep the current user shell (e.g.: zsh instead of bash)
            shellHook = "exec $SHELL";
          };
        }
      );
    };
}
