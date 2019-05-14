import gulp from 'gulp'
import {readFile} from 'mz/fs'


export async function watchGraphQLTypes() {
    await new Promise<never>((resolve, reject) => {
        gulp.watch()
    })
}

/** Generates the Typescript types for the GraphQL schema */
export async function graphQLTypes(): Promise<void> {
    const schemaStr = await readFile()
}

/**
 * Generate the Typescript types for the JSON schemas.
 */
export async function schema(): Promise<void> {

}

export async function watchSchema(): Promise<void> {
    await new Promise<never>((_resolve, reject) => {
        gulp.watch(__dirname + '/../schema/*.schema.json', schema).on('error', reject)
    })
}
