package fabman

// Package contains package configuration
type Package struct {
	ID                   int                 `json:"id,omitempty"`
	Name                 string              `json:"name"`
	Notes                string              `json:"notes,omitempty"`
	State                string              `json:"state,omitempty"`
	LockVersion          int                 `json:"lockVersion,omitempty"`
	CreatedAt            string              `json:"createdAt,omitempty"`
	UpdatedAt            string              `json:"updatedAt,omitempty"`
	UpdatedBy            int                 `json:"updatedBy,omitempty"`
	Description          string              `json:"description,omitempty"`
	Account              int                 `json:"account,omitempty"`
	AllowsBooking        bool                `json:"allowsBooking,omitempty"`
	SetupFee             string              `json:"setupFee,omitempty"`
	RecurringFee         string              `json:"recurringFee,omitempty"`
	RecurringFeePeriod   string              `json:"recurringFeePeriod,omitempty"`
	ProrateLastPeriod    bool                `json:"prorateLastPeriod,omitempty"`
	MinimumDuration      int                 `json:"minimumDuration,omitempty"`
	MinimumDurationUnit  string              `json:"minimumDurationUnit,omitempty"`
	EndsAutomatically    bool                `json:"endsAutomatically,omitempty"`
	CancellationTime     int                 `json:"cancellationTime,omitempty"`
	CancellationTimeUnit string              `json:"cancellationTimeUnit,omitempty"`
	Permissions          []PackagePermission `json:"permissions,omitempty"`
}

// Permission holds permissions configuration for a package
type PackagePermission struct {
	TimeType                  string                  `json:"timeType,omitempty"`
	Type                      string                  `json:"type,omitempty"`
	ResourceType              int                     `json:"resourceType,omitempty"`
	Resource                  int                     `json:"resource,omitempty"`
	UsageFeeDiscountPercent   string                  `json:"usageFeeDiscountPercent,omitempty"`
	BookingFeeDiscountPercent string                  `json:"bookingFeeDiscountPercent,omitempty"`
	Times                     []PackagePermissionTime `json:"times,omitempty"`
}

// PackagePermissionTime informs when permission is active
type PackagePermissionTime struct {
	DayOfWeek int    `json:"dayOfWeek,omitempty"`
	FromTime  string `json:"fromTime,omitempty"`
	UntilTime string `json:"untilTime,omitempty"`
}

// GetPackages returns all currently configured packages
func (client *Client) GetPackages() ([]Package, error) {
	var packages []Package
	err := client.get("packages", &packages)
	return packages, err
}

// NewPackage creates and returns a new package
func (client *Client) NewPackage(name string) (*Package, error) {
	pkg := &Package{
		Account: client.account,
		Name: name,
	}

	return pkg, client.create("packages", pkg)
}