package service

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/1password/onepassword-sdk-go"
	"github.com/google/uuid"
)

// CreateItemInVaultWithSection creates a new 1password item in the specified vault, from the provided environment
func CreateItemInVaultWithSection(
	client *onepassword.Client,
	vault *onepassword.VaultOverview,
	itemName,
	sectionName string,
	environment map[string]string,
) (*onepassword.Item, error) {
	section := onepassword.ItemSection{
		ID:    uuid.New().String(),
		Title: sectionName,
	}

	fields := make([]onepassword.ItemField, 0, len(environment))
	for k, v := range environment {
		field := onepassword.ItemField{
			ID:        uuid.New().String(),
			Title:     strings.TrimSpace(k),
			Value:     strings.TrimSpace(v),
			FieldType: onepassword.ItemFieldTypeConcealed,
			SectionID: &section.ID,
		}
		fields = append(fields, field)
	}

	slices.SortFunc(fields, func(a onepassword.ItemField, b onepassword.ItemField) int {
		return strings.Compare(a.Title, b.Title)
	})

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

func UpdateItem(
	client *onepassword.Client,
	item *onepassword.Item,
	sectionName string,
	environment map[string]string,
) (*onepassword.Item, error) {

	// does the section exist?
	var section = onepassword.ItemSection{}
	for _, v := range item.Sections {
		if v.Title == sectionName {
			section = v
		}
	}

	// if not, make it
	if section.ID == "" {
		section = onepassword.ItemSection{
			ID:    uuid.New().String(),
			Title: sectionName,
		}
		item.Sections = append(item.Sections, section)
	}

	l := max(len(environment), len(item.Fields))

	fieldMap := make(map[string]onepassword.ItemField, l)
	fields := make([]onepassword.ItemField, 0, l)

	// filter out the items that are in our section, vs not
	for _, field := range item.Fields {
		if *field.SectionID == section.ID {
			fieldMap[strings.TrimSpace(field.Title)] = field
		} else {
			fields = append(fields, field)
		}
	}

	for k, v := range environment {
		fieldMap[k] = onepassword.ItemField{
			ID:        uuid.New().String(),
			Title:     strings.TrimSpace(k),
			Value:     strings.TrimSpace(v),
			FieldType: onepassword.ItemFieldTypeConcealed,
			SectionID: &section.ID,
		}
	}

	for _, v := range fieldMap {
		fields = append(fields, v)
	}

	slices.SortFunc(fields, func(a onepassword.ItemField, b onepassword.ItemField) int {
		return strings.Compare(a.Title, b.Title)
	})

	item.Fields = fields

	updatedItem, err := client.Items().Put(context.Background(), *item)
	if err != nil {
		return nil, err
	}
	return &updatedItem, nil
}

func ReindexItem(client *onepassword.Client, item *onepassword.Item) (*onepassword.Item, error) {
	sectionMap := make(map[string]onepassword.ItemSection, len(item.Sections))
	for _, section := range item.Sections {
		oldId := section.ID
		section.ID = uuid.New().String()
		sectionMap[oldId] = section
	}

	item.Sections = make([]onepassword.ItemSection, 0, len(sectionMap))
	for _, section := range sectionMap {
		item.Sections = append(item.Sections, section)
	}

	for i, field := range item.Fields {
		field.ID = uuid.New().String()
		newSectionId := sectionMap[*field.SectionID].ID
		field.SectionID = &newSectionId
		item.Fields[i] = field
	}

	slices.SortFunc(item.Fields, func(a onepassword.ItemField, b onepassword.ItemField) int {
		return strings.Compare(a.Title, b.Title)
	})

	updatedItem, err := client.Items().Put(context.Background(), *item)
	if err != nil {
		return nil, err
	}

	return &updatedItem, nil
}

// FindVaultWithName retrieves a 1password vault by name
func FindVaultWithName(client *onepassword.Client, vaultName string) (*onepassword.VaultOverview, error) {
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

// FindItemWithName retrieves a 1password item from the specified vault by name
func FindItemWithName(client *onepassword.Client, vault *onepassword.VaultOverview, itemName string) (*onepassword.Item, error) {

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
	return nil, nil
}

func NewClientFromToken(token string) (*onepassword.Client, error) {
	return onepassword.NewClient(
		context.Background(),
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo("envop", "v0.0.0"),
	)
}
