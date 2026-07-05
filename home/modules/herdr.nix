{ lib, ... }:
{
  xdg.configFile."herdr/config.toml".source = ../../config/herdr/config.toml;

  home.activation.reloadHerdrConfig = lib.hm.dag.entryAfter [ "writeBoundary" ] ''
    if command -v herdr >/dev/null 2>&1; then
      herdr server reload-config >/dev/null 2>&1 || true
    fi
  '';
}
