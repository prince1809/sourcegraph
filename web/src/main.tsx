// This is the entry point for the web app

// Order is important here
// Don't remove empty lines between these imports

import React from 'react'
import {render} from 'react-dom'
import {SourcegraphWebApp} from './SourcegraphWebApp';

window.addEventListener('DOMContentLoaded', () => {
    render(
        <SourcegraphWebApp/>,
        document.querySelector('#root')
    )
})
