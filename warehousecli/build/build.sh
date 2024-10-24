#!/bin/bash

TAG_NAME="${1:-0.0.0}"
TAG_NAME="${TAG_NAME#v}"
project=warehousecli

if [ -d warehousecli/dist/ ]; then rm -r warehousecli/dist/; fi
mkdir warehousecli/dist/

mkdir -p ./warehousecli/debian/DEBIAN

cat <<EOF > ./warehousecli/debian/DEBIAN/control
Package: warehousecli
Version: ${TAG_NAME}
Section: utils
Priority: optional
Architecture: amd64
Maintainer: André Baldeweg <a@baldeweg.org>
Description: Warehouse CLI
EOF

mkdir -p ./warehousecli/debian/usr/bin
cp ./warehousecli/warehousecli ./warehousecli/debian/usr/bin

find warehousecli/debian/ -type f -exec chmod 0644 {} \;
find warehousecli/debian/ -type d -exec chmod 0755 {} \;

chmod +x warehousecli/debian/usr/bin/warehousecli

fakeroot dpkg -b ./warehousecli/debian warehousecli/dist/${project}.deb
