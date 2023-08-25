# #!/bin/bash
# # pull tu master thoi
# git pull
# export PATH=$PATH:/usr/local/go/bin
# source $HOME/.profile
# rm -rf build/
# #NIGHTLY=vmonitor make package
# #NIGHTLY=vmonitor make package include_packages="windows_amd64.zip windows_i386.zip linux_amd64.tar.gz linux_arm64.tar.gz i386.deb amd64.deb arm64.deb static_linux_amd64.tar.gz x86_64.rpm i386.rpm"
# NIGHTLY=vmonitor make package include_packages="freebsd_amd64.tar.gz freebsd_armv7.tar.gz freebsd_i386.tar.gz"

# Access token cua acc github
export GITHUB_TOKEN=$2
# Version release vmonitor metrics agent sau khi build xong
tag=$1
~/go/bin/github-release release -u anngdinh -r node-agent  --tag "${tag}"
cd ./build/dist
for file in ./*; do ~/go/bin/github-release upload -u anngdinh -r node-agent  --tag "${tag}" --name "${file##*/}" --file ${file##*/} -R; done