// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type Translation struct {
	Id          string      `json:"id"`
	Translation interface{} `json:"translation"`
}

var I18nCmd = &cobra.Command{
	Use:   "i18n",
	Short: "Management of Mattermost translations",
}

var ExtractCmd = &cobra.Command{
	Use:     "extract",
	Short:   "Extract translations",
	Long:    "Extract translations from the source code and put them into the i18n/en.json file",
	Example: "  i18n extract",
	RunE:    extractCmdF,
}

var CheckCmd = &cobra.Command{
	Use:     "check",
	Short:   "Check translations",
	Long:    "Check translations existing in the source code and compare it to the i18n/en.json file",
	Example: "  i18n list",
	RunE:    checkCmdF,
}

func init() {
	ExtractCmd.Flags().String("enterprise-dir", "../enterprise", "Path to folder with the Mattermost enterprise source code")
	ExtractCmd.Flags().String("mattermost-dir", "./", "Path to folder with the Mattermost source code")
	CheckCmd.Flags().String("enterprise-dir", "../enterprise", "Path to folder with the Mattermost enterprise source code")
	CheckCmd.Flags().String("mattermost-dir", "./", "Path to folder with the Mattermost source code")
	I18nCmd.AddCommand(
		ExtractCmd,
		CheckCmd,
	)
	RootCmd.AddCommand(I18nCmd)
}

func getCurrentTranslations(mattermostDir string) ([]Translation, error) {
	jsonFile, err := ioutil.ReadFile(path.Join(mattermostDir, "i18n", "en.json"))
	if err != nil {
		return nil, err
	}
	var translations []Translation
	json.Unmarshal(jsonFile, &translations)
	return translations, nil
}

func extractStrings(enterpriseDir, mattermostDir string) map[string]bool {
	i18nStrings := map[string]bool{}
	walkFunc := func(path string, info os.FileInfo, err error) error {
		return extractFromPath(path, info, err, &i18nStrings)
	}
	filepath.Walk(mattermostDir, walkFunc)
	filepath.Walk(enterpriseDir, walkFunc)
	return i18nStrings
}

func extractCmdF(command *cobra.Command, args []string) error {
	enterpriseDir, err := command.Flags().GetString("enterprise-dir")
	if err != nil {
		return errors.New("Invalid enterprise-dir parameter")
	}
	mattermostDir, err := command.Flags().GetString("mattermost-dir")
	if err != nil {
		return errors.New("Invalid mattermost-dir parameter")
	}

	i18nStrings := extractStrings(enterpriseDir, mattermostDir)

	i18nStringsList := []string{}
	for id := range i18nStrings {
		i18nStringsList = append(i18nStringsList, id)
	}
	sort.Strings(i18nStringsList)

	translations, err := getCurrentTranslations(mattermostDir)
	if err != nil {
		return err
	}

	translationsList := []string{}
	idx := map[string]bool{}
	resultMap := map[string]Translation{}
	for _, t := range translations {
		idx[t.Id] = true
		translationsList = append(translationsList, t.Id)
		resultMap[t.Id] = t
	}
	sort.Strings(translationsList)

	for _, translationKey := range i18nStringsList {
		if _, hasKey := idx[translationKey]; !hasKey {
			resultMap[translationKey] = Translation{Id: translationKey, Translation: ""}
		}
	}

	for _, translationKey := range translationsList {
		if _, hasKey := i18nStrings[translationKey]; !hasKey {
			delete(resultMap, translationKey)
		}
	}

	result := []Translation{}
	for _, t := range resultMap {
		result = append(result, t)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })

	f, err := os.Create(path.Join(mattermostDir, "i18n", "en.json"))

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	encoder.Encode(result)

	return nil
}

