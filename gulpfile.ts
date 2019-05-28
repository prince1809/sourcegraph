import gulp from 'gulp'
import {graphQLTypes, schema, watchSchema} from './shared/gulpfile';
import {webpack as webWebpack} from './web/gulpfile';

/**
 * Generates files needed for builds
 */
export const generate = gulp.parallel(schema, graphQLTypes)

/**
 * Builds everything.
 */
export const build = gulp.parallel(gulp.series(generate, gulp.parallel(webWebpack)))

export {schema, graphQLTypes}

/**
 * Watches everything and rebuilds on file changes.
 */
export const watch = gulp.series(
    generate,
    gulp.parallel(watchSchema)
)
