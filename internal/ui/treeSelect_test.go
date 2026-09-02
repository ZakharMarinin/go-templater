package ui

import (
	"github.com/ZakharMarinin/go-templater/internal/domain/entity"
	"testing"
)

func buildTestTree() (nodes []*entity.Node, internal, config, envFile, mainFile *entity.Node) {
	envFile = &entity.Node{Name: ".env", Path: "/.env", IsDir: false, Content: "SECRET=1"}
	config = &entity.Node{Name: "config.go", Path: "/internal/config.go", IsDir: false, Content: "package internal"}
	internal = &entity.Node{Name: "internal", Path: "/internal", IsDir: true, Children: []*entity.Node{config}}
	mainFile = &entity.Node{Name: "main.go", Path: "/main.go", IsDir: false, Content: "package main"}

	nodes = []*entity.Node{internal, envFile, mainFile}

	return nodes, internal, config, envFile, mainFile
}

func TestSeedDefaultExclusions(t *testing.T) {
	nodes, internal, config, envFile, mainFile := buildTestTree()

	var items []treeItem

	parent := make(map[*entity.Node]*entity.Node)
	flattenNodes(nodes, 0, nil, parent, &items)

	excluded := seedDefaultExclusions(items)

	if !excluded[envFile] {
		t.Error("expected dotfile .env to be excluded by default")
	}

	for _, n := range []*entity.Node{internal, config, mainFile} {
		if excluded[n] {
			t.Errorf("expected %s to be included by default", n.Name)
		}
	}
}

func TestIsExcludedCascadesFromAncestor(t *testing.T) {
	nodes, internal, config, envFile, mainFile := buildTestTree()

	var items []treeItem

	parent := make(map[*entity.Node]*entity.Node)
	flattenNodes(nodes, 0, nil, parent, &items)

	explicit := map[*entity.Node]bool{internal: true}

	if !isExcluded(config, parent, explicit) {
		t.Error("expected config.go to be excluded because parent internal/ is excluded")
	}

	if isExcluded(mainFile, parent, explicit) {
		t.Error("did not expect main.go to be excluded")
	}

	if isExcluded(envFile, parent, explicit) {
		t.Error("did not expect .env to be excluded (not seeded here, not explicit)")
	}
}

func TestApplyExclusionsClearsContentButKeepsStructure(t *testing.T) {
	nodes, internal, config, envFile, mainFile := buildTestTree()

	var items []treeItem

	parent := make(map[*entity.Node]*entity.Node)
	flattenNodes(nodes, 0, nil, parent, &items)

	explicit := map[*entity.Node]bool{internal: true, envFile: true}

	applyExclusions(items, parent, explicit)

	if config.Content != "" {
		t.Error("expected config.go content to be cleared")
	}

	if envFile.Content != "" {
		t.Error("expected .env content to be cleared")
	}

	if mainFile.Content == "" {
		t.Error("expected main.go content to be kept")
	}

	if internal.Name != "internal" || len(internal.Children) != 1 {
		t.Error("expected internal/ structure to remain untouched")
	}
}
