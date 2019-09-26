import UserOverlay from './UserOverlay'
import {connect} from 'react-redux'
import {startLogin} from './state/user'

const mapStateToProps = (state) => ({
    authenticated: !!state.authenticated,
    authenticating: !!state.authenticating,
})

const mapDispatchToProps = (dispatch) => ({
    startLogin: (email) => dispatch(startLogin(email)),
})

export default connect(mapStateToProps, mapDispatchToProps)(UserOverlay)
