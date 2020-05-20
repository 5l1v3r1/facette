/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

const merge = require("lodash/merge");
const path = require("path");

module.exports = {
    configureWebpack: config => {
        merge(config, {
            performance: {
                hints: false,
            },
            resolve: {
                alias: {
                    "@": __dirname,
                },
            },
        });

        if (process.env.NODE_ENV === "production") {
            Object.assign(config.optimization.minimizer[0].options.terserOptions, {
                keep_fnames: true,
            });
        }
    },
    devServer: {
        overlay: false,
        proxy: {
            "/api": {
                changeOrigin: true,
                target: "http://localhost:12003",
            },
        },
    },
    outputDir: path.resolve(__dirname, "../dist/assets"),
};
