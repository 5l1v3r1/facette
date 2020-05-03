/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

export default {
    name: "English",

    date: {
        long: "MMM D YYYY, HH:mm:ss",
    },

    help: {
        common: {
            name:
                "Must start and end by an alphanumerical character, and contain alphanumerical characters, hyphens " +
                "or underscores.",
            templateSupport: "This field supports template variables.",
        },
        filters: {
            action:
                "Action to be performed by the filter:\n" +
                "* discard: drops matching metrics\n" +
                "* relabel: relabels matching metrics\n" +
                "* rewrite: rewrites label of matching metrics\n" +
                "* sieve: keeps only matching metrics\n",
            into: "Replacement value to apply on the value associated with the filter label.",
            label: "Label on which the filter will be applied.",
            pattern: "Pattern to apply on the value associated with the filter label.",
        },
        providers: {
            name: "Name of the provider. @:help.common.name",
            pollInterval: "Provider metrics automatic polling interval. Disabled if empty.",
            prometheus: {
                filter: "Filter for querying metrics from upstream Prometheus service.",
            },
            rrdtool: {
                path: "Base directory from which to search for files.",
                pattern: "Pattern to apply to found files paths. Must follow the RE2 syntax.",
                daemon: "rrdcached daemon socket address.",
            },
            url: "URL to the upstream {0} service.",
        },
    },

    labels: {
        adminPanel: "Administration panel",
        cancel: "Cancel",
        catalog: "Catalog",
        charts: {
            _: "Chart|Charts",
            create: "Create chart",
            new: "New chart",
            filter: "Filter charts",
        },
        clearSelection: "Clear selection",
        clone: "Clone",
        close: "Close",
        connectors: {
            _: "Connector|Connectors",
            select: "Select a connector…",
        },
        dashboards: {
            _: "Dashboard|Dashboards",
            new: "New dashboard",
            filter: "Filter dashboards",
        },
        delete: "Delete",
        displayHelp: "Display this help",
        documentation: "Documentation",
        filters: {
            _: "Filter|Filters",
            action: {
                _: "Action",
                select: "Select an action…",
            },
            add: "Add filter",
            into: "Into",
            pattern: "Pattern",
            targets: {
                _: "Targets",
                add: "Add target",
            },
        },
        fullscreen: {
            enter: "Enter full screen",
            leave: "Leave full screen",
            toggle: "Toggle full screen",
        },
        general: "General",
        goto: {
            adminPanel: "Go to administration panel",
            home: "Go to home",
            metrics: "Go to metrics",
            providers: "Go to providers",
            settings: "Go to settings",
            settingsAlt: "Settings…",
        },
        help: "Help",
        info: {
            _: "Information",
            branch: "Branch",
            buildDate: "Build date",
            compiler: "Compiler",
            connectors: "Supported connectors",
            revision: "Revision",
            version: "Version",
        },
        keyboardShortcuts: "Keyboard shortcuts",
        labels: "Label|Labels",
        language: {
            _: "Language",
            select: "Select a language…",
        },
        lastModified: "Last modified",
        leavePage: "Leave page",
        library: "Library",
        metrics: {
            _: "Metric|Metrics",
            filter: "Filter metrics",
        },
        name: "Name",
        ok: "OK",
        openMenu: "Open menu",
        placeholders: {
            example: "e.g. {0}",
        },
        preview: "Preview",
        providers: {
            _: "Provider|Providers",
            disable: "Disable",
            disabled: "Providers is disabled",
            enable: "Enable",
            enabled: "Providers is enabled",
            filter: "Filter providers",
            name: "Provider name",
            new: "New provider",
            poll: "Poll",
            pollAlt: "Poll providers",
            pollInterval: "Poll interval",
            reset: "Reset provider",
            rrdtool: {
                path: "Path",
                pattern: "Pattern",
                daemon: "Daemon address",
            },
            save: "Save provider",
            test: "Test provider",
        },
        refresh: {
            list: "Refresh list",
        },
        reset: "Reset",
        retry: "Retry",
        saveAndGo: "Save and Go",
        settings: {
            display: {
                _: "Display",
                alt: "Display settings",
            },
            personal: "Personal settings",
            save: "Save settings",
        },
        system: "System",
        templates: {
            _: "Template|Templates",
            newFrom: "New from template",
        },
        theme: {
            _: "Theme",
            select: "Select a theme…",
        },
        timezone: {
            _: "Time zone",
            local: "Local time",
            select: "Select a time zone…",
            utc: "UTC",
        },
        tls: {
            skipVerify: "Skip server certificate verification (Insecure)",
        },
        toggleSidebar: "Toggle sidebar",
        url: "URL",
        visit: {
            documentation: "Visit documentation",
            website: "Visit website",
        },
    },

    messages: {
        charts: {
            none: "No charts defined",
        },
        confirmLeave: "All unsaved data will be lost. Are you sure?",
        dashboards: {
            none: "No dashboards defined",
        },
        data: {
            none: "No data found",
        },
        documentation:
            "Documentation regarding Facette’s architecture, its usage and REST API is available on a dedicated " +
            "website.",
        error: {
            _: "Error: {0}",
            bulk: "An error occurred during bulk execution",
            load: "An error occurred loading data",
            unhandled: "An unhandled error has occurred",
        },
        filters: {
            none: "No provider filters defined",
        },
        lastModified: "Last modified on {0}",
        metrics: {
            none: "No metrics found",
            selected: "{0} metric selected|{0} metrics selected",
        },
        notAvailable: "Not available",
        providers: {
            disable:
                "You are about to disable the “{name}” provider. Are you sure?|" +
                "You are about to disable {count} providers. Are you sure?",
            disabled: "Provider successfully disabled|Providers successfully disabled",
            enable:
                "You are about to enable the “{name}” provider. Are you sure?|" +
                "You are about to enable {count} providers. Are you sure?",
            enabled: "Provider successfully enabled|Providers successfully enabled",
            loadFailed: "Cannot load provider support: {0}",
            none: "No providers defined",
            saved: "Provider successfully saved",
            selected: "{0} provider selected|{0} providers selected",
            test: {
                error: "Provider failed to validate: {0}",
                success: "Provider successfully validated",
            },
        },
        settings: {
            saved: "Settings successfully saved",
        },
        templates: {
            none: "No templates defined",
        },
    },
};
