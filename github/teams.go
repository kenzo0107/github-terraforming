package github

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"path/filepath"
)

const (
	fpTFTeam string = "team.tf"

	tplTFTeam string = `locals {
	teams = { {{ range . }}
		{{.Name}} = {
			name        = "{{.Name}}"
			description = "{{.Description}}"
			privacy     = "{{.Privacy}}"
		}{{ end }}
	}
}

resource "github_team" "team" {
	for_each = local.teams

	name        = each.value.name
	description = each.value.description
	privacy     = each.value.privacy
}
`
)

const (
	fpTFImportTeam string = "tfimport_teams.sh"
	tplImportTeam  string = `{{ range . }}terraform import 'github_team.team["{{.Name}}"]' {{.ID}}
{{ end }}`
)

// Teams gets team.tf and shell including 'terraform import github_team.team["xxx"]'
func Teams() {
	if err := _Teams(); err != nil {
		log.Fatal(err)
	}
}

// Teams : github teams
func _Teams() error {
	cli := New(githubToken)
	teams, err := cli.GetAllTeams(githubOrganization)
	if err != nil {
		return err
	}

	// generate artifact/team.tf
	outTF := filepath.Join("artifact", fpTFTeam)
	if err := importTemplate(tplTFTeam, outTF, teams); err != nil {
		return err
	}

	// generate tfimport_teams.sh
	outTFImport := filepath.Join("artifact", fpTFImportTeam)
	if err := importTemplate(tplImportTeam, outTFImport, teams); err != nil {
		return err
	}
	return nil
}

func importTemplate(tpl, output string, a interface{}) error {
	t := template.Must(template.New("text").Parse(tpl))
	buff := new(bytes.Buffer)
	fw := io.Writer(buff)

	if err := t.Execute(fw, a); err != nil {
		return err
	}

	if err := WriteLineToFile(output, buff.String()); err != nil {
		return err
	}

	return nil
}
