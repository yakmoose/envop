/*
Copyright © 2023 John Lennard <john@yakmoo.se>
*/
package service

import (
	"github.com/1Password/connect-sdk-go/connect"
	"github.com/1Password/connect-sdk-go/onepassword"
)

// Create1PasswordItem creates a new 1password item in the specified vault, from the provided environment
func Create1PasswordItem(client connect.Client, vaultName, itemName string, environment map[string]string) (*onepassword.Item, error) {
	vault, err := client.GetVault(vaultName)
	if err != nil {
		return nil, err
	}

	fields := make([]*onepassword.ItemField, 0)
	for k, v := range environment {

		field := onepassword.ItemField{
			Label: k,
			Value: v,
			Type:  "CONCEALED",
		}
		fields = append(fields, &field)
	}

	item := &onepassword.Item{
		Title:    itemName,
		Fields:   fields,
		Vault:    onepassword.ItemVault{ID: vault.ID},
		Category: onepassword.Server,
	}

	return client.CreateItem(item, vault.ID)
}

// Get1PasswordItem retrieves a 1password item from the specified vault
func Get1PasswordItem(client connect.Client, vaultName, itemName string) (*onepassword.Item, error) {
	vault, err := client.GetVault(vaultName)
	if err != nil {
		return nil, err
	}

	return client.GetItem(itemName, vault.ID)
}
