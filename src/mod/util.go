package mod

import (
	"context"
	"fmt"
	"github.com/l1ghthouse/northstar-bootstrap/src/mod/thunderstore"
	"strings"

	"github.com/google/go-github/v42/github"
)

func latestGithubTag(ctx context.Context, repoOwner string, repoName string) (string, error) {
	client := github.NewClient(nil)
	tags, _, err := client.Repositories.ListTags(ctx, repoOwner, repoName, nil)
	if err != nil {
		return "", fmt.Errorf("error listing tags: %w", err)
	}
	if len(tags) > 0 {
		return tags[0].GetName(), nil
	}
	return "", ErrNoTagsFound
}

func cmdWgetZipBuilder(link string, zipName string) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("wget %s -O %s.zip", link, zipName))
	builder.WriteString("\n")
	return builder.String()
}

func cmdUnzipBuilder(zipName string) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("unzip %s.zip -d /", zipName))
	builder.WriteString("\n")
	return builder.String()
}

func cmdUnzipBuilderWithDst(zipName string) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("mkdir -p /%s", zipName))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("unzip %s.zip -d /%s", zipName, zipName))
	builder.WriteString("\n")
	return builder.String()
}

func dockerArgBuilder(modName string) string {
	return fmt.Sprintf("--mount \"type=bind,source=/%s,target=/mnt/mods/%s,readonly\"", modName, modName)
}

func latestThunderstoreMod(ctx context.Context, packageName string, requiredByClient bool) (string, string, string, string, bool, error) {
	pkg, err := thunderstore.GetPackageByName(ctx, packageName)
	if err != nil {
		return "", "", "", "", false, fmt.Errorf("failed to get package: %w", err)
	}
	latestVersion, err := thunderstore.GetLatestPackageVersion(pkg)
	if err != nil {
		return "", "", "", "", false, fmt.Errorf("failed to get latest package version: %w", err)
	}

	link := "https://northstar.thunderstore.io/package/download/laundmo/ParseableLogs/0.0.1/"

	builder := strings.Builder{}
	builder.WriteString(cmdWgetZipBuilder(link, packageName))
	builder.WriteString(cmdUnzipBuilderWithDst(packageName))
	return builder.String(), dockerArgBuilder(packageName), link, latestVersion.VersionNumber, requiredByClient, nil
}
