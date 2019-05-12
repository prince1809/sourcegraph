import { generateNameSpace } from '@gql2ts/from-schema'
import gulp from 'gulp'

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
