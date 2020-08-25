// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

const enterpriseKeyPrefix = "ent."

type Translation struct {
	Id          string      `json:"id"`
	Translation interface{} `json:"translation"`
}

type Item struct {
	ID          string          `json:"id"`
	Translation json.RawMessage `json:"translation"`
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
	Example: "  i18n check",
	RunE:    checkCmdF,
}

var CleanCmd = &cobra.Command{
	Use:     "clean",
	Short:   "Clean empty translations",
	Long:    "Clean empty translations in translation files other than i18n/en.json base file",
	Example: "  i18n clean",
	RunE:    cleanCmdF,
}

var CleanAllCmd = &cobra.Command{
	Use:     "clean-all",
	Short:   "Clean all empty translations",
	Long:    "Clean all empty translations in translation files other than i18n/en.json base file",
	Example: "  i18n clean-all",
	RunE:    cleanAllCmdF,
}

func init() {
	ExtractCmd.Flags().Bool("skip-dynamic", false, "Whether to skip dynamically added translations")
	ExtractCmd.Flags().String("portal-dir", "../customer-web-server", "Path to folder with the Mattermost Customer Portal source code")
	ExtractCmd.Flags().String("enterprise-dir", "../enterprise", "Path to folder with the Mattermost enterprise source code")
	ExtractCmd.Flags().String("mattermost-dir", "./", "Path to folder with the Mattermost source code")
	ExtractCmd.Flags().Bool("contributor", false, "Allows contributors safely extract translations from source code without removing enterprise messages keys")

	CheckCmd.Flags().Bool("skip-dynamic", false, "Whether to skip dynamically added translations")
	CheckCmd.Flags().String("portal-dir", "../customer-web-server", "Path to folder with the Mattermost Customer Portal source code")
	CheckCmd.Flags().String("enterprise-dir", "../enterprise", "Path to folder with the Mattermost enterprise source code")
	CheckCmd.Flags().String("mattermost-dir", "./", "Path to folder with the Mattermost source code")

	CleanCmd.Flags().Bool("dry-run", false, "Whether to not apply changes")
	CleanCmd.Flags().Bool("check", false, "Whether to throw an exit code when there are empty strings present")
	CleanCmd.Flags().String("file", "de.json", "Filename, e.g. de.json, to clean empty translation strings from")
	CleanCmd.Flags().String("portal-dir", "../customer-web-server", "Path to folder with the Mattermost Customer Portal source code")
	CleanCmd.Flags().String("enterprise-dir", "../enterprise", "Path to folder with the Mattermost enterprise source code")
	CleanCmd.Flags().String("mattermost-dir", "./", "Path to folder with the Mattermost source code")

	CleanAllCmd.Flags().Bool("dry-run", false, "Whether to not apply changes")
	CleanAllCmd.Flags().Bool("check", false, "Whether to throw an exit code when there are empty strings present")
	CleanAllCmd.Flags().String("portal-dir", "../customer-web-server", "Path to folder with the Mattermost Customer Portal source code")
	CleanAllCmd.Flags().String("enterprise-dir", "../enterprise", "Path to folder with the Mattermost enterprise source code")
	CleanAllCmd.Flags().String("mattermost-dir", "./", "Path to folder with the Mattermost source code")

	I18nCmd.AddCommand(
		ExtractCmd,
		CheckCmd,
		CleanCmd,
		CleanAllCmd,
	)
	RootCmd.AddCommand(I18nCmd)
}

func getBaseFileSrcStrings(mattermostDir string) ([]Translation, error) {
	jsonFile, err := ioutil.ReadFile(path.Join(mattermostDir, "i18n", "en.json"))
	if err != nil {
		return nil, err
	}
	var translations []Translation
	_ = json.Unmarshal(jsonFile, &translations)
	return translations, nil
}

func extractSrcStrings(enterpriseDir, mattermostDir, portalDir string) map[string]bool {
	i18nStrings := map[string]bool{}
	walkFunc := func(p string, info os.FileInfo, err error) error {
		if strings.HasPrefix(p, path.Join(mattermostDir, "vendor")) {
			return nil
		}
		return extractFromPath(p, info, err, &i18nStrings)
	}
	if portalDir != "" {
		_ = filepath.Walk(portalDir, walkFunc)
	} else {
		_ = filepath.Walk(mattermostDir, walkFunc)
		_ = filepath.Walk(enterpriseDir, walkFunc)
	}
	return i18nStrings
}

