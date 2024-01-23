#! /bin/sh
ossystem="Windows_x86_64"
currentversion="2.38.0"

rm -rf ${ossystem}
mkdir ${ossystem}
downloaduri=https://github.com/vektra/mockery/releases/download/v${currentversion}/mockery_${currentversion}_${ossystem}.tar.gz

echo "downloading from: $downloaduri"
curl --location ${downloaduri} --remote-name
tar -xf "mockery_${currentversion}_${ossystem}.tar.gz" -C ${ossystem}
rm "mockery_${currentversion}_${ossystem}.tar.gz"

read -p "Press [Enter] to exit."