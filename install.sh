#!/bin/bash
set -e

VERSION="v1.0.0-beta.1"

if [ -n "`$SHELL -c 'echo $ZSH_VERSION'`" ]; then
  shell_profile="zshrc"
elif [ -n "`$SHELL -c 'echo $BASH_VERSION'`" ]; then
  shell_profile="bashrc"
fi

if [ "$1" == "--darwin" ]; then
  DFILE="moul-darwin"
elif [ "$1" == "--linux" ]; then
  DFILE="moul-linux"
else
  echo "Please specific OS"
  exit 1
fi

if [ -d "$HOME/.local/bin" ]; then
  echo "The '.local/bin' directory already exist. Using it."
else
  echo "The '.local/bin' directory not exist. Creating it."
  mkdir "$HOME/.local/bin"
fi

echo "Downloading $DFILE ..."
wget https://github.com/sophearak/moul/releases/download/$VERSION/$DFILE.tar.gz -O /tmp/$DFILE.tar.gz

if [ $? -ne 0 ]; then
  echo "Download failed! Exiting."
  exit 1
fi

echo "Extracting File..."
tar -C "$HOME" -xzf /tmp/$DFILE.tar.gz
mv "$HOME/$DFILE" "$HOME/.local/bin/moul"

echo "Updating environment variables"
touch "$HOME/.${shell_profile}"

echo 'export PATH=$PATH:~/.local/bin' >> "$HOME/.${shell_profile}"

echo -e "\nmoul $VERSION was installed.\nMake sure to relogin into your shell or run:"
echo -e "\n\tsource $HOME/.${shell_profile}\n\nto update your environment variables."
echo "Tip: Opening a new terminal window usually just works. :)"
rm -f /tmp/$DFILE.tar.gz
