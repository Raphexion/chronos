Chronos - JIRA time tracker
===========================

[![Build Status](https://travis-ci.org/Raphexion/chronos.svg?branch=master)](https://travis-ci.org/Raphexion/chronos)
[![codecov.io](https://codecov.io/gh/Raphexion/chronos/coverage.svg?branch=master)](https://codecov.io/gh/Raphexion/chronos?branch=master)

Downloading latest release
--------------------------

[Linux](https://github.com/Raphexion/chronos/releases/latest/download/chronos)

[Windows](https://github.com/Raphexion/chronos/releases/latest/download/chronos.exe)

Install Linux
-------------

```sh
mkdir -p ~/bin
wget https://github.com/Raphexion/chronos/releases/latest/download/chronos -O ~/bin/chronos
chmod +x ~/bin/chronos
```

Building from souce
-------------------

Make sure you have working golang installation. Then run:

```sh
go install
```

Getting started
---------------

It is possible to run chronos without a configuration file but it is not recommended.
Chrons can generate a dummy config file for you, which will be placed in you home folder.

```shell
chronos --generate-config
```

If you look in the file you will see the following:

```yaml
jira:
  url: https://myJira.atlassian.net
  mail: myLogin@example.com
  apikey: 1234ABCD
  username: myUserName
```

After you have corrected the configuration, simply type

```sh
chronos
```

```sh
===========================
Week  1
===========================

2018-01-01
	AA-1234:   1.00
	AA-1235:   2.00
	-----------------
		   3.00

	Total:     3.00

===========================
Week  2
===========================

2018-01-08
	AA-1235:   3.00
	------------------
		   3.00

	Total:     3.00
```

Log work in JIRA
----------------

```sh
./chronos --logwork --issue AA-1234 --minutes 20
```
