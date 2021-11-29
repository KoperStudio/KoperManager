# KoperManager
Package manager for minecraft servers

Install minecraft server software and plugins in 1 ~~click~~ command
## Setup server
`./koper_manager setup_server <brand> <minecraft_version> <name>`\s\s
Supported brands:
+ Paper (Recommended for all)
+ Spigot (Recommenced for small testings)
+ Tuinity (Recommended for servers on 1.12.2 or 1.17 with 200+ online)
+ Airplane (Recommended for servers  1.16.5 and 1.17.1 with 200+ online)
Name is standing for folder name, where your server will be installed
## Install a plugin for an existing server
`./koper_manager install <query> [path_to_server]`tl;dr
That will install the first plugin on the list by query from spigotmc.org to plugins folder of \[path_to_server] server.
If path_to_server isn't specified, then it will search plugins folder in your current dir
## TODO
Instead of publishing all binaries, make automatic build script which will install latest version to system binaries folder
