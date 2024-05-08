/*
 * Sample Service
 *
 * This is a sample service
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// User - A User for Reference Service
type User struct {
	FirstName string `json:"first_name,omitempty"`

	LastName string `json:"last_name,omitempty"`

	Email string `json:"email,omitempty"`
}

// AssertUserRequired checks if the required fields are not zero-ed
func AssertUserRequired(obj User) error {
	return nil
}

// AssertUserConstraints checks if the values respects the defined constraints
func AssertUserConstraints(obj User) error {
	return nil
}
