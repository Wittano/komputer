{ lib
, buildGoModule
, gcc
, pkg-config
, libopus
, ffmpeg
, opusfile
}: buildGoModule {
  name = "komputer";
  version = "v1.2.0";

  src = ./.;

  vendorHash = "sha256-CThNuZ16b8SXxJAtCkDMm+mwCqaS5zrr+PbX+5N3GCc=";

  CGO_ENABLED = 1;
  proxyVendor = true;

  nativeBuildInputs = [ gcc pkg-config libopus ];
  propagatedBuildInputs = [ ffmpeg opusfile ];

  meta = with lib; {
    homepage = "https://github.com/Wittano/komputer";
    description = "Discord bot behave as like 'komputer'. One of character in Star Track parody series created by Dem3000";
    license = licenses.gpl3;
    maintainers = with maintainers; [ Wittano ];
    platforms = platforms.linux;
  };
}

