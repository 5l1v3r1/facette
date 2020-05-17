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
        charts: {
            axes: {
                label: "Label of the axis. @:help.common.templateSupport",
                max: "Maximum value of the axis.",
                min: "Minimum value of the axis.",
            },
            name: "Name of the chart. @:help.common.name",
            title: "Title of the chart. @:help.common.templateSupport",
        },
        common: {
            name:
                "Must start and end by an alphanumerical character, and contain alphanumerical characters, hyphens " +
                "or underscores.",
            templateSupport: "This field supports template variables.",
        },
        dashboards: {
            title: "Title of the dashboard. @:help.common.templateSupport",
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
        refresh: {
            interval: "Time interval for automatic refresh in seconds. Use `0` to disable.",
        },
        series: {
            alias:
                "Series name to display in both chart tooltip and summary. Use pattern like `{{__name__}}` or " +
                "`{{instance}}` to insert the value associated with the corresponding label.",
        },
    },

    labels: {
        adminPanel: "Administration panel",
        auto: "Automatic",
        basket: {
            _: "Basket",
            add: "Add to basket",
            clear: "Clear basket",
            preview: "Preview basket…",
            refresh: "Refresh basket",
            remove: "Remove from basket",
            saveDashboard: "Save as dashboard…",
        },
        cancel: "Cancel",
        catalog: "Catalog",
        charts: {
            _: "Chart|Charts",
            axes: {
                _: "Axes",
                max: "Max",
                min: "Min",
                x: "X",
                yLeft: "Left Y",
                yRight: "Right Y",
            },
            create: "Create chart",
            delete: "Delete chart|Delete charts",
            export: {
                _: "Export",
                imagePNG: "Save as PNG…",
                summaryCSV: "Summary as CSV…",
                summaryJSON: "Summary as JSON…",
            },
            name: "Chart name",
            new: "New chart",
            filter: "Filter charts",
            refresh: "Refresh chart",
            reset: "Reset chart",
            save: "Save chart",
            select: "Select a chart…",
            showLegend: "Show legend",
            showSummary: "Show summary",
            type: {
                _: "Type",
                area: "Area",
                bar: "Bar",
                line: "Line",
                select: "Select a type…",
            },
            zoom: {
                in: "Zoom in",
                out: "Zoom out",
            },
        },
        clearSelection: "Clear selection",
        clone: "Clone",
        close: "Close",
        color: "Color",
        connectors: {
            _: "Connector|Connectors",
            select: "Select a connector…",
        },
        continue: "Continue",
        custom: "Custom…",
        dashboards: {
            _: "Dashboard|Dashboards",
            delete: "Delete dashboard|Delete dashboards",
            edit: "Edit dashboard",
            filter: "Filter dashboards",
            goBack: "Go back to dashboard",
            name: "Dashboard name",
            new: "New dashboard",
            refresh: "Refresh dashboard",
            reset: "Reset dashboard",
            save: "Save dashboard",
            types: {
                chart: "Chart",
            },
        },
        delete: "Delete",
        displayHelp: "Display this help",
        documentation: "Documentation",
        edit: "Edit",
        empty: "Empty",
        filters: {
            _: "Filter|Filters",
            action: {
                _: "Action",
                select: "Select an action…",
            },
            add: "Add filter",
            edit: "Edit filter",
            into: "Into",
            pattern: "Pattern",
            remove: "Remove filter",
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
            charts: "Go to chart|Go to charts",
            dashboards: "Go to dashboard|Go to dashboards",
            home: "Go to home",
            metrics: "Go to metrics",
            providers: "Go to providers",
            settings: "Go to settings",
            settingsAlt: "Settings…",
        },
        help: "Help",
        home: "Home",
        info: {
            _: "Information",
            branch: "Branch",
            buildDate: "Build date",
            compiler: "Compiler",
            connectors: "Supported connectors",
            revision: "Revision",
            version: "Version",
        },
        items: {
            add: "Add item",
            remove: "Remove item",
            unsupported: "Unsupported item",
        },
        keyboardShortcuts: "Keyboard shortcuts",
        labels: "Label|Labels",
        language: {
            _: "Language",
            select: "Select a language…",
        },
        lastModified: "Last modified",
        layout: "Layout",
        leavePage: "Leave page",
        library: "Library",
        metrics: {
            _: "Metric|Metrics",
            filter: "Filter metrics",
        },
        moreActions: "More actions…",
        name: "Name",
        ok: "OK",
        openMenu: "Open menu",
        options: "Options",
        placeholders: {
            example: "e.g. {0}",
        },
        preview: "Preview",
        properties: "Properties",
        providers: {
            _: "Provider|Providers",
            delete: "Delete provider|Delete providers",
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
            interval: "Refresh interval",
            list: "Refresh list",
            next: "Next refresh in {0}",
            setInterval: "Set interval",
        },
        reset: "Reset",
        retry: "Retry",
        saveAndGo: "Save and Go",
        series: {
            _: "Series|Series",
            add: "Add series",
            alias: "Alias",
            edit: "Edit series",
            remove: "Remove series",
        },
        settings: {
            display: {
                _: "Display",
                alt: "Display settings",
            },
            personal: "Personal settings",
            save: "Save settings",
        },
        show: "Show",
        system: "System",
        templates: {
            _: "Template|Templates",
            edit: "Edit template",
            instance: "Template instance",
            newFrom: "New from template",
            save: "Save template",
            select: "Select a template…",
        },
        theme: {
            _: "Theme",
            select: "Select a theme…",
        },
        timeRange: {
            autoPropagate: "Automatically propagate time range",
            from: "From",
            multiple: "Multiple time ranges",
            propagate: "Propagate time range",
            reset: "Reset time range",
            set: "Set time range",
            to: "To",
            units: {
                days: "Last {count} day|Last {count} days",
                hours: "Last {count} hour|Last {count} hours",
                minutes: "Last {count} minute|Last {count} minutes",
                months: "Last {count} month|Last {count} months",
                years: "Last {count} year|Last {count} years",
            },
        },
        timezone: {
            _: "Time zone",
            local: "Local time",
            select: "Select a time zone…",
            utc: "UTC",
        },
        title: "Title",
        tls: {
            skipVerify: "Skip server certificate verification (Insecure)",
        },
        toggleSidebar: "Toggle sidebar",
        unnamed: "Unnamed",
        url: "URL",
        value: "Value",
        variables: {
            _: "Variables",
            add: "Add variable",
            clear: "Clear variable",
            dynamic: "Dynamic",
            edit: "Edit variable",
            fixed: "Fixed",
        },
        visit: {
            documentation: "Visit documentation",
            website: "Visit website",
        },
    },

    messages: {
        basket: {
            empty: "Basket is empty",
        },
        charts: {
            delete:
                "You are about to delete the “{name}” chart. Are you sure?|" +
                "You are about to delete {count} charts. Are you sure?",
            deleted: "Chart successfully deleted|Charts successfully deleted",
            none: "No charts defined",
            notFound: "Chart not found",
            saved: "Chart successfully saved",
            selected: "{0} chart selected|{0} charts selected",
        },
        confirmLeave: "All unsaved data will be lost. Are you sure?",
        dashboards: {
            delete:
                "You are about to delete the “{name}” dashboard. Are you sure?|" +
                "You are about to delete {count} dashboards. Are you sure?",
            deleted: "Dashboard successfully deleted|Dashboards successfully deleted",
            empty: "Dashboard is empty",
            none: "No dashboards defined",
            notFound: "Dashboard not found",
            saved: "Dashboard successfully saved",
            types: {
                loadFailed: "Cannot load type support: {0}",
            },
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
        items: {
            none: "No items found",
        },
        labels: {
            emptyDiscarded: "Labels having empty values will be discarded",
        },
        lastModified: "Last modified on {0}",
        metrics: {
            none: "No metrics found",
            selected: "{0} metric selected|{0} metrics selected",
        },
        notAvailable: "Not available",
        notDefined: "Not defined",
        providers: {
            delete:
                "You are about to delete the “{name}” provider. Are you sure?|" +
                "You are about to delete {count} providers. Are you sure?",
            deleted: "Provider successfully deleted|Providers successfully deleted",
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
            notFound: "Provider not found",
            saved: "Provider successfully saved",
            selected: "{0} provider selected|{0} providers selected",
            test: {
                error: "Provider failed to validate: {0}",
                success: "Provider successfully tested",
            },
        },
        series: {
            none: "No series defined",
        },
        settings: {
            saved: "Settings successfully saved",
        },
        templates: {
            none: "No templates defined",
        },
        variables: {
            none: "No variables defined",
        },
    },
};
