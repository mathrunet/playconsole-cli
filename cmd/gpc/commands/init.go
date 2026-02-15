package commands

import (
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/apks"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/apps"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/auth"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/bundles"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/edits"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/images"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/listings"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/products"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/purchases"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/reviews"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/subscriptions"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/testing"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/tracks"
	"github.com/AndroidPoet/gpc/cmd/gpc/commands/users"
)

func init() {
	// Add all command groups to root
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(apps.AppsCmd)
	rootCmd.AddCommand(tracks.TracksCmd)
	rootCmd.AddCommand(bundles.BundlesCmd)
	rootCmd.AddCommand(apks.APKsCmd)
	rootCmd.AddCommand(listings.ListingsCmd)
	rootCmd.AddCommand(images.ImagesCmd)
	rootCmd.AddCommand(reviews.ReviewsCmd)
	rootCmd.AddCommand(products.ProductsCmd)
	rootCmd.AddCommand(subscriptions.SubscriptionsCmd)
	rootCmd.AddCommand(purchases.PurchasesCmd)
	rootCmd.AddCommand(edits.EditsCmd)
	rootCmd.AddCommand(users.UsersCmd)
	rootCmd.AddCommand(testing.TestingCmd)
}
