import React, { PropTypes } from 'react'
import { connect } from 'react-redux'

const Project = () => (
    <div className='project'>

    </div>
)

const mapStateToProps = state => {
    return {
        description: state.project.description
    }
}

export default connect(mapStateToProps)(Project)
