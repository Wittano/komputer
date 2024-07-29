{ lib
, buildGoModule
, gcc
, pkg-config
, libopus
, ffmpeg
, opusfile
}: buildGoModule {
  name = "komputer";
  version = "v1.2.1";

  src = ./.;

  vendorHash = "sha256-B/kII44/cuzGwO/pWCamFl7clHbz/qon4YgtUHYWV30=";

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

