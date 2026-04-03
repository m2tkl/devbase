{ lib, ... }:
let
  localGitConfigPath = ".config/devbase/git/local.gitconfig";
in
{
  home.file.".gitignore_global".source = ../../config/git/.gitignore_global;

  programs.git = {
    enable = true;
    settings = {
      color.ui = true;
      core = {
        autocrlf = false;
        editor = "vim";
        excludesfile = "~/.gitignore_global";
        quotepath = false;
      };
      ghq.root = "~/repos";
      include.path = "~/${localGitConfigPath}";
      init.defaultBranch = "main";
      merge.ff = false;
      pull.ff = "only";
    };
  };

  home.activation.createLocalGitConfig = lib.hm.dag.entryAfter [ "writeBoundary" ] ''
    local_git_config="$HOME/${localGitConfigPath}"
    if [ ! -e "$local_git_config" ]; then
      mkdir -p "$(dirname "$local_git_config")"
      cp ${../../config/git/local.gitconfig.example} "$local_git_config"
      chmod 600 "$local_git_config"
      echo "Created local Git config template: $local_git_config"
    fi
  '';
}
