#!/bin/bash

TAG_NAME="${1:-0.0.0}"
TAG_NAME="${TAG_NAME#v}"
project=warehousecli

if [ -d dist/ ]; then rm -r dist/; fi
mkdir dist/

mkdir -p ./debian/DEBIAN

cat <<EOF > ./debian/DEBIAN/control
Package: warehousecli
Version: ${TAG_NAME}
Section: utils
Priority: optional
Architecture: amd64
Maintainer: André Baldeweg <a@baldeweg.org>
Description: Warehouse CLI
EOF

mkdir -p ./debian/usr/bin
cp ./warehousecli/warehousecli ./debian/usr/bin

find debian/ -type f -exec chmod 0644 {} \;
find debian/ -type d -exec chmod 0755 {} \;

chmod +x debian/usr/bin/warehousecli

fakeroot dpkg -b ./debian dist/${project}.deb
