#!/bin/sh
# This code should (try to) follow Google's Shell Style Guide
# (https://google-styleguide.googlecode.com/svn/trunk/shell.xml)
wget -o /usr/local/bin/weave \
    https://github.com/zettio/weave/release/download/latest_release/weave
chmod a+x /usr/local/bin/weave

apt-get install bridge-utils
weave launch && weave launch-dns && weave launch-proxy
eval $(weave proxy-env)
docker run --name a1 -ti ubuntu
