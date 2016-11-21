import React, { PropTypes } from 'react'
import { connect } from 'react-redux'


const TopBar = () => (
    <ul>
        <ProfileLink
            profileLocation={ this.props.link }
            name={this.props.user.name} />
        <ProjectDropdown projects={ this.props.projects } />
        <SettingsDropdown />
    </ul>
)

const mapStateToProps = state => {
    return {
        user: state.user,
        link: './' + state.user,
        projects: state.projects
    }
}

export default connect(mapStateToProps)(TopBar)


const ProfileLink = ({ profileLocation, name }) => (
    <li>
        <a href={ profileLocation }>{ name }</a>
    </li>
)

ProfileLink.propTypes = {
    profileLocation: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired
}


const ProjectDropdown = ({ projects }) => (
    <li>
        <ul className='projects-dropdown'>
            { projects.map((project, i) =>
                <li><a href={ project.link }>{ project.name }</a></li>
            )}
        </ul>
    </li>
)

ProjectsDropdown.propTypes = {
    projects: PropTypes.arrayOf(PropTypes.shape({
        link: PropTypes.string.isRequired,
        name: PropTypes.string.isRequired,
    })).isRequired
}


const SettingsDropdown = () => (
    <li>
        settings
    </li>
)
