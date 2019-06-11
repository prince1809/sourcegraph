import * as React from "react";

export class App extends React.Component<{}, {}> {
    public render(): JSX.Element | null {
        return (
            <div>
                <h1 className="app__title">Sourcegraph management console</h1>
                <p className="app_subtitle">
                    View and edit critical Sourcegraph configuration. See{' '}
                    <a target="_blank" href="https://docs.sourcegraph.com/admin/management_console">
                        documentation
                    </a>{' '}
                    for more information.
                </p>
            </div>
        )
    }
}
