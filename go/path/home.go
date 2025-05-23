// Copyright The KCL Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Reference: k8s.io/client-go/util/homedir
package path

import (
	"os"
	"os/user"
	"runtime"
)

// HomeDirWithError returns the home directory for the current user along with an error if any.
// On Windows:
// 1. if none of those locations are writeable, the first of %HOME%, %USERPROFILE%, %HOMEDRIVE%%HOMEPATH% that exists is returned.
// 2. if none of those locations exists, the first of %HOME%, %USERPROFILE%, %HOMEDRIVE%%HOMEPATH% that is set is returned.
func HomeDirWithError() (string, error) {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOME")
		homeDriveHomePath := ""
		if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); len(homeDrive) > 0 && len(homePath) > 0 {
			homeDriveHomePath = homeDrive + homePath
		}
		userProfile := os.Getenv("USERPROFILE")

		firstSetPath := ""
		firstExistingPath := ""

		// Prefer %USERPROFILE% over %HOMEDRIVE%/%HOMEPATH% for compatibility with other auth-writing tools
		for _, p := range []string{home, userProfile, homeDriveHomePath} {
			if len(p) == 0 {
				continue
			}
			if len(firstSetPath) == 0 {
				// remember the first path that is set
				firstSetPath = p
			}
			info, err := os.Stat(p)
			if err != nil {
				continue
			}
			if len(firstExistingPath) == 0 {
				// remember the first path that exists
				firstExistingPath = p
			}
			if info.IsDir() && info.Mode().Perm()&(1<<(uint(7))) != 0 {
				// return first path that is writeable
				return p, nil
			}
		}

		// If none are writeable, return first location that exists
		if len(firstExistingPath) > 0 {
			return firstExistingPath, nil
		}

		// If none exist, return first location that is set
		if len(firstSetPath) > 0 {
			return firstSetPath, nil
		}

		// We've got nothing
		return "", nil
	}
	home := os.Getenv("HOME")
	if home != "" {
		return home, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

// HomeDir returns the home directory for the current user.
// On Windows:
// 1. if none of those locations are writeable, the first of %HOME%, %USERPROFILE%, %HOMEDRIVE%%HOMEPATH% that exists is returned.
// 2. if none of those locations exists, the first of %HOME%, %USERPROFILE%, %HOMEDRIVE%%HOMEPATH% that is set is returned.
func HomeDir() string {
	homeDir, _ := HomeDirWithError()
	return homeDir
}
