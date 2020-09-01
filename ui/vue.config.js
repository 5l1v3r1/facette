module.exports = {
    chainWebpack: config => {
        config.plugin("html").tap(args => {
            args[0].minify = {
                collapseBooleanAttributes: true,
                collapseWhitespace: true,
                minifyCSS: true,
                minifyJS: true,
                removeComments: true,
                removeRedundantAttributes: true,
                removeScriptTypeAttributes: true,
                removeStyleLinkTypeAttributes: true,
            };

            return args;
        });
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
};
