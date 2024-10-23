#!/bin/bash

# Define TAG_NAME as a parameter with a default value of 0.0.0
TAG_NAME="${1:-0.0.0}"
project=warehouse

if [ -d dist/ ]; then rm -r dist/; fi
mkdir dist/

mkdir -p ./debian/DEBIAN

cat <<EOF > ./debian/DEBIAN/control
Package: baldeweg-desktop
Version: 3.1.0
Section: admin
Priority: optional
Architecture: all
Depends: sudo, curl, apt-transport-https, dirmngr, wget, unzip, ffmpeg, clamav, clamav-base, clamav-freshclam, clamtk, clamtk-gnome, gimp, firefox, vlc, synaptic, gthumb, ubuntu-restricted-extras, gnome-tweaks, dconf-editor, chrome-gnome-shell, shotwell, baobab, simple-scan, usb-creator-common, usb-creator-gtk, gnome-clocks, ansible, duplicity, software-properties-common, python3-openssl, rsync, nano, sysstat, flatpak, gnome-software-plugin-flatpak, git, borgbackup, keepassxc
Conflicts:
Installed-Size: 5400
Maintainer: André Baldeweg <a@baldeweg.org>
Homepage: http://baldeweg.org/
Description: Installs packages for desktop
 This package installs packages for the desktop.
EOF

mkdir -p ./debian/usr/bin
cp warehousecli ./debian/usr/bin

find debian/ -type f -exec chmod 0644 {} \;
find debian/ -type d -exec chmod 0755 {} \;

chmod +x debian/usr/bin/warehousecli

fakeroot dpkg -b ./debian dist/${project}.deb
