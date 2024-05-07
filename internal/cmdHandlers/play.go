package cmdhandlers
//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
//
// 	"github.com/SimonVillalonIT/go-music/internal/services"
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// )
//
// func PlayHandler(cmd *cobra.Command, args []string) {
// 	rawFile, err := os.ReadFile(viper.GetString(services.STORE_PATH))
//
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
//
// 	var jsonFile services.JsonFile
//
// 	err = json.Unmarshal(rawFile, &jsonFile)
//
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
//
// 	idxs, err := services.DisplayFuzzyFind(jsonFile)
//
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
//
// 	// Create a temporary playlist file
// 	playlistFile, err := os.CreateTemp("", "playlist*.txt")
// 	if err != nil {
// 		fmt.Println("Error creating temporary file:", err)
// 		return
// 	}
// 	defer os.Remove(playlistFile.Name()) // Remove the temporary file when done
//
// 	for _, idx := range idxs {
// 		itemPath := getItemPath(jsonFile[idx])
// 		_, err := fmt.Fprintf(playlistFile, "%s\n", itemPath)
// 		if err != nil {
// 			fmt.Println("Error writing to playlist file:", err)
// 			return
// 		}
// 	}
//
// 	playlistFile.Close()
//
// 	if err := services.Play(playlistFile.Name()); err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
//
//
