{ lib
, stdenv
, alsa-lib
, buildGoModule
, darwin
, pkg-config
, version ? "unknown"
, withSound ? true
, withStatic ? false
}:

# Darwin builds always have sound since it doesn't depend in CGO, and darwin
# builds always depends on CGO anyway because gioui
assert stdenv.isDarwin -> withSound;
# No sound builds are always static
assert withStatic -> withSound;

buildGoModule {
  pname = "twenty-twenty-twenty";
  inherit version;
  src = lib.cleanSource ./.;
  vendorHash = "sha256-3RtdnS4J7JbdU+jMTEzClSlDDPh6bWqbjchvrtS8HUc";

  CGO_ENABLED = if withSound then "1" else "0";

  nativeBuildInputs = lib.optionals (withSound && stdenv.hostPlatform.isLinux) [
    pkg-config
  ];

  buildInputs = lib.optionals (withSound && stdenv.hostPlatform.isLinux) [
    alsa-lib
  ] ++
  lib.optionals stdenv.hostPlatform.isDarwin [
    darwin.apple_sdk_11_0.frameworks.MetalKit
    darwin.apple_sdk_11_0.frameworks.UserNotifications
  ];

  # Tests are mostly useful for development, not to ensure that
  # program is running correctly.
  doCheck = false;

  ldflags = [ "-X=main.version=${version}" "-s" "-w" ]
    ++ lib.optionals withStatic [ "-linkmode external" ''-extldflags "-static"'' ];

  meta = with lib; {
    description = "Alerts every 20 minutes to look something at 20 feet away for 20 seconds";
    homepage = "https://github.com/thiagokokada/twenty-twenty-twenty";
    license = licenses.mit;
    mainProgram = "twenty-twenty-twenty";
  };
}
