package bg

// MigrateAllSettingsMOTDToNotices migrates the deprecated "motd" setting property to the new
// "notices" settings property.
//
// TODO: Remove this migration code in Sourcegraph 3.4, which is one minor release after Sourcegraph
// 3.3. (which introduces the "motd" deprecation and this migration)
