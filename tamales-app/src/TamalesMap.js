import React from 'react'
import MapView from 'react-native-maps'
import * as Location from 'expo-location'

export default class TamalesMap extends React.Component {

    state = {
        cameraCenter: undefined,
    }

    async componentDidMount() {
        await Location.requestPermissionsAsync()
        await Location.watchPositionAsync({accuracy: 1}, this.setUserLocation)
    }

    setUserLocation = (userLocation) => {
        const {latitude, longitude, speed} = userLocation.coords
        console.log('speed on user location update ' + speed)
        this.setState({cameraCenter: {latitude, longitude}})
    }

    render() {
        if (!this.state.cameraCenter) {
            return null
        }
        const mapProps = {}
        mapProps.camera = {center: this.state.cameraCenter, zoom: 0, pitch: 8, heading: 0, altitude: 5000}
        mapProps.zoomControlEnabled = false
        mapProps.loadingEnabled = true
        mapProps.showsTraffic = false
        mapProps.showsIndoors = false
        mapProps.showsScale = false
        mapProps.showsPointsOfInterest = false
        mapProps.showsMyLocationButton = false
        mapProps.showsCompass = false
        return (
            <MapView style={{flex: 1}} camera={mapProps.camera}>
            </MapView>
        )
    }
}
