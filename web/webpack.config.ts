// tslint:disable-next-line:no-reference
// <reference path="../shared/src/types/terser-webpack-plugin/index.d.ts" />

import MiniCssExtractPlugin from 'mini-css-extract-plugin'
import * as path from 'path'
import * as webpack from 'webpack'
// @ts-ignore

const mode = process.env.NODE_ENV === 'production' ? 'production' : 'development'
console.log('Using mode', mode)

const devtool = mode === 'production' ? 'source-map' : 'cheap-module-eval-source-map'

const rootDir = path.resolve(__dirname, '..')
const nodeModulesPath = path.resolve(__dirname, '..', 'node_modules')
// const monacoEditorPaths = [path.resolve(nodeModulePath, 'monaco-editor')]

const babelLoader: webpack.RuleSetUseItem = {
    loader: 'babel-loader',
    options: {
        cacheDirectory: true,
        configFile: path.join(__dirname, 'babel.config.js')
    },
}

const typescriptLoader: webpack.RuleSetUseItem = {
    loader: 'ts-loader',
    options: {
        compilerOptions: {
            target: 'es6',
            module: 'esnext',
            noEmit: false,
        },
        experimentalWatchApi: true,
        happyPackMode: true, // Typechecking is done by a separate sc process, disable here for performance
    }
}

const isEnterpriseBuild = !!process.env.ENTERPRISE
const enterpriseDir = path.resolve(__dirname, 'src', 'enterprise')
const sourceRoots = [path.resolve(__dirname, 'src'), path.resolve(rootDir, 'shared')]

const config: webpack.Configuration = {
    context: __dirname, // needed when running `gulp webpackDevServer` from the root dir
    mode,
    optimization: {},
    entry: {
        // Enterprise vs OSS builds use different entrypoints. For app (Typescript), a single enterypoint is used
        // (enterprise or OSS). For style (SCSS), the OSS entrypoint is always used, and the enterprise entrypoint
        // is appended for enterprise builds.
        app: isEnterpriseBuild ? path.join(enterpriseDir, 'main.tsx') : path.join(__dirname, 'src', 'main.tsx'),
        style: [
            path.join(__dirname, 'src', 'main.scss'),
            isEnterpriseBuild ? path.join(__dirname, 'src', 'enterprise.scss') : null,
        ].filter((path): path is string => !!path),
        'editor.worker': 'monaco-editor/esm/vs/editor/editor.worker.js',
        'json.worker': 'monaco-editor/esm/vs/language/json/json.worker',
    },
    output: {
        path: path.join(rootDir, 'ui', 'assets'),
        filename: 'script/[name].bundle.js',
        chunkFilename: 'scripts/[id]-[contenthash].chunk.js',
        publicPath: 'self',
        pathinfo: false,
    },
    devtool,
    plugins: [],
    resolve: {},
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                include: sourceRoots,
                use: [babelLoader, typescriptLoader],
            },
            {
                test: /\.(sass|scss)$/,
                use: [
                    MiniCssExtractPlugin.loader,
                    {
                        loader: 'css-loader',
                    },
                    {
                        loader: 'postcss-loader',
                        options: {
                            config: {
                                path: __dirname,
                            }
                        }
                    },
                    {
                        loader: 'sass-loader',
                        options: {
                            includePaths: [nodeModulesPath],
                        },
                    },
                ],
            },
        ],
    },
}

export default config
