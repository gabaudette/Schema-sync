package tsgenerator

// TODO Rename this file and refactor the code to make it more readable and maintainable

import (
	"fmt"
	"github.com/pterm/pterm"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ProcessModelFiles(paths []string, modelExt string) {
	targetSuffix := modelExt + ".ts"

	spinnerInfo, _ := pterm.DefaultSpinner.Start("Process models files.")
	for _, path := range paths {
		filename := filepath.Base(path)

		// TODO Add support for file that are not in src directory
		// If the file is not in the src directory, skip it
		if !strings.Contains(path, "src") {
			continue
		}

		if strings.HasSuffix(filename, targetSuffix) {
			err := readModelFile(path)
			if err != nil {
				continue
			}
		}
	}

	spinnerInfo.Success("Model files listed successfully.")
}

func ProcessEnumFiles(paths []string, enumExt string) {
	spinnerInfo, _ := pterm.DefaultSpinner.Start("Process enum files.")
	targetSuffix := enumExt + ".ts"
	for _, path := range paths {
		filename := filepath.Base(path)
		if strings.HasSuffix(filename, targetSuffix) {
			err := readEnumFile(path)
			if err != nil {
				continue
			}
		}
	}

	spinnerInfo.Success("Enum files listed successfully.")
}

func ListInterfaceFiles(paths []string, interfaceExt string) {
	spinnerInfo, _ := pterm.DefaultSpinner.Start("Process interface files.")
	targetSuffix := interfaceExt + ".ts"
	for _, path := range paths {
		filename := filepath.Base(path)
		if strings.HasSuffix(filename, targetSuffix) {
			err := readInterfaceFile(path)
			if err != nil {
				continue
			}
		}
	}

	spinnerInfo.Success("Interface files listed successfully.")
}

func readModelFile(filePath string) error {
	//println("Reading file: ", filePath)

	// Read the file content
	content, err := os.ReadFile(filePath)

	filename := filepath.Base(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filename, err)
		return err
	}

	// Convert the content to a string for easier processing
	strContent := string(content)

	// Initialize variables to store class name, properties, and relationships
	var className string
	var properties []map[string]string
	var relationships []map[string]string

	err = addDefaultHelperClassOrInterface()
	if err != nil {
		return err
	}

	// Split file content into lines for easier processing
	lines := strings.Split(strContent, "\n")
	for i := 0; i < len(lines); i++ {
		trimmedLine := strings.TrimSpace(lines[i])

		// Find the class name
		if strings.HasPrefix(trimmedLine, "export class") {
			parts := strings.Fields(trimmedLine)
			if len(parts) >= 3 {
				className = parts[2]
			}
		}

		// Check for properties annotated with @Column
		if strings.HasPrefix(trimmedLine, string(COLUMN_TOKEN)) {
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				propParts := strings.Split(nextLine, ":")
				if len(propParts) == 2 {
					propertyName := strings.TrimSpace(propParts[0])

					// If the propertyName is "id" "createdAt" or "updatedAt" or "deletedAt" continue
					// These are the properties that are automatically added by the DatabaseEntities class
					if propertyName == "id" || propertyName == "createdAt" || propertyName == "updatedAt" || propertyName == "deletedAt" {
						continue
					}

					// If the propertyName start with a get or set continue
					if strings.HasPrefix(propertyName, "get") || strings.HasPrefix(propertyName, "set") {
						continue
					}

					// If the propertyName start with a _ continue
					if strings.HasPrefix(propertyName, "_") {
						continue
					}

					// If the propertyName start with a @ continue
					if strings.HasPrefix(propertyName, "@") {
						continue
					}

					propertyType := strings.TrimSpace(strings.TrimSuffix(propParts[1], ";"))
					properties = append(properties, map[string]string{"name": propertyName, "type": propertyType})
				}
			}
		}

		// Check for relationships annotated with @BelongsTo or @HasMany
		if strings.HasPrefix(trimmedLine, string(BELONGS_TO_TOKEN)) || strings.HasPrefix(trimmedLine, string(HAS_MANY_TOKEN)) {
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				relParts := strings.Split(nextLine, ":")
				if len(relParts) == 2 {
					relationshipName := strings.TrimSpace(relParts[0])
					relationshipType := strings.TrimSpace(strings.TrimSuffix(relParts[1], ";"))
					relationships = append(relationships, map[string]string{"name": relationshipName, "type": relationshipType})
				}
			}
		}
	}

	model := TSModel{ClassName: className, Properties: []TSProperty{}, Relationships: []TSRelationship{}}

	for _, property := range properties {
		model.Properties = append(model.Properties, TSProperty{Name: property["name"], Type: TSTypeToken(property["type"])})
	}

	for _, relationship := range relationships {
		model.Relationships = append(model.Relationships, TSRelationship{Name: relationship["name"], Type: APIModelToken(relationship["type"])})
	}

	model.ApiFilePath = filePath

	err = generateTSModelsFile(model)

	return nil
}

