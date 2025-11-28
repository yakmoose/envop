package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/1password/onepassword-sdk-go"
	"github.com/google/uuid"
	"github.com/yakmoose/envop/collection"
)

func anyToStringish(val any) string {
	switch val.(type) {
	case string:
		return val.(string)

	default:
		vv, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		return string(vv)
	}
}

func stringishToAny(stringish string) any {

	stringish = strings.TrimSpace(stringish)

	// if it's an int, make sure it returns as an int, since the json unmarshal treats em as floats if
	// no type is specified... and we can lose precision...
	matched, _ := regexp.MatchString("^[+-]?[0-9]+$", stringish)
	if matched {
		parseInt, _ := strconv.ParseInt(stringish, 10, 64)
		return parseInt
	}

	if json.Valid([]byte(stringish)) {
		var anyish any
		err := json.Unmarshal([]byte(stringish), &anyish)
		if err != nil {
			return nil
		}

		return anyish
	}
	return stringish
}

func NewClientFromToken(token string) (*onepassword.Client, error) {
	return onepassword.NewClient(
		context.Background(),
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo("envop", "v0.0.0"),
	)
}

func EnvironmentToFields(
	environment *map[string]any,
	section *onepassword.ItemSection,
) *[]onepassword.ItemField {

	fields := make([]onepassword.ItemField, 0, len(*environment))
	for k, v := range *environment {

		vv := anyToStringish(v)

		field := onepassword.ItemField{
			ID:        uuid.New().String(),
			Title:     strings.TrimSpace(k),
			Value:     strings.TrimSpace(vv),
			FieldType: onepassword.ItemFieldTypeConcealed,
			SectionID: &section.ID,
		}
		fields = append(fields, field)
	}

	return &fields
}

// CreateItem creates a new 1password item in the specified vault, from the provided environment
func CreateItem(
	client *onepassword.Client,
	vault *onepassword.VaultOverview,
	itemName string,
	sectionName string,
	// environment *map[string]any,
) (*onepassword.Item, error) {
	section := onepassword.ItemSection{
		ID:    uuid.New().String(),
		Title: sectionName,
	}
	//
	//fields := EnvironmentToFields(environment, &section)
	//
	//slices.SortFunc(*fields, func(a onepassword.ItemField, b onepassword.ItemField) int {
	//	return strings.Compare(a.Title, b.Title)
	//})

	sections := append([]onepassword.ItemSection{}, section)

	itemParams := onepassword.ItemCreateParams{
		Title:    itemName,
		Sections: sections,
		//		Fields:   *fields,
		VaultID:  vault.ID,
		Category: onepassword.ItemCategoryServer,
	}

	item, err := client.Items().Create(context.Background(), itemParams)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateItem updates them
func UpdateItem(
	client *onepassword.Client,
	item *onepassword.Item,
	sectionName string,
	environment *map[string]any,
) (*onepassword.Item, error) {

	// does the section exist?
	var section = onepassword.ItemSection{}
	for _, v := range item.Sections {
		if v.Title == sectionName {
			section = v
			break
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

	l := max(len(*environment), len(item.Fields))

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

	for _, v := range *EnvironmentToFields(environment, &section) {
		fieldMap[v.Title] = v
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

func FindSection(item *onepassword.Item, sectionName string) *onepassword.ItemSection {
	for _, v := range item.Sections {
		if v.Title == sectionName {
			return &v
		}
	}
	return nil
}

func CopySection(
	client *onepassword.Client,
	sourceItem *onepassword.Item,
	sourceSectionName string,
	destinationItem *onepassword.Item,
	destinationSectionName string,

) error {

	sourceSection := FindSection(sourceItem, sourceSectionName)
	if sourceSection == nil {
		return fmt.Errorf("sourceSection %s not found in item %s", sourceSectionName, sourceItem.Title)
	}

	destinationSection := FindSection(destinationItem, destinationSectionName)
	if destinationSection == nil {
		destinationSection = &onepassword.ItemSection{
			ID:    uuid.New().String(),
			Title: destinationSectionName,
		}
		destinationItem.Sections = append(destinationItem.Sections, *destinationSection)
	}

	// find fields in sourceSection...
	// and grab them...
	for _, v := range sourceItem.Fields {
		if *v.SectionID == sourceSection.ID {
			v.ID = uuid.New().String()
			*v.SectionID = destinationSection.ID
			destinationItem.Fields = append(destinationItem.Fields, v)
		}
	}

	slices.SortFunc(destinationItem.Fields, func(a onepassword.ItemField, b onepassword.ItemField) int {
		return strings.Compare(a.Title, b.Title)
	})

	updatedItem, err := client.Items().Put(context.Background(), *destinationItem)
	if err != nil {
		return err
	}

	destinationItem = &updatedItem

	return nil
}

func RemoveSection(
	client *onepassword.Client,
	item *onepassword.Item,
	sectionName string,
) error {

	section := FindSection(item, sectionName)
	if section == nil {
		return fmt.Errorf("sourceSection %s not found in item %s", sectionName, item.Title)
	}

	fields := make([]onepassword.ItemField, 0, len(item.Fields))
	for _, field := range item.Fields {
		if *field.SectionID != section.ID {
			fields = append(fields, field)
		}
	}

	item.Fields = fields

	updatedItem, err := client.Items().Put(context.Background(), *item)
	if err != nil {
		return err
	}

	item = &updatedItem

	return nil
}

func MoveSection(
	client *onepassword.Client,
	sourceItem *onepassword.Item,
	sourceSectionName string,
	destinationItem *onepassword.Item,
	destinationSectionName string,
) error {

	err := CopySection(client, sourceItem, sourceSectionName, destinationItem, destinationSectionName)
	if err != nil {
		return err
	}

	// we need to read back the source section... incase it's changed...
	refreshedSourceItem, err := client.Items().Get(context.Background(), sourceItem.VaultID, sourceItem.ID)

	err = RemoveSection(client, &refreshedSourceItem, sourceSectionName)
	if err != nil {
		return err
	}

	return nil
}

func ReadOnePassword(
	client *onepassword.Client,
	vaultName string,
	itemName string,
	sectionName string,
) (map[string]any, error) {
	vault, err := FindVaultWithName(client, vaultName)
	if err != nil {
		return nil, err
	}

	item, err := FindItemWithName(client, vault, itemName)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, fmt.Errorf("item %s not found in vault", itemName)
	}

	var section *onepassword.ItemSection
	if sectionName != "" {

	}

	var environment map[string]any

	if sectionName == "" {

	} else {

		section = FindSection(item, sectionName)
		environment = collection.Reduce(item.Fields, func(env map[string]any, v onepassword.ItemField) map[string]any {
			if section.ID == *v.SectionID {
				env[strings.TrimSpace(v.Title)] = stringishToAny(v.Value)
			}
			return env
		}, make(map[string]any))
	}

	return environment, nil
}
