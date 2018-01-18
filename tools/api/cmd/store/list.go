// Copyright © 2017 Osiloke Emoekpere <me@osiloke.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/osiloke/dostow-contrib/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all stores",
	Long:  `list all stores.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := viper.GetString("api.url")
		apiGroup := viper.GetString("api.group")
		apiAuth := viper.GetString("api.auth")

		client := api.NewAdminClient(apiURL, apiGroup, apiAuth)
		var result struct {
			Data []struct {
				Name string `json:"name"`
			} `json:"data"`
		}
		raw, _, err := client.Schema.List(&api.PaginationParams{Size: 100})
		if err == nil {
			err = json.Unmarshal(*raw, &result)
			if err == nil {
				log.Println(fmt.Sprintf("%v", result))
			}
		}
	},
}

func init() {
	StoreCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
