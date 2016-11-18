# edvcfs

![](https://img.shields.io/badge/coverage-88%25-brightgreen.svg)

This is an *encrypted distributed version control file system* built on top of [fossil](http://fossil-scm.org/index.html/doc/trunk/www/index.wiki).

# Guide

This filesystem (fs) stores documents. A document is composed of entries. A single entry has:

- text content, main data of entry
- timestamp, date of that entry
- document name which refers to which it belongs to
- entry name which is a unique identifier of this entry

The fs stores an entry by writing a file, `data.aes` (AES encrypted). This file is committed to a new branch which has the name `documentname-==-entryname` (AES encrypted), and overrides the commit date with the specified timestamp. The commit message will take four possibilities:

- "new" - designates new entry
- "update" - designates editing of entry
- "ignore-document" - designates deletion of documents, i.e. ignoring of this document henceforth
- "ignore-entry" - designates deletion of entries, i.e. ignoring this entry henceforth

## Use

The API is simple for use:

- `Put(text,document,entry,timestamp)`: `text` and `document` cannot be empty. If `entry` is empty, a new entry is made, otherwise the entry is updated. If `timestamp` is empty, the current datetime used, otherwise it will override.
- `Delete(document,entry)`: deletes `document` if `entry` is empty, otherwise deletes respective `entry`
- `Get(document,entry)`: if `entry` is empty, returns array of entries for latest version of document, in order of timestamp. If `entry` is not empty, it returns an array of entries for each version of that `entry`




# My notes

Binaries This repository contains builds for fossil that are on non-major architectures / systems. Built using the following:

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