func extractCmdF(command *cobra.Command, args []string) error {
	skipDynamic, err := command.Flags().GetBool("skip-dynamic")
	if err != nil {
		return errors.New("invalid skip-dynamic parameter")
	}
	enterpriseDir, err := command.Flags().GetString("enterprise-dir")
	if err != nil {
		return errors.New("invalid enterprise-dir parameter")
	}
	mattermostDir, err := command.Flags().GetString("mattermost-dir")
	if err != nil {
		return errors.New("invalid mattermost-dir parameter")
	}
	contributorMode, err := command.Flags().GetBool("contributor")
	if err != nil {
		return errors.New("invalid contributor parameter")
	}
	portalDir, err := command.Flags().GetString("portal-dir")
	if err != nil {
		return errors.New("invalid portal-dir parameter")
	}
	translationDir := mattermostDir
	if portalDir != "" {
		if enterpriseDir != "" || mattermostDir != "" {
			return errors.New("please specify EITHER portal-dir or enterprise-dir/mattermost-dir")
		}
		skipDynamic = true // dynamics are not needed for portal
		translationDir = portalDir
	}
	i18nStrings := extractSrcStrings(enterpriseDir, mattermostDir, portalDir)
	if !skipDynamic {
		addDynamicallyGeneratedStrings(&i18nStrings)
	}
	var i18nStringsList []string
	for id := range i18nStrings {
		i18nStringsList = append(i18nStringsList, id)
	}
	sort.Strings(i18nStringsList)

	sourceStrings, err := getBaseFileSrcStrings(translationDir)
	if err != nil {
		return err
	}

	var baseFileList []string
	idx := map[string]bool{}
	resultMap := map[string]Translation{}
	for _, t := range sourceStrings {
		idx[t.Id] = true
		baseFileList = append(baseFileList, t.Id)
		resultMap[t.Id] = t
	}
	sort.Strings(baseFileList)

	for _, translationKey := range i18nStringsList {
		if _, hasKey := idx[translationKey]; !hasKey {
			resultMap[translationKey] = Translation{Id: translationKey, Translation: ""}
		}
	}

	for _, translationKey := range baseFileList {
		if _, hasKey := i18nStrings[translationKey]; !hasKey {
			if contributorMode && strings.HasPrefix(translationKey, enterpriseKeyPrefix) {
				continue
			}
			delete(resultMap, translationKey)
		}
	}

	var result []Translation
	for _, t := range resultMap {
		result = append(result, t)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })

	f, err := os.Create(path.Join(mattermostDir, "i18n", "en.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(result)
	if err != nil {
		return err
	}

	return nil
}

func checkCmdF(command *cobra.Command, args []string) error {
	skipDynamic, err := command.Flags().GetBool("skip-dynamic")
	if err != nil {
		return errors.New("invalid skip-dynamic parameter")
	}
	enterpriseDir, err := command.Flags().GetString("enterprise-dir")
	if err != nil {
		return errors.New("invalid enterprise-dir parameter")
	}
	mattermostDir, err := command.Flags().GetString("mattermost-dir")
	if err != nil {
		return errors.New("invalid mattermost-dir parameter")
	}
	portalDir, err := command.Flags().GetString("portal-dir")
	if err != nil {
		return errors.New("invalid portal-dir parameter")
	}
	translationDir := mattermostDir
	if portalDir != "" {
		if enterpriseDir != "" || mattermostDir != "" {
			return errors.New("please specify EITHER portal-dir or enterprise-dir/mattermost-dir")
		}
		translationDir = portalDir
		skipDynamic = true // dynamics are not needed for portal
	}
	extractedSrcStrings := extractSrcStrings(enterpriseDir, mattermostDir, portalDir)
	if !skipDynamic {
		addDynamicallyGeneratedStrings(&extractedSrcStrings)
	}
	var extractedList []string
	for id := range extractedSrcStrings {
		extractedList = append(extractedList, id)
	}
	sort.Strings(extractedList)

	srcStrings, err := getBaseFileSrcStrings(translationDir)
	if err != nil {
		return err
	}

	var baseFileList []string
	idx := map[string]bool{}
	for _, t := range srcStrings {
		idx[t.Id] = true
		baseFileList = append(baseFileList, t.Id)
	}
	sort.Strings(baseFileList)

	changed := false
	for _, translationKey := range extractedList {
		if _, hasKey := idx[translationKey]; !hasKey {
			fmt.Println("Added:", translationKey)
			changed = true
		}
	}

	for _, translationKey := range baseFileList {
		if _, hasKey := extractedSrcStrings[translationKey]; !hasKey {
			fmt.Println("Removed:", translationKey)
			changed = true
		}
	}
	if changed {
		command.SilenceUsage = true
		return errors.New("translation source strings file out of date")
	}
	return nil
}

func addDynamicallyGeneratedStrings(i18nStrings *map[string]bool) {
	(*i18nStrings)["model.user.is_valid.pwd.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_lowercase.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_lowercase_number.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_lowercase_number_symbol.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_lowercase_symbol.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_lowercase_uppercase.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_lowercase_uppercase_number.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_lowercase_uppercase_number_symbol.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_lowercase_uppercase_symbol.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_number.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_number_symbol.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_symbol.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_uppercase.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_uppercase_number.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_uppercase_number_symbol.app_error"] = true
	(*i18nStrings)["model.user.is_valid.pwd_uppercase_symbol.app_error"] = true
	(*i18nStrings)["model.user.is_valid.id.app_error"] = true
	(*i18nStrings)["model.user.is_valid.create_at.app_error"] = true
	(*i18nStrings)["model.user.is_valid.update_at.app_error"] = true
	(*i18nStrings)["model.user.is_valid.username.app_error"] = true
	(*i18nStrings)["model.user.is_valid.email.app_error"] = true
	(*i18nStrings)["model.user.is_valid.nickname.app_error"] = true
	(*i18nStrings)["model.user.is_valid.position.app_error"] = true
	(*i18nStrings)["model.user.is_valid.first_name.app_error"] = true
	(*i18nStrings)["model.user.is_valid.last_name.app_error"] = true
	(*i18nStrings)["model.user.is_valid.auth_data.app_error"] = true
	(*i18nStrings)["model.user.is_valid.auth_data_type.app_error"] = true
	(*i18nStrings)["model.user.is_valid.auth_data_pwd.app_error"] = true
	(*i18nStrings)["model.user.is_valid.password_limit.app_error"] = true
	(*i18nStrings)["model.user.is_valid.locale.app_error"] = true
	(*i18nStrings)["January"] = true
	(*i18nStrings)["February"] = true
	(*i18nStrings)["March"] = true
	(*i18nStrings)["April"] = true
	(*i18nStrings)["May"] = true
	(*i18nStrings)["June"] = true
	(*i18nStrings)["July"] = true
	(*i18nStrings)["August"] = true
	(*i18nStrings)["September"] = true
	(*i18nStrings)["October"] = true
	(*i18nStrings)["November"] = true
	(*i18nStrings)["December"] = true
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
	} else if name == "newAppError" {
		if len(args) < 1 {
			return nil
		}
		key, ok := args[0].(*ast.BasicLit)
		if !ok {
			return nil
		}
		return &key.Value
	} else if name == "NewUserFacingError" {
		if len(args) < 1 {
			return nil
		}
		key, ok := args[0].(*ast.BasicLit)
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

func extractForConstants(name string, valueNode ast.Expr) *string {
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
	value, ok := valueNode.(*ast.BasicLit)

	if !ok {
		return nil
	}
	return &value.Value

}

func extractFromPath(path string, info os.FileInfo, err error, i18nStrings *map[string]bool) error {
	if strings.HasSuffix(path, "model/client4.go") {
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
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					if len(valueSpec.Names) == 0 {
						continue
					}
					if len(valueSpec.Values) == 0 {
						continue
					}
					id = extractForConstants(valueSpec.Names[0].Name, valueSpec.Values[0])
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

func cleanCmdF(command *cobra.Command, args []string) error {
	dryRun, err := command.Flags().GetBool("dry-run")
	if err != nil {
		return errors.New("invalid dry-run parameter")
	}
	check, err := command.Flags().GetBool("check")
	if err != nil {
		return errors.New("invalid check parameter")
	}
	file, err := command.Flags().GetString("file")
	if err != nil {
		return errors.New("invalid file parameter")
	}
	enterpriseDir, err := command.Flags().GetString("enterprise-dir")
	if err != nil {
		return errors.New("invalid enterprise-dir parameter")
	}
	mattermostDir, err := command.Flags().GetString("mattermost-dir")
	if err != nil {
		return errors.New("invalid mattermost-dir parameter")
	}
	portalDir, err := command.Flags().GetString("portal-dir")
	if err != nil {
		return errors.New("invalid portal-dir parameter")
	}
	translationDir := path.Join(mattermostDir, "i18n")
	if portalDir != "" {
		if enterpriseDir != "" || mattermostDir != "" {
			return errors.New("please specify EITHER portal-dir or enterprise-dir/mattermost-dir")
		}
		translationDir = portalDir
	}

	if filepath.Ext(file) == ".json" && file != "en.json" {
		i, err2 := clean(translationDir, file, dryRun)
		if err2 != nil {
			return err2
		}
		if check {
			log.Fatalf("%v has %v empty translations\n", file, i)
		}
	}

	return nil
}

func cleanAllCmdF(command *cobra.Command, args []string) error {
	dryRun, err := command.Flags().GetBool("dry-run")
	if err != nil {
		return errors.New("invalid dry-run parameter")
	}
	check, err := command.Flags().GetBool("check")
	if err != nil {
		return errors.New("invalid check parameter")
	}
	enterpriseDir, err := command.Flags().GetString("enterprise-dir")
	if err != nil {
		return errors.New("invalid enterprise-dir parameter")
	}
	mattermostDir, err := command.Flags().GetString("mattermost-dir")
	if err != nil {
		return errors.New("invalid mattermost-dir parameter")
	}
	portalDir, err := command.Flags().GetString("portal-dir")
	if err != nil {
		return errors.New("invalid portal-dir parameter")
	}
	translationDir := path.Join(mattermostDir, "i18n")
	if portalDir != "" {
		if enterpriseDir != "" || mattermostDir != "" {
			return errors.New("please specify EITHER portal-dir or enterprise-dir/mattermost-dir")
		}
		translationDir = portalDir
	}

	var shippedFs []string
	files, err := ioutil.ReadDir(translationDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".json" && f.Name() != "en.json" {
			shippedFs = append(shippedFs, f.Name())
		}
	}

	rs := "\n"
	for _, f := range shippedFs {
		r, err2 := clean(translationDir, f, dryRun)
		if err2 != nil {
			return err2
		}
		rs += *r
	}
	if check && rs != "\n" {
		log.Fatalf(rs)
	}
	return nil
}

func clean(translationDir string, f string, dryRun bool) (*string, error) {
	oldJ, err := ioutil.ReadFile(path.Join(translationDir, f))
	if err != nil {
		return nil, err
	}

	var ts []Item
	if err = json.Unmarshal(oldJ, &ts); err != nil {
		return nil, err
	}
	cts, i := removeEmptyTranslations(ts)
	r := ""
	if i > 0 {
		r = fmt.Sprintf("%v has %v empty translations\n", f, i)
	}

	if !dryRun {
		newJ, err := JSONMarshal(cts)
		if err != nil {
			return nil, err
		}
		filename := path.Join(translationDir, f)
		fi, err := os.Lstat(filename)
		if err != nil {
			return nil, err
		}
		if err = ioutil.WriteFile(filename, newJ, fi.Mode().Perm()); err != nil {
			return nil, err
		}
	}
	return &r, nil
}

func removeEmptyTranslations(ts []Item) ([]Item, int) {
	var k int
	var cts []Item
	for i, t := range ts {
		if string(t.Translation) != "\"\"" {
			cts = append(cts, ts[i])
		} else {
			k++
		}

	}
	return cts, k
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
