/*
Copyright © 2025 John Lennard <john@yakmoo.se>
*/
package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/1password/onepassword-sdk-go"
)

// Create1PasswordItem creates a new 1password item in the specified vault, from the provided environment
func Create1PasswordItem(
	client *onepassword.Client,
	vaultName,
	itemName,
	sectionName string,
	environment map[string]string,
) (*onepassword.Item, error) {
	vault, err := Get1PasswordVault(client, vaultName)
	if err != nil {
		return nil, err
	}

	fields := make([]onepassword.ItemField, 0)
	for k, v := range environment {
		field := onepassword.ItemField{
			ID:        k,
			Title:     k,
			Value:     v,
			FieldType: onepassword.ItemFieldTypeConcealed,
		}
		fields = append(fields, field)
	}

	itemParams := onepassword.ItemCreateParams{
		Title: itemName,
		Sections: []onepassword.ItemSection{
			{
				Title: sectionName,
			},
		},
		Fields:   fields,
		VaultID:  vault.ID,
		Category: onepassword.ItemCategoryServer,
	}
	item, err := client.Items().Create(context.Background(), itemParams)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Get1PasswordVault retrieves a 1password vault by name
func Get1PasswordVault(client *onepassword.Client, vaultName string) (*onepassword.VaultOverview, error) {
	vaults, err := client.Vaults().ListAll(context.Background())
	if err != nil {
		return nil, err
	}

	for {
		vault, err := vaults.Next()
		if errors.Is(err, onepassword.ErrorIteratorDone) {
			break
		} else if err != nil {
			panic(err)
		}
		if vault.Title == vaultName {
			return vault, nil
		}
	}

	return nil, fmt.Errorf("vault %s not found", vaultName)
}

// Get1PasswordItem retrieves a 1password item from the specified vault by name
func Get1PasswordItem(client *onepassword.Client, vaultName string, itemName string) (*onepassword.Item, error) {

	vault, err := Get1PasswordVault(client, vaultName)
	if err != nil {
		return nil, err
	}

	items, err := client.Items().ListAll(context.Background(), vault.ID)
	if err != nil {
		return nil, err
	}

	for {
		item, err := items.Next()
		if errors.Is(err, onepassword.ErrorIteratorDone) {
			break
		} else if err != nil {
			panic(err)
		}
		if item.Title == itemName {
			item, err := client.Items().Get(context.Background(), vault.ID, item.ID)
			if err != nil {
				return nil, err
			}

			return &item, nil

		}
	}
	return nil, fmt.Errorf("item %s not found", itemName)
}
