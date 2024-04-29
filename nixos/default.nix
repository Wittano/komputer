{ pkgs, lib, config, ... }:
with lib;
with builtins;
let
  cfg = config.komputer;

  komputer = pkgs.callPackage ./../default.nix { };
in
{
  options = {
    komputer = {
      enable = mkEnableOption "Enable komputer discord bot";
      package = mkOption {
        type = types.package;
        default = komputer;
        description = "komputer package";
      };
      guildID = mkOption {
        type = types.str;
        description = "Discord server id, that you deploy bot";
      };
      applicationID = mkOption {
        type = types.str;
        description = "Application ID for you local version of komputer bot";
      };
      token = mkOption {
        type = types.str;
        description = ''
          Discord token for bot. 
          <REMEMBER!>
          Your token never shouldn't be publish on any public git repository e.g. Github or Gitlab
        '';
      };
    };
  };

  config = mkIf (cfg.enable) {
    assertions = [
      {
        assertion = cfg.token != "";
        message = "Option komputer.token is empty";
      }
      {
        assertion = cfg.guildID != "";
        message = "Option komputer.guildID is empty";
      }
      {
        assertion = cfg.applicationID != "";
        message = "Option komputer.applicationID is empty";
      }
    ];

    systemd.services.komputer = {
      description = "Komputer - Discord bot behave as like 'komputer'. One of character in Star Track parody series created by Dem3000";
      wantedBy = [ "multi-user.target" ];
      environment = {
        DISCORD_BOT_TOKEN = cfg.token;
        APPLICATION_ID = cfg.applicationID;
        SERVER_GUID = cfg.guildID;
      };
      script = "${cfg.package}/bin/komputer";
    };
  };
}

