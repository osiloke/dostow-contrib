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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// storeCmd represents the store command
var StoreCmd = &cobra.Command{
	Use:   "store",
	Short: "Store commands",
	Long:  `Store commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("store called")
	},
}

func init() {

	StoreCmd.Flags().StringP("api-url", "u", "", "source api url")
	StoreCmd.Flags().StringP("api-group", "g", "", "source api group")
	StoreCmd.Flags().StringP("api-auth", "t", "", "source api authorization")

	viper.BindPFlag("api.url", StoreCmd.Flags().Lookup("api-url"))
	viper.BindPFlag("api.group", StoreCmd.Flags().Lookup("api-group"))
	viper.BindPFlag("api.auth", StoreCmd.Flags().Lookup("api-auth"))
}
