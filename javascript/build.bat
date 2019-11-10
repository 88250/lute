: Lute - A structured markdown engine.
: Copyright (c) 2019-present, b3log.org
:
: Lute is licensed under the Mulan PSL v1.
: You can use this software according to the terms and conditions of the Mulan PSL v1.
: You may obtain a copy of Mulan PSL v1 at:
:     http:#license.coscl.org.cn/MulanPSL
: THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
: IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
: PURPOSE.
: See the Mulan PSL v1 for more details.

SET GOOS=linux

RD /S /Q %GOPATH%\pkg\linux_js

: go list --tags "!sm"  -f {{.Deps}}
gopherjs build . --tags "javascript !sm" -o lute.min.js
RD /S /Q %GOPATH%\pkg\linux_js

: go list --tags "sm"  -f {{.Deps}}
gopherjs build . --tags "javascript sm" -o lute-sm.min.js
RD /S /Q %GOPATH%\pkg\linux_js
