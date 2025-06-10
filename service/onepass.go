/*
Copyright © 2025 John Lennard <john@yakmoo.se>
*/
package service

import (
	"context"
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

	section := onepassword.ItemSection{
		ID:    sectionName,
		Title: sectionName,
	}

	fields := make([]onepassword.ItemField, 0)
	for k, v := range environment {
		field := onepassword.ItemField{
			ID:        k,
			Title:     k,
			Value:     v,
			FieldType: onepassword.ItemFieldTypeConcealed,
			SectionID: &sectionName,
		}
		fields = append(fields, field)
	}

	itemParams := onepassword.ItemCreateParams{
		Title:    itemName,
		Sections: append([]onepassword.ItemSection{}, section),
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
	vaults, err := client.Vaults().List(context.Background())
	if err != nil {
		return nil, err
	}

	for i := range vaults {
		if vaults[i].Title == vaultName {
			return &vaults[i], nil
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

	items, err := client.Items().List(context.Background(), vault.ID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].Title == itemName {
			item, err := client.Items().Get(context.Background(), vault.ID, items[i].ID)
			if err != nil {
				return nil, err
			}
			return &item, nil
		}
	}
	return nil, fmt.Errorf("item %s not found", itemName)
}

func NewClientFromToken(token string) (*onepassword.Client, error) {
	return onepassword.NewClient(
		context.Background(),
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo("envop", "v0.0.0"),
	)
}
