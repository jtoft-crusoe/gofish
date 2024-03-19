//
// SPDX-License-Identifier: BSD-3-Clause
//

package redfish

import (
	"encoding/json"

	"github.com/jtoft-crusoe/gofish/common"
)

// UpdateService is used to represent the update service offered by the redfish API
type UpdateService struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// Description provides a description of this resource.
	Description string
	// FirmwareInventory points towards the firmware store endpoint
	FirmwareInventory string
	// SoftwareInventory points towards the firmware store endpoint
	SoftwareInventory string
	// HTTPPushURI endpoint is used to push (POST) firmware updates
	HTTPPushURI string `json:"HttpPushUri"`
	// MultipartHTTPPushURI endpoint is used to perform a multipart push (POST) updates
	MultipartHTTPPushURI string `json:"MultiPartHttpPushUri"`
	// ServiceEnabled indicates whether this service isenabled.
	ServiceEnabled bool
	// Status describes the status and health of a resource and its children.
	Status common.Status
	// TransferProtocol is the list of network protocols used by the UpdateService to retrieve the software image file
	TransferProtocol []string
	// UpdateServiceTarget indicates where theupdate image is to be applied.
	UpdateServiceTarget string
	// StartUpdateTarget is the endpoint which starts updating images that have been previously
	// invoked using an OperationApplyTime value of OnStartUpdateRequest.
	StartUpdateTarget string
	// OemActions contains all the vendor specific actions. It is vendor responsibility to parse this field accordingly
	OemActions json.RawMessage
	// Oem shall contain the OEM extensions. All values for properties that
	// this object contains shall conform to the Redfish Specification
	// described requirements.
	Oem json.RawMessage
	// rawData holds the original serialized JSON so we can compare updates.
	rawData []byte
}

// UnmarshalJSON unmarshals a UpdateService object from the raw JSON.
func (updateService *UpdateService) UnmarshalJSON(b []byte) error {
	type temp UpdateService
	type actions struct {
		SimpleUpdate struct {
			AllowableValues []string `json:"TransferProtocol@Redfish.AllowableValues"`
			Target          string
		} `json:"#UpdateService.SimpleUpdate"`

		// This action starts updating all images that have been previously
		// invoked using an OperationApplyTime value of OnStartUpdateRequest.
		StartUpdate struct {
			Target string
		} `json:"#UpdateService.StartUpdate"`

		Oem json.RawMessage // OEM actions will be stored here
	}
	var t struct {
		temp
		Actions           actions
		FirmwareInventory common.Link
		SoftwareInventory common.Link
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	// Extract the links to other entities for later
	*updateService = UpdateService(t.temp)
	updateService.FirmwareInventory = t.FirmwareInventory.String()
	updateService.SoftwareInventory = t.SoftwareInventory.String()
	updateService.TransferProtocol = t.Actions.SimpleUpdate.AllowableValues
	updateService.UpdateServiceTarget = t.Actions.SimpleUpdate.Target
	updateService.StartUpdateTarget = t.Actions.StartUpdate.Target
	updateService.OemActions = t.Actions.Oem
	updateService.rawData = b

	return nil
}

// GetUpdateService will get a UpdateService instance from the service.
func GetUpdateService(c common.Client, uri string) (*UpdateService, error) {
	var updateService UpdateService
	return &updateService, updateService.Get(c, uri, &updateService)
}

// SoftwareInventories gets the collection of software inventories of this update service
func (updateService *UpdateService) SoftwareInventories() ([]*SoftwareInventory, error) {
	return ListReferencedSoftwareInventories(updateService.GetClient(), updateService.SoftwareInventory)
}

// FirmwareInventories gets the collection of firmware inventories of this update service
func (updateService *UpdateService) FirmwareInventories() ([]*SoftwareInventory, error) {
	return ListReferencedSoftwareInventories(updateService.GetClient(), updateService.FirmwareInventory)
}
