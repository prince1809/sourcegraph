import gulp from 'gulp'


/**
 * Generates files needed for builds
 */
export const generate = gulp.parallel(schema, graphQLTypes)

