package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/schollz/progressbar/v3"
)

type ConnectorContext struct {
	source        *SourceRequest
	target        *TargetRequest
	sourceOrg     string
	sourceProject string
	targetOrg     string
	targetProject string
}

func NewConnectorOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) ConnectorContext {
	return ConnectorContext{
		source:        sourceApi,
		target:        targetApi,
		sourceOrg:     st.SourceOrg,
		sourceProject: st.SourceProject,
		targetOrg:     st.TargetOrg,
		targetProject: st.TargetProject,
	}
}

func (c ConnectorContext) Move() error {

	var connectors []*nextgen.ConnectorInfo
	var err error
	if c.sourceOrg == "" && c.sourceProject == "" {
		// If there is no source org or project pull account level connectors
		connectors, err = c.listAccountConnectors()
		if err != nil {
			return err
		}
	} else {
		connectors, err = c.listConnectors(c.sourceOrg, c.sourceProject)
		if err != nil {
			return err
		}

	}

	bar := progressbar.Default(int64(len(connectors)), "Connectors")
	var failed []string

	for _, conn := range connectors {
		conn.OrgIdentifier = c.targetOrg
		conn.ProjectIdentifier = c.targetProject

		strConn, _ := json.Marshal(conn)
		newConn := strings.ReplaceAll(string(strConn), "account.", "org.")
		json.Unmarshal([]byte(newConn), &conn)

		err = c.createConnector(&model.CreateConnectorRequest{
			Connector: conn,
		})
		if err != nil {
			failed = append(failed, fmt.Sprintln(conn.Name, "-", err.Error()))
		}
		bar.Add(1)
	}
	bar.Finish()

	reportFailed(failed, "connectors")
	return nil
}

func (c ConnectorContext) listConnectors(org, project string) ([]*nextgen.ConnectorInfo, error) {

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier":                    api.Account,
			"orgIdentifier":                        org,
			"projectIdentifier":                    project,
			"size":                                 "1000",
			"includeAllConnectorsAvailableAtScope": "false",
		}).
		SetBody(model.ListRequestBody{
			FilterType: "Connector",
		}).
		Post(api.Url + "/ng/api/connectors/listV2")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.ListConnectorResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	connectors := []*nextgen.ConnectorInfo{}
	for _, conn := range result.Data.Content {
		if !conn.HarnessManaged {
			connectors = append(connectors, conn.Connector)
		}
	}

	return connectors, nil
}

func (c ConnectorContext) listAccountConnectors() ([]*nextgen.ConnectorInfo, error) {

	api := c.source
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"accountIdentifier":                    api.Account,
			"size":                                 "1000",
			"includeAllConnectorsAvailableAtScope": "false",
		}).
		SetBody(model.ListRequestBody{
			FilterType: "Connector",
		}).
		Post(api.Url + "/ng/api/connectors/listV2")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, handleErrorResponse(resp)
	}

	result := model.ListConnectorResponse{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	connectors := []*nextgen.ConnectorInfo{}
	for _, conn := range result.Data.Content {
		if !conn.HarnessManaged {
			connectors = append(connectors, conn.Connector)
		}
	}

	return connectors, nil
}

// processStruct recursively processes a struct value to replace account identifiers with org values
func (c ConnectorContext) processStruct(v reflect.Value, orgValue string) {
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return // Only process valid struct values
	}

	t := v.Type() // Get the struct type

	// Process each field in the struct
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)      // Get field metadata
		fieldValue := v.Field(i) // Get field value
		fieldKind := fieldValue.Kind()

		// Check if this is a string field with "account" in the name
		if strings.Contains(strings.ToLower(field.Name), "account.") &&
			fieldKind == reflect.String &&
			fieldValue.CanSet() {
			// Replace the string value with org value
			fieldValue.SetString(orgValue)
		}

		// Recursively process nested structs and pointers
		switch fieldKind {
		case reflect.Struct:
			// Process nested struct
			c.processStruct(fieldValue, orgValue)
		case reflect.Ptr:
			// Process pointer to struct
			if !fieldValue.IsNil() {
				elem := fieldValue.Elem()
				if elem.Kind() == reflect.Struct {
					c.processStruct(elem, orgValue)
				}
			}
		case reflect.Slice, reflect.Array:
			// Process each element in slice/array
			for j := 0; j < fieldValue.Len(); j++ {
				elem := fieldValue.Index(j)
				// Handle pointers and structs in arrays/slices
				if elem.Kind() == reflect.Ptr && !elem.IsNil() && elem.Elem().Kind() == reflect.Struct {
					c.processStruct(elem.Elem(), orgValue)
				} else if elem.Kind() == reflect.Struct {
					c.processStruct(elem, orgValue)
				}
			}
		case reflect.Map:
			// Process map values if the field is settable
			if fieldValue.CanSet() && fieldValue.Len() > 0 {
				// Get map key and value types
				mapKeyType := fieldValue.Type().Key()
				mapValueType := fieldValue.Type().Elem()
				
				// Create a new map with same types
				newMap := reflect.MakeMap(reflect.MapOf(mapKeyType, mapValueType))
				
				// Iterate through existing map entries
				for _, key := range fieldValue.MapKeys() {
					originalValue := fieldValue.MapIndex(key)
					
					// Process the map value based on its type
					switch originalValue.Kind() {
					case reflect.String:
						// If key contains "account", replace the string value
						strKey, ok := key.Interface().(string)
						if ok && strings.Contains(strings.ToLower(strKey), "account") {
							newMap.SetMapIndex(key, reflect.ValueOf(orgValue))
						} else {
							// Otherwise keep the original value
							newMap.SetMapIndex(key, originalValue)
						}
						
					case reflect.Ptr:
						// Handle pointers to structs
						if !originalValue.IsNil() && originalValue.Elem().Kind() == reflect.Struct {
							// Create a new pointer to a new struct value
							newStruct := reflect.New(originalValue.Elem().Type())
							// Copy the struct data
							newStruct.Elem().Set(originalValue.Elem())
							// Process the new struct
							c.processStruct(newStruct.Elem(), orgValue)
							// Add to new map
							newMap.SetMapIndex(key, newStruct)
						} else {
							// Keep original for non-struct pointers
							newMap.SetMapIndex(key, originalValue)
						}
						
					case reflect.Struct:
						// For struct values, we need special handling since map values aren't addressable
						// Create a new struct, copy the value, process it, then add to map
						newStruct := reflect.New(originalValue.Type()).Elem()
						newStruct.Set(originalValue)
						
						// The problem is we can't modify newStruct directly since it's a value not a pointer
						// We'll need a temporary solution to handle this case
						// For now, just copy the original value
						newMap.SetMapIndex(key, originalValue)
						
					default:
						// For other types, just copy the value
						newMap.SetMapIndex(key, originalValue)
					}
				}
				
				// Replace the original map with our processed map
				fieldValue.Set(newMap)
			}
		}
	}
}

func (c ConnectorContext) createConnector(connector *model.CreateConnectorRequest) error {

	api := c.target
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(connector).
		SetQueryParams(map[string]string{
			"accountIdentifier": api.Account,
		}).
		Post(api.Url + "/ng/api/connectors")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return handleErrorResponse(resp)
	}

	return nil
}
