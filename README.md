# KoperManager
Package manager for minecraft servers

Install minecraft server software and plugins in 1 ~~click~~ command
## Installation
### For Windows (use linux pls)
+ Download code as zip file
+ Unpack it
+ Win + R and enter cmd.exe (or other another way to open command prompt)
+ Go to folder that you unpacked. Example `d:` `cd Project\KoperManager`
+ Run installation script: `install.bat` **you must run script with admin perms!**
### For other OSes
+ `git clone https://github.com/KoperStudio/KoperManager.git`
+ `cd KoperManager`
+ `chmod +777 install.sh`
+ `./install.sh`

Now you're able to use koper_manager command on windows and ./koper_manager on other OSes
## Setup server
`sudo koper_manager setup_server <brand> <minecraft_version> <name>`

Supported brands:
+ Paper (Recommended for all)
+ Spigot (Recommenced for small testings)
+ Tuinity (Recommended for servers on 1.12.2 or 1.17 with 200+ online)
+ Airplane (Recommended for servers  1.16.5 and 1.17.1 with 200+ online)

Name is standing for folder name, where your server will be installed
## Install a plugin for an existing server
`sudo koper_manager install <query> [path_to_server]`

That will install the first plugin on the list by query from spigotmc.org to plugins folder of \[path_to_server] server.
If path_to_server isn't specified, then it will search plugins folder in your current dir

For windows all commands are same but without `sudo` word
