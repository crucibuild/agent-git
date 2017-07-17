// Copyright (C) 2016 Christophe Camel, Jonathan Pigrée
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schema

import (
	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/crucibuild/sdk-agent-go/agentimpl"
	"github.com/crucibuild/sdk-agent-go/util"
)

// GitCloneCommand represents a "git clone command" as defined in the avro file.
type GitCloneCommand struct {
	Src string
	Dst string
}

// GitCloneCommandType represents the type of a GitCloneCommand.
var GitCloneCommandType agentiface.Type

func init() {
	t, err := util.GetStructType(&GitCloneCommand{})
	if err != nil {
		panic(err)
	}
	GitCloneCommandType = agentimpl.NewTypeFromType("crucibuild/agent-git#git-clone-command", t)
}
