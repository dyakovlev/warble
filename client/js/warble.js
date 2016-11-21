import React from 'react'
import render from 'react-dom'
import { Provider } from 'react-redux'
import { createStore } from 'redux'

import { TopBar, NotificationArea, Project } from './components/index'
import reducer from './reducers/index'

let store = createStore(reducer, {
    user: 'dyakovlev',
    projects: [
        { name: 'test project 1', link: 'test-project-1' },
        { name: 'test project 1', link: 'test-project-2' }
    ],
    activeProject: 0,
    notification: {level: null, msg: null},
    project: {
        description: 'long form text description of project'
    }
})

render(
    <Provider store={createStore(reducer)}>
        <div>
            <TopBar />
            <NotificationArea />
            <Project />
        </div>
    </Provider>
    , document.querySelector('.root')
)

