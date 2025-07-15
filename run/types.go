package run

import (
	"fmt"
	"os"
	"strings"
)

type PlistFile struct {
}

type BuildDetails struct {
	Progress string
	Action   string

	Error error
}

func (b *BuildDetails) String() string {
	if b.Error != nil {
		return fmt.Sprintf("%v complete", b.Progress)
	}
	return fmt.Sprintf("[%v]: %v\n%v", b.Progress, b.Action, b.Error.Error())
}

type Bundle struct {
	ProjectPath string

	Bundle struct {
		Name string
		Path string
	}
}

func (b *Bundle) Initialized() bool {
	if _, err := os.ReadDir(b.Bundle.Path); err == nil {
		return true
	}

	return false
}

type Certificate struct {
	Hash  string
	Title string
}

func (c *Certificate) String() string {
	return fmt.Sprintf("%v \"%v\"", c.Hash, c.Title)
}

func decodeCertificate(v string) *Certificate {
	cert := &Certificate{}

	pieces := strings.Split(v, " ")
	if len(pieces) > 2 {
		// Example: 0) HASH "Apple Development: ..."

		index := pieces[0]
		hash := pieces[1]
		title := strings.Join(pieces[2:], " ")

		// validate signature
		if strings.HasSuffix(index, ")") {
			if strings.HasPrefix(title, "\"") && strings.HasSuffix(title, "\"") {
				cert.Title = strings.TrimPrefix(strings.TrimSuffix(title, "\""), "\"")
				cert.Hash = hash

				return cert
			}
		}
	}

	return nil
}