func checkCmdF(command *cobra.Command, args []string) error {
	enterpriseDir, err := command.Flags().GetString("enterprise-dir")
	if err != nil {
		return errors.New("Invalid enterprise-dir parameter")
	}
	mattermostDir, err := command.Flags().GetString("mattermost-dir")
	if err != nil {
		return errors.New("Invalid mattermost-dir parameter")
	}

	i18nStrings := extractStrings(enterpriseDir, mattermostDir)

	i18nStringsList := []string{}
	for id := range i18nStrings {
		i18nStringsList = append(i18nStringsList, id)
	}
	sort.Strings(i18nStringsList)

	translations, err := getCurrentTranslations(mattermostDir)
	if err != nil {
		return err
	}

	translationsList := []string{}
	idx := map[string]bool{}
	for _, t := range translations {
		idx[t.Id] = true
		translationsList = append(translationsList, t.Id)
	}
	sort.Strings(translationsList)

	for _, translationKey := range i18nStringsList {
		if _, hasKey := idx[translationKey]; !hasKey {
			fmt.Println("Added:", translationKey)
		}
	}

	for _, translationKey := range translationsList {
		if _, hasKey := i18nStrings[translationKey]; !hasKey {
			fmt.Println("Removed:", translationKey)
		}
	}
	return nil
}

func extractByFuncName(name string, args []ast.Expr) *string {
	if name == "T" {
		if len(args) == 0 {
			return nil
		}

		key, ok := args[0].(*ast.BasicLit)
		if !ok {
			return nil
		}
		return &key.Value
	} else if name == "NewAppError" {
		if len(args) < 2 {
			return nil
		}

		key, ok := args[1].(*ast.BasicLit)
		if !ok {
			return nil
		}
		return &key.Value
	} else if name == "translateFunc" {
		if len(args) < 1 {
			return nil
		}

		key, ok := args[0].(*ast.BasicLit)
		if !ok {
			return nil
		}
		return &key.Value
	} else if name == "TranslateAsHtml" {
		if len(args) < 2 {
			return nil
		}

		key, ok := args[1].(*ast.BasicLit)
		if !ok {
			return nil
		}
		return &key.Value
	} else if name == "userLocale" {
		if len(args) < 1 {
			return nil
		}

		key, ok := args[0].(*ast.BasicLit)
		if !ok {
			return nil
		}
		return &key.Value
	} else if name == "localT" {
		if len(args) < 1 {
			return nil
		}

		key, ok := args[0].(*ast.BasicLit)
		if !ok {
			return nil
		}
		return &key.Value
	}
	return nil
}

func extractForCostants(name string, value_node ast.Expr) *string {
	validConstants := map[string]bool{
		"MISSING_CHANNEL_ERROR":        true,
		"MISSING_CHANNEL_MEMBER_ERROR": true,
		"CHANNEL_EXISTS_ERROR":         true,
		"MISSING_STATUS_ERROR":         true,
		"TEAM_MEMBER_EXISTS_ERROR":     true,
		"MISSING_AUTH_ACCOUNT_ERROR":   true,
		"MISSING_ACCOUNT_ERROR":        true,
		"EXPIRED_LICENSE_ERROR":        true,
		"INVALID_LICENSE_ERROR":        true,
	}

	if _, ok := validConstants[name]; !ok {
		return nil
	}
	value, ok := value_node.(*ast.BasicLit)

	if !ok {
		return nil
	}
	return &value.Value

}

func extractFromPath(path string, info os.FileInfo, err error, i18nStrings *map[string]bool) error {
	if strings.HasPrefix(path, "./vendor") {
		return nil
	}
	if strings.HasSuffix(path, "_test.go") {
		return nil
	}
	if !strings.HasSuffix(path, ".go") {
		return nil
	}

	src, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		var id *string = nil

		switch expr := n.(type) {
		case *ast.CallExpr:
			switch fun := expr.Fun.(type) {
			case *ast.SelectorExpr:
				id = extractByFuncName(fun.Sel.Name, expr.Args)
				if id == nil {
					return true
				}
				break
			case *ast.Ident:
				id = extractByFuncName(fun.Name, expr.Args)
				break
			default:
				return true
			}
			break
		case *ast.GenDecl:
			if expr.Tok == token.CONST {
				for _, spec := range expr.Specs {
					value_spec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					if len(value_spec.Names) == 0 {
						continue
					}
					if len(value_spec.Values) == 0 {
						continue
					}
					id = extractForCostants(value_spec.Names[0].Name, value_spec.Values[0])
					if id == nil {
						continue
					}
					(*i18nStrings)[strings.Trim(*id, "\"")] = true
				}
			}
			return true
		default:
			return true
		}

		if id != nil {
			(*i18nStrings)[strings.Trim(*id, "\"")] = true
		}

		return true
	})
	return nil
}
