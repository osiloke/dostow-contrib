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
	"github.com/osiloke/dostow-contrib/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Copy a store",
	Long:  `Copy a store.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := viper.GetString("api.url")
		apiGroup := viper.GetString("api.group")
		apiAuth := viper.GetString("api.auth")

		destURL := viper.GetString("dest.url")
		destGroup := viper.GetString("dest.group")
		destAuth := viper.GetString("dest.auth")

		scon := &api.ConnectionParams{
			apiURL,
			apiGroup,
			apiAuth,
		}
		dcon := &api.ConnectionParams{
			destURL,
			destGroup,
			destAuth,
		}
		name, _ := cmd.Flags().GetString("store-name")
		dest, _ := cmd.Flags().GetString("dest-name")
		copyData, _ := cmd.Flags().GetBool("copy-data")
		if copyData {
			var extraFields map[string]interface{}
			api.CloneStore(name, dest, extraFields, scon, dcon)
		} else {
			api.CopySchema(name, dest, scon, dcon)
		}
	},
}

func init() {
	StoreCmd.AddCommand(cloneCmd)

	cloneCmd.Flags().StringP("store-name", "n", "", "store name")
	cloneCmd.Flags().StringP("dest-name", "e", "", "destination store name")
	cloneCmd.Flags().StringP("extra-fields", "x", "", "extra fields")
	cloneCmd.Flags().BoolP("copy-data", "c", false, "copy data")

	cloneCmd.Flags().StringP("dest-api-url", "s", "", "dest api url")
	cloneCmd.Flags().StringP("dest-api-group", "o", "", "dest api group")
	cloneCmd.Flags().StringP("dest-api-auth", "i", "", "dest api authorization")

	viper.BindPFlag("dest.url", cloneCmd.Flags().Lookup("dest-api-url"))
	viper.BindPFlag("dest.group", cloneCmd.Flags().Lookup("dest-api-group"))
	viper.BindPFlag("dest.auth", cloneCmd.Flags().Lookup("dest-api-auth"))
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
