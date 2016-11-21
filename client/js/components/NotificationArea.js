import React, { PropTypes } from 'react'
import { connect } from 'react-redux'

const NotificationArea = () => (
    <div className={`notification-area ${this.props.level}`}>
        { this.props.msg }
    </div>
)

const mapStateToProps = state => {
    return {
        level: state.notification.level,
        msg: state.notification.msg
    }
}

export default connect(mapStateToProps)(NotificationArea)
