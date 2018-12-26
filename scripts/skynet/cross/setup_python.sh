#!/usr/bin/env bash

set -ex

if [ -z "${IN_SKYNET_CLI+yes}" ]; then
    mkdir /python
    tar -xzf ${SKYNET_APPLICATION_FRUGAL_PYPI} -C /python
    cd /python/frugal*
else
    cd $GOPATH/src/github.com/Workiva/frugal/lib/python
fi

pip install -e ".[tornado]"
python3 /usr/bin/pip3 install Cython==0.27.2
python3 /usr/bin/pip3 install -e ".[asyncio]"
python3 /usr/bin/pip3 install aiostomp==1.4.0
