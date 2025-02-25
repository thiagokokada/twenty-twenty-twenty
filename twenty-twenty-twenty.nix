{
  lib,
  stdenv,
  alsa-lib,
  buildGo124Module,
  darwin,
  pkg-config,
  version ? "unknown",
  withSound ? true,
  withStatic ? false,
  withSystray ? true,
}:

# Darwin builds always have sound since it doesn't depend in CGO, and darwin
# builds always depends on CGO anyway because gioui
assert stdenv.isDarwin -> withSound;
# No sound builds are always static
assert withStatic -> withSound;

buildGo124Module {
  pname = "twenty-twenty-twenty";
  inherit version;
  src = lib.cleanSource ./.;
  vendorHash = "sha256-rZoL4WQIqOwmEnopo05AksKuIKAcE3aCzlr9D5hmoz4=";

  env.CGO_ENABLED = if withSound then "1" else "0";

  nativeBuildInputs =
    lib.optionals (withSound && stdenv.hostPlatform.isLinux) [ pkg-config ]
    ++ lib.optionals stdenv.hostPlatform.isDarwin [ darwin.autoSignDarwinBinariesHook ];

  buildInputs = lib.optionals (withSound && stdenv.hostPlatform.isLinux) [ alsa-lib ];

  preBuild = lib.optionalString stdenv.isDarwin ''
    export MACOSX_DEPLOYMENT_TARGET=11.0
  '';

  preFixup = lib.optionalString stdenv.isDarwin ''
    OUT_APP="$out/Applications/TwentyTwentyTwenty.app"
    mkdir -p "$OUT_APP/Contents/MacOS"
    cp -r assets/macos/TwentyTwentyTwenty.app/* "$OUT_APP"
    cp $out/bin/twenty-twenty-twenty "$OUT_APP/Contents/MacOS/TwentyTwentyTwenty"
  '';

  # Disable tests that can't run in CI, since they require desktop environment
  preCheck = ''
    export CI=1
  '';

  ldflags =
    [
      "-X=main.version=${version}"
      "-s"
      "-w"
    ]
    ++ lib.optionals withStatic [
      "-linkmode external"
      ''-extldflags "-static"''
    ]
    ++ lib.optionals (!withSystray) [ "-tags=nosystray" ];

  meta = with lib; {
    description = "Alerts every 20 minutes to look something at 20 feet away for 20 seconds";
    homepage = "https://github.com/thiagokokada/twenty-twenty-twenty";
    license = licenses.mit;
    mainProgram = "twenty-twenty-twenty";
  };
}
