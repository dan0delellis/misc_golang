install vm
set vim to disable mouse, set tabstop and width to 4, clear trailing whitepsaec at exit
install keychain
Change the hostname to something unique
Chane the default password
Make another user and disable the pi user, add user to sudo
Make sudo require a password
Write the rpi pubkey to authorized_keys
Add AllowUsers to sshdconf
enable sshd service
copy the tlstc private key
Disable ssh password auth and enable pubkey auth
set hosts.allow to have [sshd]=10.0.0.
add hosts to /etc/hosts file (maybe just set up a fucking dhcp server)
