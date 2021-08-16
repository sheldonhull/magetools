#!/usr/bin/env bash
# Run with elevated privileges
set -e

apt-get -yqq update --fix-missing && apt-get -yqq install pv
mkdir -p ./tmpinstall && curl --silent "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "./tmpinstall/awscliv2.zip"
COUNT=`unzip -q -l "./tmpinstall/awscliv2.zip" | wc -l`
mkdir -p ./tmpinstall/aws
unzip "./tmpinstall/awscliv2.zip" -d "./tmpinstall/"  | pv -l -s $COUNT >/dev/null
./tmpinstall/aws/install --update | (pv --timer --name "🤖 awscli")
rm -rf ./tmpinstall/
apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts

echo "✅ installing session manager plugin" && mkdir -p ./tmpinstall
curl --silent "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb" -o "./tmpinstall/session-manager-plugin.deb"
dpkg -i ./tmpinstall/session-manager-plugin.deb
rm -rf ./tmpinstall/
