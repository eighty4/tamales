import React from 'react'
import {StyleSheet, TouchableOpacity} from 'react-native'
import PropTypes from 'prop-types'

const styles = StyleSheet.create({
    container: {
        borderRadius: 75/2,
        backgroundColor: 'red',
        width: 75,
        height: 75,
        top: 50,
        right: 25,
        position: 'absolute',
        opacity: .5,
        zIndex: 1,
    },
})

class UserOverlay extends React.Component {

    onPress = () => {
        if (this.props.authenticated || this.props.authenticating) {
            return
        }
        this.props.startLogin('adam.mckee84@gmail.com')
    }

    render() {

        const containerStyles = this.props.authenticated ? {...styles.container, backgroundColor: 'green'} : this.props.authenticating ? {...styles.container, backgroundColor: 'orange'} : styles.container
        return (
            <TouchableOpacity style={containerStyles} onPress={this.onPress}>

            </TouchableOpacity>
        )
    }
}

UserOverlay.propTypes = {
    startLogin: PropTypes.func.isRequired,
    authenticated: PropTypes.bool.isRequired,
    authenticating: PropTypes.bool.isRequired,
}

export default UserOverlay
