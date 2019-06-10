interface PageError {
    statusCode: number
    statusText: string
    error: string
    errorID: string
}

interface Window {
    pageError?: PageError
    context: SourcegraphContext
    MonacoEnvironment: {
        getWorkerUrl(moduleId: string, label: string): string
    }
}

interface ImmutableUser {
    readonly UID: number
}

type DeployType = 'cluster' | 'docker-container' | 'dev'

interface SourcegraphContext {
    xhrHeaders: { [key: string]: string }
    csrfToken: string
    userAgentIsBot: boolean

    readonly isAuthenticatedUser: boolean

    readonly sentryDSN: string | null

    externalURL: string

    assetsRoot: string

    version: string

    debug: boolean

    sourcegraphDotComMode: boolean
}
