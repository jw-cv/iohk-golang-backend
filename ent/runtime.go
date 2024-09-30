// Code generated by ent, DO NOT EDIT.

package ent

import (
	"iohk-golang-backend/ent/customer"
	"iohk-golang-backend/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	customerFields := schema.Customer{}.Fields()
	_ = customerFields
	// customerDescName is the schema descriptor for name field.
	customerDescName := customerFields[1].Descriptor()
	// customer.NameValidator is a validator for the "name" field. It is called by the builders before save.
	customer.NameValidator = func() func(string) error {
		validators := customerDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// customerDescSurname is the schema descriptor for surname field.
	customerDescSurname := customerFields[2].Descriptor()
	// customer.SurnameValidator is a validator for the "surname" field. It is called by the builders before save.
	customer.SurnameValidator = func() func(string) error {
		validators := customerDescSurname.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(surname string) error {
			for _, fn := range fns {
				if err := fn(surname); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// customerDescNumber is the schema descriptor for number field.
	customerDescNumber := customerFields[3].Descriptor()
	// customer.NumberValidator is a validator for the "number" field. It is called by the builders before save.
	customer.NumberValidator = customerDescNumber.Validators[0].(func(int) error)
	// customerDescCountry is the schema descriptor for country field.
	customerDescCountry := customerFields[5].Descriptor()
	// customer.CountryValidator is a validator for the "country" field. It is called by the builders before save.
	customer.CountryValidator = func() func(string) error {
		validators := customerDescCountry.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(country string) error {
			for _, fn := range fns {
				if err := fn(country); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// customerDescDependants is the schema descriptor for dependants field.
	customerDescDependants := customerFields[6].Descriptor()
	// customer.DefaultDependants holds the default value on creation for the dependants field.
	customer.DefaultDependants = customerDescDependants.Default.(int)
	// customer.DependantsValidator is a validator for the "dependants" field. It is called by the builders before save.
	customer.DependantsValidator = customerDescDependants.Validators[0].(func(int) error)
	// customerDescID is the schema descriptor for id field.
	customerDescID := customerFields[0].Descriptor()
	// customer.IDValidator is a validator for the "id" field. It is called by the builders before save.
	customer.IDValidator = customerDescID.Validators[0].(func(int) error)
}
