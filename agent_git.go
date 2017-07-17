// Copyright (C) 2016 Camel Christophe, Christophe Camel, Jonathan Pigrée
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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/crucibuild/agent-git/schema"
	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/crucibuild/sdk-agent-go/agentimpl"
)

// Resources represents an handler on the various data files
// Used by the agent(avro files, manifest, etc...).
var Resources http.FileSystem

// AgentGit is an implementation over the Agent implementation available in sdk-agent-go.
// Its goal is to permit Git actions via the agent.
type AgentGit struct {
	*agentimpl.Agent
}

func mustOpenResources(path string) []byte {
	file, err := Resources.Open(path)

	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(file)

	if err != nil {
		panic(err)
	}

	return content
}

// NewAgentGit creates a new instance of AgentGit.
func NewAgentGit() (agentiface.Agent, error) {
	var agentSpec map[string]interface{}

	manifest := mustOpenResources("/resources/manifest.json")

	err := json.Unmarshal(manifest, &agentSpec)

	if err != nil {
		return nil, err
	}

	impl, err := agentimpl.NewAgent(agentimpl.NewManifest(agentSpec))

	if err != nil {
		return nil, err
	}

	agent := &AgentGit{
		impl,
	}

	if err := agent.init(); err != nil {
		return nil, err
	}

	return agent, nil
}

func (a *AgentGit) init() (err error) {
	// register schemas:
	schemas := []string{
		"/schema/git-clone-command.avro",
		"/schema/git-cloned-event.avro",
	}
	if err = a.registerSchemas(schemas); err != nil {
		return err
	}

	// register types:
	types := []agentiface.Type{
		schema.GitCloneCommandType,
		schema.GitClonedEventType,
	}
	err = a.registerTypes(types)

	// register state callback
	a.RegisterStateCallback(a.onStateChange)

	return err
}

func (a *AgentGit) registerSchemas(pathes []string) error {
	for _, path := range pathes {
		content := mustOpenResources(path)

		s, err := agentimpl.LoadAvroSchema(string(content[:]), a)
		if err != nil {
			return fmt.Errorf("Failed to load schema %s: %s", path, err.Error())
		}

		_, err = a.SchemaRegister(s)

		if err != nil {
			return fmt.Errorf("Failed to register schema %s: %s", path, err.Error())
		}
	}
	return nil
}

func (a *AgentGit) registerTypes(types []agentiface.Type) error {
	for _, t := range types {
		if _, err := a.TypeRegister(t); err != nil {
			return fmt.Errorf("Failed to register type %s (which is a %s): %s", t.Name(), t.Type().Name(), err.Error())
		}
	}
	return nil
}

func (a *AgentGit) onStateChange(state agentiface.State) error {
	switch state {
	case agentiface.StateConnected:
		if _, err := a.RegisterCommandCallback("crucibuild/agent-git#clone-command", a.onGitCloneCommand); err != nil {
			return err
		}
	}
	return nil

}

func (a *AgentGit) onGitCloneCommand(ctx agentiface.CommandCtx) error {
	cmd := ctx.Message().(*schema.GitCloneCommand)

	a.Info(fmt.Sprintf("Received git-clone-command: '%s' '%s' '%d' ", cmd.Src, cmd.Dst))

	message := fmt.Sprintf("Repository '%s' cloned into '%s'", cmd.Src, cmd.Dst)

	return ctx.SendEvent(&schema.GitClonedEvent{Rcode: 0, Message: message})
}
