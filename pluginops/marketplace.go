package main

import (
	"github.com/mattermost/mattermost-server/v6/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.AddCommand(marketplaceCmd)

	marketplaceCmd.AddCommand(
		marketplaceInstallCmd,
		marketplaceUnnstallCmd,
	)
}

var marketplaceCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "TODO",
}

var marketplaceInstallCmd = &cobra.Command{
	Use:   "install-all",
	Short: "Install all Marketplace plugins",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			log.WithError(err).Fatal("Failed to create client")
		}

		filter := &model.MarketplacePluginFilter{
			Page:                 0,
			PerPage:              perPage,
			EnterprisePlugins:    true,
			BuildEnterpriseReady: true,
		}

		plugins, _, err := client.GetMarketplacePlugins(filter)
		if err != nil {
			log.WithError(err).Fatal("Failed to get marketplace plugin")
		}

		log.Info("Successfully fetched all plugins")

		var g errgroup.Group

		var i int

		for _, p := range plugins {
			p := p

			g.Go(func() error {
				plugin := &model.InstallMarketplacePluginRequest{
					Id:      p.Manifest.Id,
					Version: p.Manifest.Version,
				}

				log.Infof("Requesting install of %s", p.Manifest.Name)

				_, _, err = client.InstallMarketplacePlugin(plugin)
				if err != nil {
					log.WithError(err).WithField("plugin", p.Manifest.Name).Error("Failed to install plugin")
					return err
				}

				i++

				log.Infof("Successfully installed %s. (%v of %v)", p.Manifest.Name, i, len(plugins))

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			log.Error("Some plugins failed to get installed")
			return
		}

		log.WithField("number of plugins", len(plugins)).Info("Successfully installed all Marketplace plugin")
	},
}

var marketplaceUnnstallCmd = &cobra.Command{
	Use:   "uninstall-all",
	Short: "Uninstall all Marketplace plugins",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			log.WithError(err).Fatal("Failed to create client")
		}

		filter := &model.MarketplacePluginFilter{
			Page:              0,
			PerPage:           perPage,
			EnterprisePlugins: true,
		}

		plugins, _, err := client.GetMarketplacePlugins(filter)
		if err != nil {
			log.WithError(err).Fatal("Failed to get marketplace plugin")
		}

		log.Info("Successfully fetched all plugins")

		var g errgroup.Group

		var i int

		for _, p := range plugins {
			p := p

			g.Go(func() error {
				log.Infof("Requesting uninstall of %s", p.Manifest.Name)
				_, err := client.RemovePlugin(p.Manifest.Id)
				if err != nil {
					log.WithError(err).WithField("plugin", p.Manifest.Name).Error("Failed to uninstall plugin")
					return err
				}

				i++

				log.Infof("Successfully uninstall %s. (%v of %v)", p.Manifest.Name, i, len(plugins))

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			log.Error("Some plugins failed to get uninstall")
			return
		}

		log.WithField("number of plugins", len(plugins)).Info("Successfully uninstall all Marketplace plugin")
	},
}
