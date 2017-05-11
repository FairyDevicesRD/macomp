
# macomp: Morphological Analyzer Comparator

[![CircleCI](https://circleci.com/gh/FairyDevicesRD/macomp.svg?style=svg)](https://circleci.com/gh/FairyDevicesRD/macomp)
[![Report card](http://goreportcard.com/badge/FairyDevicesRD/macomp)](http://goreportcard.com/report/FairyDevicesRD/macomp)
[![GoDoc](https://godoc.org/github.com/FairyDevicesRD/macomp?status.svg)](https://godoc.org/github.com/FairyDevicesRD/macomp)
[![Apache License](http://img.shields.io/badge/license-APACHE2-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

## Usage

### Initialize

First, initialize ``~/.macomp.json`` to run the following command.

```sh
macomp --init
```

Modify it if you need.


### macomp

See ``macomp -h`` for the detail.

```sh
$ macomp
宇宙にあるいて座
juman          |宇 宙|に|あ る い て|座|
jumanpp        |宇 宙|に|あ る い て|座|
mecab-ipa-NE   |宇 宙|に|あ る|い て 座|
mecab-unidic   |宇 宙|に|あ る い|て|座|

$ macomp --pos --check
宇宙にある|いて?座|
X juman          |宇 宙|に|あ る_い て|座|
                 |名   |助|動         |名|
X jumanpp        |宇 宙|に|あ る_い て|座|
                 |名   |助|動         |名|
O mecab-ipa-NE   |宇 宙|に|あ る|い て 座|
                 |名   |助|動   |名      |
X mecab-unidic   |宇 宙|に|あ る_い/て|座|
                 |名   |助|動      |助|名|

$ macomp -t jumanpp -t mecab-ipa-NE
柱で食べるジャパリまんは美味しい
jumanpp      |柱|で|食 べ る|ジ ャ パ リ|ま ん|は|美 味 し い|
mecab-ipa-NE |柱|で|食 べ る|ジ ャ|パ リ|ま|ん|は|美 味 し い|
```

### macomp-server

```sh
macomp-server
```

Open [localhost:5000](http://localhost:5000) and use the form, or use commands.

```sh
curl http://localhost:5000/api/v1/ma/食べたい
curl -X POST -d 'text=食べたい' http://localhost:5000/api/v1/ma
curl -X POST -d 'text=食べたい' -d 'callback=jsonp123' http://localhost:5000/api/v1/ma
```

See ``macomp-server -h`` for the detail.


## Install

Macomp requires the following packages.

- [golang](https://golang.org/)
- [MeCab](https://github.com/taku910/mecab)

JUMAN and JUMAN++ are optional.

- [JUMAN](http://nlp.ist.i.kyoto-u.ac.jp/?JUMAN)
- [JUMAN++](http://nlp.ist.i.kyoto-u.ac.jp/?JUMAN++)

### Golang

```sh
sudo apt-get install golang
echo 'export GOPATH=~/.go; export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
soruce ~/.bashrc
```

### MeCab

```sh
sudo apt-get install mecab libmecab-dev
```

You can also install from the source codes.

```sh
git clone git@github.com:taku910/mecab.git --depth 1
cd mecab/mecab
autoreconf -i
./configure --with-charset=utf8 --enable-utf8-only
make
sudo make install
sudo ldconfig

cd ../mecab-ipadic
./configure --with-charset=utf8
make
sudo make install
```

### macomp

```sh
export CGO_LDFLAGS="`mecab-config --libs`"
export CGO_CFLAGS="-I`mecab-config --inc-dir`"
go get github.com/FairyDevicesRD/macomp/cmd/macomp
go get github.com/FairyDevicesRD/macomp/cmd/macomp-server
```


## Configure

```sh
macomp --init
vi ~/.macomp.json
```

- ``type`` is a mandatory field
- ``path``, ``aliases`` and ``options`` are not mandatory fields.
- Set ``disable`` to ``true`` if you want to disable the setting
- Be aware to use absolute path


## License

- [Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0) by [Fairy Devices Inc](http://www.fairydevices.jp/)
