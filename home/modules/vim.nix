{ pkgs, ... }:
{
  programs.vim = {
    enable = true;
    defaultEditor = true;
    plugins = with pkgs.vimPlugins; [
      vim-easymotion
      lightline-vim
      vim-lsp
      vim-lsp-settings
      asyncomplete-vim
      asyncomplete-lsp-vim
      vim-easy-align
      vim-commentary
      vim-fern
      emmet-vim
      vim-gitgutter
      ctrlp-vim
    ];
    extraConfig = ''
      ${builtins.readFile ../../common/vim/.vimrc}

      ${builtins.readFile ../../common/vim/plugins.vim}
    '';
  };
}
