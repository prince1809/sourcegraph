import * as React from 'react'
import { Redirect, RouteComponentProps } from 'react-router'
import { LayoutProps } from './Layout'
import { parseSearchURLQuery } from './search'

export interface LayoutRouteComponentProps extends RouteComponentProps<any>, LayoutProps {
}

export interface LayoutRouteProps {
    path: string
    exact?: boolean
    render: (props: LayoutRouteComponentProps) => React.ReactNode

    forceNarrowWidth?: boolean
}


export const routes: ReadonlyArray<LayoutRouteProps> = [
    {
        path: '/',
        render: (props: any) =>
            window.context.sourcegraphDotComMode && !props.user ? (
                <Redirect to="/welcome"/>
            ) : (
                <Redirect to="/search"/>
            ),
        exact: true,
    },
    {
        path: '/search',
        render: (props: any) =>
            parseSearchURLQuery(props.location.search)
    }
]
