/*
 * Copyright (C) 2016 Red Hat, Inc.
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package client

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/skydive-project/skydive/api/client"
	shttp "github.com/skydive-project/skydive/http"
	"github.com/skydive-project/skydive/js"
	"github.com/skydive-project/skydive/logging"
)

// ShellCmd skydive shell root command
var ShellCmd = &cobra.Command{
	Use:          "shell",
	Short:        "Shell Command Line Interface",
	Long:         "Skydive Shell Command Line Interface, yet another shell",
	SilenceUsage: false,
	Run: func(cmd *cobra.Command, args []string) {
		shellMain()
	},
}

var (
	// ErrContinue parser error continue input
	ErrContinue = errors.New("<continue input>")
	// ErrQuit parser error quit session
	ErrQuit = errors.New("<quit session>")
)

func (s *Session) completeWord(line string, pos int) (string, []string, string) {
	if len(line) == 0 || pos == 0 {
		return "", nil, ""
	}

	// Chunck data to relevant part for autocompletion
	// E.g. in case of nested lines eth.getBalance(eth.coinb<tab><tab>
	start := pos - 1
	for ; start > 0; start-- {
		// Skip all methods and namespaces (i.e. including the dot)
		if line[start] == '.' || (line[start] >= 'a' && line[start] <= 'z') || (line[start] >= 'A' && line[start] <= 'Z') {
			continue
		}
		// We've hit an unexpected character, autocomplete form here
		start++
		break
	}
	return line[:start], s.jsre.CompleteKeywords(line[start:pos]), line[pos:]
}

func shellMain() {
	s, err := NewSession()
	if err != nil {
		panic(err)
	}

	rl := newContLiner()
	defer rl.Close()

	var historyFile string
	home, err := homeDir()
	if err != nil {
		logging.GetLogger().Errorf("home: %s", err)
	} else {
		historyFile = filepath.Join(home, "history")

		f, err := os.Open(historyFile)
		if err != nil {
			if !os.IsNotExist(err) {
				logging.GetLogger().Errorf("%s", err)
			}
		} else {
			_, err := rl.ReadHistory(f)
			if err != nil {
				logging.GetLogger().Errorf("while reading history: %s", err)
			}
		}
	}

	rl.SetWordCompleter(s.completeWord)

	for {
		in, err := rl.Prompt()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "fatal: %s", err)
			os.Exit(1)
		}

		if in == "" {
			continue
		}

		err = s.Eval(in)
		if err != nil {
			if err == ErrContinue {
				continue
			} else if err == ErrQuit {
				break
			}
			fmt.Println(err)
		}
		rl.Accepted()
	}

	if historyFile != "" {
		err := os.MkdirAll(filepath.Dir(historyFile), 0755)
		if err != nil {
			logging.GetLogger().Errorf("%s", err)
		} else {
			f, err := os.Create(historyFile)
			if err != nil {
				logging.GetLogger().Errorf("%s", err)
			} else {
				_, err := rl.WriteHistory(f)
				if err != nil {
					logging.GetLogger().Errorf("while saving history: %s", err)
				}
			}
		}
	}
}

func homeDir() (home string, err error) {
	home = os.Getenv("SKYDIVE_HOME")
	if home != "" {
		return
	}

	home, err = homedir.Dir()
	if err != nil {
		return
	}

	home = filepath.Join(home, ".skydive")
	return
}

// Session describes a shell session
type Session struct {
	authenticationOpts shttp.AuthenticationOpts
	jsre               *js.JSRE
}

// NewSession creates a new shell session
func NewSession() (*Session, error) {
	s := &Session{
		authenticationOpts: AuthenticationOpts,
		jsre:               js.NewJSRE(),
	}

	client, err := client.NewCrudClientFromConfig(&AuthenticationOpts)
	if err != nil {
		return nil, err
	}

	s.jsre.Start()
	s.jsre.RegisterAPIClient(client)

	return s, nil
}

// Eval evaluation a input expression
func (s *Session) Eval(in string) error {
	_, err := s.jsre.Exec(in)
	if err != nil {
		return fmt.Errorf("Error while executing Javascript '%s': %s", in, err.Error())
	}
	return nil
}
