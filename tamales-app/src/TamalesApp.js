import React from 'react'
import {StyleSheet, View} from 'react-native'
import {Provider} from 'react-redux'
import store from './state/store'
import TamalesMap from './TamalesMap'
import UserOverlay from './UserOverlayContainer'

const styles = StyleSheet.create({
    container: {
        ...StyleSheet.absoluteFillObject,
    },
})

export default class TamalesApp extends React.Component {

    render() {
        return (
            <Provider store={store}>
                <View style={styles.container}>
                    <UserOverlay/>
                    <TamalesMap/>
                </View>
            </Provider>
        )
    }
}
