DynConf [![Travis CI](https://travis-ci.org/hahnjo/dynconf.svg?branch=master)](https://travis-ci.org/hahnjo/dynconf) [![Coverage Status](https://coveralls.io/repos/github/hahnjo/dynconf/badge.svg?branch=master)](https://coveralls.io/github/hahnjo/dynconf?branch=master)
=======

DynConf is a small program to apply recipes to configuration files.
This can be used to dynamically alter a configuration without diverging from the defaults:
When there is an update you simply re-run DynConf to produce a new configuration file.
In this case the recipes describe which lines should be deleted, replaced, or appended.

The executable expects one of the following subcommands:
 * `apply` takes a recipe, produces a configuration file and writes the result.
   (You might need to run this subcommand as `root` to modify files in `/etc/`.)
 * `check` validates the given recipe.
 * `show` produces a configuration file, but outputs the result for inspection.

Recipes are written in YAML and look like this:
```yaml
file: "/etc/test.conf"

delete:
  -
    context:
      begin: "begin"
      end: "end"
    search: "remove"

replace:
  -
    context:
      begin: "begin"
      end: "end"
    search: "pattern"
    replace: "substitution"

append: "last line"
```
`delete` and `replace` are arrays and their `search` key is interpreted as regular expression.

`context` is optional and allows to restrict `delete` and `replace` to a subset of the file.
`begin` and `end` are interpreted as regular expressions and matched to input file before deleting a line or replacing its contents.
If `begin` or `end` is omitted, the context begins in the first line or ends at the last.
`begin` and `end` do not match the same substring, ie. `end` can only match from the position where the match of `begin` ended.
However, if `begin` and `end` still match at the same line the context will not be enabled.

`file` names the configuration file that should be produced.
The unmodified input is taken from (in this order):
1. An updated configuration file installed by the distribution's package manager.
   For example, on Arch Linux these are called `.pacnew`, `rpm` calls them `.rmpnew`.
2. A saved copy of the unmodified configuration file, suffixed with `.orig`.
3. If all else fails DynConf will modify the configuration file itself.

To ensure idempotence DynConf will create an unmodified copy and suffix it with `.orig`.
As seen in the previous paragraph, updates have higher priority and an invocation of `apply` will work with the new configuration file.
In addition, DynConf will update the `.orig` file and remove the new file installed by the package manager.

License
-------

The code is released under the GNU General Public License v3:

    Copyright (C) 2018 - 2019  Jonas Hahnfeld

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

Individual files contain the following tag instead of the full license text.

	SPDX-License-Identifier:	GPL-3.0-or-later

This enables machine processing of license information based on the SPDX
License Identifiers that are here available: http://spdx.org/licenses/
