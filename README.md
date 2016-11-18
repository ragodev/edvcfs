# edvcfs

![](https://img.shields.io/badge/coverage-88%25-brightgreen.svg)

This repository contains builds for fossil that are on non-major architectures / systems.

# Binaries

Built using the following:

```
wget https://www.fossil-scm.org/download/fossil-src-1.X.tar.gz
mkdir build
cd build
sudo apt-get install openssl libssl-dev
../configure
make

uname -a
fossil version
tar -czvf fossil-VERSION-UNAMEINFO.tar.gz
```
