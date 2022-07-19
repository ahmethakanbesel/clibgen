package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/laureanray/clibgen/pkg/api"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func truncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndex(s[:max], " ")] + " ..."
}

var (
	selectedSite    string
	numberOfResults = 10

	searchCmd = &cobra.Command{
		Use:   "search",
		Short: "search for a book, paper or article",
		Long: `search for a book, paper or article
	example: clibgen search "Eloquent JavaScript"`,
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please enter a search query!")
				return
			}

			var libgenType = api.LibgenNew

			if selectedSite == "old" {
				libgenType = api.LibgenOld
			} else if selectedSite == "new" {
				libgenType = api.LibgenNew
			}

			books, err := api.SearchBookByTitle(args[0], numberOfResults, libgenType)
			if err != nil {
				log.Fatalln(err)
			}

			if err != nil {
				log.Fatal(err)
				return
			}

			var titles []string

			for _, book := range books {
				parsedTitle := truncateText(book.Title, 42)
				parsedAuthor := truncateText(book.Author, 24)
				titles = append(titles, fmt.Sprintf("[%5s %4s] %-45s %s", book.FileSize, book.Extension, parsedTitle, parsedAuthor))
			}

			prompt := promptui.Select{
				Label: "Select Title",
				Items: titles,
			}

			resultInt, _, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			api.DownloadSelection(books[resultInt], libgenType)
		},
	}
)

func init() {
	searchCmd.
		PersistentFlags().
		StringVarP(&selectedSite, "site", "s", "old", `select which site to use
		options: 
			"old" -> libgen.is
			"new" -> liggen.li 
	`)

	searchCmd.
		PersistentFlags().
		IntVarP(&numberOfResults, "number of results", "n", 10, `number of result(s) to be displayed maximum: 25`)

	rootCmd.AddCommand(searchCmd)
}
