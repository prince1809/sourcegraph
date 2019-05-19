declare module 'terser-webpack-plugin' {
    import {Plugin} from 'webpack';

    declare class TerserPlugin extends Plugin {
        constructor(options?: TerserPlugin.TerserPluginOptions)
    }

    declare namespace TerserPlugin {
        interface TerserPluginOptions {
            test?: RegExp | RegExp[]
            sourceMap?: boolean
            terserOptions?: TerserOptions
        }

        interface TerserOptions {
            ie8?: boolean
            compress?: boolean | object
        }

        interface ExtractCommentOptions {

        }
    }

    export = TerserPlugin
}
