import {combineEpics} from 'redux-observable'
import {epics as userEpics} from './user'

export default combineEpics(...userEpics)
