package fabman

import (
  "fmt"
  "net/url"
  "strings"
)

type Gender string

const (
  GenderFemale Gender = "female"
  GenderMale   Gender = "male"
  GenderOther  Gender = "other"
)

type MemberState string

const (
  MemberStateActive  MemberState = "active"
  MemberStateLocked  MemberState = "locked"
  MemberStateDeleted MemberState = "deleted"
)

// ErrMemberNotFound error returned when requested member wasn't found in the system
var ErrMemberNotFound = fmt.Errorf("member not found")

// Member information from FabMan
type Member struct {
  ID                    int               `json:"id,omitempty"`
  Account               int               `json:"account,omitempty"`
  Space                 int               `json:"space,omitempty"`
  LockVersion           int               `json:"lockVersion,omitempty"`
  UpdatedBy             int               `json:"updatedBy,omitempty"`
  CreatedAt             string            `json:"createdAt,omitempty"`
  UpdatedAt             string            `json:"updatedAt,omitempty"`
  State                 MemberState       `json:"state,omitempty"`
  MemberNumber          string            `json:"memberNumber,omitempty"`
  FirstName             string            `json:"firstName,omitempty"`
  LastName              string            `json:"lastName,omitempty"`
  Gender                Gender            `json:"gender,omitempty"`
  DateOfBirth           string            `json:"dateOfBirth,omitempty"`
  EmailAddress          string            `json:"emailAddress,omitempty"`
  Company               string            `json:"company,omitempty"`
  Phone                 string            `json:"phone,omitempty"`
  Address               string            `json:"address,omitempty"`
  Address2              string            `json:"address2,omitempty"`
  City                  string            `json:"city,omitempty"`
  Zip                   string            `json:"zip,omitempty"`
  CountryCode           string            `json:"countryCode,omitempty"`
  Region                string            `json:"region,omitempty"`
  Notes                 string            `json:"notes,omitempty"`
  BillingFirstName      string            `json:"billingFirstName,omitempty"`
  BillingLastName       string            `json:"billingLastName,omitempty"`
  BillingCompany        string            `json:"billingCompany,omitempty"`
  BillingAddress        string            `json:"billingAddress,omitempty"`
  BillingAddress2       string            `json:"billingAddress2,omitempty"`
  BillingCity           string            `json:"billingCity,omitempty"`
  BillingZip            string            `json:"billingZip,omitempty"`
  BillingCountryCode    string            `json:"billingCountryCode,omitempty"`
  BillingRegion         string            `json:"billingRegion,omitempty"`
  BillingInvoiceText    string            `json:"billingInvoiceText,omitempty"`
  Metadata              map[string]string `json:"metadata,omitempty"`
  StripeCustomer        string            `json:"stripeCustomer,omitempty"`
  UpfrontMinimumBalance string            `json:"upfrontMinimumBalance,omitempty"`
  TaxExempt             bool              `json:"taxExempt,omitempty"`
  HasBillingAddress     bool              `json:"hasBillingAddress,omitempty"`
  RequireUpfrontPayment bool              `json:"requireUpfrontPayment,omitempty"`
  AllowLogin            bool              `json:"allowLogin,omitempty"`
}

// MemberPackage holds information about package assigned to a member
type MemberPackage struct {
  ID                 int    `json:"id,omitempty"`
  PackageID          int    `json:"package"`
  LockVersion        int    `json:"lockVersion,omitempty"`
  FromDate           string `json:"fromDate,omitempty"`
  UntilDate          string `json:"untilDate,omitempty"`
  Notes              string `json:"notes,omitempty"`
  CreatedAt          string `json:"createdAt,omitempty"`
  UpdatedAt          string `json:"updatedAt,omitempty"`
  UpdatedBy          int    `json:"updatedBy,omitempty"`
  CustomFees         bool   `json:"customFees,omitempty"`
  SetupFee           string `json:"setupFee,omitempty"`
  RecurringFee       string `json:"recurringFee,omitempty"`
  RecurringFeePeriod string `json:"recurringFeePeriod,omitempty"`
  ChargedUntilDate   string `json:"chargedUntilDate,omitempty"`
}

// GetMemberByEmail returns a member with
func (client *Client) GetMemberByEmail(email string) (*Member, error) {
  params := url.Values{}
  params.Add("q", email)

  var members []*Member
  if err := client.get("members?"+params.Encode(), &members); err != nil {
    return nil, err
  }

  for _, member := range members {
    if strings.EqualFold(member.EmailAddress, email) {
      return member, nil
    }
  }

  return nil, ErrMemberNotFound
}

// CreateMember creates a new member
func (client *Client) NewMember(email, name, customerID string) (*Member, error) {
  member := &Member{
    Account:        client.account,
    FirstName:      name,
    EmailAddress:   email,
    StripeCustomer: customerID,
    Metadata: map[string]string{
      "m4m": "1",
    },
  }

  err := client.create("members", member)
  return member, err
}

// SendInvitation sends invitation email to the newly added member
func (client *Client) SendInvitation(memberID int) error {
  path := fmt.Sprintf("members/%d/invitation", memberID)
  return client.create(path, nil)
}

// UpdateMember updates member information in FabMan
func (client *Client) UpdateMember(member *Member) (*Member, error) {
  path := fmt.Sprintf("members/%d", member.ID)
  err := client.update(path, member)
  return member, err
}

// GetMemberPackages gets member's packages from FabMan
func (client *Client) GetMemberPackages(memberID int) ([]MemberPackage, error) {
  var memberPackages []MemberPackage
  path := fmt.Sprintf("members/%d/packages", memberID)
  return memberPackages, client.get(path, &memberPackages)
}

// DeleteMemberPackage deletes existing member's package
func (client *Client) DeleteMemberPackage(memberID, memberPackageID int) error {
  path := fmt.Sprintf("members/%d/packages/%d", memberID, memberPackageID)
  return client.delete(path)
}

// AddMemberPackage adds a new package to a member
func (client *Client) AddMemberPackage(memberID, packageID int, startDate, endDate string) (*MemberPackage, error) {
  path := fmt.Sprintf("members/%d/packages", memberID)

  memberPackage := &MemberPackage{
    PackageID: packageID,
    FromDate:  startDate,
    UntilDate: endDate,
  }
  return memberPackage, client.create(path, memberPackage)
}

// UpdateMemberPackage updates existing member's package
func (client *Client) UpdateMemberPackage(memberID int, memberPackage *MemberPackage) (*MemberPackage, error) {
  path := fmt.Sprintf("members/%d/packages/%d", memberID, memberPackage.ID)
  return memberPackage, client.update(path, memberPackage)
}
