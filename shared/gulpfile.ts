import { generateNamespace } from '@gql2ts/from-schema'
import { DEFAULT_OPTIONS, DEFAULT_TYPE_MAP } from '@gql2ts/language-typescript'
import {buildSchema, graphql, IntrospectionQuery, introspectionQuery} from 'graphql';
import gulp from 'gulp'
import {readFile} from 'mz/fs'
import path from 'path';
import {resolveConfig} from 'prettier'

const GRAPHQL_SCHEMA_PATH = path.join(__dirname, '../cmd/frontend/graphqlbackend/schema.graphql')

export async function watchGraphQLTypes(): Promise<void> {
    await new Promise<never>((resolve, reject) => {
        gulp.watch(GRAPHQL_SCHEMA_PATH, graphQLTypes).on('error', reject)
    })
}

/** Generates the Typescript types for the GraphQL schema */
export async function graphQLTypes(): Promise<void> {
    const schemaStr = await readFile(GRAPHQL_SCHEMA_PATH, 'utf8')
    const schema = buildSchema(schemaStr)
    const result = (await graphql(schema, introspectionQuery)) as { data: IntrospectionQuery }

    const formatOptions = (await resolveConfig(__dirname, {config: __dirname + '/../prettier.config.js'}))
    const typings =
        'export type ID = string\n\n' +
        generateNamespace(
            '',
            result,
            {
                typeMap: {
                    ...DEFAULT_TYPE_MAP,
                    ID: 'ID',
                },
            },
            {

            }
        )
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
