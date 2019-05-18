import log from 'fancy-log'
import gulp from 'gulp'
import createWebpackCompiler, {Stats} from 'webpack'
import webpackConfig from './webpack.config'

const WEBPACK_STATS_OPTIONS: Stats.ToStringOptions & { colors?: boolean } = {}

const logWebpackStats = (stats: Stats) => log(stats.toString(WEBPACK_STATS_OPTIONS))

export async function webpack(): Promise<void> {
    const compiler = createWebpackCompiler(webpackConfig)
    const stats = await new Promise<Stats>((resolve, reject) => {
        compiler.run((err, stats) => (err ? reject(err) : resolve(stats)))
    })
    logWebpackStats(stats)
    if (stats.hasErrors()) {
        throw Object.assign(new Error('Failed to compile'), {showStack: false})
    }
}

/**
 * Builds everything.
 */
export const build = gulp.parallel(
    gulp.series(gulp.parallel(webpack))
)
