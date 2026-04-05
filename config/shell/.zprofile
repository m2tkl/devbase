# devbase common zprofile

if [ -f /etc/profile ]; then
  source /etc/profile
fi

if [ -f "$HOME/.nix-profile/etc/profile.d/hm-session-vars.sh" ]; then
  source "$HOME/.nix-profile/etc/profile.d/hm-session-vars.sh"
fi
