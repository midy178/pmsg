// Copyright 2022-2023 The pmsg Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version

import (
	"fmt"
	"runtime"
)

var (
	AppName     = "pmsg"             // 名称
	Version     = "0.5.0"            // 版本
	BuildCommit = "latest"           // git commit
	BuildTime   = "2023-20-26 15:55" // 编译时间

	OpenSource = "https://github.com/lenye/pmsg" // 开发人
)

const versionTemplate = `  %s
  Version:     %s
  Commit:      %s
  Built:       %s
  Go version:  %s
  OS/Arch:     %s/%s
  Open source: %s
`

func Print() string {
	return fmt.Sprintf(versionTemplate,
		AppName, Version, BuildCommit, BuildTime,
		runtime.Version(), runtime.GOOS, runtime.GOARCH,
		OpenSource)
}