func readEnumFile(filePath string) error {
	//println("Reading file: ", filePath)

	// Read the file content
	content, err := os.ReadFile(filePath)

	filename := filepath.Base(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filename, err)
	}

	// Convert the content to a string for easier processing
	strContent := string(content)

	// Initialize variables to store enum name and values
	var enumName string
	var values []string

	// Split file content into lines for easier processing
	lines := strings.Split(strContent, "\n")
	for i := 0; i < len(lines); i++ {
		enumName = strings.TrimSpace(lines[i])

		if strings.HasPrefix(enumName, "export enum") {
			parts := strings.Fields(enumName)
			if len(parts) >= 3 {
				enumName = parts[2]
			}

			// Check for values
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if strings.HasPrefix(nextLine, "{") {
					for j := i + 2; j < len(lines); j++ {
						value := strings.TrimSpace(lines[j])
						if strings.HasPrefix(value, "}") {
							break
						}
						values = append(values, value)
					}
				}
			}

			enum := TSEnum{EnumName: enumName, Values: values}

			enum.ApiFilePath = filePath

			err = generateTSEnumsFile(enum)

			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func readInterfaceFile(filePath string) error {
	//println("Reading file: ", filePath)

	// Read the file content
	content, err := os.ReadFile(filePath)

	filename := filepath.Base(filePath)

	if err != nil {
		log.Printf("Error reading file %s: %v", filename, err)
	}

	// Convert the content to a string for easier processing
	strContent := string(content)

	// Initialize variables to store interface name and properties

	var interfaceName string
	var properties []map[string]string

	// Split file content into lines for easier processing
	lines := strings.Split(strContent, "\n")

	for i := 0; i < len(lines); i++ {
		interfaceName = strings.TrimSpace(lines[i])

		if strings.HasPrefix(interfaceName, "export interface") {
			parts := strings.Fields(interfaceName)
			if len(parts) >= 3 {
				interfaceName = parts[2]
			}

			// Check for properties
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if strings.HasPrefix(nextLine, "{") {
					for j := i + 2; j < len(lines); j++ {
						prop := strings.TrimSpace(lines[j])
						if strings.HasPrefix(prop, "}") {
							break
						}
						propParts := strings.Split(prop, ":")
						if len(propParts) == 2 {
							propertyName := strings.TrimSpace(propParts[0])
							propertyType := strings.TrimSpace(strings.TrimSuffix(propParts[1], ";"))
							properties = append(properties, map[string]string{"name": propertyName, "type": propertyType})
						}
					}
				}
			}

			tsInterface := TSInterface{InterfaceName: interfaceName, Properties: []TSProperty{}}

			for _, property := range properties {
				tsInterface.Properties = append(tsInterface.Properties, TSProperty{Name: property["name"], Type: TSTypeToken(property["type"])})
			}

			tsInterface.ApiFilePath = filePath

			err = generateTSInterfacesFile(tsInterface)

			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func generateTSModelsFile(model TSModel) error {
	file, err := os.OpenFile("types.d.ts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	fileContent := parseModelToTS(model)
	if fileContent == "" {
		return nil
	}

	if _, err := file.WriteString(fileContent); err != nil {
		log.Fatal(err)
	}

	if err := file.Sync(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func generateTSEnumsFile(enum TSEnum) error {
	file, err := os.OpenFile("types.d.ts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	fileContent := parseEnumToTS(enum)
	if fileContent == "" {
		return nil
	}

	if _, err := file.WriteString(fileContent); err != nil {
		log.Fatal(err)
	}

	if err := file.Sync(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func generateTSInterfacesFile(tsInterface TSInterface) error {
	file, err := os.OpenFile("types.d.ts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	fileContent := parseInterfaceToTS(tsInterface)
	if fileContent == "" {
		return nil
	}

	if _, err := file.WriteString(fileContent); err != nil {
		log.Fatal(err)
	}

	if err := file.Sync(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func parseModelToTS(model TSModel) string {
	if model.ClassName == "" {
		return ""
	}

	if len(model.Properties) == 0 && len(model.Relationships) == 0 {
		fileContent := "// tslint:disable-next-line:no-empty-interface\n"

		// Add a formatted comment to the file that describes the model and tells the user that the model is empty
		// and that it might be caused by the model being deprecated, model in backend is empty (not yet implemented) or no properties or relationships were found or no columns were annotated with @Column
		// or that the user need to sync again
		fileContent += fmt.Sprintf("/**\n * Model [%s]:\n\n", model.ClassName)
		fileContent += fmt.Sprintf("API file path %s \n", model.ApiFilePath)
		fileContent += " * This model is empty\n"
		fileContent += " * This might be caused by a empty backend  model (eg. not yet implemented), no properties or relationships were found in the backend model (eg. no columns annoted with  @Column)\n"
		fileContent += " * or that you need to sync again\n"
		fileContent += " */\n"

		fileContent += fmt.Sprintf("export interface %s extends DatabaseEntities {\n", model.ClassName)
		fileContent += "    // This model is empty\n"
		fileContent += "}\n\n"
		return fileContent
	}

	fileContent := fmt.Sprintf("/**\n * Model [%s]:\n\n", model.ClassName)
	fileContent += fmt.Sprintf("API file path %s \n", model.ApiFilePath)

	if len(model.Properties) > 0 {
		fileContent += " * Properties:\n\n"
		for _, property := range model.Properties {
			propType := removeNonAttributeFromType(string(property.Type))
			fileContent += fmt.Sprintf(" * - %s: %s\n", property.Name, propType)
		}
		fileContent += "\n\n"
	}

	if len(model.Relationships) > 0 {
		fileContent += " * Relationships:\n\n"
		for _, relationship := range model.Relationships {

			relationshipType := removeNonAttributeFromType(string(relationship.Type))
			fileContent += fmt.Sprintf(" * - %s: %s\n", relationship.Name, relationshipType)
		}
	}

	fileContent += " */\n"
	fileContent += fmt.Sprintf("export interface %s extends DatabaseEntities {\n", model.ClassName)

	for _, property := range model.Properties {
		if property.Type == "GeoPoint" {
			property.Type = "GeoJson"
		}

		if property.Type == "Date" {
			property.Type = "Date | string"
		}

		if property.Type == "Date[]" {
			property.Type = "Date[] | string[]"
		}

		if property.Type == "GeoPoint[]" {
			property.Type = "GeoJson[]"
		}

		propType := removeNonAttributeFromType(string(property.Type))

		fileContent += fmt.Sprintf("    %s: %s;\n", property.Name, propType)
	}

	for _, relationship := range model.Relationships {
		relationshipType := removeNonAttributeFromType(string(relationship.Type))
		fileContent += fmt.Sprintf("    %s: %s;\n", relationship.Name, relationshipType)
	}

	fileContent += "}\n\n"

	//print("TS Content: " + fileContent)

	return fileContent
}

func parseEnumToTS(enum TSEnum) string {
	if enum.EnumName == "" {
		return ""
	}

	if len(enum.Values) == 0 {
		fileContent := fmt.Sprintf("/**\n * Enum [%s]:\n\n", enum.EnumName)
		fileContent += fmt.Sprintf("API file path %s \n", enum.ApiFilePath)
		fileContent += " * This enum is empty\n"
		fileContent += " * This might be caused by asa empty backend enum (eg. not yet implemented), no values were found in the backend enum\n"
		fileContent += " * or that you need to sync again\n"
		fileContent += " */\n"

		fileContent += fmt.Sprintf("export enum %s {\n", enum.EnumName)
		fileContent += "    // This enum is empty\n"
		fileContent += "}\n\n"
		return fileContent
	}

	fileContent := fmt.Sprintf("/**\n * Enum [%s]:\n\n", enum.EnumName)
	fileContent += fmt.Sprintf("API file path %s \n", enum.ApiFilePath)

	// FIXME: This is a bug, the values are not being printed
	if len(enum.Values) > 0 {
		fileContent += " * Values:\n\n"
		for _, value := range enum.Values {
			fileContent += fmt.Sprintf(" * - %s\n", value)
		}
		fileContent += "\n\n"
	}

	fileContent += " */\n"

	fileContent += fmt.Sprintf("export enum %s {\n", enum.EnumName)
	for _, value := range enum.Values {
		fileContent += fmt.Sprintf("    %s,\n", value)
	}

	fileContent += "}\n\n"

	//print("TS Content: " + fileContent)

	return fileContent
}

func parseInterfaceToTS(tsInterface TSInterface) string {
	if tsInterface.InterfaceName == "" {
		return ""
	}

	if len(tsInterface.Properties) == 0 {
		fileContent := fmt.Sprintf("/**\n * Interface [%s]:\n\n", tsInterface.InterfaceName)
		fileContent += fmt.Sprintf("API file path %s \n", tsInterface.ApiFilePath)
		fileContent += " * This interface is empty\n"
		fileContent += " * This might be caused by a empty backend interface (eg. not yet implemented), no properties were found in the backend interface\n"
		fileContent += " * or that you need to sync again\n"
		fileContent += " */\n"

		fileContent += fmt.Sprintf("export interface %s {\n", tsInterface.InterfaceName)
		fileContent += "    // This interface is empty\n"
		fileContent += "}\n\n"
		return fileContent
	}

	fileContent := fmt.Sprintf("/**\n * Interface [%s]:\n\n", tsInterface.InterfaceName)
	fileContent += fmt.Sprintf("API file path %s \n", tsInterface.ApiFilePath)

	if len(tsInterface.Properties) > 0 {
		fileContent += " * Properties:\n\n"
		for _, property := range tsInterface.Properties {
			if property.Type == "GeoPoint" {
				property.Type = "GeoJson"
			}

			if property.Type == "Date" {
				property.Type = "Date | string"
			}

			if property.Type == "Date[]" {
				property.Type = "Date[] | string[]"
			}

			if property.Type == "GeoPoint[]" {
				property.Type = "GeoJson[]"
			}

			// If the property is wrapped like this NonAttribute<Prop> or NonAttribute<Prop>[] remove the NonAttribute<> part and keep the Prop
			propType := removeNonAttributeFromType(string(property.Type))
			fileContent += fmt.Sprintf(" * - %s: %s\n", property.Name, propType)
		}
		fileContent += "\n\n"
	}

	fileContent += " */\n"
	fileContent += fmt.Sprintf("export interface %s {\n", tsInterface.InterfaceName)

	for _, property := range tsInterface.Properties {
		propType := removeNonAttributeFromType(string(property.Type))
		fileContent += fmt.Sprintf("    %s: %s;\n", property.Name, propType)
	}

	fileContent += "}\n\n"

	//print("TS Content: " + fileContent)

	return fileContent
}

func addDefaultHelperClassOrInterface() error {
	file, err := os.OpenFile("types.d.ts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	fc, err := os.ReadFile(file.Name())

	// Check if the DatabaseEntities class has already been added to the file
	// If it has been added, return
	if strings.Contains(string(fc), "export class DatabaseEntities") {
		return nil
	}

	// Check if the GeoJson interface has already been added to the file
	// If it has been added, return
	if strings.Contains(string(fc), "export interface GeoJson") {
		return nil
	}

	fileContent := "import { Type } from \"class-transformer\";\n\n"
	fileContent += fmt.Sprintf("export class DatabaseEntities {\n")
	fileContent += "    id: number;\n"
	fileContent += "    @Type(() => Date)\n"
	fileContent += "    createdAt: Date;\n"
	fileContent += "    @Type(() => Date)\n"
	fileContent += "    updatedAt: Date;\n"
	fileContent += "    @Type(() => Date)\n"
	fileContent += "    deletedAt: Date;\n"
	fileContent += "    constructor(value: Partial<DatabaseEntities>) {\n"
	fileContent += "        if (value) {\n"
	fileContent += "            Object.assign(this, value);\n"
	fileContent += "        }\n"
	fileContent += "    }\n"
	fileContent += "}\n\n"

	fileContent += "export interface GeoJson {\n"
	fileContent += "    type: \"Point\";\n"
	fileContent += "    crs: {\n"
	fileContent += "        type: string;\n"
	fileContent += "        properties: {\n"
	fileContent += "            name: string;\n"
	fileContent += "        };\n"
	fileContent += "    };\n"
	fileContent += "    coordinates: [number, number];\n"
	fileContent += "}\n\n"

	if _, err := file.WriteString(fileContent); err != nil {
		log.Fatal(err)
	}

	if err := file.Sync(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func removeNonAttributeFromType(typeString string) string {
	pattern := `(\w+)\s+NonAttribute<([^<>]+)>`
	re := regexp.MustCompile(pattern)

	output := re.ReplaceAllStringFunc(typeString, func(match string) string {
		matches := re.FindStringSubmatch(match)
		if len(matches) > 2 {
			prefix := matches[1]      // Captured group before "NonAttribute"
			dynamicPart := matches[2] // Captured group inside "<>"
			return fmt.Sprintf("%s %s", prefix, dynamicPart)
		}
		return match
	})

	return output
}
