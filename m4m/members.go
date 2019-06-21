package m4m

import (
  "time"

  "github.com/stripe/stripe-go"
  stripeCustomers "github.com/stripe/stripe-go/customer"
  stripeProducts "github.com/stripe/stripe-go/product"

  "github.com/Marineholmen-Makerspace/members-management/fabman"
)

const packageExpirationGracePeriod = 60 * 60 * 24 * 30

// syncMemberState synchronizes member state to reflect subscription in Stripe
func syncMemberState(customerID string) error {
  customer, err := stripeCustomers.Get(customerID, nil)
  if err != nil {
    return err
  }

  if customer.Subscriptions == nil || len(customer.Subscriptions.Data) == 0 {
    LogError("customer %s does not have any subscriptions", customerID)
    return nil
  }

  subscription := customer.Subscriptions.Data[0]
  productID := subscription.Plan.Product.ID
  product, err := stripeProducts.Get(productID, nil)
  if err != nil {
    return err
  }

  cfg := GetConfig()
  fm := fabman.NewClient(cfg.FabMan.Account, cfg.FabMan.Token)

  pkg, err := getOrCreatePackage(fm, product)
  if err != nil {
    return err
  }

  member, err := getOrCreateMember(fm, customer)
  if err != nil {
    LogError("cannot get a member for customer %s: %v", customerID, err)
    return nil
  }

  memberPackages, err := fm.GetMemberPackages(member.ID)
  if err != nil {
    LogError("cannot get a member's packages for %d: %v", member.ID, err)
    return nil
  }

  startDate := time.Unix(subscription.CurrentPeriodStart, 0).Format("2006-01-02")
  endDateTs := subscription.CurrentPeriodEnd + packageExpirationGracePeriod
  endDate := time.Unix(endDateTs, 0).Format("2006-01-02")

  packageDoesNotExist := true
  for _, memberPkg := range memberPackages {
    if memberPkg.PackageID == pkg.ID {
      packageDoesNotExist = false

      if memberPkg.UntilDate != endDate {
        memberPkg.UntilDate = endDate
        LogInfo("Extending package expiration %d for member %d", pkg.ID, member.ID)
        if _, err := fm.UpdateMemberPackage(member.ID, &memberPkg); err != nil {
          return err
        }
      }
    }
  }

  if packageDoesNotExist {
    LogInfo("Adding package %d to member %d", pkg.ID, member.ID)
    if _, err := fm.AddMemberPackage(member.ID, pkg.ID, startDate, endDate); err != nil {
      return err
    }
  }

  return nil
}

// getOrCreatePackage tries to find existing package with matching name, otherwise returns a newly created one
func getOrCreatePackage(fm *fabman.Client, product *stripe.Product) (*fabman.Package, error) {
  packages, err := fm.GetPackages()
  if err != nil {
    return nil, err
  }

  for _, pkg := range packages {
    if pkg.Name == product.Name {
      return &pkg, nil
    }
  }

  return fm.NewPackage(product.Name)

}

// getOrCreateMember tries to find existing member with matching email, otherwise returns a newly created one
func getOrCreateMember(fm *fabman.Client, customer *stripe.Customer) (*fabman.Member, error) {

  // get member information from FabMan
  member, err := fm.GetMemberByEmail(customer.Email)
  if err == fabman.ErrMemberNotFound {
    member, err = fm.NewMember(customer.Email, customer.Description, customer.ID)
    if err != nil {
      return nil, err
    }

    if err = fm.SendInvitation(member.ID); err != nil {
      return nil, err
    }
  } else if err != nil {
    return nil, err
  }

  // update stripe's customer ID if not set
  if member.StripeCustomer == "" {
    member.StripeCustomer = customer.ID
    if member, err = fm.UpdateMember(member); err != nil {
      return nil, err
    }
  }

  return member, nil
}
