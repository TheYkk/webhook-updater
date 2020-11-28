package main

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"os"
)


func main() {
	log.Print("Init Webhook updater ")


	token := getenv("GITHUB_TOKEN", "")
	if token == "" {
		log.Fatal("Github token not set")
	}

	url := getenv("URL", "")
	if token == "" {
		log.Fatal("URL not set")
	}

	secret := getenv("SECRET", "")
	if token == "" {
		log.Fatal("SECRET not set")
	}

	ctx := context.Background()

	// ? Login to github
	client := login(ctx, token)
	listOpt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}
	starHook := &github.Hook{
		Events: []string{"watch"},
		Config: map[string]interface{}{
			"url":          url,
			"content_type": "json",
			"secret": secret,
		},
		Active: github.Bool(true),
	}

	for {
		repository,_,err := client.Repositories.List(ctx, "TheYkk" ,listOpt)
		if err != nil {
			log.Fatal(err)
		}

		for repoID,repo := range repository {
			if !*repo.Fork && !*repo.Archived{
				log.Printf("ID: %v Fork degil %s %s", (listOpt.ListOptions.Page * listOpt.ListOptions.PerPage) + repoID +1 , *repo.Owner.Login, *repo.Name)

				// ? List repo hooks
				existHooks, _, errListHooks := client.Repositories.ListHooks(ctx, *repo.Owner.Login, *repo.Name, &github.ListOptions{})
				if errListHooks != nil {
					log.Fatal(errListHooks)
				}
				updated := false

				// ? If the hook to be added is exist, delete and re add.

				for _,hook := range existHooks {

					if hook.Config["url"] ==  starHook.Config["url"]{
						log.Print("Hook bulundu, Silinip yeninded olusturulacak")
						_, _ = client.Repositories.DeleteHook(ctx, *repo.Owner.Login, *repo.Name, *hook.ID)

						_, _, errHook := client.Repositories.CreateHook(ctx, *repo.Owner.Login, *repo.Name, starHook)
						if errHook != nil {
							log.Fatal(errHook)
						}
						updated = true
					}
				}

				// ? If newer added hooks , add the hook.
				if !updated{
					_, _, errHook := client.Repositories.CreateHook(ctx, *repo.Owner.Login, *repo.Name, starHook)
					if errHook != nil {
						log.Fatal(errHook)
					}
				}

			}
		}

		if len(repository) < listOpt.ListOptions.PerPage || len(repository) == 0 {
			break
		}
		listOpt.ListOptions.Page += 1
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func login(ctx context.Context, accessToken string) *github.Client {
	if len(accessToken) == 0 {
		return github.NewClient(nil)
	}

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tokenClient := oauth2.NewClient(ctx, tokenSource)

	client := github.NewClient(tokenClient)
	return client
}