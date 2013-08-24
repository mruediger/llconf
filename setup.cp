<html>

(done (and (vim) (bash)))

(vim (and (installed "vim") (configured_vim)))

(bash (and (installed "bash") (configured_bash)))

(configured_bash
 (sync_template "$HOME/.bashrc" "bashrc"))
)

(configured_vim
 (sync_template "$HOME/.vimrc" "vimrc"))
)

(sync_template
  (exec [env:HOME] "cp ~/templates/" [arg: 1] " " [arg:0]))

(installed
 (or
  (and (debian) (apt_installed) (apt [arg: 0]))
  (and (fedora) (yum_installed) (yum [arg: 0]))
 )
)

(apt_installed (exec [env:HOME] "/usr/bin/test -x /usr/bin/apt"))
(yum_installed (exec [env:HOME] "/usr/bin/test -x /usr/bin/yum"))
(debian (exec [env:HOME] "/usr/bin/grep -q Debian /etc/system-release"))
(fedora (exec [env:HOME] "/usr/bin/grep -q Fedora /etc/system-release"))
(yum (exec [env:HOME] "yum install " [arg: 0]))
(apt (exec [env:HOME] "apt-get install " [arg: 0]))

(git
  (or (git update [arg:0]) (git clone [arg:0] [arg:1]))
)

(git update
  (exec [arg:0] "git fetch --all")))

(git clone
  (exec [arg:0] "git clone " [arg:1])))

</html>

