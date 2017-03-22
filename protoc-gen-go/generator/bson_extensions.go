package generator

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	bsonTagPattern        = "@bson_tag: (.*)"
	bsonCompatiblePattern = "@bson_compatible"
	bsonUpsertablePattern = "@bson_upsertable"
	goInjectPattern       = `(?s)@go_inject\s(.+)`
	// TODO: This import pattern doesn't support aliased imports
	goImportPattern = `@import[ \t]+"([^"]+)"`
)

var bsonTagRegex, bsonCompatibleRegex, bsonUpsertableRegex, goInjectRegex, goImportRegex *regexp.Regexp

func init() {
	bsonTagRegex = regexp.MustCompile(bsonTagPattern)
	bsonCompatibleRegex = regexp.MustCompile(bsonCompatiblePattern)
	bsonUpsertableRegex = regexp.MustCompile(bsonUpsertablePattern)
	goInjectRegex = regexp.MustCompile(goInjectPattern)
	goImportRegex = regexp.MustCompile(goImportPattern)
}

func (g *Generator) IsMessageBsonCompatible(message *Descriptor) bool {
	if loc, ok := g.file.comments[message.path]; ok {
		preMessageComments := strings.TrimSuffix(loc.GetLeadingComments(), "\n")
		return bsonCompatibleRegex.Match([]byte(preMessageComments))
	}

	return false
}

func (g *Generator) IsMessageBsonUpsertable(message *Descriptor) bool {
	if loc, ok := g.file.comments[message.path]; ok {
		preMessageComments := strings.TrimSuffix(loc.GetLeadingComments(), "\n")
		return bsonUpsertableRegex.Match([]byte(preMessageComments))
	}

	return false
}

func (g *Generator) GetBsonTagForField(message *Descriptor, fieldNumber int) string {
	fieldPath := fmt.Sprintf("%s,%d,%d", message.path, messageFieldPath, fieldNumber)
	if loc, ok := g.file.comments[fieldPath]; ok {
		allFieldComments := []string{loc.GetTrailingComments(), loc.GetLeadingComments()}

		for _, fieldComment := range allFieldComments {
			fieldComment = strings.TrimSuffix(fieldComment, "\n")
			matchedGroups := bsonTagRegex.FindStringSubmatch(fieldComment)
			if matchedGroups != nil {
				return strings.TrimSpace(matchedGroups[1])
			}
		}
	}

	return ""
}

func (g *Generator) GetGoInjectForMessage(message *Descriptor) string {
	if loc, ok := g.file.comments[message.path]; ok {
		allLeadingComments := loc.GetLeadingDetachedComments()
		allLeadingComments = append(allLeadingComments, loc.GetLeadingComments())

		for _, leadingComment := range allLeadingComments {
			matchedGroups := goInjectRegex.FindStringSubmatch(leadingComment)
			if matchedGroups != nil {
				return matchedGroups[1]
			}
		}
	}

	return ""
}

func goImportsFromGoInject(injectBlock string) []string {
	matches := goImportRegex.FindAllStringSubmatch(injectBlock, -1)
	imports := []string{}
	for _, match := range matches {
		imports = append(imports, match[1])
	}

	return imports
}

func goCodeFromGoInject(injectBlock string) string {
	return goImportRegex.ReplaceAllString(injectBlock, "")
}
