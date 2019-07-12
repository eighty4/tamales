
export const USER_LOCATION_UPDATED = 'USER_LOCATION_UPDATED'

export const reducers = (prevState = {location: null}, action) => {
    switch (action.type) {
        case USER_LOCATION_UPDATED:
            const {location} = action
            return {...prevState, location}
        default:
            return prevState
    }
}

export const epics = []
