package main

import (
    "os"
    "github.com/google/go-github/github"
    "github.com/xoebus/statham"
    "golang.org/x/oauth2"
    "io"
    "net/http"
)

import log "github.com/palette-software/insight-tester/common/logging"

// This mapping thing is a workaround an issue in the google/go-github package
// https://github.com/google/go-github/issues/246
// When that is solved this can be much more simple
func connect(token string) *github.Client {
    githubTransport := &oauth2.Transport{
        Source: oauth2.StaticTokenSource(
            &oauth2.Token{AccessToken: token},
        ),
    }
    defaultTransport := &http.Transport{}
    mappedTransport := statham.NewTransport(defaultTransport, statham.Mapping{
        "api.github.com": githubTransport,
    })
    httpClient := &http.Client{Transport: mappedTransport}
    return github.NewClient(httpClient)
}

func main() {
	log.AddTarget(os.Stdout, log.InfoLevel)
    if len(os.Args) < 4 {
        log.Error("Usage: ", os.Args[0], "<owner> <repository> <access token>")
        os.Exit(1)
    }
    owner := os.Args[1]
    repo := os.Args[2]
    token := os.Args[3]

    client := connect(token)

    // Get the latest release
    latest, _, err := client.Repositories.GetLatestRelease(owner, repo)
    if err != nil {
        log.Error("Error while getting latest release: ", err)
        os.Exit(1)
    }

    // Get assets for the latest release
    var opt github.ListOptions
    assets, _, err := client.Repositories.ListReleaseAssets(owner, repo, *(latest.ID), &opt)

    // Download each asset
    for _, asset := range assets {

        // Create file on file system
        outFile, err := os.Create(*(asset.Name))
        if err != nil {
            log.Error("Error creating file: ", asset.Name, err)
            continue
        }
        defer outFile.Close()

        // Download the asset
        content, redirectURL, err := client.Repositories.DownloadReleaseAsset(owner, repo, *(asset.ID))
        if err != nil {
            log.Error("Error while downlaoding asset: ", asset.Name, err)
            continue
        }

		if redirectURL != "" {
			// We are going to overwrite the contents acquired before, but as per the documentation the
			// content must be nil. So we are free to overwrite it
			response, err := http.Get(redirectURL)
			if err != nil {
				log.Error("Failed to download redirected content from URL:", redirectURL)
				continue
			}
			content = response.Body
		}

		// Do not move this deferred call before the redirect URL check, since content might be
		// changed in there!
		defer content.Close()

        // Write contents to file
        _, err = io.Copy(outFile, content)
        if err != nil {
            log.Error("Error while writing file: ", asset.Name, err)
        }
    }
}
