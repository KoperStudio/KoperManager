#!/bin/bash

ARCH=$(arch)
echo "Installing latest version to your BIN path"
echo "Detected platform type: $ARCH"
if [[ "$ARCH" == *"x64"* ]]; then
  wget -O /usr/bin/koper_manager https://github.com/KoperStudio/KoperManager/releases/prerelease/koper_manager
else
  wget -O /usr/bin/koper_manager https://github.com/KoperStudio/KoperManager/releases/prerelease/koper_manager_32bit
fi
echo "Done! Type 'koper_manager -h' to test installation"