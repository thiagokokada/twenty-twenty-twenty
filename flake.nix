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
      apps.default = forAllSystems (system: {
        type = "app";
        program = "${self.packages.${system}.twenty-twenty-twenty}/bin/twenty-twenty-twenty";
      });

      packages = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          default = self.packages.${system}.twenty-twenty-twenty;
          twenty-twenty-twenty = pkgs.buildGoModule {
            pname = "twenty-twenty-twenty";
            inherit version;
            src = ./.;
            vendorHash = "sha256-3RtdnS4J7JbdU+jMTEzClSlDDPh6bWqbjchvrtS8HUc";

            nativeBuildInputs = with pkgs; lib.optionals stdenv.hostPlatform.isLinux [
              pkg-config
            ];

            buildInputs = with pkgs;
              lib.optionals stdenv.hostPlatform.isLinux [
                alsa-lib
              ] ++
              lib.optionals stdenv.hostPlatform.isDarwin [
                darwin.apple_sdk_11_0.frameworks.MetalKit
                darwin.apple_sdk_11_0.frameworks.UserNotifications
              ];

            # Tests are mostly useful for development, not to ensure that
            # program is running correctly.
            doCheck = false;

            ldflags = [ "-X=main.version=${version}" ];

            meta = with pkgs.lib; {
              description = "Alerts every 20 minutes to look something at 20 feet away for 20 seconds";
              homepage = "https://github.com/thiagokokada/twenty-twenty-twenty";
              license = licenses.mit;
              mainProgram = "twenty-twenty-twenty";
            };
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
            ]
            ++ lib.optionals stdenv.hostPlatform.isLinux [
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
