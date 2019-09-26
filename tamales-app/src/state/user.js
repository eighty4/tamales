import {filter, mergeMap} from 'rxjs/operators'
import {waitForAuthToken} from './auth'

const USER_LOCATION_UPDATED = 'USER_LOCATION_UPDATED'
const LOGIN_START = 'LOGIN_START'
const LOGIN_SUCCESS = 'LOGIN_SUCCESS'
const LOGIN_ERROR = 'LOGIN_ERROR'

export const startLogin = (email) => ({type: LOGIN_START, email})

export const reducers = (prevState = {location: null, authenticated: false, authenticating: false}, action) => {
    switch (action.type) {
        case USER_LOCATION_UPDATED:
            const {location} = action
            return {...prevState, location}
        case LOGIN_START:
            const {email} = action
            return {...prevState, authenticating: email}
        case LOGIN_SUCCESS:
            const {authenticated} = action
            return {...prevState, authenticated, authenticating: false}
        default:
            return prevState
    }
}

export const startLoginEpic = action$ => action$.pipe(
    filter(action => action.type === LOGIN_START),
    mergeMap(async action => {
        const {email} = action
        try {
            const authToken = await waitForAuthToken(email)
            return {type: LOGIN_SUCCESS, authenticated: {authToken, email}}
        } catch (error) {
            return {
                type: LOGIN_ERROR,
                error,
            }
        }
    }),
)

export const epics = [
    startLoginEpic,
]
